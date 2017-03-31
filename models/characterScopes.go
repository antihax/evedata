package models

import "github.com/antihax/goesi"

type ScopeGroup struct {
	Scope string
	Group string
}

var characterScopes = []ScopeGroup{
	{"esi-assets.read_assets.v1", "assets"},

	{"esi-characters.read_contacts.v1", "contacts"},
	{"esi-characters.write_contacts.v1", "contacts"},
	{goesi.ScopeCharacterContactsRead, "contacts"},
	{goesi.ScopeCharacterContactsWrite, "contacts"},

	{goesi.ScopeCharacterMarketOrdersRead, "market"},
	{"esi-universe.read_structures.v1", "market"},
	{"esi-search.search_structures.v1", "market"},
	{"esi-markets.structure_markets.v1", "market"},

	{"esi-ui.open_window.v1", "ui-control"},
	{"esi-ui.write_waypoint.v1", "ui-control"},

	{goesi.ScopeCharacterWalletRead, "wallet"},
	{"esi-wallet.read_character_wallet.v1", "wallet"},
}

var groupReasons = map[string]string{
	"market":     "Reporting and Market tools",
	"contacts":   "War Contact Synchronization",
	"wallet":     "Profit and Loss Tools",
	"assets":     "Asset Value Tools",
	"ui-control": "Control of in-game UI",
}

func GetCharacterScopes() []string {
	var s []string
	for _, scope := range characterScopes {
		s = append(s, scope.Scope)
	}
	return s
}

func GetCharacterScopeGroups() map[string]string {
	return groupReasons
}

func GetCharacterScopesByGroups(groups []string) []string {
	var s []string
	for _, group := range groups {
		for _, scope := range characterScopes {
			if scope.Group == group {
				s = append(s, scope.Scope)
			}
		}
	}
	return s
}
