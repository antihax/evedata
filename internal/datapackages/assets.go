package datapackages

import (
	"github.com/antihax/goesi/esi"
)

type CharacterAssets struct {
	Assets           []esi.GetCharactersCharacterIdAssets200Ok
	CharacterID      int32
	TokenCharacterID int32
}
