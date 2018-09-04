package hammer

import (
	"context"
	"database/sql"
	"log"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
)

func init() {
	registerConsumer("charSearch", charSearchConsumer)
	registerConsumer("alliance", allianceConsumer)
	registerConsumer("corporation", corporationConsumer)
	registerConsumer("corporationHistory", corporationHistoryConsumer)
	registerConsumer("allianceHistory", allianceHistoryConsumer)

	registerConsumer("character", characterConsumer)

	registerConsumer("loyaltyStore", loyaltyStoreConsumer)
}

// AddAlliance adds an alliance to queue
func (s *Hammer) AddAlliance(allianceID int32) error {
	if !s.inQueue.CheckWorkExpired("evedata_entity", int64(allianceID)) {
		return s.inQueue.QueueWork([]redisqueue.Work{
			{Operation: "alliance", Parameter: allianceID},
		}, redisqueue.Priority_Low)
	}
	return nil
}

// AddCorporation adds a corporation to queue
func (s *Hammer) AddCorporation(corporationID int32) error {
	if corporationID > 1000000 {
		if !s.inQueue.CheckWorkExpired("evedata_entity", int64(corporationID)) {
			return s.inQueue.QueueWork([]redisqueue.Work{
				{Operation: "corporation", Parameter: corporationID},
				{Operation: "allianceHistory", Parameter: corporationID},
			}, redisqueue.Priority_Low)
		}
	}
	return nil
}

// AddCharacter adds a character to queue
func (s *Hammer) AddCharacter(characterID int32) error {
	if characterID > 3000000 {
		if !s.inQueue.CheckWorkExpired("evedata_entity", int64(characterID)) {
			return s.inQueue.QueueWork([]redisqueue.Work{
				{Operation: "character", Parameter: characterID},
				{Operation: "corporationHistory", Parameter: characterID},
			}, redisqueue.Priority_Low)
		}
	}
	return nil
}

// BulkLookup looks up ID to entities in bulk and adds them to the queue
func (s *Hammer) BulkLookup(ids []int32) error {
	if len(ids) > 0 {
		for start := 0; start < len(ids); start = start + 1000 {
			end := min(start+1000, len(ids))
			if len(ids[start:end]) == 0 {
				break
			}
			resolved, _, err := s.esi.ESI.UniverseApi.PostUniverseNames(context.Background(), ids[start:end], nil)
			if err != nil {
				log.Printf("%s %+v\n", err, ids[start:end])
				return err
			}
			for _, r := range resolved {
				switch r.Category {
				case "alliance":
					s.AddAlliance(r.Id)
				case "corporation":
					s.AddCorporation(r.Id)
				case "character":
					s.AddCharacter(r.Id)
				}
			}
		}
	}
	return nil
}

func sliceUniq(s []int32) []int32 {
	for i := 0; i < len(s); i++ {
		for i2 := i + 1; i2 < len(s); i2++ {
			if s[i] == s[i2] {
				// delete
				s = append(s[:i2], s[i2+1:]...)
				i2--
			}
		}
	}
	return s
}

func charSearchConsumer(s *Hammer, parameter interface{}) {
	char := parameter.(string)

	// Check if we know this character already
	id, err := s.GetCharacterIDByName(char)
	if err != nil {
		log.Println(err)
		return
	}

	if id == 0 {
		search, _, err := s.esi.ESI.SearchApi.GetSearch(nil, []string{"character"}, char,
			&esi.GetSearchOpts{
				Strict: optional.NewBool(true),
			})
		if err != nil {
			log.Println(err)
			return
		}
		if len(search.Character) > 0 {
			for _, newid := range search.Character {
				s.AddCharacter(newid)
			}
		}
	} else { // add the character to the queue so we get latest data.
		s.AddCharacter(id)
	}
	return
}

// GetCharacterIDByName checks if a character exists in the database
func (s *Hammer) GetCharacterIDByName(character string) (int32, error) {
	var id int32
	if err := s.db.Get(&id, `
		SELECT characterID 
		FROM evedata.characters C
		WHERE C.name = ? LIMIT 1;`, character); err != nil && err != sql.ErrNoRows {
		return id, err
	}
	return id, nil
}

func allianceConsumer(s *Hammer, parameter interface{}) {
	allianceID := int32(parameter.(int))

	alliance, _, err := s.esi.ESI.AllianceApi.GetAlliancesAllianceId(context.Background(), allianceID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	allianceCorporations, _, err := s.esi.ESI.AllianceApi.GetAlliancesAllianceIdCorporations(context.Background(), allianceID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.Alliance{
		AllianceID:           allianceID,
		Alliance:             alliance,
		AllianceCorporations: allianceCorporations,
	}, "alliance")
	if err != nil {
		log.Println(err)
		return
	}

	err = s.inQueue.SetWorkExpire("evedata_entity", int64(allianceID), 43200)
	if err != nil {
		log.Println(err)
		return
	}

	// Grab intel from meta data
	// Add all known corporations
	for _, corp := range allianceCorporations {
		err = s.AddCorporation(corp)
		if err != nil {
			log.Println(err)
			return
		}
	}
	return
}

func loyaltyStoreConsumer(s *Hammer, parameter interface{}) {
	corporationID := int32(parameter.(int))
	store, _, err := s.esi.ESI.LoyaltyApi.GetLoyaltyStoresCorporationIdOffers(context.Background(), corporationID, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if len(store) > 0 {
		// Send out the result
		err = s.QueueResult(&datapackages.Store{
			CorporationID: corporationID,
			Store:         store},
			"loyaltyStore")
		if err != nil {
			log.Println(err)
			return
		}
	}
	return
}

func corporationConsumer(s *Hammer, parameter interface{}) {
	corporationID := int32(parameter.(int))
	corporation, _, err := s.esi.ESI.CorporationApi.GetCorporationsCorporationId(context.Background(), corporationID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.Corporation{
		CorporationID: corporationID,
		Corporation:   corporation,
	}, "corporation")
	if err != nil {
		log.Println(err)
		return
	}

	s.inQueue.SetWorkExpire("evedata_entity", int64(corporationID), 43200)

	// Grab intel from meta data
	err = s.AddCharacter(corporation.CeoId)
	if err != nil {
		log.Println(err)
		return
	}

	if corporation.CeoId != corporation.CreatorId {
		err := s.AddCharacter(corporation.CreatorId)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func corporationHistoryConsumer(s *Hammer, parameter interface{}) {
	characterID := int32(parameter.(int))

	corporationHistory, _, err := s.esi.ESI.CharacterApi.GetCharactersCharacterIdCorporationhistory(context.Background(), characterID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.CorporationHistory{
		CharacterID:        characterID,
		CorporationHistory: corporationHistory,
	}, "corporationHistory")
	if err != nil {
		log.Println(err)
		return
	}

	// Add all known corporations
	for _, corp := range corporationHistory {
		err = s.AddCorporation(corp.CorporationId)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func allianceHistoryConsumer(s *Hammer, parameter interface{}) {
	corporationID := int32(parameter.(int))

	allianceHistory, _, err := s.esi.ESI.CorporationApi.GetCorporationsCorporationIdAlliancehistory(context.Background(), corporationID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.AllianceHistory{
		CorporationID:   corporationID,
		AllianceHistory: allianceHistory,
	}, "allianceHistory")
	if err != nil {
		log.Println(err)
		return
	}
}

func characterConsumer(s *Hammer, parameter interface{}) {
	characterID := int32(parameter.(int))
	character, _, err := s.esi.ESI.CharacterApi.GetCharactersCharacterId(context.Background(), characterID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.Character{
		CharacterID: characterID,
		Character:   character,
	}, "character")
	if err != nil {
		log.Println(err)
		return
	}

	s.inQueue.SetWorkExpire("evedata_entity", int64(characterID), 43200)

	// Grab intel from meta data
	err = s.AddCorporation(character.CorporationId)
	if err != nil {
		log.Println(err)
		return
	}

}
