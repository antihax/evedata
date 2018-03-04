package datapackages

import "github.com/antihax/goesi/esi"

// Alliance contains alliance information and alliance corporations
type Alliance struct {
	Alliance             esi.GetAlliancesAllianceIdOk
	AllianceCorporations []int32
	AllianceID           int32
}

// Character contains corp history and character information
type Character struct {
	Character          esi.GetCharactersCharacterIdOk
	CorporationHistory []esi.GetCharactersCharacterIdCorporationhistory200Ok
	CharacterID        int32
}

// Character contains corp history and character information
type CharacterRoles struct {
	Roles            esi.GetCharactersCharacterIdRolesOk
	CharacterID      int32
	TokenCharacterID int32
}

// Corporation contains corp history and character information
type Corporation struct {
	CorporationID int32
	Corporation   esi.GetCorporationsCorporationIdOk
}

// Store contains NPC Corp loyalty point offers
type Store struct {
	CorporationID int32
	Store         []esi.GetLoyaltyStoresCorporationIdOffers200Ok
}

// Killmail contains Killmail data
type Killmail struct {
	Hash string
	Kill esi.GetKillmailsKillmailIdKillmailHashOk
}
