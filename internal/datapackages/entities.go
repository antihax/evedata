package datapackages

import "github.com/antihax/goesi/esi"

// Alliance contains alliance information and alliance corporations
type Alliance struct {
	Alliance             esi.GetAlliancesAllianceIdOk
	AllianceCorporations []int32
	AllianceID           int32
	ETag                 string
}

// Character contains corp history and character information
type Character struct {
	Character   esi.GetCharactersCharacterIdOk
	CharacterID int32
	ETag        string
}

type CorporationHistory struct {
	CharacterID        int32
	CorporationHistory []esi.GetCharactersCharacterIdCorporationhistory200Ok
	ETag               string
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
	ETag          string
}

type AllianceHistory struct {
	CorporationID   int32
	AllianceHistory []esi.GetCorporationsCorporationIdAlliancehistory200Ok
	ETag            string
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

// AllianceContacts contains AllianceContacts data
type AllianceContacts struct {
	AllianceID int32
	Contacts   []esi.GetAlliancesAllianceIdContacts200Ok
}

// CorporationContacts contains CorporationContacts data
type CorporationContacts struct {
	CorporationID int32
	Contacts      []esi.GetCorporationsCorporationIdContacts200Ok
}
