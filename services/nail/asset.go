package nail

import (
	"log"

	"github.com/antihax/evedata/internal/sqlhelper"

	sq "github.com/Masterminds/squirrel"
	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("characterAssets", func(s *Nail, consumer *nsq.Consumer) {
		consumer.AddConcurrentHandlers(s.wait(nsq.HandlerFunc(s.characterAssetsConsumer)), 5)
	})
}

func (s *Nail) characterAssetsConsumer(message *nsq.Message) error {
	assets := datapackages.CharacterAssets{}
	err := gobcoder.GobDecoder(message.Body, &assets)
	if err != nil {
		log.Println(err)
		return err
	}

	// Start a new transaction to keep the delete an inserts together
	tx, err := s.db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	// Delete everything because we will replace with the new assets.
	_, err = tx.Exec("DELETE FROM evedata.assets WHERE characterID = ?;", assets.TokenCharacterID)
	if err != nil {
		log.Println(err)
		return err
	}

	// early out if there are no assets.
	if len(assets.Assets) == 0 {
		return tx.Commit()
	}

	// Dump all assets into the DB.
	assetSQL := sq.Insert("evedata.assets").
		Columns("locationID", "typeID", "quantity", "characterID",
			"locationFlag", "itemID", "locationType", "isSingleton")

	// Build a multi-insert statement
	for _, asset := range assets.Assets {
		assetSQL = assetSQL.Values(asset.LocationId, asset.TypeId, asset.Quantity, assets.TokenCharacterID,
			asset.LocationFlag, asset.ItemId, asset.LocationType, asset.IsSingleton)
	}

	// Retry the transaction until we succeed.
	err = sqlhelper.RetrySquirrelInsertTransaction(tx, assetSQL)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
