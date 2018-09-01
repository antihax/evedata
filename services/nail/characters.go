package nail

import (
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	nsq "github.com/nsqio/go-nsq"
)

var (
	characterSQLQueue chan datapackages.Character
)

func init() {
	characterSQLQueue = make(chan datapackages.Character, 500)
	AddHandler("character", func(s *Nail, consumer *nsq.Consumer) {
		consumer.AddConcurrentHandlers(s.wait(nsq.HandlerFunc(s.characterHandler)), 100)
		go s.characterSQLPost()
	})
}

func (s *Nail) characterSQLPost() {
	for {
		count := 0
		cacheUntil := time.Now().UTC().Add(time.Hour * 24 * 14)
		sql := sq.Insert("evedata.characters").Columns("characterID", "name", "bloodlineID", "ancestryID", "corporationID", "allianceID", "race", "gender", "securityStatus", "updated", "cacheUntil", "birthDate")
		for c := range characterSQLQueue {
			count++
			sql = sql.Values(c.CharacterID, c.Character.Name, c.Character.BloodlineId, c.Character.AncestryId, c.Character.CorporationId, c.Character.AllianceId, s.characterRaces[c.Character.RaceId], c.Character.Gender, c.Character.SecurityStatus, time.Now(), cacheUntil, c.Character.Birthday)
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
		err = s.doSQL(sqlq+` ON DUPLICATE KEY UPDATE  corporationID=VALUES(corporationID), gender=VALUES(gender), allianceID=VALUES(allianceID),birthDate=VALUES(birthDate),
				securityStatus=VALUES(securityStatus), updated = UTC_TIMESTAMP(), cacheUntil=VALUES(cacheUntil)`, args...)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Nail) characterHandler(message *nsq.Message) error {
	c := datapackages.Character{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	characterSQLQueue <- c

	return s.addEntity(c.CharacterID, "character")
}
