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

	err := UpdateCharacter(1001, "dude", 1, 1, 147035273, 0, 1, "male", -10, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	err = UpdateCharacter(1002, "dude 2", 1, 1, 147035273, 0, 2, "Female", -10, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	err = UpdateCharacter(1003, "dude 3", 1, 1, 147035273, 0, 2, "Female", -10, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	err = UpdateCorporationHistory(1001, 147035273, 100000222, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}
	err = UpdateCorporationHistory(1002, 147035273, 100000223, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}
	err = UpdateCorporationHistory(1003, 147035273, 100000224, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	err = BuildRelationships()
	if err != nil {
		t.Error(err)
	}
}
