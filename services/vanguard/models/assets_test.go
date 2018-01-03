package models

import (
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestAssetSetup(t *testing.T) {
	_, err := database.Exec(`
			INSERT IGNORE INTO evedata.assets VALUES
			 (60012526, 39, 2, 1001, "station", 21, "station", 0)
			,(60012526, 179, 2, 1001, "station", 22, "station", 0)
			,(60012526, 182, 2, 1001, "station", 23, "station", 0)
			,(23, 196, 2, 1001, "other", 24, "other", 0)
			,(24, 199, 2, 1001, "other", 25, "other", 0)
			,(25, 199, 2, 1001, "other", 26, "other", 0)
			,(26, 199, 2, 1001, "other", 27, "other", 0)
			,(23, 39, 2, 1001, "other", 28, "other", 0);
		`)
	if err != nil {
		t.Error(err)
		return
	}

	tok := oauth2.Token{
		AccessToken:  "FAKE",
		RefreshToken: "So Fake",
		Expiry:       time.Now().UTC().Add(time.Hour * 100000),
		TokenType:    "Bearer"}

	err = AddCRESTToken(1001, 1001, "Dude", &tok, "")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetAssetLocations(t *testing.T) {
	_, err := GetAssetLocations(1001, 1)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = GetAssetLocations(1001, 0)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetAssetCharacters(t *testing.T) {
	_, err := GetAssetCharacters(1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetAsset(t *testing.T) {
	data, err := GetAssets(1001, 0, 60012526)
	if err != nil {
		t.Error(err)
		return
	}

	if data[0].CharacterID != 1001 {
		t.Error("wrong character ID")
		return
	}

	subs := false
	for i := range data {
		if data[i].SubCount > 0 {
			subs = true
		}
	}

	if !subs {
		t.Error("failed recursion. should have some sub items.")
		return
	}

	data, err = GetAssets(1001, 1001, 60012526)
	if err != nil {
		t.Error(err)
		return
	}

	if data[0].CharacterID != 1001 {
		t.Error("wrong character ID")
		return
	}

	subs = false
	for i := range data {
		if data[i].SubCount > 0 {
			subs = true
		}
	}

	if !subs {
		t.Error("failed recursion. should have some sub items.")
		return
	}
}
