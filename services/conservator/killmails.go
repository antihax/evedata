package conservator

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/evedata/internal/datapackages"

	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/goesi/esi"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	addHandler("killmail", spawnKillmailConsumer)
}

func spawnKillmailConsumer(s *Conservator, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.killmailHandler)))
}

type atWarWith struct {
	ID     int32     `db:"id"`
	Start  time.Time `db:"timeStarted"`
	Finish time.Time `db:"timeFinished"`
}

func (s *Conservator) atWarWithKillmail(entityID int32, mail *esi.GetKillmailsKillmailIdKillmailHashOk) bool {
	var entity atWarWith
	for _, a := range mail.Attackers {
		_, ok := s.warsMap[entityID]
		if ok {
			if i, ok := s.warsMap[entityID].Load(a.AllianceId); ok {
				v := i.(atWarWith)
				if v.Start.Before(time.Now().UTC()) && (v.Finish.IsZero() || v.Finish.After(time.Now().UTC())) {
					entity = v
					break
				}
			} else if i, ok := s.warsMap[entityID].Load(a.CorporationId); ok {
				v := i.(atWarWith)
				if v.Start.Before(time.Now().UTC()) && (v.Finish.IsZero() || v.Finish.After(time.Now().UTC())) {
					entity = v
					break
				}
			}
		}
	}
	return entity.ID != 0
}

func (s *Conservator) reportKillmail(mail *esi.GetKillmailsKillmailIdKillmailHashOk) error {
	s.channels.Range(func(ki, vi interface{}) bool {
		channel := vi.(Channel)

		if !inSlice("kill", strings.Split(channel.Services, ",")) {
			return true // Don't have killmails on this channel
		}

		// Don't duplicate
		if s.outQueue.CheckWorkCompleted(fmt.Sprintf("evedata-bot-killmail-sent:%s", channel.ChannelID), int64(mail.KillmailId)) {
			return true
		}

		// Get the service
		si, ok := s.services.Load(channel.BotServiceID)
		if !ok {
			log.Printf("Missing Bot ID %d\n", channel.BotServiceID)
			return true
		}
		service := si.(Service)

		// filters
		if channel.Options.Killmail.IgnoreWorthless && isWorthlessTypeID(mail.Victim.ShipTypeId) {
			return true
		}
		if channel.Options.Killmail.IgnoreHighSec && s.solarSystems[mail.SolarSystemId] >= 0.5 {
			return true
		}
		if channel.Options.Killmail.IgnoreLowSec &&
			s.solarSystems[mail.SolarSystemId] > 0.0 && s.solarSystems[mail.SolarSystemId] < 0.5 {
			return true
		}
		if channel.Options.Killmail.IgnoreNullSec && s.solarSystems[mail.SolarSystemId] <= 0.0 {
			return true
		}

		// Determine if we send the mail
		sendMail := false
		if channel.Options.Killmail.War && s.atWarWithKillmail(service.EntityID, mail) {
			sendMail = true
		} else if channel.Options.Killmail.FactionWar && mail.Victim.FactionId > 0 {
			sendMail = true
		} else if channel.Options.Killmail.SendAll {
			sendMail = true
		}

		if sendMail {
			if err := service.Server.SendMessageToChannel(channel.ChannelID,
				fmt.Sprintf("https://zkillboard.com/kill/%d/", mail.KillmailId)); err != nil {
				log.Println(err)
			}
		}

		s.outQueue.SetWorkCompleted(fmt.Sprintf("evedata-bot-killmail-sent:%s", channel.ChannelID), int64(mail.KillmailId))
		return true
	})

	return nil
}

func (s *Conservator) killmailHandler(message *nsq.Message) error {
	killmail := datapackages.Killmail{}
	if err := gobcoder.GobDecoder(message.Body, &killmail); err != nil {
		log.Println(err)
		return err
	}

	mail := killmail.Kill
	// Skip killmails more than an hour old
	if mail.KillmailTime.Before(time.Now().UTC().Add(-time.Hour * 6)) {
		return nil
	}

	if err := s.reportKillmail(&mail); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func isWorthlessTypeID(typeID int32) bool {
	types := []int32{
		670,   // capsule
		33328, // capsule
		672,   // shuttle
		11129, // shuttle
		11132, // shuttle
		21097, // shuttle
		11134, // shuttle
		21628, // shuttle
		30842, // shuttle
		588,   // rookie
		596,   // rookie
		601,   // rookie
		606,   // rookie
	}

	for _, id := range types {
		if id == typeID {
			return true
		}
	}
	return false
}
