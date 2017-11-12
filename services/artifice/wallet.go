package artifice

import (
	"time"

	"github.com/antihax/evedata/internal/redisqueue"
)

func init() {
	registerTrigger("characterTransactions", characterTransactions, time.NewTicker(time.Second*3600))
}

func characterTransactions(s *Artifice) error {
	entities, err := s.db.Query(
		`SELECT characterID, tokenCharacterID FROM evedata.crestTokens T
		WHERE scopes LIKE "%esi-assets.read_assets.v1%"`)
	if err != nil {
		return err
	}

	work := []redisqueue.Work{}

	// Loop the entities
	for entities.Next() {
		var cid, tcid int32

		err = entities.Scan(&cid, &tcid)
		if err != nil {
			return err
		}

		work = append(work, redisqueue.Work{Operation: "characterWalletTransactions", Parameter: []int32{cid, tcid}})
		work = append(work, redisqueue.Work{Operation: "characterWalletJournal", Parameter: []int32{cid, tcid}})
	}
	s.QueueWork(work)
	entities.Close()

	return nil
}
