package artifice

import (
	"time"

	"github.com/antihax/evedata/internal/redisqueue"
)

func init() {
	registerTrigger("characterTransactions", characterTransactions, time.NewTicker(time.Second*3600))
	registerTrigger("characterAssets", characterAssets, time.NewTicker(time.Second*3600))
}

func characterTransactions(s *Artifice) error {
	entities, err := s.db.Query(
		`SELECT characterID, tokenCharacterID FROM evedata.crestTokens T
		WHERE scopes LIKE "%read_character_wallet%"`)
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

func characterAssets(s *Artifice) error {
	entities, err := s.db.Query(
		`SELECT characterID, tokenCharacterID FROM evedata.crestTokens T
		WHERE scopes LIKE "%assets.read_assets%"`)
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

		work = append(work, redisqueue.Work{Operation: "characterAssets", Parameter: []int32{cid, tcid}})

	}
	s.QueueWork(work)
	entities.Close()

	return nil
}
