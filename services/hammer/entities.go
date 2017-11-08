package hammer

import (
	"context"
	"log"

	"encoding/gob"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/goesi"
)

func init() {
	registerConsumer("alliance", allianceConsumer)
	registerConsumer("corporation", corporationConsumer)
	registerConsumer("character", characterConsumer)
	gob.Register(datapackages.Corporation{})
	gob.Register(datapackages.Alliance{})
	gob.Register(datapackages.Character{})

}

// AddAlliance adds an alliance to queue
func (s *Hammer) AddAlliance(allianceID int32) error {
	if allianceID > 99000000 { // Skip NPC Alliances
		if !s.inQueue.CheckWorkExpired("evedata_entity", int64(allianceID)) {
			return s.inQueue.QueueWork([]redisqueue.Work{
				{Operation: "alliance", Parameter: allianceID},
			})
		}
	}
	return nil
}

// AddCorporation adds a corporation to queue
func (s *Hammer) AddCorporation(corporationID int32) error {
	if corporationID > 98000000 { // Skip NPC Corporations
		if !s.inQueue.CheckWorkExpired("evedata_entity", int64(corporationID)) {
			return s.inQueue.QueueWork([]redisqueue.Work{
				{Operation: "corporation", Parameter: corporationID},
			})
		}
	}
	return nil
}

// AddCharacter adds a character to queue
func (s *Hammer) AddCharacter(characterID int32) error {
	if characterID > 90000000 { // Skip NPC Characters
		if !s.inQueue.CheckWorkExpired("evedata_entity", int64(characterID)) {
			return s.inQueue.QueueWork([]redisqueue.Work{
				{Operation: "character", Parameter: characterID},
			})
		}
	}
	return nil
}

func allianceConsumer(s *Hammer, parameter interface{}) {
	allianceID := parameter.(int32)

	alliance, r, err := s.esi.ESI.AllianceApi.GetAlliancesAllianceId(context.TODO(), allianceID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	allianceCorporations, _, err := s.esi.ESI.AllianceApi.GetAlliancesAllianceIdCorporations(context.TODO(), allianceID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	b, err := gobcoder.GobEncoder(&datapackages.Alliance{AllianceID: allianceID, Alliance: alliance, AllianceCorporations: allianceCorporations})
	if err != nil {
		log.Println(err)
		return
	}

	err = s.nsq.Publish("alliance", b)
	if err != nil {
		log.Println(err)
		return
	}

	err = s.inQueue.SetWorkExpire("evedata_entity", int64(allianceID), 10800)
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

func corporationConsumer(s *Hammer, parameter interface{}) {
	corporationID := parameter.(int32)

	corporation, r, err := s.esi.ESI.CorporationApi.GetCorporationsCorporationId(context.TODO(), corporationID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	b, err := gobcoder.GobEncoder(&datapackages.Corporation{CorporationID: corporationID, Corporation: corporation})
	if err != nil {
		log.Println(err)
		return
	}

	err = s.nsq.Publish("corporation", b)
	if err != nil {
		log.Println(err)
		return
	}

	s.inQueue.SetWorkExpire("evedata_entity", int64(corporationID), 10800)

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
	return
}

func characterConsumer(s *Hammer, parameter interface{}) {
	characterID := parameter.(int32)

	character, r, err := s.esi.ESI.CharacterApi.GetCharactersCharacterId(context.TODO(), characterID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	corporationHistory, _, err := s.esi.ESI.CharacterApi.GetCharactersCharacterIdCorporationhistory(context.TODO(), characterID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	b, err := gobcoder.GobEncoder(&datapackages.Character{CharacterID: characterID, Character: character, CorporationHistory: corporationHistory})
	if err != nil {
		log.Println(err)
		return
	}

	err = s.nsq.Publish("character", b)
	if err != nil {
		log.Println(err)
		return
	}

	s.inQueue.SetWorkExpire("evedata_entity", int64(characterID), 10800)

	// Grab intel from meta data
	err = s.AddCorporation(character.CorporationId)
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

	return
}
