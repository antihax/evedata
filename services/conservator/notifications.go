package conservator

import (
	"fmt"
	"log"
	"time"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
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

	// Only process contracts.
	if notifications.TokenCharacterID == 1962167517 || notifications.TokenCharacterID == 94135910 {
		for _, n := range notifications.Notifications {
			// Skip the notification if if is more than three hours old

			if n.Timestamp.Before(time.Now().Add(-time.Hour * 6)) {
				continue
			}

			if s.outQueue.CheckWorkCompleted(fmt.Sprintf("evedata-bot-notification-sent:%d", 99002974), n.NotificationId) {
				continue
			}

			err := s.checkNotification(n.Type_, n.Text, n.Timestamp)
			if err != nil {
				continue
			}

			s.outQueue.SetWorkCompleted(fmt.Sprintf("evedata-bot-notification-sent:%d", 99002974), n.NotificationId)
		}
	}
	return nil
}

// AllWarDeclaredMsg message
type AllWarDeclaredMsg struct {
	AgainstID    int64   `yaml:"againstID"`
	Cost         float64 `yaml:"cost"`
	DeclaredByID int64   `yaml:"declaredByID"`
	DelayHours   int64   `yaml:"delayHours"`
	HostileState int64   `yaml:"hostileState"`
}

// OrbitalAttacked message
type OrbitalAttacked struct {
	AggressorAllianceID int64   `yaml:"aggressorAllianceID"`
	AggressorCorpID     int64   `yaml:"aggressorCorpID"`
	PlanetID            int64   `yaml:"planetID"`
	MoonID              int64   `yaml:"moonID"`
	ShieldLevel         float64 `yaml:"shieldLevel"`
	ArmorValue          float64 `yaml:"armorValue"`
	HullValue           float64 `yaml:"hullValue"`
	TypeID              int64   `yaml:"typeID"`
	SolarSystemID       int64   `yaml:"solarSystemID"`
}

// OrbitalReinforced message
type OrbitalReinforced struct {
	AggressorAllianceID int64 `yaml:"aggressorAllianceID"`
	AggressorCorpID     int64 `yaml:"aggressorCorpID"`
	PlanetID            int64 `yaml:"planetID"`
	MoonID              int64 `yaml:"moonID"`
	ReinforceExitTime   int64 `yaml:"reinforceExitTime"`
	TypeID              int64 `yaml:"typeID"`
	SolarSystemID       int64 `yaml:"solarSystemID"`
}

func (s *Conservator) checkNotification(notificationType, text string, timestamp time.Time) error {

	switch notificationType {

	case "AllWarDeclaredMsg", "CorpWarDeclaredMsg":
		l := AllWarDeclaredMsg{}
		err := yaml.Unmarshal([]byte(text), &l)
		if err != nil {
			return err
		}

		defender, err := s.getEntityName(l.AgainstID)
		if err != nil {
			return err
		}
		attacker, err := s.getEntityName(l.DeclaredByID)
		if err != nil {
			return err
		}

		sendNotificationMessage(fmt.Sprintf("@everyone [%s] [%s](https://www.evedata.org/%s?id=%d) just declared war on [%s](https://www.evedata.org/%s?id=%d)\n",
			timestamp.UTC().String(), attacker.Name, attacker.EntityType, l.DeclaredByID, defender.Name, defender.EntityType, l.AgainstID))
	case "StructureUnderAttack", "OrbitalAttacked", "TowerAlertMsg":
		l := OrbitalAttacked{}
		yaml.Unmarshal([]byte(text), &l)

		location := int64(0)
		if l.MoonID > 0 {
			location = l.MoonID
		} else if l.PlanetID > 0 {
			location = l.PlanetID
		}

		attacker := int64(0)
		attackerType := ""
		if l.AggressorAllianceID > 0 {
			attacker = l.AggressorAllianceID
			attackerType = "alliance"
		} else if l.AggressorCorpID > 0 {
			attacker = l.AggressorCorpID
			attackerType = "corporation"
		}

		locationName, err := s.getCelestialName(location)
		if err != nil {
			return err
		}
		systemName, err := s.getCelestialName(l.SolarSystemID)
		if err != nil {
			return err
		}
		structureType, err := s.getTypeName(l.TypeID)
		if err != nil {
			return err
		}
		attackerName, err := s.getEntityName(attacker)
		if err != nil {
			return err
		}

		return sendNotificationMessage(fmt.Sprintf("@everyone [%s] %s is under attack at %s in %s by [%s](https://www.evedata.org/%s?id=%d) S: %.1f%%  A: %.1f%%  H: %.1f%% \n",
			timestamp.UTC().String(), structureType, locationName, systemName, attackerName.Name, attackerType, attacker, l.ShieldLevel*100, l.ArmorValue*100, l.HullValue*100))

	case "OrbitalReinforced":
		l := OrbitalReinforced{}
		yaml.Unmarshal([]byte(text), &l)

		location := int64(0)
		if l.MoonID > 0 {
			location = l.MoonID
		} else if l.PlanetID > 0 {
			location = l.PlanetID
		}

		attacker := int64(0)
		attackerType := ""
		if l.AggressorAllianceID > 0 {
			attacker = l.AggressorAllianceID
			attackerType = "alliance"
		} else if l.AggressorCorpID > 0 {
			attacker = l.AggressorCorpID
			attackerType = "corporation"
		}

		locationName, err := s.getCelestialName(location)
		if err != nil {
			return err
		}
		systemName, err := s.getCelestialName(l.SolarSystemID)
		if err != nil {
			return err
		}
		structureType, err := s.getTypeName(l.TypeID)
		if err != nil {
			return err
		}
		attackerName, err := s.getEntityName(attacker)
		if err != nil {
			return err
		}

		return sendNotificationMessage(fmt.Sprintf("@everyone [%s] %s was reinforced at %s in %s by [%s](https://www.evedata.org/%s?id=%d). Timer expires at %s\n",
			timestamp.UTC().String(), structureType, locationName, systemName, attackerName.Name, attackerType, attacker,
			time.Unix(datapackages.WintoUnixTimestamp(l.ReinforceExitTime), 0).String()))
	}
	return nil
}

type EntityName struct {
	Name       string `db:"name" json:"name"`
	EntityType string `db:"type" json:"type"`
}

// Obtain entity name and type by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func (s *Conservator) getEntityName(id int64) (*EntityName, error) {
	ref := EntityName{}
	if err := s.db.QueryRowx(`
		SELECT name, 'corporation' AS type FROM evedata.corporations WHERE corporationID = ?
		UNION
		SELECT name, 'alliance' AS type FROM evedata.alliances WHERE allianceID = ?
		LIMIT 1`, id, id).StructScan(&ref); err != nil {
		return nil, err
	}
	return &ref, nil
}

// Obtain type name.
// [BENCHMARK] 0.000 sec / 0.000 sec
func (s *Conservator) getTypeName(id int64) (string, error) {
	ref := ""
	if err := s.db.QueryRowx(`
		SELECT typeName FROM invTypes WHERE typeID = ?
		LIMIT 1`, id).Scan(&ref); err != nil {
		return "", err
	}
	return ref, nil
}

// Obtain SolarSystem name.
// [BENCHMARK] 0.000 sec / 0.000 sec
func (s *Conservator) getSystemName(id int64) (string, error) {
	ref := ""
	if err := s.db.QueryRowx(`
		SELECT solarSystemName FROM mapSolarSystems WHERE solarSystemID = ?
		LIMIT 1`, id).Scan(&ref); err != nil {
		return "", err
	}
	return ref, nil
}

// Obtain Celestial name.
// [BENCHMARK] 0.000 sec / 0.000 sec
func (s *Conservator) getCelestialName(id int64) (string, error) {
	ref := ""
	if err := s.db.QueryRowx(`
		SELECT itemName FROM mapDenormalize WHERE itemID = ?
		LIMIT 1`, id).Scan(&ref); err != nil {
		return "", err
	}
	return ref, nil
}
