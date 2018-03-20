package hammer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/antihax/goesi/esi"

	"github.com/antihax/goesi"
	"golang.org/x/oauth2"
)

func init() {
	registerConsumer("characterContactSync", characterContactSyncConsumer)
}

func characterContactSyncConsumer(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]interface{})
	characterID := int32(parameters[0].(int))
	source := int32(parameters[1].(int))
	destinations := parameters[2].(string)

	char, _, err := s.esi.ESI.CharacterApi.GetCharactersCharacterId(nil, int32(source), nil)
	if err != nil {
		log.Println(err)
		return
	}

	corp, _, err := s.esi.ESI.CorporationApi.GetCorporationsCorporationId(nil, char.CorporationId, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Find the Entity ID to search for wars.
	var searchID int32
	if corp.AllianceId > 0 {
		searchID = corp.AllianceId
	} else {
		searchID = char.CorporationId
	}

	// Map of tokens
	type characterToken struct {
		token *oauth2.TokenSource
		cid   int32
	}
	tokens := make(map[int64]characterToken)

	// Get the tokens for our destinations
	for _, cidS := range strings.Split(destinations, ",") {
		cid, _ := strconv.ParseInt(cidS, 10, 64)
		a, err := s.tokenStore.GetTokenSource(int32(characterID), int32(cid))
		if err != nil {
			log.Println(err, characterID, cidS)
			return
		}
		// Save the token.
		tokens[cid] = characterToken{token: &a, cid: int32(cid)}
	}

	// Active Wars
	activeWars, err := s.GetActiveWarsByID((int64)(searchID))
	if err != nil {
		log.Println(err)
		return
	}

	// Pending Wars
	pendingWars, err := s.GetPendingWarsByID((int64)(searchID))
	if err != nil {
		log.Println(err)
		return
	}

	// Faction Wars
	var factionWars []FactionWarEntities
	if corp.FactionId > 0 {
		factionWars, err = s.GetFactionWarEntitiesForID(corp.FactionId)
		if err != nil {
			log.Println(err)
			return
		}
	}

	// Loop through all the destinations
	for _, token := range tokens {
		// authentication token context for destination char
		contacts := []esi.GetCharactersCharacterIdContacts200Ok{}
		auth := context.WithValue(context.Background(), goesi.ContextOAuth2, *token.token)
		page := int32(1)
		for {
			c, _, err := s.esi.ESI.ContactsApi.GetCharactersCharacterIdContacts(auth, (int32)(token.cid), map[string]interface{}{"page": int32(page)})
			if err != nil {
				s.tokenStore.CheckSSOError(characterID, token.cid, err)
				log.Println(err)
				return
			}
			if len(c) == 0 {
				break
			}
			contacts = append(contacts, c...)
			page++
		}

		var erase []int32
		var active []int32
		var pending []int32
		var pendingMove []int32
		var activeMove []int32
		var untouchableContacts int

		// Figure out how many contacts they have outside of ours
		for _, contact := range contacts {
			if contact.Standing > -0.4 {
				untouchableContacts++
			}
		}

		// Faction wars can get over the 1024 contact limit so we need to trim
		// real wars will take precedence.
		trim := len(activeWars) + len(pendingWars)

		activeCheck := make(map[int32]bool)
		pendingCheck := make(map[int32]bool)

		// Build a map of active wars
		for _, war := range activeWars {
			activeCheck[(int32)(war.ID)] = true
		}

		// Add faction wars to the active list
		maxFactionWarLength := min(1023-trim-untouchableContacts, len(factionWars))
		if maxFactionWarLength > 0 {
			for _, war := range factionWars[:maxFactionWarLength] {
				activeCheck[(int32)(war.ID)] = true
			}
		}

		// Build a map of pending wars
		for _, war := range pendingWars {
			pendingCheck[(int32)(war.ID)] = true
		}

		// Loop through all current contacts and figure out needed moves
		for _, contact := range contacts {
			// skip anything > -0.4
			if contact.Standing > -0.4 {
				continue
			}

			pend := pendingCheck[contact.ContactId]
			act := activeCheck[contact.ContactId]

			// Is this existing contact in the active list
			if !act {
				// Is this existing contact in the pending list
				if !pend { // Not in either list. delete it.
					erase = append(erase, (int32)(contact.ContactId))
				} else if pend && contact.Standing > -5.0 { // in pending list but wrong standing
					// Take it out of the active list and put into pending move.
					delete(pendingCheck, contact.ContactId)
					pendingMove = append(pendingMove, (int32)(contact.ContactId))
				} else if pend && contact.Standing == -5.0 { // Contact correct, do nothing.
					delete(pendingCheck, contact.ContactId)
				}
			} else if act && contact.Standing != -10.0 { // in active list, but wrong standing
				delete(activeCheck, contact.ContactId)
				activeMove = append(activeMove, (int32)(contact.ContactId))
			} else if act && contact.Standing == -10.0 { // Contact correct, do nothing.
				delete(activeCheck, contact.ContactId)
			}
		}

		// Build a list of active wars to add
		for con := range activeCheck {
			active = append(active, con)

		}

		// Build a list of pending wars to add
		for con := range pendingCheck {
			pending = append(pending, con)
		}

		for i := 0; i < 2; i++ {
			// Erase contacts which have no wars.
			if len(erase) > 0 {
				for start := 0; start < len(erase); start = start + 20 {
					end := min(start+20, len(erase))
					if len(erase[start:end]) == 0 {
						break
					}
					if _, err := s.esi.ESI.ContactsApi.DeleteCharactersCharacterIdContacts(auth, token.cid, erase[start:end], nil); err != nil {
						s.tokenStore.CheckSSOError(characterID, token.cid, err)
						log.Println(err)
					}
				}
			}
		}

		// Add contacts for active wars
		if len(active) > 0 {
			for start := 0; start < len(active); start = start + 100 {
				end := min(start+100, len(active))
				if len(active[start:end]) == 0 {
					break
				}
				if _, _, err := s.esi.ESI.ContactsApi.PostCharactersCharacterIdContacts(auth, (int32)(token.cid), active[start:end], -10, nil); err != nil {
					s.tokenStore.CheckSSOError(characterID, token.cid, err)
					log.Println(err, active[start:end])
				}
			}
		}

		// Add contacts for pending wars
		if len(pending) > 0 {
			for start := 0; start < len(pending); start = start + 100 {
				end := min(start+100, len(pending))
				if len(pending[start:end]) == 0 {
					break
				}
				if _, _, err := s.esi.ESI.ContactsApi.PostCharactersCharacterIdContacts(auth, (int32)(token.cid), pending[start:end], -5, nil); err != nil {
					s.tokenStore.CheckSSOError(characterID, token.cid, err)
					log.Println(err)
				}
			}
		}

		// Move contacts to active wars
		if len(activeMove) > 0 {
			for start := 0; start < len(activeMove); start = start + 100 {
				end := min(start+100, len(activeMove))
				if len(activeMove[start:end]) == 0 {
					break
				}
				if _, err := s.esi.ESI.ContactsApi.PutCharactersCharacterIdContacts(auth, (int32)(token.cid), activeMove[start:end], -10, nil); err != nil {
					s.tokenStore.CheckSSOError(characterID, token.cid, err)
					log.Println(err)
				}
			}
		}

		// Move contacts to pending wars
		if len(pendingMove) > 0 {
			for start := 0; start < len(pendingMove); start = start + 100 {
				end := min(start+100, len(pendingMove))
				if len(pendingMove[start:end]) == 0 {
					break
				}
				if _, err := s.esi.ESI.ContactsApi.PutCharactersCharacterIdContacts(auth, (int32)(token.cid), pendingMove[start:end], -5, nil); err != nil {
					s.tokenStore.CheckSSOError(characterID, token.cid, err)
					log.Println(err)
				}
			}
		}
	}
}

func (s *Hammer) deleteContactsCREST(auth context.Context, characterID int32, contacts []int32) error {
	names, _, err := s.esi.ESI.UniverseApi.PostUniverseNames(context.Background(), contacts, nil)
	if err != nil {
		return err
	}

	tokenSource, ok := auth.Value(goesi.ContextOAuth2).(oauth2.TokenSource)
	if ok {
		for _, erase := range names {
			ref := fmt.Sprintf("https://crest-tq.eveonline.com/%ss/%d/", erase.Category, erase.Id)
			err := s.esi.EVEAPI.ContactDeleteV1(tokenSource, int64(characterID), int64(erase.Id), ref)
			if err != nil {
				return err
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

// Entity denormalizes corporations, alliance, and characters
type ContactEntity struct {
	ID   int64
	Type string
}

// GetActiveWarsByID gets active wars for an entityID
// [BENCHMARK] 0.000 sec / 0.000 sec
func (s *Hammer) GetActiveWarsByID(id int64) ([]ContactEntity, error) {
	w := []ContactEntity{}
	if err := s.db.Select(&w, `
			SELECT K.id, type FROM
			(SELECT defenderID AS id FROM evedata.wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = ?
			UNION
			SELECT aggressorID AS id FROM evedata.wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND defenderID = ?
			UNION
			SELECT aggressorID  AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND allyID = ?
			UNION
			SELECT allyID AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = ?) AS K
			INNER JOIN evedata.entities C ON C.id = K.id
		`, id, id, id, id); err != nil {
		return nil, err
	}
	return w, nil
}

// GetPendingWarsByID gets pending wars for an entityID
// [BENCHMARK] 0.000 sec / 0.000 sec
func (s *Hammer) GetPendingWarsByID(id int64) ([]ContactEntity, error) {
	w := []ContactEntity{}
	if err := s.db.Select(&w, `
			SELECT K.id, type FROM
			(SELECT defenderID AS id FROM evedata.wars WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND aggressorID = ?
			UNION
			SELECT aggressorID AS id FROM evedata.wars WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND defenderID = ?
			UNION
			SELECT aggressorID  AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND allyID = ?
			UNION
			SELECT allyID AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND aggressorID = ?) AS K
			INNER JOIN evedata.entities C ON C.id = K.id
		`, id, id, id, id); err != nil {
		return nil, err
	}
	return w, nil
}

// FactionsByName resolves faction name to ID
var FactionsByName = map[string]int32{"Caldari": 500001, "Minmatar": 500002, "Amarr": 500003, "Gallente": 500004}

// FactionsByID resolves faction ID to Name
var FactionsByID = map[int32]string{500001: "Caldari", 500002: "Minmatar", 500003: "Amarr", 500004: "Gallente"}

// FactionsAtWar resolves two enemy parties for each factionID
var FactionsAtWar = map[int32][]int32{
	500001: {500002, 500004}, // Caldari  : Minmatar, Gallente
	500003: {500002, 500004}, // Amarr    : Minmatar, Gallente
	500002: {500001, 500003}, // Minmatar : Caldari, Amarr
	500004: {500001, 500003}, // Gallente : Caldari, Amarr
}

type FactionWarEntities struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Type string `db:"type" json:"type"`
}

// [BENCHMARK] 0.031 sec / 0.000 sec
func (s *Hammer) GetFactionWarEntitiesForID(factionID int32) ([]FactionWarEntities, error) {
	if FactionsByID[factionID] == "" {
		return nil, errors.New("Unknown FactionID")
	}

	// Due to CCP limitation, make sure count is under 1024, cut stuff off until it is.

	wars := FactionsAtWar[factionID]
	w := []FactionWarEntities{}
	if err := s.db.Select(&w, `
		SELECT 
			DISTINCT IF(C.allianceID > 0, C.allianceID, corporationID) AS id,
			IF(C.allianceID > 0, A.name, C.name) AS name,
			IF(C.allianceID > 0, "alliance", "corporation") AS type 
			FROM evedata.corporations C 
			LEFT OUTER JOIN evedata.alliances A ON C.allianceID = A.allianceID
			INNER JOIN evedata.entityKillStats K ON K.id = IF(C.allianceID > 0, C.allianceID, C.corporationID)
			WHERE factionID IN (?, ?) AND C.memberCount > 0
			ORDER BY K.kills + K.losses + C.memberCount DESC, name ASC;
			`, wars[0], wars[1]); err != nil {
		return nil, err
	}

	return w, nil
}
