package artifice

import (
	"context"
	"log"
	"time"

	"github.com/antihax/evedata/internal/redisqueue"
)

func init() {
	registerTrigger("characterTransactions", characterTransactions, time.NewTicker(time.Second*3600))
	registerTrigger("characterAssets", characterAssets, time.NewTicker(time.Second*3600))
	registerTrigger("characterNotifications", characterNotifications, time.NewTicker(time.Second*600))
	registerTrigger("characterContactSync", characterContactSync, time.NewTicker(time.Second*300))
	registerTrigger("characterAuthOwners", characterAuthOwners, time.NewTicker(time.Second*3600))
}

func characterTransactions(s *Artifice) error {
	entities, err := s.db.Query(
		`SELECT characterID, tokenCharacterID FROM evedata.crestTokens T
		WHERE scopes LIKE "%read_character_wallet%"`)
	if err != nil {
		return err
	}
	defer entities.Close()

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

	return s.QueueWork(work, redisqueue.Priority_Normal)
}

func characterAssets(s *Artifice) error {
	entities, err := s.db.Query(
		`SELECT characterID, tokenCharacterID FROM evedata.crestTokens T
		WHERE scopes LIKE "%assets.read_assets%"`)
	if err != nil {
		return err
	}
	defer entities.Close()

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

	return s.QueueWork(work, redisqueue.Priority_Normal)
}

func characterNotifications(s *Artifice) error {
	entities, err := s.db.Query(
		`SELECT characterID, tokenCharacterID FROM evedata.crestTokens T
		WHERE scopes LIKE "%esi-characters.read_notifications%"`)
	if err != nil {
		return err
	}
	defer entities.Close()

	work := []redisqueue.Work{}

	// Loop the entities
	for entities.Next() {
		var cid, tcid int32

		err = entities.Scan(&cid, &tcid)
		if err != nil {
			return err
		}

		work = append(work, redisqueue.Work{Operation: "characterNotifications", Parameter: []int32{cid, tcid}})

	}

	return s.QueueWork(work, redisqueue.Priority_High)
}

func characterContactSync(s *Artifice) error {
	entities, err := s.db.Query(
		`SELECT S.characterID, source, group_concat(destination) AS destinations
		FROM evedata.contactSyncs S
		INNER JOIN evedata.crestTokens T ON T.tokenCharacterID = destination
		WHERE lastStatus NOT LIKE "%400 Bad Request%"
		GROUP BY source
		HAVING max(nextSync) < UTC_TIMESTAMP();`)
	if err != nil {
		return err
	}
	defer entities.Close()

	work := []redisqueue.Work{}

	// Loop the entities
	for entities.Next() {
		var (
			cid, source  int32
			destinations string
		)

		err = entities.Scan(&cid, &source, &destinations)
		if err != nil {
			return err
		}

		work = append(work, redisqueue.Work{Operation: "characterContactSync", Parameter: []interface{}{
			cid,
			source,
			destinations,
		}})
	}

	return s.QueueWork(work, redisqueue.Priority_Normal)
}

func characterAuthOwners(s *Artifice) error {
	entities, err := s.db.Query(
		`SELECT characterID, tokenCharacterID FROM evedata.crestTokens T
		WHERE scopes LIKE "%characters.read_corporation_roles%"`)
	if err != nil {
		return err
	}
	defer entities.Close()

	work := []redisqueue.Work{}

	// Update character corps/alliances
	crestCharacters(s)

	// Loop the entities
	for entities.Next() {
		var cid, tcid int32

		err = entities.Scan(&cid, &tcid)
		if err != nil {
			return err
		}

		work = append(work, redisqueue.Work{Operation: "characterAuthOwner", Parameter: []int32{cid, tcid}})

	}

	return s.QueueWork(work, redisqueue.Priority_Urgent)
}

// This is sensitive so we will do it here to prevent mixing it with public data.
// figure out character alliance and corp for our members
func crestCharacters(s *Artifice) error {
	var chars []int32
	err := s.db.Select(&chars,
		`SELECT DISTINCT tokenCharacterID FROM evedata.crestTokens`)
	if err != nil {
		log.Println(err)
		return err
	}

	for start := 0; start < len(chars); start = start + 1000 {
		end := min(start+1000, len(chars))
		if affiliation, _, err := s.esi.ESI.CharacterApi.PostCharactersAffiliation(context.Background(), chars[start:end], nil); err != nil {
			log.Println(err)
			continue
		} else {
			for _, c := range affiliation {
				if err := s.doSQL(`UPDATE evedata.crestTokens
					SET corporationID = ?, allianceID = ?, factionID = ?
					WHERE tokenCharacterID = ?;
					`, c.CorporationId, c.AllianceId, c.FactionId, c.CharacterId); err != nil {
					log.Println(err)
					continue
				}
			}
		}

	}
	return nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
