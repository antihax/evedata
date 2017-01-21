package models

import "testing"

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
