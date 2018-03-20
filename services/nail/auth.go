package nail

import (
	"fmt"
	"log"
	"strings"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("corporationContacts", spawnCorporationContactsConsumer)
	AddHandler("allianceContacts", spawnAllianceContactsConsumer)
}

func spawnCorporationContactsConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.corporationContactsHandler)))
}

func spawnAllianceContactsConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.allianceContactsHandler)))
}

func (s *Nail) corporationContactsHandler(message *nsq.Message) error {
	c := datapackages.CorporationContacts{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	defer tx.Rollback()

	var values []string
	for _, e := range c.Contacts {
		values = append(values, fmt.Sprintf("(%d,%d,%f)", c.CorporationID, e.ContactId, e.Standing))
	}

	_, err = tx.Exec("DELETE FROM evedata.entityContacts WHERE entityID = ?", c.CorporationID)
	if err != nil {
		log.Println(err)
		return err
	}

	stmt := fmt.Sprintf(`INSERT INTO evedata.entityContacts
		(entityID,contactID,standing)
		VALUES %s
		ON DUPLICATE KEY UPDATE standing=VALUES(standing)
		`, strings.Join(values, ",\n"))

	_, err = tx.Exec(stmt)

	if err != nil {
		log.Println(err)
		return err

	}

	return retryTransaction(tx)
}

func (s *Nail) allianceContactsHandler(message *nsq.Message) error {
	c := datapackages.AllianceContacts{}
	err := gobcoder.GobDecoder(message.Body, &c)
	if err != nil {
		log.Println(err)
		return err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	defer tx.Rollback()

	var values []string
	for _, e := range c.Contacts {
		values = append(values, fmt.Sprintf("(%d,%d,%f)", c.AllianceID, e.ContactId, e.Standing))
	}

	_, err = tx.Exec("DELETE FROM evedata.entityContacts WHERE entityID = ?", c.AllianceID)
	if err != nil {
		log.Println(err)
		return err
	}

	stmt := fmt.Sprintf(`INSERT INTO evedata.entityContacts
		(entityID,contactID,standing)
		VALUES %s
		ON DUPLICATE KEY UPDATE standing=VALUES(standing)
		`, strings.Join(values, ",\n"))

	_, err = tx.Exec(stmt)

	if err != nil {
		log.Println(err)
		return err

	}

	return retryTransaction(tx)
}
