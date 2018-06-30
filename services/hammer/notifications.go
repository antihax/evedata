package hammer

import (
	"context"
	"log"

	"github.com/antihax/goesi/esi"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/ghodss/yaml"
)

func init() {
	registerConsumer("characterNotifications", characterNotificationsConsumer)
}

func addID(e *[]int32, n int32) {
	if n > 0 {
		*e = append(*e, n)
	}
}

func characterNotificationsConsumer(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]interface{})
	characterID := int32(parameters[0].(int))
	tokenCharacterID := int32(parameters[1].(int))

	ctx, err := s.GetTokenSourceContext(context.Background(), characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	notifications, _, err := s.esi.ESI.CharacterApi.GetCharactersCharacterIdNotifications(ctx, tokenCharacterID, nil)
	if err != nil {
		s.tokenStore.CheckSSOError(characterID, tokenCharacterID, err)
		log.Println(err)
		return
	}
	if len(notifications) == 0 {
		return
	}
	// see what we can learn about these notifications so our alerts do not fail
	s.learnFromNotifications(notifications)

	// Send out the result
	err = s.QueueResult(&datapackages.CharacterNotifications{
		CharacterID:      characterID,
		TokenCharacterID: tokenCharacterID,
		Notifications:    notifications,
	}, "characterNotifications")
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Hammer) learnFromNotifications(notifications []esi.GetCharactersCharacterIdNotifications200Ok) {
	type stripEntities struct {
		CharacterID         int32 `yaml:"characterID,omitempty"`
		AggressorCorpID     int32 `yaml:"aggressorCorpID,omitempty"`
		AggressorAllianceID int32 `yaml:"aggressorAllianceID,omitempty"`
		DeclaredByID        int32 `yaml:"declaredByID,omitempty"`
		AgainstID           int32 `yaml:"againstID,omitempty"`
		CharID              int32 `yaml:"charID,omitempty"`
		CorpID              int32 `yaml:"corpID,omitempty"`
		AggressorID         int32 `yaml:"aggressorID,omitempty"`
	}

	lookup := []int32{}
	for _, n := range notifications {
		l := stripEntities{}
		yaml.Unmarshal([]byte(n.Text), &l)
		addID(&lookup, l.AgainstID)
		addID(&lookup, l.AggressorAllianceID)
		addID(&lookup, l.AggressorCorpID)
		addID(&lookup, l.AgainstID)
		addID(&lookup, l.CharacterID)
		addID(&lookup, l.DeclaredByID)
		addID(&lookup, l.CharID)
		addID(&lookup, l.CorpID)
		addID(&lookup, l.AggressorID)
	}

	// Lookup every thing we learned
	if err := s.BulkLookup(sliceUniq(lookup)); err != nil {
		log.Println(err)
	}
}
