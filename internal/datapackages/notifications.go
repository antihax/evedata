package datapackages

import "github.com/antihax/goesi/esi"

type CharacterNotifications struct {
	Notifications    []esi.GetCharactersCharacterIdNotifications200Ok
	CharacterID      int32
	TokenCharacterID int32
}
