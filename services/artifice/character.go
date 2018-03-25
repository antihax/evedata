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
	registerTrigger("characterContactSync", characterContactSync, time.NewTicker(time.Second*360))
	registerTrigger("characterAuthOwners", characterAuthOwners, time.NewTicker(time.Second*3600))
	registerTrigger("crestCharacters", crestCharacters, time.NewTicker(time.Second*3600))

	registerTrigger("allianceContacts", allianceContacts, time.NewTicker(time.Second*360))
	registerTrigger("corporationContacts", corporationContacts, time.NewTicker(time.Second*360))
}

func characterTransactions(s *Artifice) error {
	work := []redisqueue.Work{}
	if pairs, err := s.GetCharactersForScope("read_character_wallet"); err != nil {
		return err
	} else {
		for _, p := range pairs {
			work = append(work, redisqueue.Work{Operation: "characterWalletTransactions", Parameter: []int32{p.CharacterID, p.TokenCharacterID}})
			work = append(work, redisqueue.Work{Operation: "characterWalletJournal", Parameter: []int32{p.CharacterID, p.TokenCharacterID}})
		}
	}

	return s.QueueWork(work, redisqueue.Priority_Normal)
}

func characterAssets(s *Artifice) error {
	work := []redisqueue.Work{}
	if pairs, err := s.GetCharactersForScope("read_assets"); err != nil {
		return err
	} else {
		for _, p := range pairs {
			work = append(work, redisqueue.Work{Operation: "characterAssets", Parameter: []int32{p.CharacterID, p.TokenCharacterID}})
		}
	}

	return s.QueueWork(work, redisqueue.Priority_Normal)
}

func characterNotifications(s *Artifice) error {
	work := []redisqueue.Work{}
	if pairs, err := s.GetCharactersForScope("read_notifications"); err != nil {
		return err
	} else {
		for _, p := range pairs {
			work = append(work, redisqueue.Work{Operation: "characterNotifications", Parameter: []int32{p.CharacterID, p.TokenCharacterID}})
		}
	}

	return s.QueueWork(work, redisqueue.Priority_High)
}

func allianceContacts(s *Artifice) error {
	work := []redisqueue.Work{}
	if pairs, err := s.GetAllianceForScope("alliances.read_contacts"); err != nil {
		return err
	} else {
		for _, p := range pairs {
			work = append(work, redisqueue.Work{Operation: "allianceContacts", Parameter: []int32{p.CharacterID, p.TokenCharacterID, p.AllianceID}})
		}
	}

	return s.QueueWork(work, redisqueue.Priority_High)
}

func corporationContacts(s *Artifice) error {
	work := []redisqueue.Work{}
	if pairs, err := s.GetCorporationForScope("corporations.read_contacts"); err != nil {
		return err
	} else {
		for _, p := range pairs {
			work = append(work, redisqueue.Work{Operation: "corporationContacts", Parameter: []int32{p.CharacterID, p.TokenCharacterID, p.CorporationID}})
		}
	}

	return s.QueueWork(work, redisqueue.Priority_High)
}

func characterContactSync(s *Artifice) error {
	entities, err := s.db.Query(
		`SELECT S.characterID, source, group_concat(destination) AS destinations
		FROM evedata.contactSyncs S
		INNER JOIN evedata.crestTokens T ON T.tokenCharacterID = destination
		WHERE lastStatus NOT LIKE "%400 Bad Request%" AND scopes LIKE ?
		GROUP BY source
		HAVING max(nextSync) < UTC_TIMESTAMP();`, "%characters.read_contacts%")
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
	work := []redisqueue.Work{}
	if pairs, err := s.GetCharactersForScope("read_corporation_roles"); err != nil {
		return err
	} else {
		for _, p := range pairs {
			work = append(work, redisqueue.Work{Operation: "characterAuthOwner", Parameter: []int32{p.CharacterID, p.TokenCharacterID}})
		}
	}

	return s.QueueWork(work, redisqueue.Priority_High)
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

	// Get a list of characters sharing data to check for changes
	sharing := make(map[int32]int32)
	rows, err := s.db.Query(`
		SELECT DISTINCT T.tokenCharacterID, T.corporationID FROM evedata.crestTokens T
			INNER JOIN evedata.sharing S ON S.tokenCharacterID = T.tokenCharacterID`)
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var char, corp int32
		err := rows.Scan(&char, &corp)
		if err != nil {
			log.Println(err)
			return err
		}
		sharing[char] = corp
	}

	for start := 0; start < len(chars); start = start + 50 {
		end := min(start+50, len(chars))
		if affiliation, _, err := s.esi.ESI.CharacterApi.PostCharactersAffiliation(context.Background(), chars[start:end], nil); err != nil {
			log.Println(err)
			continue
		} else {
			tx, err := s.db.Beginx()
			defer tx.Rollback()
			if err != nil {
				return err
			}
			for _, c := range affiliation {
				// See if they changed corporation, if they have shares, warn them they are still sharing.
				if check, ok := sharing[c.CharacterId]; ok {
					if check != c.CorporationId {
						s.mailCorporationChangeWithShares(c.CharacterId)
					}
				}

				if _, err := tx.Exec(`UPDATE evedata.crestTokens
					SET corporationID = ?, allianceID = ?, factionID = ?
					WHERE tokenCharacterID = ?;
					`, c.CorporationId, c.AllianceId, c.FactionId, c.CharacterId); err != nil {
					log.Println(err)
					continue
				}
			}
			if err = retryTransaction(tx); err != nil {
				return err
			}

		}

		// Update integrations services with factionID.
		if err = s.doSQL(`
			UPDATE evedata.integrations S
			INNER JOIN evedata.crestTokens C ON C.allianceID = S.entityID OR C.corporationID = S.entityID 
				AND expiry > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 20 MINUTE)
			SET S.factionID = C.factionID`); err != nil {
			return err
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
