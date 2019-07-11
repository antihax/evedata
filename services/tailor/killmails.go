package tailor

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/antihax/eve-axiom/attributes"
	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/goesi/esi"
	"github.com/prometheus/client_golang/prometheus"

	nsq "github.com/nsqio/go-nsq"
)

// KillmailAttributes provides a package of killmail, attributes, and id->name lookup
type KillmailAttributes struct {
	Attributes *attributes.Attributes                    `json:"attributes"`
	Killmail   *esi.GetKillmailsKillmailIdKillmailHashOk `json:"killmail"`
	NameMap    map[int32]string                          `json:"nameMap"`
	PriceMap   map[int32]float64                         `json:"priceMap"`
	SystemInfo *SystemInformation                        `json:"systemInfo"`
	DNA        string                                    `json:"dna"`
}

// SystemInformation for the killmail
type SystemInformation struct {
	CelestialID     int32   `db:"celestialID" json:"celestialID"`
	CelestialName   string  `db:"celestialName" json:"celestialName"`
	RegionID        int32   `db:"regionID" json:"regionID"`
	RegionName      string  `db:"regionName" json:"regionName"`
	SolarSystemName string  `db:"solarSystemName" json:"solarSystemName"`
	SolarSystemID   int32   `db:"solarSystemID" json:"solarSystemID"`
	Security        float32 `db:"security" json:"security"`
}

var chanKillmailAttributes chan KillmailAttributes

func init() {
	chanKillmailAttributes = make(chan KillmailAttributes, 1000)
}

// saveKillmail saves the data to backblaze b2 in json.gz format
func (s *Tailor) saveKillmail(pack *KillmailAttributes) error {
	b, err := json.Marshal(pack)
	if err != nil {
		log.Println(err)
		return err
	}

	var gzb bytes.Buffer
	gz, err := gzip.NewWriterLevel(&gzb, gzip.BestCompression)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err := gz.Write(b); err != nil {
		log.Println(err)
		return err
	}
	if err := gz.Flush(); err != nil {
		log.Println(err)
		return err
	}
	if err := gz.Close(); err != nil {
		log.Println(err)
		return err
	}
	if len(gzb.Bytes()) > 0 {
		metadata := make(map[string]string)
		_, err = s.bucket.UploadFile(
			fmt.Sprintf("%d.json.gz", pack.Killmail.KillmailId),
			metadata,
			bytes.NewReader(gzb.Bytes()),
		)
		if err != nil {
			log.Println(err)

			metricStatus.With(
				prometheus.Labels{"status": "fail"},
			).Inc()
			return err
		}
	}
	metricStatus.With(
		prometheus.Labels{"status": "ok"},
	).Inc()
	return nil
}

func timeStage(s string, t time.Time) {
	metricStageTime.With(
		prometheus.Labels{"stage": s},
	).Observe(float64(time.Since(t).Nanoseconds() / 1000000))
}

func (s *Tailor) killmailHandler(message *nsq.Message) error {
	killmail := datapackages.Killmail{}
	start := time.Now()
	if err := gobcoder.GobDecoder(message.Body, &killmail); err != nil {
		log.Println(err)
		return err
	}
	timeStage("decode", start)
	start = time.Now()

	attr, err := getAttributesForKillmail(&killmail.Kill)
	if err != nil {
		log.Println(err)
		return err
	}
	timeStage("axiom", start)
	start = time.Now()

	names, err := s.resolveNames(&killmail.Kill)
	if err != nil {
		log.Println(err, " ", killmail.Kill.KillmailId)
		return err
	}
	timeStage("names", start)
	start = time.Now()

	prices, err := s.getPrices(&killmail.Kill)
	if err != nil {
		log.Println(err, " ", killmail.Kill.KillmailId)
		return err
	}
	timeStage("prices", start)
	start = time.Now()

	pos := killmail.Kill.Victim.Position
	sysinfo, err := s.getSystemInformation(killmail.Kill.SolarSystemId, pos.X, pos.Y, pos.Z)
	if err != nil {
		log.Println(err, ": missing system:", killmail.Kill.SolarSystemId)
		return err
	}
	timeStage("system", start)
	start = time.Now()

	dna, err := s.getDNA(killmail.Kill.Victim.ShipTypeId)
	if err != nil {
		log.Println(err)
		dna = ""
	}
	timeStage("dna", start)
	start = time.Now()

	pack := KillmailAttributes{
		Attributes: attr,
		Killmail:   &killmail.Kill,
		NameMap:    names,
		PriceMap:   prices,
		SystemInfo: sysinfo,
		DNA:        dna,
	}
	// Add the package to the list
	chanKillmailAttributes <- pack

	err = s.saveKillmail(&pack)
	if err != nil {
		return err
	}
	timeStage("save", start)
	start = time.Now()

	return nil
}

func joinInt32(a []int32) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), ","), "[]")
}

func (s *Tailor) getDNA(typeID int32) (string, error) {
	var dna string
	err := s.db.QueryRow(
		`SELECT concat(sofHullName, ":", sofFactionName, ":", sofRaceName) 
		FROM eve.invTypes T
		INNER JOIN eve.eveGraphics G ON T.graphicID = G.graphicID
		WHERE typeID = ?;`, typeID).Scan(&dna)

	return dna, err
}

func (s *Tailor) getFallbackSystemInformation(system int32) (*SystemInformation, error) {
	sys := &SystemInformation{CelestialName: "Unknown", CelestialID: 0}
	err := s.db.QueryRow(
		`SELECT  S.solarSystemID, S.regionID, solarSystemName, regionName, security
		FROM  eve.mapSolarSystems S 
		INNER JOIN eve.mapRegions R ON R.regionID = S.regionID
		WHERE solarSystemID = ?;`, system).Scan(
		&sys.SolarSystemID, &sys.RegionID, &sys.SolarSystemName, &sys.RegionName, &sys.Security)
	if err != nil {
		return nil, err
	}
	return sys, err
}

func (s *Tailor) getSystemInformation(system int32, x, y, z float64) (*SystemInformation, error) {
	sys := &SystemInformation{}
	err := s.db.QueryRow(
		`SELECT itemName AS celestialName, itemID AS celestialID, S.solarSystemID, S.regionID, solarSystemName, regionName, S.security
		FROM  eve.mapDenormalize D
		INNER JOIN eve.mapSolarSystems S ON S.solarSystemID = D.solarSystemID
		INNER JOIN eve.mapRegions R ON R.regionID = D.regionID
		WHERE itemID = closestCelestial(?,?,?,?);`, system, x, y, z).Scan(
		&sys.CelestialName, &sys.CelestialID, &sys.SolarSystemID, &sys.RegionID, &sys.SolarSystemName, &sys.RegionName, &sys.Security)
	if err != nil {
		return s.getFallbackSystemInformation(system)
	}
	return sys, err
}

func (s *Tailor) getPrices(kill *esi.GetKillmailsKillmailIdKillmailHashOk) (map[int32]float64, error) {

	// Make a list of lost items
	idList := []int32{}
	idList = append(idList, kill.Victim.ShipTypeId)
	for _, a := range kill.Victim.Items {
		idList = append(idList, a.ItemTypeId)
		for _, a := range a.Items {
			idList = append(idList, a.ItemTypeId)
		}
	}

	// Lookup the information
	rows, err := s.db.Query(`
		SELECT typeID, mean
		FROM evedata.typePricesMonthly 
		WHERE typeID IN (`+joinInt32(idList)+`) AND 
			month = MONTH(DATE_SUB(?, INTERVAL 28 DAY)) AND
			year = YEAR(DATE_SUB(?, INTERVAL 28 DAY));
		`, kill.KillmailTime, kill.KillmailTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Make a map of the lost items
	ids := make(map[int32]float64)
	for rows.Next() {
		var (
			id    int32
			price float64
		)
		err := rows.Scan(&id, &price)
		if err != nil {
			return nil, err
		}
		ids[id] = price
	}

	return ids, nil
}

func (s *Tailor) resolveNames(kill *esi.GetKillmailsKillmailIdKillmailHashOk) (map[int32]string, error) {
	// make a unique list of ids
	ids := make(map[int32]string)
	for _, a := range kill.Attackers {
		ids[a.AllianceId] = ""
		ids[a.CorporationId] = ""
		ids[a.CharacterId] = ""
		ids[a.ShipTypeId] = ""
		ids[a.WeaponTypeId] = ""
		ids[a.FactionId] = ""
	}
	for _, a := range kill.Victim.Items {
		ids[a.ItemTypeId] = ""
		for _, a := range a.Items {
			ids[a.ItemTypeId] = ""
		}
	}

	ids[kill.Victim.AllianceId] = ""
	ids[kill.Victim.CorporationId] = ""
	ids[kill.Victim.CharacterId] = ""
	ids[kill.Victim.FactionId] = ""
	ids[kill.Victim.ShipTypeId] = ""

	// Delete 0 if we picked up no alliance
	delete(ids, 0)

	idList := []int32{}
	for id := range ids {
		if id > 0 {
			idList = append(idList, id)
		}
	}

	// Lookup the information
	rows, err := s.db.Query(`
		SELECT allianceID as id, name FROM evedata.alliances FORCE INDEX(PRIMARY) WHERE allianceID IN (` + joinInt32(idList) + `)
		UNION
		SELECT corporationID as id, name FROM evedata.corporations WHERE corporationID IN (` + joinInt32(idList) + `)
		UNION
		SELECT characterID as id, name FROM evedata.characters WHERE characterID IN (` + joinInt32(idList) + `)
		UNION
		SELECT typeID as id, typeName as name FROM eve.invTypes WHERE typeID IN (` + joinInt32(idList) + `)
		UNION
		SELECT itemID as ID, itemName as name FROM eve.eveNames WHERE itemID IN (` + joinInt32(idList) + `)
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   int32
			name string
		)
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		ids[id] = name
	}

	// Double check we are not missing any
	for id, n := range ids {
		if id > 0 {
			if n == "" {
				return nil, fmt.Errorf("cannot find name for %d", id)
			}
		}
	}

	return ids, nil
}

// killmailConsumer receives killmails from NSQ and dumps attributes to SQL
func (s *Tailor) killmailConsumer() {
	for {
		a := <-chanKillmailAttributes
		b := a.Attributes.Ship
		pm := a.PriceMap

		// figure out the mean security of the attackers.
		meanSecurity := float32(0.0)
		count := 0
		for _, attacker := range a.Killmail.Attackers {
			if attacker.CharacterId > 0 {
				meanSecurity += attacker.SecurityStatus
				count++
			}
		}
		if count > 0 {
			meanSecurity = meanSecurity / float32(count)
		}

		value := pm[a.Attributes.TypeID]
		totalValue := pm[a.Attributes.TypeID]
		if value > 0 {
			for _, e := range a.Attributes.Modules {
				value += pm[int32(e["typeID"])]
				if pm[int32(e["chargeID"])] > 0 {
					value += pm[int32(e["chargeID"])]
				}
			}
			for _, e := range a.Attributes.Drones {
				value += pm[int32(e["typeID"])] * e["quantity"]
			}
		}

		for _, e := range a.Killmail.Victim.Items {
			for _, e := range e.Items {
				totalValue += pm[e.ItemTypeId] * float64(e.QuantityDestroyed+e.QuantityDropped)
			}
			totalValue += pm[e.ItemTypeId] * float64(e.QuantityDestroyed+e.QuantityDropped)
		}

		sq := squirrel.Insert("evedata.killmailAttributes").Columns(
			"id", "eHP", "DPS", "Alpha", "scanResolution", "signatureRadiusNoMWD",
			"signatureRadius", "agility", "warpSpeed", "speedNoMWD", "speed", "remoteArmorRepair",
			"remoteShieldRepair", "remoteEnergyTransfer", "energyNeutralization", "sensorStrength",
			"RPS", "CPURemaining", "powerRemaining", "capacitorNoMWD", "capacitor",
			"capacitorTimeNoMWD", "capacitorTime", "fittedValue", "totalValue", "stasisWebifierStrength", "totalWarpScrambleStrength",
			"meanSecurity",
		)

		sq = sq.Values(
			a.Killmail.KillmailId, b["avgEHP"], b["totalDPS"], b["totalAlphaDamage"], b["scanResolution"], b["signatureRadius"],
			b["signatureRadiusMWD"], b["agility"], b["warpSpeedMultiplier"], b["maxVelocity"], b["maxVelocityMWD"], b["remoteArmorRepairPerSecond"],
			b["remoteShieldBonusAmountPerSecond"], b["remotePowerTransferAmountPerSecond"], b["energyNeutralizerAmountPerSecond"],
			b["scanRadarStrength"]+b["scanLadarStrength"]+b["scanMagnetometricStrength"]+b["scanGravimetricStrength"],
			b["avgRPS"], b["cpuRemaining"], b["powerRemaining"], b["capacitorFraction"], b["capacitorFractionMWD"],
			b["capacitorDuration"]*100000, b["capacitorDurationMWD"]*100000, value, totalValue, b["stasisWebifierStrength"], b["totalWarpScrambleStrength"],
			meanSecurity,
		)

		sqlq, args, err := sq.ToSql()
		if err != nil {
			log.Println(err)
			continue
		}

		err = s.doSQL(sqlq+" ON DUPLICATE KEY UPDATE id = id, speed = VALUES(speed), meanSecurity = VALUES(meanSecurity)", args...)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func getAttributesForKillmail(km *esi.GetKillmailsKillmailIdKillmailHashOk) (*attributes.Attributes, error) {
	j, err := json.Marshal(km)
	if err != nil {
		return nil, err
	}
	ctx, cncl := context.WithTimeout(context.Background(), time.Second*30)
	defer cncl()
	req, err := http.NewRequest("POST", "http://axiom.evedata.svc.cluster.local:3005/killmail", bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Printf("exceeded %d", km.KillmailId)
		return nil, err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	attributes := &attributes.Attributes{}
	if err = dec.Decode(attributes); err != nil {
		return nil, err
	}

	return attributes, nil
}
