package nail

import (
	"fmt"
	"log"
	"strings"

	"github.com/antihax/goesi/notification"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/evedata/services/vanguard/models"
	nsq "github.com/nsqio/go-nsq"
	yaml "gopkg.in/yaml.v2"
)

func init() {
	AddHandler("characterNotifications", spawnCharacterNotificationsConsumer)
}

func spawnCharacterNotificationsConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.characterNotificationsHandler)))
}

func (s *Nail) characterNotificationsHandler(message *nsq.Message) error {
	notifications := datapackages.CharacterNotifications{}
	err := gobcoder.GobDecoder(message.Body, &notifications)
	if err != nil {
		log.Println(err)
		return err
	}
	if len(notifications.Notifications) == 0 {
		return nil
	}

	done := false
	var locatorValues, allValues []string

	// Dump all locators into the DB.
	for _, n := range notifications.Notifications {
		if n.Type_ == "LocateCharMsg" {
			l := notification.LocateCharMsg{}
			err = yaml.Unmarshal([]byte(n.Text), &l)
			if err == nil { // Ignore old locator responses
				done = true
				locatorValues = append(locatorValues, fmt.Sprintf("(%d,%d,%d,%d,%d,%d,%d,%q)",
					n.NotificationId, notifications.CharacterID, l.TargetLocation.SolarSystem, l.TargetLocation.Constellation,
					l.TargetLocation.Region, l.TargetLocation.Station, l.CharacterID, n.Timestamp.Format(models.SQLTimeFormat)))
			}
		}
		allValues = append(allValues, fmt.Sprintf("(%d,%d,%d,%d,%q,%q,%q,%q)",
			n.NotificationId, notifications.CharacterID, notifications.TokenCharacterID, n.SenderId, n.SenderType,
			n.Timestamp.Format(models.SQLTimeFormat), n.Type_, models.Escape(n.Text)))
	}

	if done {
		stmt := fmt.Sprintf(`INSERT INTO evedata.locatedCharacters
									(notificationID, characterID, solarSystemID, constellationID, 
										regionID, stationID, locatedCharacterID, time)
				VALUES %s ON DUPLICATE KEY UPDATE characterID = characterID;`, strings.Join(locatorValues, ",\n"))

		err = s.doSQL(stmt)
		if err != nil {
			return err
		}
	}

	stmt := fmt.Sprintf(`INSERT INTO evedata.notifications
		(notificationID,characterID,notificationCharacterID,senderID,senderType,timestamp,type,text)
		VALUES %s ON DUPLICATE KEY UPDATE characterID = characterID;`, strings.Join(allValues, ",\n"))

	err = s.doSQL(stmt)
	if err != nil {
		return err
	}
	return nil
}
