package discord

import (
	"errors"
	"evedata/appContext"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/garyburd/redigo/redis"
)

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

	rate := time.Second * 60
	throttle := time.Tick(rate)

	for {

		r := ctx.Cache.Get()
		defer r.Close()

		// Skip this entity if we have touched it recently
		startID, err := redis.Int64(r.Do("GET", "EVEDATA_killqueue:99006652"))
		if err != nil {
			startID = 0
		}

		// [BENCHMARK] 0.016 sec / 0.000 sec
		rows, err := ctx.Db.Query(`
			SELECT K.id FROM  killmails K 
            INNER JOIN killmailAttackers A ON K.id = A.id
            INNER JOIN mapSolarSystems M ON K.solarSystemID = M.solarSystemID
            WHERE
            (
				A.allianceID IN (
					SELECT defenderID AS id FROM wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = 99006652
					UNION
					SELECT aggressorID AS id FROM wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND defenderID = 99006652
					UNION
					SELECT aggressorID  AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND allyID = 99006652
					UNION
					SELECT allyID AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = 99006652
				) OR
				A.corporationID IN (
					SELECT defenderID AS id FROM wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = 99006652
					UNION
					SELECT aggressorID AS id FROM wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND defenderID = 99006652
					UNION
					SELECT aggressorID  AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND allyID = 99006652
					UNION
					SELECT allyID AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = 99006652
				)
			) AND
            K.id > ? AND 
            killTime > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 10 MINUTE) AND
            M.security > 0.5
            GROUP BY K.id
            ORDER BY K.id ASC`, startID)

		for rows.Next() {
			var killID int64
			if err := rows.Scan(&killID); err != nil {
				log.Printf("EMDRCrestBridge: ", err)
			}
			// Say we touched the entity and expire after one day
			r.Do("SET", "EVEDATA_killqueue:99006652", killID)
			SendMessage(fmt.Sprintf("https://zkillboard.com/kill/%d/", killID))
		}
		rows.Close()
		<-throttle
	}
}

// [TODO] Temporary Hack... test feasibility
func SendMessage(message string) error {
	if dg == nil {
		return errors.New("Not Connected")
	}
	_, err := dg.ChannelMessageSend("229342742399025154", message)
	return err
}
