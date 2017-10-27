package hammer

import (
	"log"

	"encoding/gob"

	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/goesi/esi"
)

func init() {
	registerConsumer("killmail", killmailConsumer)
	gob.Register(esi.GetKillmailsKillmailIdKillmailHashOk{})
}

func killmailConsumer(s *Hammer, parameter interface{}) {
	parameters := parameter.([]interface{})

	hash := parameters[0].(string)
	id := parameters[1].(int32)

	// We know this kill, early out
	if s.inQueue.CheckWorkCompleted("evedata_known_kills", int64(id)) {
		return
	}

	kill, _, err := s.esi.ESI.KillmailsApi.GetKillmailsKillmailIdKillmailHash(nil, hash, id, nil)
	if err != nil {
		log.Println(err)
		return
	}

	err = s.inQueue.SetWorkCompleted("evedata_known_kills", int64(id))
	if err != nil {
		log.Println(err)
	}

	b, err := gobcoder.GobEncoder(kill)
	if err != nil {
		log.Println(err)
		return
	}

	err = s.nsq.Publish("killmail", b)
	if err != nil {
		log.Println(err)
		return
	}

	return
}
