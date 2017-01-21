package models

import "testing"

func TestGetAllianceAssetsInSpace(t *testing.T) {
	_, err := GetAllianceAssetsInSpace(1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetCorporationAssetsInSpace(t *testing.T) {
	_, err := GetCorporationAssetsInSpace(1)
	if err != nil {
		t.Error(err)
		return
	}
}
