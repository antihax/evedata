package models

import "testing"

func TestAssetSetup(t *testing.T) {
	database.Exec(`
			INSERT INTO evedata.assets (10, 44, 2, 1001, "station", 1, "somewhere", 0);
			INSERT INTO evedata.assets (11, 44, 2, 1001, "station", 10, "somewhere", 0);		
			INSERT INTO evedata.assets (12, 44, 2, 1001, "station", 11, "somewhere", 0);`,
	)
}

func TestGetAssetLocations(t *testing.T) {
	_, err := GetAssetLocations(1, 1)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = GetAssetLocations(1, 0)
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
	_, err := GetAssets(1, 1, 1)
	if err != nil {
		t.Error(err)
		return
	}
}
