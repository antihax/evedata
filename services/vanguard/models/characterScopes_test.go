package models

import "testing"

func TestGetCharacterScopes(t *testing.T) {
	scopes := GetCharacterScopes()

	if len(scopes) != len(characterScopes) {
		t.Errorf("Scope count does not match")
		return
	}
}

func TestGetCharacterScopeGroups(t *testing.T) {
	scopeGroups := GetCharacterScopeGroups()

	if len(scopeGroups) != len(groupReasons) {
		t.Errorf("Scope count does not match")
		return
	}
}

func TestGetCharacterScopesByGroups(t *testing.T) {
	scopes := GetCharacterScopesByGroups([]string{"market"})

	expected := 0
	for _, scope := range characterScopes {
		if scope.Group == "market" {
			expected++
		}
	}

	if len(scopes) != expected {
		t.Errorf("Scope count does not match")
		return
	}
}
