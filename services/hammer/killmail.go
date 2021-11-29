package hammer

import (
	"context"
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

	known := s.inQueue.CheckWorkCompleted("evedata_known_kills", id)
	if known {
		return
	}

	kill, _, err := s.esi.ESI.KillmailsApi.GetKillmailsKillmailIdKillmailHash(context.Background(), hash, id, nil)
	if err != nil {
		log.Println(err)
		return
	}

	s.inQueue.SetWorkCompleted("evedata_known_kills", id)
	if err != nil {
		log.Println(err)
		return
	}

	// Send out the result, but ignore DUST stuff.
	if kill.Victim.ShipTypeId < 65535 {
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

	err = s.AddAlliance(kill.Victim.AllianceId)
	if err != nil {
		log.Println(err)
		return
	}

	err = s.AddCorporation(kill.Victim.CorporationId)
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
		err = s.AddAlliance(a.AllianceId)
		if err != nil {
			log.Println(err)
			return
		}

		err = s.AddCorporation(a.CorporationId)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
