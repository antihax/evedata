package hammer

import (
	"context"
	"log"
	"strings"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/goesi"
)

func init() {
	registerConsumer("structure", structureConsumer)
}

func structureConsumer(s *Hammer, parameter interface{}) {
	structureID := parameter.(int64)

	if s.inQueue.CheckWorkExpired("evedata_structure_failure", structureID) {
		return
	}

	ctx := context.WithValue(context.Background(), goesi.ContextOAuth2, *s.token)
	struc, _, err := s.esi.ESI.UniverseApi.GetUniverseStructuresStructureId(ctx, structureID, nil)

	if err != nil {
		if strings.Contains(err.Error(), "403") {
			s.inQueue.SetWorkExpire("evedata_structure_failure", structureID, 86400)
		}
		log.Println(err)
		return
	}
	// Send out the result
	err = s.QueueResult(&datapackages.Structure{Structure: struc, StructureID: structureID}, "structure")
	if err != nil {
		log.Println(err)
		return
	}
}
