package nail

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/goesi"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("character", spawnCharacterConsumer)
	AddHandler("alliance", spawnAllianceConsumer)
	AddHandler("corporation", spawnCorporationConsumer)

}

func spawnCharacterConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.characterHandler)))
}

func spawnAllianceConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.allianceHandler)))
}

func spawnCorporationConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.corporationHandler)))
}

func (s *Nail) corporationHandler(message *nsq.Message) error {
	c := datapackages.Corporation{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	cacheUntil := time.Now().UTC().Add(time.Hour * 24 * 7)

	return s.DoSQL(`INSERT INTO evedata.corporations
		(corporationID,name,ticker,ceoID,allianceID,factionID,memberCount,updated,cacheUntil)
		VALUES(?,?,?,?,?,?,?,UTC_TIMESTAMP(),?) 
		ON DUPLICATE KEY UPDATE 
		ceoID=VALUES(ceoID), name=VALUES(name), ticker=VALUES(ticker), allianceID=VALUES(allianceID), 
		factionID=VALUES(factionID), memberCount=VALUES(memberCount),  
		updated=UTC_TIMESTAMP(), cacheUntil=VALUES(cacheUntil)
	`, c.CorporationID, c.Corporation.CorporationName, c.Corporation.Ticker, c.Corporation.CeoId, c.Corporation.AllianceId, goesi.FactionNameToID(c.Corporation.Faction), c.Corporation.MemberCount, cacheUntil)

}

func (s *Nail) allianceHandler(message *nsq.Message) error {
	c := datapackages.Alliance{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	cacheUntil := time.Now().UTC().Add(time.Hour * 24 * 7)

	return s.DoSQL(`
		INSERT INTO evedata.alliances 
			(
				allianceID,
				name,
				shortName,
				executorCorpID,
				startDate,
				corporationsCount,
				updated,
				cacheUntil
			)
			VALUES(?,?,?,?,?,?,UTC_TIMESTAMP(),?) 
			ON DUPLICATE KEY UPDATE 
				executorCorpID = VALUES(executorCorpID),
				corporationsCount = VALUES(corporationsCount), 
				updated = UTC_TIMESTAMP(), 
				cacheUntil=VALUES(cacheUntil)
	`, c.AllianceID, c.Alliance.AllianceName, c.Alliance.Ticker, c.Alliance.ExecutorCorp, c.Alliance.DateFounded.UTC().Format("2006-01-02 15:04:05"), len(c.AllianceCorporations), cacheUntil)
}

func (s *Nail) characterHandler(message *nsq.Message) error {
	c := datapackages.Character{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	cacheUntil := time.Now().UTC().Add(time.Hour * 24 * 7)

	err = s.DoSQL(`
		INSERT INTO evedata.characters (characterID,name,bloodlineID,ancestryID,corporationID,allianceID,race,gender,securityStatus,updated,cacheUntil)
			VALUES(?,?,?,?,?,?,evedata.raceByID(?),?,?,UTC_TIMESTAMP(),?) 
			ON DUPLICATE KEY UPDATE 
			corporationID=VALUES(corporationID), gender=VALUES(gender), allianceID=VALUES(allianceID), securityStatus=VALUES(securityStatus), updated = UTC_TIMESTAMP(), cacheUntil=VALUES(cacheUntil)
	`, c.CharacterID, c.Character.Name, c.Character.BloodlineId, c.Character.AncestryId, c.Character.CorporationId, c.Character.AllianceId, c.Character.RaceId, c.Character.Gender, c.Character.SecurityStatus, cacheUntil)
	if err != nil {
		log.Println(err)
		return err
	}

	var values []string
	for _, e := range c.CorporationHistory {

		values = append(values, fmt.Sprintf("(%d,'%s',%d,%d)",
			c.CharacterID, e.StartDate.UTC().Format("2006-01-02 15:04:05"), e.RecordId, e.CorporationId))
	}

	stmt := fmt.Sprintf(`INSERT INTO evedata.corporationHistory (characterID,startDate,recordID,corporationID)
			VALUES %s
			ON DUPLICATE KEY UPDATE 
			startDate=VALUES(startDate);
				`, strings.Join(values, ",\n"))

	return s.DoSQL(stmt)
}
