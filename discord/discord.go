package discord

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/models"
	"github.com/bwmarrin/discordgo"
	"github.com/garyburd/redigo/redis"
	yaml "gopkg.in/yaml.v2"
)

// This is all one massive hack until we get microservices going for this.
// currently locked into stuff i need personally.

var dg *discordgo.Session

func GoDiscordBot(ctx *appContext.AppContext) {
	var err error
	log.Printf("DiscordBot: Starting \n")
	dg, err = discordgo.New("Bot " + ctx.Conf.Discord.Token)
	if err != nil {
		log.Fatal("DiscordBot: Failure ", err)
	}

	err = dg.Open()
	if err != nil {
		log.Fatal("DiscordBot: Open Socket ", err)
	}
	go goKillmailHunter(ctx)
}

// Excuse the mess.. this is a temporary test to determine it
// we are capable of handling this traffic. Interface will be
// developed around it later.
func goKillmailHunter(ctx *appContext.AppContext) {
	rate := time.Second * 60 * 3
	throttle := time.Tick(rate)

	for {
		<-throttle
		r := ctx.Cache.Get()
		defer r.Close()

		checkNotifications(ctx)
		// Skip this entity if we have touched it recently
		startID, err := redis.Int64(r.Do("GET", "EVEDATA_killqueue:99002974"))
		if err != nil {
			startID = 0
		}

		// [BENCHMARK] 0.016 sec / 0.000 sec
		rows, err := ctx.Db.Query(`
			SELECT K.id FROM  evedata.killmails K 
            INNER JOIN evedata.killmailAttackers A ON K.id = A.id
            INNER JOIN mapSolarSystems M ON K.solarSystemID = M.solarSystemID
            WHERE
            (
				A.allianceID IN (
					SELECT defenderID AS id FROM evedata.wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = 99002974
					UNION
					SELECT aggressorID AS id FROM evedata.wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND defenderID = 99002974
					UNION
					SELECT aggressorID  AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND allyID = 99002974
					UNION
					SELECT allyID AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = 99002974
				) OR
				A.corporationID IN (
					SELECT defenderID AS id FROM evedata.wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = 99002974
					UNION
					SELECT aggressorID AS id FROM evedata.wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND defenderID = 99002974
					UNION
					SELECT aggressorID  AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND allyID = 99002974
					UNION
					SELECT allyID AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = 99002974
				)
			) AND
            K.id > ? AND 
            killTime > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 MINUTE) AND
            ROUND(M.security, 1) >= 0.5
            GROUP BY K.id
            ORDER BY K.id ASC`, startID)
		if err != nil {
			continue
		}

		defer rows.Close()

		for rows.Next() {
			var killID int64
			if err := rows.Scan(&killID); err != nil {
				log.Printf("Discord: %v", err)
			}
			// Say we touched the entity and expire after one day
			r.Do("SET", "EVEDATA_killqueue:99002974", killID)
			sendKillMessage(fmt.Sprintf("https://zkillboard.com/kill/%d/", killID))
		}
	}
}

type AllWarDeclaredMsg struct {
	AgainstID    int64   `yaml:"againstID"`
	Cost         float64 `yaml:"cost"`
	DeclaredByID int64   `yaml:"declaredByID"`
	DelayHours   int64   `yaml:"delayHours"`
	HostileState int64   `yaml:"hostileState"`
}

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

func checkNotifications(ctx *appContext.AppContext) error {
	r := ctx.Cache.Get()
	defer r.Close()
	startID, err := redis.Int64(r.Do("GET", "EVEDATA_notificationqueue:99002974"))
	if err != nil {
		startID = 0
	}

	// [BENCHMARK] 0.016 sec / 0.000 sec
	rows, err := ctx.Db.Query(`
		SELECT notificationID, type, text FROM evedata.notifications
		WHERE type IN ('TowerAlertMsg', 'StructureUnderAttack', 'OrbitalReinforced', 'OrbitalAttacked', 'CorpWarDeclaredMsg', 'AllWarDeclaredMsg')
		AND timestamp > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 1440 MINUTE) AND notificationCharacterID IN (1962167517,94135910) AND notificationID > ? ORDER BY notificationID ASC
		`, startID)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			notificationID         int64
			notificationType, text string
		)
		if err := rows.Scan(&notificationID, &notificationType, &text); err != nil {
			log.Printf("Discord: %v", err)
		}
		fmt.Println(startID)

		switch notificationType {

		case "AllWarDeclaredMsg", "CorpWarDeclaredMsg":
			l := AllWarDeclaredMsg{}
			err = yaml.Unmarshal([]byte(text), &l)

			defender, _ := models.GetEntityName(l.AgainstID)
			attacker, _ := models.GetEntityName(l.DeclaredByID)

			fmt.Printf("@everyone [%s](https://www.evedata.org/%s?id=%d) just declared war on [%s](https://www.evedata.org/%s?id=%d)\n", attacker.Name,
				attacker.EntityType, l.DeclaredByID, defender.Name, defender.EntityType, l.AgainstID)

			sendNotificationMessage(fmt.Sprintf("@everyone [%s](https://www.evedata.org/%s?id=%d) just declared war on [%s](https://www.evedata.org/%s?id=%d)\n", attacker.Name,
				attacker.EntityType, l.DeclaredByID, defender.Name, defender.EntityType, l.AgainstID))
		case "StructureUnderAttack", "OrbitalAttacked", "TowerAlertMsg":
			l := OrbitalAttacked{}
			err = yaml.Unmarshal([]byte(text), &l)

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

			locationName, _ := models.GetCelestialName(location)
			systemName, _ := models.GetCelestialName(l.SolarSystemID)
			structureType, _ := models.GetTypeName(l.TypeID)
			attackerName, _ := models.GetEntityName(attacker)
			fmt.Printf("@everyone %s is under attack at %s in %s by [%s](https://www.evedata.org/%s?id=%d) S: %.1f%%  A: %.1f%%  H: %.1f%% \n",
				structureType, locationName, systemName, attackerName.Name, attackerType, attacker, l.ShieldLevel*100, l.ArmorValue*100, l.HullValue*100)

			sendNotificationMessage(fmt.Sprintf("@everyone %s is under attack at %s in %s by [%s](https://www.evedata.org/%s?id=%d) S: %.1f%%  A: %.1f%%  H: %.1f%% \n",
				structureType, locationName, systemName, attackerName.Name, attackerType, attacker, l.ShieldLevel*100, l.ArmorValue*100, l.HullValue*100))
		}
		r.Do("SET", "EVEDATA_notificationqueue:99002974", notificationID)
	}
	return nil
}

// [TODO] Temporary Hack... test feasibility
func sendKillMessage(message string) error {
	if dg == nil {
		return errors.New("Not Connected")
	}
	_, err := dg.ChannelMessageSend("369208842443292675", message)
	return err
}

// [TODO] Temporary Hack... test feasibility
func sendNotificationMessage(message string) error {
	if dg == nil {
		return errors.New("Not Connected")
	}
	_, err := dg.ChannelMessageSend("369620236019695616", message)
	return err
}
