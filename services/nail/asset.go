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
		consumer.AddConcurrentHandlers(s.wait(nsq.HandlerFunc(s.characterAssetsConsumer)), 10)
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
		return sqlhelper.RetryTransaction(tx)
	}

	// Dump all assets into the DB.
	assetSQL := sq.Insert("evedata.assets").
		Columns("locationID", "typeID", "quantity", "characterID",
			"locationFlag", "itemID", "locationType", "isSingleton")

	count := 0
	// Build a multi-insert statement
	for _, asset := range assets.Assets {
		count++
		assetSQL = assetSQL.Values(asset.LocationId, asset.TypeId, asset.Quantity, assets.TokenCharacterID,
			asset.LocationFlag, asset.ItemId, asset.LocationType, asset.IsSingleton)

		// create statement if >80 or at end of array
		if count%80 == 0 || len(assets.Assets) == count {
			sqlq, args, err := assetSQL.ToSql()
			if err != nil {
				return err
			}

			_, err = tx.Exec(sqlq+" ON DUPLICATE KEY UPDATE quantity=VALUES(quantity) ", args...)
			if err != nil {
				return err
			}

			assetSQL = sq.Insert("evedata.assets").
				Columns("locationID", "typeID", "quantity", "characterID",
					"locationFlag", "itemID", "locationType", "isSingleton")
		}
	}

	err = sqlhelper.RetryTransaction(tx)
	if err != nil {
		return err
	}

	return nil
}
