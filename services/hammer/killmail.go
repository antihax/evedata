package hammer

import (
	"log"
	"time"

	"github.com/antihax/evedata/internal/datapackages"
)

func init() {
	registerConsumer("killmail", killmailConsumer)
}

func killmailConsumer(s *Hammer, parameter interface{}) {
	parameters := parameter.([]interface{})
	hash := parameters[0].(string)
	id := int32(parameters[1].(int))

	// don't duplicate calls
	if s.inQueue.CheckWorkCompleted("evedata_known_kills", int64(id)) {
		return
	}

	kill, _, err := s.esi.ESI.KillmailsApi.GetKillmailsKillmailIdKillmailHash(nil, hash, id, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Send out the result, but ignore DUST stuff.
	if kill.Victim.ShipTypeId < 65535 {
		if kill.KillmailTime.After(time.Now().UTC().Add(-time.Hour * 24 * 365)) {
			err = s.QueueResult(&datapackages.Killmail{Hash: hash, Kill: kill}, "killmail")
			if err != nil {
				log.Println(err)
				return
			}
		}

		err = s.AddCharacter(kill.Victim.CharacterId)
		if err != nil {
			log.Println(err)
			return
		}
	}

	// We don't want to ignore attackers on DUST mails.
	for _, a := range kill.Attackers {
		err = s.AddCharacter(a.CharacterId)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
