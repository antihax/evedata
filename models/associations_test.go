package models

import (
	"log"
	"testing"
	"time"
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
	err := UpdateCorporationHistory(1001, 147035273, 100000222, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}
	err = UpdateCorporationHistory(1002, 147035273, 100000222, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}
	err = UpdateCorporationHistory(1003, 147035273, 100000222, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	err = BuildRelationships()
	if err != nil {
		t.Error(err)
	}
}
