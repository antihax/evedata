package hammer

import (
	"context"
	"log"

	"github.com/antihax/evedata/internal/datapackages"
)

func init() {
	// Resolve mutaplasmids dogma from contracts and things
	registerConsumer("mutatedItem", func(s *Hammer, parameter interface{}) {
		itemID := int64(parameter.([]interface{})[0].(int64))
		typeID := int32(parameter.([]interface{})[1].(int))
		h, _, err := s.esi.ESI.DogmaApi.GetDogmaDynamicItemsTypeIdItemId(context.Background(), itemID, typeID, nil)
		if err != nil {
			log.Println(err)
			return
		}

		// Send out the result
		err = s.QueueResult(&datapackages.ResolveItems{
			Item:   h,
			ItemID: itemID,
			TypeID: typeID}, "mutatedItem")
		if err != nil {
			log.Println(err)
			return
		}
	})
}
