package datapackages

import (
	"github.com/antihax/goesi/esi"
)

type CharacterOrders struct {
	Orders           []esi.GetCharactersCharacterIdOrders200Ok
	CharacterID      int32
	TokenCharacterID int32
}

type CharacterWalletTransactions struct {
	Transactions     []esi.GetCharactersCharacterIdWalletTransactions200Ok
	CharacterID      int32
	TokenCharacterID int32
}

type CharacterJournal struct {
	Journal          []esi.GetCharactersCharacterIdWalletJournal200Ok
	CharacterID      int32
	TokenCharacterID int32
}
