package hammer

import (
	"log"

	"encoding/gob"

	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/goesi/v1"
	"github.com/garyburd/redigo/redis"
)

func init() {
	registerConsumer("killmail", killmailConsumer)
	gob.Register(goesiv1.GetKillmailsKillmailIdKillmailHashOk{})
}

func killmailConsumer(s *Hammer, parameter interface{}) {
	parameters := parameter.([]interface{})

	hash := parameters[0].(string)
	id := parameters[1].(int32)

	// We know this kill, early out
	if s.knownKill(id) {
		return
	}

	kill, _, err := s.esi.V1.KillmailsApi.GetKillmailsKillmailIdKillmailHash(hash, id, nil)
	if err != nil {
		log.Println(err)
		return
	}
	s.setKnownKill(id)

	err = s.outQueue.QueueWork([]redisqueue.Work{
		redisqueue.Work{Operation: "killmail", Parameter: kill},
	})
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (c *Hammer) knownKill(id int32) bool {
	conn := c.redis.Get()
	defer conn.Close()

	found, _ := redis.Bool(conn.Do("SISMEMBER", "evedata_known_kills", id))
	return found
}

func (c *Hammer) setKnownKill(id int32) {
	conn := c.redis.Get()
	defer conn.Close()

	conn.Do("SADD", "evedata_known_kills", id)
}
