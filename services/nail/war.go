package nail

import (
	"log"
	"time"

	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi/esi"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("war", spawnWarConsumer)
}

func spawnWarConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.warHandler)))
}

func (s *Nail) warHandler(message *nsq.Message) error {
	war := esi.GetWarsWarIdOk{}
	err := gobcoder.GobDecoder(message.Body, &war)
	if err != nil {
		log.Println(err)
		return err
	}

	// save the aggressor id
	var aggressor, defender int32
	if war.Aggressor.AllianceId > 0 {
		aggressor = war.Aggressor.AllianceId
	} else {
		aggressor = war.Aggressor.CorporationId
	}

	// save the defender id
	if war.Defender.AllianceId > 0 {
		defender = war.Defender.AllianceId
	} else {
		defender = war.Defender.CorporationId
	}

	err = s.DoSQL(`INSERT INTO evedata.wars
		(id, timeFinished,timeStarted,timeDeclared,openForAllies,cacheUntil,aggressorID,defenderID,mutual)
		VALUES(?,?,?,?,?,?,?,?,?)
		ON DUPLICATE KEY UPDATE 
			timeFinished=VALUES(timeFinished), 
			openForAllies=VALUES(openForAllies), 
			mutual=VALUES(mutual), 
			cacheUntil=VALUES(cacheUntil);`,
		war.Id, war.Finished.Format(models.SQLTimeFormat), war.Started, war.Declared,
		war.OpenForAllies, time.Now().UTC().Format(models.SQLTimeFormat), aggressor,
		defender, war.Mutual)
	if err != nil {
		return err
	}

	// Add information on allies in the war
	for _, a := range war.Allies {
		var ally int32
		if a.AllianceId > 0 {
			ally = a.AllianceId
		} else {
			ally = a.CorporationId
		}
		err = s.DoSQL(`INSERT INTO evedata.warAllies (id, allyID) VALUES(?,?) ON DUPLICATE KEY UPDATE id = id;`, war.Id, ally)
		if err != nil {
			return err
		}
	}

	return nil
}
