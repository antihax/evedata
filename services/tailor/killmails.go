package tailor

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/antihax/eve-axiom/dogma"
	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/goesi/esi"

	nsq "github.com/nsqio/go-nsq"
)

type killmailAttributes struct {
	Attributes *dogma.Attributes
	KillmailID int32
}

var chanKillmailAttributes chan killmailAttributes

func init() {
	chanKillmailAttributes = make(chan killmailAttributes, 100)
}

func (s *Tailor) killmailHandler(message *nsq.Message) error {
	killmail := datapackages.Killmail{}
	if err := gobcoder.GobDecoder(message.Body, &killmail); err != nil {
		log.Println(err)
		return err
	}

	attr, err := getAttributesForKillmail(&killmail.Kill)
	if err != nil {
		log.Println(err)
		return err
	}

	// Add the package to the list
	chanKillmailAttributes <- killmailAttributes{
		Attributes: attr,
		KillmailID: killmail.Kill.KillmailId,
	}
	return nil
}

func (s *Tailor) killmailConsumer() {
	for {
		a := <-chanKillmailAttributes
		b := a.Attributes

		sq := squirrel.Insert("evedata.killmailAttributes").Columns(
			"id", "eHP", "DPS", "Alpha", "scanResolution", "signatureRadiusNoMWD",
			"signatureRadius", "agility", "warpSpeed", "speedNoMWD", "speed", "remoteArmorRepair",
			"remoteShieldRepair", "remoteEnergyTransfer", "energyNeutralization", "sensorStrength",
			"RPS", "CPURemaining", "powerRemaining", "capacitorNoMWD", "capacitor",
			"capacitorTimeNoMWD", "capacitorTime",
		)
		sq = sq.Values(
			a.KillmailID, b.AvgEHP, b.TotalDPS, b.TotalAlpha, b.ScanResolution, b.WithoutMWD.SignatureRadius,
			b.WithMWD.SignatureRadius, b.Agility, b.WarpSpeed, b.WithoutMWD.MaxVelocity, b.WithMWD.MaxVelocity, b.RemoteArmorRepairPerSecond,
			b.RemoteShieldTransferPerSecond, b.RemoteEnergyTransferPerSecond, b.EnergyNeutralizerPerSecond,
			b.GravStrenth+b.RadarStrength+b.LadarStrength+b.MagStrength,
			b.AvgRPS, b.CPURemaining, b.PGRemaining, b.WithoutMWD.Capacitor.Fraction, b.WithMWD.Capacitor.Fraction,
			b.WithoutMWD.Capacitor.Duration.Nanoseconds(), b.WithMWD.Capacitor.Duration.Nanoseconds(),
		)
		sqlq, args, err := sq.ToSql()
		if err != nil {
			log.Println(err)
			continue
		}

		err = s.doSQL(sqlq+" ON DUPLICATE KEY UPDATE id = id", args...)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func getAttributesForKillmail(km *esi.GetKillmailsKillmailIdKillmailHashOk) (*dogma.Attributes, error) {
	j, err := json.Marshal(km)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "http://axiom.evedata:3005/killmail", bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	attributes := &dogma.Attributes{}
	if err = dec.Decode(attributes); err != nil {
		return nil, err
	}

	return attributes, nil
}
