package artifice

import (
	"time"

	"github.com/antihax/evedata/internal/redisqueue"
)

func init() {
	registerTrigger("mutatedItem", resolveItems, time.NewTicker(time.Second*30))
}

func resolveItems(s *Artifice) error {
	entities, err := s.db.Query(
		`SELECT I.itemID, I.typeID FROM evedata.contractItems I
			INNER JOIN eve.invTypes T ON T.typeID = I.typeID
			LEFT OUTER JOIN evedata.mutations M ON M.itemID = I.itemID
			WHERE typeName LIKE "%Abyssal%" AND M.itemID IS NULL`)
	if err != nil {
		return err
	}
	defer entities.Close()

	work := []redisqueue.Work{}

	// Loop the entities
	for entities.Next() {
		var (
			itemID int64
			typeID int32
		)

		err = entities.Scan(&itemID, &typeID)
		if err != nil {
			return err
		}

		work = append(work, redisqueue.Work{
			Operation: "mutatedItem",
			Parameter: []interface{}{
				itemID, typeID,
			}})
	}

	return s.QueueWork(work, redisqueue.Priority_High)
}
