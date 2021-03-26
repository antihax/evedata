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
	Character   esi.GetCharactersCharacterIdOk
	CharacterID int32
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
