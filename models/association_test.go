package models

import (
	"testing"
)

func TestGetCharacterKnownAssociates(t *testing.T) {

	_, err := database.Exec(`
			INSERT IGNORE INTO evedata.characterAssociations VALUES
			 (1001, 1002,3);
		`)
	if err != nil {
		t.Error(err)
		return
	}

	alts, err := GetCharacterKnownAssociates(1001)
	if err != nil {
		t.Error(err)
		return
	}

	if alts[0].CharacterID != 1002 {
		t.Error("CharacterID does not match")
		return
	}
}

func TestGetCharacterKnownKillmailAssociates(t *testing.T) {

	_, err := database.Exec(`
			INSERT IGNORE INTO evedata.characterAssociations VALUES
			 (1001, 1002,3);
		`)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = GetCharacterKnownKillmailAssociates(1001)
	if err != nil {
		t.Error(err)
		return
	}
}
