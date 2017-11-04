package artifice

import (
	"context"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/antihax/evedata/internal/redisqueue"
)

func init() {
	registerTrigger("wars", warsTrigger, time.NewTicker(time.Second*3600))
}

func warsTrigger(s *Artifice) error {
	maxWarID := int32(math.MaxInt32)
	cycle := 0
	for {
		wars, _, err := s.esi.ESI.WarsApi.GetWars(context.TODO(), map[string]interface{}{"maxWarId": maxWarID})
		if err != nil {
			return err
		}

		wars64 := make([]int64, len(wars))
		for i := range wars64 {
			if maxWarID > wars[i] {
				maxWarID = wars[i]
			}
			wars64[i] = int64(wars[i])
		}

		work := []redisqueue.Work{}

		known, err := s.inQueue.CheckWorkCompletedInBulk("evedata_war_finished", wars64)
		if err != nil {
			return err
		}

		for i := range known {
			if !known[i] {
				work = append(work, redisqueue.Work{Operation: "war", Parameter: wars[i]})
				warChan <- wars[i]
			}
		}

		s.QueueWork(work)

		if maxWarID < 100 {
			return nil
		}

		if cycle > 300 {
			return nil
		}

		cycle++
	}
}

var warChan chan int32

// zkillboardPost posts killmails to zkillboard from zkillChan
func (s *Artifice) warKillmails() {
	// Create the channel for feeding kills
	warChan = make(chan int32, 1000000)

	for {
		// pop a war off the channel
		war := <-warChan
		err := getWarKills(s, war)
		if err != nil {
			log.Println(err)
		}
	}
}

func getWarKills(s *Artifice, id int32) error {
	page := int32(1)
	for {
		kills, r, err := s.esi.ESI.WarsApi.GetWarsWarIdKillmails(context.TODO(), id, map[string]interface{}{"page": int32(page)})
		if err != nil {
			log.Println(err)
			return err
		}

		kills64 := make([]int64, len(kills))
		for i := range kills64 {
			kills64[i] = int64(kills[i].KillmailId)
		}

		known, err := s.inQueue.CheckWorkCompletedInBulk("evedata_known_kills", kills64)
		if err != nil {
			log.Println(err)
			return err
		}

		work := []redisqueue.Work{}
		for i := range known {
			if !known[i] {
				work = append(work, redisqueue.Work{Operation: "killmail", Parameter: []interface{}{kills[i].KillmailHash, kills[i].KillmailId}})

				// Send to zkillboard
				zkillChan <- killmail{ID: kills[i].KillmailId, Hash: kills[i].KillmailHash}
			}
		}

		s.QueueWork(work)

		xpagesS := r.Header.Get("X-Pages")
		xpages, _ := strconv.Atoi(xpagesS)
		if int32(xpages) == page || len(kills) == 0 {
			return nil
		}
		page++
	}
}
