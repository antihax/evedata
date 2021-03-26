package nail

import (
	"log"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	nsq "github.com/nsqio/go-nsq"
)

type entityIDType struct {
	ID   int32
	Type string
}

var entitySQLQueue chan entityIDType

func init() {
	entitySQLQueue = make(chan entityIDType, 500)
	AddHandler("alliance", func(s *Nail, consumer *nsq.Consumer) {
		consumer.AddHandler(s.wait(nsq.HandlerFunc(s.allianceHandler)))
		go s.entitySQLPost()
	})
	AddHandler("corporation", func(s *Nail, consumer *nsq.Consumer) {
		consumer.AddHandler(s.wait(nsq.HandlerFunc(s.corporationHandler)))
	})

}

func (s *Nail) corporationHandler(message *nsq.Message) error {
	c := datapackages.Corporation{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	cacheUntil := time.Now().UTC().Add(time.Hour * 24 * 30)

	err = s.doSQL(`INSERT INTO evedata.corporations
		(corporationID,name,ticker,ceoID,allianceID,factionID,memberCount,updated,cacheUntil)
		VALUES(?,?,?,?,?,?,?,UTC_TIMESTAMP(),?) 
		ON DUPLICATE KEY UPDATE ceoID=VALUES(ceoID), name=VALUES(name), ticker=VALUES(ticker), allianceID=VALUES(allianceID), 
		factionID=VALUES(factionID), memberCount=VALUES(memberCount),  
		updated=UTC_TIMESTAMP(), cacheUntil=VALUES(cacheUntil)
	`, c.CorporationID, c.Corporation.Name, c.Corporation.Ticker, c.Corporation.CeoId, c.Corporation.AllianceId, c.Corporation.FactionId, c.Corporation.MemberCount, cacheUntil)
	if err != nil {
		log.Println(err)
		return err
	}

	return s.addEntity(c.CorporationID, "corporation")
}

func (s *Nail) allianceHandler(message *nsq.Message) error {
	c := datapackages.Alliance{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	cacheUntil := time.Now().UTC().Add(time.Hour * 24 * 31)

	err = s.doSQL(`
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
			ON DUPLICATE KEY UPDATE executorCorpID = VALUES(executorCorpID),
				corporationsCount = VALUES(corporationsCount), 
				updated = UTC_TIMESTAMP(), 
				cacheUntil=VALUES(cacheUntil)
	`, c.AllianceID, c.Alliance.Name, c.Alliance.Ticker, c.Alliance.ExecutorCorporationId, c.Alliance.DateFounded.UTC().Format("2006-01-02 15:04:05"), len(c.AllianceCorporations), cacheUntil)
	if err != nil {
		log.Println(err)
		return err
	}

	return s.addEntity(c.AllianceID, "alliance")
}

func (s *Nail) addEntity(id int32, entityType string) error {
	entitySQLQueue <- entityIDType{id, entityType}
	return nil
}

func (s *Nail) entitySQLPost() {
	for {
		count := 0
		sql := sq.Insert("evedata.entities").Columns("id", "type")
		for c := range entitySQLQueue {
			count++
			sql = sql.Values(c.ID, c.Type)
			if count > 80 {
				break
			}
		}
		if count == 0 {
			break
		}
		sqlq, args, err := sql.ToSql()
		if err != nil {
			log.Println(err)
		}
		err = s.doSQL(sqlq+` ON DUPLICATE KEY UPDATE id = id`, args...)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Nail) characterAuthOwnerHandler(message *nsq.Message) error {
	c := datapackages.CharacterRoles{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.doSQL(`
		UPDATE evedata.crestTokens SET roles = ?
		WHERE characterID = ? AND tokenCharacterID = ?
	`, strings.Join(c.Roles.Roles, ","), c.CharacterID, c.TokenCharacterID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
