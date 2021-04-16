package hammer

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/goesi"
)

func init() {
	//registerConsumer("structure", structureConsumer)
	//registerConsumer("characterStructures", characterStructuresConsumer)
}

func structureConsumer(s *Hammer, parameter interface{}) {
	structureID := parameter.(int64)

	if s.inQueue.CheckWorkExpired("evedata_structure_failure", structureID) {
		return
	}

	ctx := context.WithValue(context.Background(), goesi.ContextOAuth2, *s.token)
	structure, _, err := s.esi.ESI.UniverseApi.GetUniverseStructuresStructureId(ctx, structureID, nil)

	if err != nil {
		if strings.Contains(err.Error(), "403") {
			s.inQueue.SetWorkExpire("evedata_structure_failure", structureID, 86400)
		}
		log.Println(err)
		return
	}
	// Send out the result
	err = s.QueueResult(&datapackages.Structure{Structure: structure, StructureID: structureID}, "structure")
	if err != nil {
		log.Println(err)
		return
	}
}

// Handle character structures separately since they should remain private
func characterStructuresConsumer(s *Hammer, parameter interface{}) {
	parameters := parameter.([]interface{})
	characterID := int32(parameters[0].(int))
	tokenCharacterID := int32(parameters[1].(int))
	structureID := parameters[2].(int64)

	if s.inQueue.CheckWorkExpired("evedata_structurechar_failure",
		fmt.Sprintf("%d%d", structureID, tokenCharacterID)) {
		log.Printf("failed structure ignored %d\n", structureID)
		return
	}

	ctx, err := s.GetTokenSourceContext(context.Background(), characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	// [TODO] tick failure to database
	structure, r, err := s.esi.ESI.UniverseApi.GetUniverseStructuresStructureId(ctx, structureID, nil)
	if err != nil {
		s.tokenStore.CheckSSOError(characterID, tokenCharacterID, err)
		if r != nil && r.StatusCode == 403 {
			err := s.inQueue.SetWorkExpire("evedata_structurechar_failure", fmt.Sprintf("%d%d", structureID, tokenCharacterID), 86400*3)
			if err != nil {
				log.Printf("failed setting failure: %s %d\n", err, structureID)
			}
		}
		log.Println(err)
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.CharacterStructure{Structure: structure, StructureID: structureID, CharacterID: tokenCharacterID}, "characterStructure")
	if err != nil {
		log.Println(err)
		return
	}
}
