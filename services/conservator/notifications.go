package conservator

import (
	"fmt"
	"log"
	"time"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/goesi/notification"
	"github.com/bradfitz/slice"
	nsq "github.com/nsqio/go-nsq"
	yaml "gopkg.in/yaml.v2"
)

func init() {
	addHandler("characterNotifications", spawnCharacterNotificationsConsumer)
}

func spawnCharacterNotificationsConsumer(s *Conservator, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.characterNotificationsHandler)))
}

func (s *Conservator) characterNotificationsHandler(message *nsq.Message) error {
	notifications := datapackages.CharacterNotifications{}
	err := gobcoder.GobDecoder(message.Body, &notifications)
	if err != nil {
		log.Println(err)
		return err
	}

	if len(notifications.Notifications) == 0 {
		return nil
	}

	// sort by time
	slice.Sort(notifications.Notifications[:], func(i, j int) bool {
		return notifications.Notifications[i].Timestamp.Unix() < notifications.Notifications[j].Timestamp.Unix()
	})

	if len(notifications.Notifications) == 0 {
		log.Println("we broke the notifications")
		return nil
	}

	for _, n := range notifications.Notifications {
		// Skip the notification if if is old
		if n.Timestamp.Before(time.Now().UTC().Add(-time.Hour * 12)) {
			continue
		}

		err := s.checkNotification(notifications.TokenCharacterID, n.NotificationId, n.Type_, n.Text, n.Timestamp)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (s *Conservator) checkNotification(characterID int32, notificationID int64, notificationType, text string, timestamp time.Time) error {
	switch notificationType {

	case "LocateCharMsg":
		l := notification.LocateCharMsg{}
		if err := yaml.Unmarshal([]byte(text), &l); err != nil {
			log.Println(err)
			return err
		}
		systemName, err := s.getCelestialName(l.TargetLocation.SolarSystem)
		if err != nil {
			log.Println(err)
			return err
		}
		regionName, _ := s.getCelestialName(l.TargetLocation.Region)

		character, _ := s.getEntityName(l.CharacterID)

		stationName, _ := s.getStationName(l.TargetLocation.Station)

		message := fmt.Sprintf(" %s (https://www.evedata.org/character?id=%d) has been located in %s of %s",
			character.Name, l.CharacterID, systemName, regionName)

		if stationName != "" {
			message = fmt.Sprintf("%s docked at %s", message, stationName)
		}

		return s.sendNotificationMessage("locator", characterID, notificationID, message)

	case "AllWarDeclaredMsg", "CorpWarDeclaredMsg":
		l := notification.AllWarDeclaredMsg{}
		if err := yaml.Unmarshal([]byte(text), &l); err != nil {
			log.Println(err)
		}

		defender, err := s.getEntityName(l.AgainstID)
		if err != nil {
			log.Println(err)
		}
		attacker, err := s.getEntityName(l.DeclaredByID)
		if err != nil {
			log.Println(err)
		}

		message := fmt.Sprintf("[%s](https://www.evedata.org/%s?id=%d) just declared war on [%s](https://www.evedata.org/%s?id=%d)\n",
			attacker.Name, attacker.EntityType, l.DeclaredByID, defender.Name, defender.EntityType, l.AgainstID)

		return s.sendNotificationMessage("war", characterID, notificationID, message)

	case "CorpAppNewMsg":
		{
			l := notification.CorpAppNewMsg{}
			yaml.Unmarshal([]byte(text), &l)
			character, err := s.getEntityName(l.CharID)
			if err != nil {
				character.Name = "!! No fricking clue !!"
			}
			corporation, _ := s.getEntityName(l.CorpID)

			message := fmt.Sprintf("New corporation application from [%s](https://www.evedata.org/character?id=%d) to %s. Application Comment: %s\n",
				character.Name, l.CharID, corporation.Name, l.ApplicationText)

			return s.sendNotificationMessage("application", characterID, notificationID, message)
		}

	case "StructureUnderAttack":
		l := notification.StructureUnderAttack{}
		yaml.Unmarshal([]byte(text), &l)

		attacker := int32(0)
		attackerType := ""
		attackerName := ""
		if l.AllianceID > 0 {
			attacker = l.AllianceID
			attackerType = "alliance"
			attackerName = l.AllianceName
		} else if l.CharID > 0 { // There is no corporationID? WTF CCPls?
			attacker = l.CharID
			attackerType = "character"
			attackerName = l.CorpName
		}

		systemName, err := s.getCelestialName(l.SolarsystemID) // -.-
		if err != nil {
			log.Println(err)
		}

		message := fmt.Sprintf("@everyone structure is under attack in %s by [%s](https://www.evedata.org/%s?id=%d) S: %.1f%%  A: %.1f%%  H: %.1f%% \n",
			systemName, attackerName, attackerType, attacker, l.ShieldPercentage, l.ArmorPercentage, l.HullPercentage)

		return s.sendNotificationMessage("structure", characterID, notificationID, message)

	case "OrbitalAttacked":
		l := notification.OrbitalAttacked{}
		yaml.Unmarshal([]byte(text), &l)

		attacker := int32(0)
		attackerType := ""
		if l.AggressorAllianceID > 0 {
			attacker = l.AggressorAllianceID
			attackerType = "alliance"
		} else if l.AggressorCorpID > 0 {
			attacker = l.AggressorCorpID
			attackerType = "corporation"
		}

		locationName, err := s.getCelestialName(l.PlanetID)
		if err != nil {
			log.Println(err)
		}
		systemName, err := s.getCelestialName(l.SolarSystemID)
		if err != nil {
			log.Println(err)
		}
		structureType, err := s.getTypeName(l.TypeID)
		if err != nil {
			log.Println(err)
		}
		attackerName, err := s.getEntityName(attacker)
		if err != nil {
			log.Println(err)
		}

		message := fmt.Sprintf("@everyone %s is under attack at %s in %s by [%s](https://www.evedata.org/%s?id=%d) S: %.1f%%\n",
			structureType, locationName, systemName, attackerName.Name, attackerType, attacker, l.ShieldLevel*100)

		return s.sendNotificationMessage("structure", characterID, notificationID, message)

	case "TowerAlertMsg":
		l := notification.TowerAlertMsg{}
		yaml.Unmarshal([]byte(text), &l)

		attacker := int32(0)
		attackerType := ""
		if l.AggressorAllianceID > 0 {
			attacker = l.AggressorAllianceID
			attackerType = "alliance"
		} else if l.AggressorCorpID > 0 {
			attacker = l.AggressorCorpID
			attackerType = "corporation"
		}

		locationName, err := s.getCelestialName(l.MoonID)
		if err != nil {
			log.Println(err)
		}
		systemName, err := s.getCelestialName(l.SolarSystemID)
		if err != nil {
			log.Println(err)
		}
		structureType, err := s.getTypeName(l.TypeID)
		if err != nil {
			log.Println(err)
		}
		attackerName, err := s.getEntityName(attacker)
		if err != nil {
			log.Println(err)
		}
		message := fmt.Sprintf("@everyone %s is under attack at %s in %s by [%s](https://www.evedata.org/%s?id=%d) S: %.1f%%  A: %.1f%%  H: %.1f%% \n",
			structureType, locationName, systemName, attackerName.Name, attackerType, attacker, l.ShieldValue*100, l.ArmorValue*100, l.HullValue*100)

		return s.sendNotificationMessage("structure", characterID, notificationID, message)

	case "OrbitalReinforced":
		l := notification.OrbitalReinforced{}
		yaml.Unmarshal([]byte(text), &l)

		attacker := int32(0)
		attackerType := ""
		if l.AggressorAllianceID > 0 {
			attacker = l.AggressorAllianceID
			attackerType = "alliance"
		} else if l.AggressorCorpID > 0 {
			attacker = l.AggressorCorpID
			attackerType = "corporation"
		}

		locationName, err := s.getCelestialName(l.PlanetID)
		if err != nil {
			log.Println(err)
		}
		systemName, err := s.getCelestialName(l.SolarSystemID)
		if err != nil {
			log.Println(err)
		}
		structureType, err := s.getTypeName(l.TypeID)
		if err != nil {
			log.Println(err)
		}
		attackerName, err := s.getEntityName(attacker)
		if err != nil {
			log.Println(err)
		}

		message := fmt.Sprintf("@everyone %s was reinforced at %s in %s by [%s](https://www.evedata.org/%s?id=%d).\n\n Timer expires at %s\n",
			structureType, locationName, systemName, attackerName.Name, attackerType, attacker,
			time.Unix(datapackages.WintoUnixTimestamp(l.ReinforceExitTime), 0).UTC().String())

		return s.sendNotificationMessage("structure", characterID, notificationID, message)

	case "StructureLostShields", "StructureLostArmor":
		l := notification.StructureLostShields{}
		yaml.Unmarshal([]byte(text), &l)

		systemName, err := s.getCelestialName(l.SolarsystemID)
		if err != nil {
			log.Println(err)
		}

		structureType, err := s.getTypeName(l.StructureTypeID)
		if err != nil {
			log.Println(err)
		}

		message := fmt.Sprintf("@everyone %s entered %s in %s: Timer expires at %s\n",
			structureType, notificationType, systemName,
			time.Unix(datapackages.WintoUnixTimestamp(l.Timestamp), 0).UTC().String())

		return s.sendNotificationMessage("structure", characterID, notificationID, message)
	}
	return nil
}

func (s *Conservator) sendNotificationMessage(messageType string, characterID int32, notificationID int64, message string) error {
	shares, ok := s.notifications[messageType][characterID]
	if !ok {
		return nil
	}
	for _, share := range shares {
		channelData := unpackChannelData(share.Packed)
		for _, channel := range channelData {
			if inSlice(messageType, channel.Services) {
				if s.outQueue.CheckWorkCompleted(fmt.Sprintf("evedata-bot-notification-sent:%s", channel.ChannelID), notificationID) {
					continue
				}

				c, err := s.getChannel(channel.ChannelID)
				if err != nil {
					continue
				}

				// Get the service
				service, err := s.getService(c.IntegrationID)
				if err != nil {
					log.Println(err)
					return err
				}

				if err := service.Server.SendMessageToChannel(channel.ChannelID, message); err != nil {
					log.Println(err)
					return err
				}
				s.outQueue.SetWorkCompleted(fmt.Sprintf("evedata-bot-notification-sent:%s", channel.ChannelID), notificationID)
			}
		}
	}
	return nil
}

type EntityName struct {
	Name       string `db:"name" json:"name"`
	EntityType string `db:"type" json:"type"`
}

// Obtain entity name and type by ID.

func (s *Conservator) getEntityName(id int32) (*EntityName, error) {
	ref := EntityName{}
	if err := s.db.QueryRowx(`
		SELECT name, 'corporation' AS type FROM evedata.corporations WHERE corporationID = ?
		UNION
		SELECT name, 'alliance' AS type FROM evedata.alliances WHERE allianceID = ?
		UNION
		SELECT name, 'character' AS type FROM evedata.characters WHERE characterID = ?
		LIMIT 1`, id, id, id).StructScan(&ref); err != nil {
		return &ref, err
	}
	return &ref, nil
}

// Obtain type name.

func (s *Conservator) getTypeName(id int32) (string, error) {
	ref := ""
	if err := s.db.QueryRowx(`
		SELECT typeName FROM invTypes WHERE typeID = ?
		LIMIT 1`, id).Scan(&ref); err != nil {
		return "", err
	}
	return ref, nil
}

// Obtain SolarSystem name.

func (s *Conservator) getSystemName(id int32) (string, error) {
	ref := ""
	if err := s.db.QueryRowx(`
		SELECT solarSystemName FROM mapSolarSystems WHERE solarSystemID = ?
		LIMIT 1`, id).Scan(&ref); err != nil {
		return "", err
	}
	return ref, nil
}

// Obtain Celestial name.

func (s *Conservator) getCelestialName(id int32) (string, error) {
	ref := ""
	if err := s.db.QueryRowx(`
		SELECT itemName FROM mapDenormalize WHERE itemID = ?
		LIMIT 1`, id).Scan(&ref); err != nil {
		return "", err
	}
	return ref, nil
}

// Obtain Station name.

func (s *Conservator) getStationName(id int32) (string, error) {
	ref := ""
	if err := s.db.QueryRowx(`
		SELECT stationName FROM staStations WHERE stationID = ?
		LIMIT 1`, id).Scan(&ref); err != nil {
		return "", err
	}
	return ref, nil
}
