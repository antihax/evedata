package nail

import (
	"fmt"
	"log"
	"strings"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("characterAssets", spawnCharacterAssetsConsumer)
}

func spawnCharacterAssetsConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.characterAssetsConsumer)))
}

func (s *Nail) characterAssetsConsumer(message *nsq.Message) error {
	assets := datapackages.CharacterAssets{}
	err := gobcoder.GobDecoder(message.Body, &assets)
	if err != nil {
		log.Println(err)
		return err
	}
	if len(assets.Assets) == 0 {
		return nil
	}
	var values []string

	err = s.doSQL("DELETE FROM evedata.assets WHERE characterID = ?;", assets.TokenCharacterID)
	if err != nil {
		log.Println(err)
		return err
	}

	// Dump all assets into the DB.
	count := 0
	for _, asset := range assets.Assets {
		count++
		values = append(values, fmt.Sprintf("(%d,%d,%d,%d,%q,%d,%q,%v)",
			asset.LocationId, asset.TypeId, asset.Quantity, assets.TokenCharacterID,
			asset.LocationFlag, asset.ItemId, asset.LocationType, asset.IsSingleton))

		if count > 100 {
			err := s.doSQL(doAssets(values))
			if err != nil {
				return err
			}
			values = values[:0]
			count = 0
		}
	}

	return s.doSQL(doAssets(values))
}

func doAssets(values []string) string {
	return fmt.Sprintf(`INSERT INTO evedata.assets
		(locationID, typeID, quantity, characterID, 
		locationFlag, itemID, locationType, isSingleton)
		VALUES %s 
		ON DUPLICATE KEY UPDATE typeID = typeID;`, strings.Join(values, ",\n"))
}
