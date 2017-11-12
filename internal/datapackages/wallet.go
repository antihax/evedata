package datapackages

import (
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/eveapi"
)

type CharacterWalletTransactions struct {
	Transactions     []esi.GetCharactersCharacterIdWalletTransactions200Ok
	CharacterID      int32
	TokenCharacterID int32
}

type CharacterJournal struct {
	Journal          eveapi.WalletJournalXML
	CharacterID      int32
	TokenCharacterID int32
}
