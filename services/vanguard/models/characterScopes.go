package models

import (
	"sort"
	"strings"
)

type ScopeGroup struct {
	Scope string
	Group string
}

var characterScopes = []ScopeGroup{
	{"esi-assets.read_assets.v1", "assets"},

	{"esi-characters.read_contacts.v1", "contacts"},
	{"esi-characters.write_contacts.v1", "contacts"},

	{"esi-characters.read_notifications.v1", "notifications"},

	{"esi-universe.read_structures.v1", "market"},
	{"esi-search.search_structures.v1", "market"},
	{"esi-markets.structure_markets.v1", "market"},

	{"esi-ui.open_window.v1", "ui-control"},
	{"esi-ui.write_waypoint.v1", "ui-control"},

	{"esi-wallet.read_character_wallet.v1", "wallet"},

	{"esi-characters.read_corporation_roles.v1", "roles"},
	{"esi-characters.read_titles.v1", "roles"},
	{"esi-alliances.read_contacts.v1", "roles"},
	{"esi-corporations.read_contacts.v1", "roles"},
}

var groupReasons = map[string]string{
	"market":        "Reporting and Market tools",
	"contacts":      "War Contact Synchronization",
	"wallet":        "Profit and Loss Tools",
	"assets":        "Asset Value Tools",
	"ui-control":    "Control of in-game UI",
	"notifications": "Notification tools (locators, structures, wars)",
	"roles":         "Corporation services based on corp roles",
}

// shareReasons for data shares between characters and entities
var shareReasons = map[string]string{
	"application": "Corporation Applications",
	"locator":     "Locator Responses",
	"structure":   "Corporation structures under attack",
	"war":         "War Declared on Corporation",
}

func GetCharacterShareGroups() map[string]string {
	return shareReasons
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

// GetCharacterGroupsByScopesString takes a space seperated string of scopes and returns the groups
func GetCharacterGroupsByScopesString(scopes string) string {
	groups := make(map[string]bool)
	for _, scope := range strings.Split(scopes, " ") {
		for _, charScope := range characterScopes {
			if charScope.Scope == scope {
				groups[charScope.Group] = true
			}
		}
	}

	m := make([]string, len(groups))
	i := 0
	for k := range groups {
		m[i] = k
		i++
	}
	sort.Strings(m)
	return strings.Join(m, ", ")
}
