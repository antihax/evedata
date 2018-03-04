package hammer

import (
	"log"

	"github.com/antihax/evedata/internal/datapackages"
)

func init() {
	registerConsumer("killmail", killmailConsumer)
}

func killmailConsumer(s *Hammer, parameter interface{}) {
	parameters := parameter.([]interface{})
	hash := parameters[0].(string)
	id := int32(parameters[1].(int))

	kill, _, err := s.esi.ESI.KillmailsApi.GetKillmailsKillmailIdKillmailHash(nil, hash, id, nil)
	if err != nil {
		log.Println(err)
		return
	}

	err = s.inQueue.SetWorkCompleted("evedata_known_kills", int64(id))
	if err != nil {
		log.Println(err)
	}

	// Send out the result
	err = s.QueueResult(&datapackages.Killmail{Hash: hash, Kill: kill}, "killmail")
	if err != nil {
		log.Println(err)
		return
	}

	err = s.AddCharacter(kill.Victim.CharacterId)
	if err != nil {
		log.Println(err)
		return
	}

	for _, a := range kill.Attackers {
		err = s.AddCharacter(a.CharacterId)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
