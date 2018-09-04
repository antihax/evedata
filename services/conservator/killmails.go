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
	addHandler("killmail", func(s *Conservator, consumer *nsq.Consumer) {
		consumer.AddConcurrentHandlers(s.wait(nsq.HandlerFunc(s.killmailHandler)), 20)
	})
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
		service, err := s.getService(channel.IntegrationID)
		if err != nil {
			return true
		}
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
		} else if channel.Options.Killmail.SendAllAbyssalT4 && isAbyssalCruiser(mail.Victim.ShipTypeId) {
			for _, item := range mail.Victim.Items {
				if isAbyssalT4TypeID(item.ItemTypeId) {
					sendMail = true
				}
			}
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

func isAbyssalT4TypeID(typeID int32) bool {
	types := []int32{
		47891,
		47895,
		47899,
		47903,
		47907,
		47890,
		47894,
		47898,
		47902,
		47906,
		47906,
		47701,
		47731,
		47735,
		47739,
		47741,
		47744,
		47748,
		47752,
		47756,
		47768,
		47772,
		47776,
		47780,
		47784,
		47788,
		47792,
		47799,
		47803,
		47807,
		47811,
		47815,
		47819,
		47823,
		47827,
		47831,
		47830,
		47826,
		47822,
		47818,
		47814,
		47810,
		47806,
		47802,
		47798,
		47791,
		47787,
		47783,
		47779,
		47775,
		47771,
		47767,
		47755,
		47751,
		47747,
		47743,
		47738,
		47734,
		47730,
		47700,
	}

	for _, id := range types {
		if id == typeID {
			return true
		}
	}
	return false
}

func isAbyssalCruiser(typeID int32) bool {
	types := []int32{
		2836,
		3518,
		11993,
		11999,
		12003,
		12005,
		12011,
		12015,
		12019,
		12023,
		32209,
		34477,
		34479,
		620,
		621,
		622,
		623,
		624,
		11959,
		11961,
		11971,
		20125,
		625,
		626,
		627,
		628,
		11957,
		11963,
		11965,
		11969,
		33395,
		33675,
		44995,
		45531,
		629,
		630,
		631,
		632,
		633,
		634,
		635,
		1904,
		2006,
		11011,
		17634,
		17709,
		17713,
		17715,
		17718,
		17720,
		17722,
		17843,
		17922,
		25560,
		29336,
		29337,
		29340,
		29344,
		33470,
		33553,
		33639,
		33641,
		33643,
		33645,
		33647,
		33649,
		33651,
		33653,
		33818,
		34445,
		34475,
		34590,
		47270,
	}

	for _, id := range types {
		if id == typeID {
			return true
		}
	}
	return false
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
