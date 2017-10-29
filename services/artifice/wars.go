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
	for {
		wars, _, err := s.esi.ESI.WarsApi.GetWars(context.TODO(), map[string]interface{}{"max_war_id": int32(maxWarID)})
		if err != nil {
			return err
		}

		work := []redisqueue.Work{}
		for _, war := range wars {
			if !s.inQueue.CheckWorkCompleted("evedata-war-finished", int64(war)) {
				if maxWarID > war {
					maxWarID = war
				}
				work = append(work, redisqueue.Work{Operation: "war", Parameter: war})
				err = getWarKills(s, war)
				if err != nil {
					log.Println(err)
				}
			}
		}
		s.QueueWork(work)

		if maxWarID < 100 {
			return nil
		}
	}
}

func getWarKills(s *Artifice, id int32) error {
	page := int32(1)
	for {
		kills, r, err := s.esi.ESI.WarsApi.GetWarsWarIdKillmails(context.TODO(), id, map[string]interface{}{"page": int32(page)})
		if err != nil {
			return err
		}

		work := []redisqueue.Work{}
		for _, kill := range kills {
			if !s.inQueue.CheckWorkCompleted("evedata_known_kills", int64(kill.KillmailId)) {
				work = append(work, redisqueue.Work{Operation: "killmail", Parameter: []interface{}{kill.KillmailHash, kill.KillmailId}})

				// Send to zkillboard
				zkillChan <- killmail{ID: kill.KillmailId, Hash: kill.KillmailHash}
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
