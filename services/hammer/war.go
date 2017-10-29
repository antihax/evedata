package hammer

import (
	"context"
	"log"
	"time"

	"encoding/gob"

	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/goesi/esi"
)

func init() {
	registerConsumer("war", warConsumer)
	gob.Register(esi.GetWarsWarIdOk{})
}

func warConsumer(s *Hammer, parameter interface{}) {
	id := parameter.(int32)

	// We know this kill, early out
	if s.inQueue.CheckWorkCompleted("evedata_war_finished", int64(id)) {
		return
	}

	war, _, err := s.esi.ESI.WarsApi.GetWarsWarId(context.TODO(), id, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// if the war ended, market it finished
	if war.Finished.IsZero() == false && war.Finished.Before(time.Now().UTC()) {
		err = s.inQueue.SetWorkCompleted("evedata_known_kills", int64(id))
		if err != nil {
			log.Println(err)
		}
	}

	b, err := gobcoder.GobEncoder(war)
	if err != nil {
		log.Println(err)
		return
	}

	err = s.nsq.Publish("war", b)
	if err != nil {
		log.Println(err)
		return
	}

	return
}
