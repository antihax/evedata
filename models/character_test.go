package models

import (
	"log"
	"testing"
	"time"

	"github.com/antihax/goesi"
)

func TestUpdateCharacter(t *testing.T) {
	err := UpdateCorporation(147035273, "Dude Corp", "TEST2", 10,
		0, 0, 50, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	err = UpdateCharacter(1001, "dude", 1, 1, 147035273, 0, 1, "male", -10, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}
}

func TestAddCRESTToken(t *testing.T) {
	tok := goesi.CRESTToken{
		AccessToken:  "FAKE",
		RefreshToken: "So Fake",
		Expiry:       time.Now().UTC().Add(time.Hour * 100000),
		TokenType:    "Bearer"}

	err := AddCRESTToken(1, 1, "Dude", &tok, "")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetCRESTTokens(t *testing.T) {
	dude, err := GetCRESTTokens(1)
	if err != nil {
		t.Error(err)
		return
	}
	if dude[0].CharacterName != "Dude" {
		t.Error("Character name is not as stored")
		return
	}
}

func TestGetCRESTToken(t *testing.T) {
	dude, err := GetCRESTToken(1, 1)
	if err != nil {
		t.Error(err)
		return
	}
	if dude.CharacterName != "Dude" {
		t.Error("Character name is not as stored")
		return
	}
}

func TestSetTokenError(t *testing.T) {
	err := SetTokenError(1, 1, 200, "OK", nil, nil)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteCRESTToken(t *testing.T) {
	err := DeleteCRESTToken(1, 1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetCharacter(t *testing.T) {
	char, err := GetCharacter(1001)
	if err != nil {
		t.Error(err)
		return
	}
	if char.CorporationID != 147035273 {
		t.Error("Character corporationID does not match")
		return
	}
}

func TestGetCharacterIDByName(t *testing.T) {
	char, err := GetCharacterIDByName("dude")
	if err != nil {
		t.Error(err)
		return
	}
	if char != 1001 {
		t.Error("CharacterID does not match")
		return
	}
}
