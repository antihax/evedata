package discordbottemp

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/goesi/esi"
	nsq "github.com/nsqio/go-nsq"
)

type atWarWith struct {
	ID     int32     `db:"id"`
	Start  time.Time `db:"timeStarted"`
	Finish time.Time `db:"timeFinished"`
}

var warsMap sync.Map
var highsecSystems map[int32]bool

func init() {
	addHandler("killmail", spawnKillmailConsumer)
	highsecSystems = make(map[int32]bool)
}

func spawnKillmailConsumer(s *DiscordBot, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.killmailHandler)))
}

func (s *DiscordBot) killmailHandler(message *nsq.Message) error {
	mail := esi.GetKillmailsKillmailIdKillmailHashOk{}
	err := gobcoder.GobDecoder(message.Body, &mail)
	if err != nil {
		log.Println(err)
		return err
	}

	// Skip killmails more than an hour old
	if mail.KillmailTime.Before(time.Now().UTC().Add(-time.Hour * 6)) {
		return nil
	}

	var entity atWarWith
	for _, a := range mail.Attackers {
		if i, ok := warsMap.Load(a.AllianceId); ok {
			v := i.(atWarWith)
			if v.Start.Before(time.Now().UTC()) && (v.Finish.IsZero() || v.Finish.After(time.Now().UTC())) {
				entity = v
			}
		} else if i, ok := warsMap.Load(a.CorporationId); ok {
			v := i.(atWarWith)
			if v.Start.Before(time.Now().UTC()) && (v.Finish.IsZero() || v.Finish.After(time.Now().UTC())) {
				entity = v
			}
		}
	}

	// didn't match
	if entity.ID == 0 {
		return nil
	}

	// Don't duplicate notifications
	if !s.outQueue.CheckWorkCompleted(fmt.Sprintf("evedata-bot-killmail-sent:%d", 99002974), int64(mail.KillmailId)) {
		if highsecSystems[mail.SolarSystemId] { // is it in highsec?
			err = sendKillMessage(fmt.Sprintf("https://zkillboard.com/kill/%d/", mail.KillmailId))
			if err != nil {
				return err
			}

			s.outQueue.SetWorkCompleted(fmt.Sprintf("evedata-bot-killmail-sent:%d", 99002974), int64(mail.KillmailId))
		}
	}
	return nil
}

func (s *DiscordBot) getSystems() error {
	var systems []int32
	err := s.db.Select(&systems, "SELECT solarSystemID FROM mapSolarSystems WHERE round(security, 1) > 0.4")
	if err != nil {
		return err
	}
	for _, s := range systems {
		highsecSystems[s] = true
	}
	return nil
}

func (s *DiscordBot) updateWars() {
	throttle := time.Tick(time.Second * 120)
	for {
		warlist := []atWarWith{}
		err := s.db.Select(&warlist, "CALL evedata.atWarWith(99002974);")
		if err != nil {
			log.Println(err)
		}

		for _, war := range warlist {
			warsMap.Store(war.ID, war)
		}
		<-throttle
	}
}
