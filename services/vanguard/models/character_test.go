package models

import (
	"log"
	"testing"
	"time"

	"golang.org/x/oauth2"
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

	err = UpdateCharacter(1002, "dude 2", 1, 1, 147035273, 0, 2, "Female", -10, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	err = UpdateCorporationHistory(1001, 147035273, 100000222, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	hist, err := GetCorporationHistory(1001)
	if err != nil {
		log.Fatal(err)
		return
	}

	if hist[0].CorporationID != 147035273 {
		log.Fatal("wrong corporationID in history")
	}
}

func TestCursor(t *testing.T) {
	err := SetCursorCharacter(1001, 1001)
	if err != nil {
		log.Fatal(err)
		return
	}

	cursor, err := GetCursorCharacter(1001)
	if err != nil {
		log.Fatal(err)
		return
	}

	if cursor.CursorCharacterID != 1001 {
		log.Fatal("Wrong cursor returned")
	}
}

func TestAddCRESTToken(t *testing.T) {
	tok := oauth2.Token{
		AccessToken:  "FAKE",
		RefreshToken: "So Fake",
		Expiry:       time.Now().UTC().Add(time.Hour * 100000),
		TokenType:    "Bearer",
	}

	err := AddCRESTToken(1, 1, "Dude", &tok, "", "ownerhash", 0, 0, 0)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAddIntegrationToken(t *testing.T) {
	tok := oauth2.Token{
		AccessToken:  "FAKE",
		RefreshToken: "So Fake",
		Expiry:       time.Now().UTC().Add(time.Hour * 100000),
		TokenType:    "Bearer",
	}

	if err := AddIntegrationToken("mock", 1, "1232142345345", "duuuuuuuude", &tok, "identify"); err != nil {
		t.Error(err)
	}
}

func TestGetIntegrationTokens(t *testing.T) {
	dude, err := GetIntegrationTokens(1)
	if err != nil {
		t.Error(err)
		return
	}
	if dude[0].IntegrationUserID != "1232142345345" {
		t.Error("IntegrationUserID is not as stored")
		return
	}
}
func TestDeleteIntegrationToken(t *testing.T) {
	err := DeleteIntegrationToken("mock", 1, "1232142345345")
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
