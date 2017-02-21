package models

import (
	"testing"
	"time"

	"github.com/antihax/goesi"
)

func TestAssetSetup(t *testing.T) {
	database.Exec(`
			INSERT INTO evedata.assets (60012526, 1373, 2, 1001, "station", 11, "somewhere", 0);
			INSERT INTO evedata.assets (60012526, 1373, 2, 1001, "station", 12, "somewhere", 0);		
			INSERT INTO evedata.assets (60012526, 1373, 2, 1001, "station", 13, "somewhere", 0);
			INSERT INTO evedata.assets (13, 1373, 2, 1001, "other", 14, "somewhere", 0);
			INSERT INTO evedata.assets (14, 1373, 2, 1001, "other", 15, "somewhere", 0);
		`)

	tok := goesi.CRESTToken{
		AccessToken:  "FAKE",
		RefreshToken: "So Fake",
		Expiry:       time.Now().UTC().Add(time.Hour * 100000),
		TokenType:    "Bearer"}

	err := AddCRESTToken(1001, 1001, "Dude", &tok, "")
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
	_, err := GetAssets(1001, 0, 0)
	if err != nil {
		t.Error(err)
		return
	}
}
