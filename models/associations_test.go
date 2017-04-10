package models

import (
	"testing"
)

func TestGetCharacterKnownAssociates(t *testing.T) {
	_, err := database.Exec(`
			INSERT IGNORE INTO evedata.characterAssociations VALUES
			 (1001, 1002, 3, 1, UTC_TIMESTAMP);
		`)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = GetCharacterKnownAssociates(1001)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetCorporationKnownAssociates(t *testing.T) {
	_, err := GetCorporationKnownAssociates(1001)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetAllianceKnownAssociates(t *testing.T) {
	_, err := GetAllianceKnownAssociates(1001)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestBuildRelationships(t *testing.T) {
	err := BuildRelationships()
	if err != nil {
		t.Error(err)
	}
}
