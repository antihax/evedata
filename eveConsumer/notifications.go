package eveConsumer

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
	yaml "gopkg.in/yaml.v2"
)

func init() {
	addConsumer("notifications", notificationsConsumer, "EVEDATA_notificationQueue")
	addTrigger("notifications", notificationsTrigger)
}

type Locator struct {
	AgentLocation struct {
		Region        int `yaml:"3"`
		Constellation int `yaml:"4"`
		SolarSystem   int `yaml:"5"`
		Station       int `yaml:"15"`
	} `yaml:"agentLocation"`
	TargetLocation struct {
		Region        int `yaml:"3"`
		Constellation int `yaml:"4"`
		SolarSystem   int `yaml:"5"`
		Station       int `yaml:"15"`
	} `yaml:"targetLocation"`
	CharacterID  int `yaml:"characterID"`
	MessageIndex int `yaml:"messageIndex"`
}

// notificationsConsumer handles gathering notifications from API
func notificationsConsumer(c *EVEConsumer, redisPtr *redis.Conn) (bool, error) {
	// Dereference the redis pointer.
	r := *redisPtr

	// POP some work of the queue
	ret, err := r.Do("SPOP", "EVEDATA_notificationQueue")
	if err != nil {
		return false, err
	} else if ret == nil {
		return false, nil
	}

	// Get the notification string
	v, err := redis.String(ret, err)
	if err != nil {
		return false, err
	}

	// Split our char:tokenChar string
	dest := strings.Split(v, ":")

	// Quick sanity check
	if len(dest) != 2 {
		return false, errors.New("Invalid notification string")
	}

	char, err := strconv.ParseInt(dest[0], 10, 64)
	if err != nil {
		return false, err
	}
	tokenChar, err := strconv.ParseInt(dest[1], 10, 64)
	if err != nil {
		return false, err
	}

	// Get the OAuth2 Token from the database.
	token, err := c.ctx.TokenStore.GetTokenSource(char, tokenChar)
	if err != nil {
		return false, err
	}

	// Put the token into a context for the API client
	auth := context.WithValue(context.TODO(), goesi.ContextOAuth2, token)

	notifications, res, err := c.ctx.ESI.ESI.CharacterApi.GetCharactersCharacterIdNotifications(auth, (int32)(tokenChar), nil)
	if err != nil {
		tokenError(char, tokenChar, res, err)
		return false, err
	}
	tokenSuccess(char, tokenChar, 200, "OK")

	tx, err := models.Begin()
	if err != nil {
		return false, err
	}

	// Skip if we have no notifications
	if len(notifications) != 0 {
		done := false
		var locatorValues, allValues []string

		// Dump all locators into the DB.
		for _, n := range notifications {
			if n.Type_ == "LocateCharMsg" {
				done = true
				l := Locator{}
				err = yaml.Unmarshal([]byte(n.Text), &l)
				if err != nil {
					return false, err
				}
				locatorValues = append(locatorValues, fmt.Sprintf("(%d,%d,%d,%d,%d,%d,%d,'%s')",
					n.NotificationId, char, l.TargetLocation.SolarSystem, l.TargetLocation.Constellation,
					l.TargetLocation.Region, l.TargetLocation.Station, l.CharacterID, n.Timestamp.Format(models.SQLTimeFormat)))
			}
			allValues = append(allValues, fmt.Sprintf("(%d,%d,%d,%d,'%s','%s','%s','%s')",
				n.NotificationId, char, tokenChar, n.SenderId, n.SenderType,
				n.Timestamp.Format(models.SQLTimeFormat), n.Type_, models.Escape(n.Text)))
		}

		if done {
			stmt := fmt.Sprintf(`INSERT INTO evedata.locatedCharacters
										(notificationID, characterID, solarSystemID, constellationID, 
											regionID, stationID, locatedCharacterID, time)
					VALUES %s ON DUPLICATE KEY UPDATE characterID = characterID;`, strings.Join(locatorValues, ",\n"))

			_, err = tx.Exec(stmt)
			if err != nil {
				return false, err
			}
		}

		stmt := fmt.Sprintf(`INSERT INTO evedata.notifications
			(notificationID,characterID,notificationCharacterID,senderID,senderType,timestamp,type,text)
			VALUES %s ON DUPLICATE KEY UPDATE characterID = characterID;`, strings.Join(allValues, ",\n"))

		_, err = tx.Exec(stmt)
		if err != nil {
			tx.Rollback()
			return false, err
		}

		// Update our cacheUntil flag
		tx.Exec(`UPDATE evedata.crestTokens SET notificationCacheUntil = ? 
						WHERE characterID = ? AND tokenCharacterID = ?`,
			goesi.CacheExpires(res), char, tokenChar)

		// Retry the transaction if we get deadlocks
		err = models.RetryTransaction(tx)
		if err != nil {
			return false, err
		}
	}

	return true, err
}

// Update character notifications
func notificationsTrigger(c *EVEConsumer) (bool, error) {

	// Gather characters for update. Group for optimized updating.
	rows, err := c.ctx.Db.Query(
		`SELECT characterID, tokenCharacterID FROM evedata.crestTokens WHERE 
		notificationCacheUntil < UTC_TIMESTAMP() AND lastStatus NOT LIKE "%400 Bad Request%" AND 
		scopes LIKE "%esi-characters.read_notifications.v1%";`)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	r := c.ctx.Cache.Get()
	defer r.Close()

	// Loop updatable characters
	for rows.Next() {
		var (
			char      int64 // Source char
			tokenChar int64 // Token Char
		)

		err = rows.Scan(&char, &tokenChar)
		if err != nil {
			return false, err
		}

		// Add the job to the queue
		_, err = r.Do("SADD", "EVEDATA_notificationQueue", fmt.Sprintf("%d:%d", char, tokenChar))
		if err != nil {
			return false, err
		}
	}
	return true, err
}
