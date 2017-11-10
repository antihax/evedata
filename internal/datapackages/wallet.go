package datapackages

import "github.com/antihax/goesi/esi"

type CharacterWalletTransactions struct {
	Transcations     []esi.GetCharactersCharacterIdWalletTransactions200Ok
	CharacterID      int32
	TokenCharacterID int32
}

type CharacterJournal struct {
	Journal          []esi.GetCharactersCharacterIdWalletJournal200Ok
	CharacterID      int32
	TokenCharacterID int32
}
