package zkillboard

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/antihax/evedata/internal/redisqueue"
)

// Collect killmails from RedisQ (zkillboard live feed)
func (s *ZKillboard) redisQ() error {
	// ZKillboard redisq json format
	type redisqkill struct {
		Package struct {
			KillID int32
			ZKB    struct {
				Hash string
			}
		}
	}

	k := redisqkill{}
	err := s.getJSON(fmt.Sprintf("https://redisq.zkillboard.com/listen.php?queueID=croakroach"), &k)
	if err != nil {
		return err
	}
	if k.Package.KillID > 0 {
		if !s.outQueue.CheckWorkCompleted("evedata_known_kills", int64(k.Package.KillID)) {
			err = s.outQueue.QueueWork(
				[]redisqueue.Work{
					{Operation: "killmail", Parameter: []interface{}{k.Package.ZKB.Hash, k.Package.KillID}}},
				redisqueue.Priority_High,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Go Routine to collect killmails from ZKill API.
// Loops collecting one year of kill mails.
func (s *ZKillboard) apiConsumer() error {
	// Start from where we left off.
	nextCheck, _ := time.Parse("20060102", "20071201")

	rate := time.Second / 2
	throttle := time.Tick(rate)

	for {
		<-throttle

		// Move to the next day
		date := nextCheck.Format("20060102")
		nextCheck = nextCheck.Add(time.Hour * 24)

		// If we are at today, restart from 90 days
		if nextCheck.Sub(time.Now().UTC()) > 0 {
			nextCheck = time.Now().UTC().Add(time.Hour * 24 * -365)
			log.Printf("Restart zKill Consumer to %s", nextCheck.String())
		}

		// Get the kill history from ZKill for this day.
		k := make(map[string]interface{})
		err := s.getJSON(fmt.Sprintf("https://zkillboard.com/api/history/%s/", date), &k)
		if err != nil {
			log.Println(err)
			continue
		}

		kills := []redisqueue.Work{}
		// Loop through the killmails
		for idS, hash := range k {
			id, err := strconv.ParseInt(idS, 10, 64)
			if err != nil {
				log.Println(err)
				continue
			}

			// Add to the killmail queue
			kills = append(kills, redisqueue.Work{Operation: "killmail", Parameter: []interface{}{hash.(string), (int32)(id)}})
			if err != nil {
				log.Println(err)
				continue
			}
		}

		err = s.outQueue.QueueWork(kills, redisqueue.Priority_Low)
		if err != nil {
			log.Println(err)
		}
	}
}
