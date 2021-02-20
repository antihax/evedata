package models

import (
	"testing"
)

func TestGetKnownKillmails(t *testing.T) {
	mails, err := GetKnownKillmails()
	if err != nil {
		t.Error(err)
		return
	}
	if len(mails) == 0 {
		t.Error("No killmail ids returned")
		return
	}
}

func TestGetConstellationActivity(t *testing.T) {
	_, err := GetConstellationActivity(1, "character")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = GetConstellationActivity(1, "corporation")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = GetConstellationActivity(1, "alliance")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetKnownShipTypes(t *testing.T) {
	_, err := GetKnownShipTypes(1, "character")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = GetKnownShipTypes(1, "corporation")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = GetKnownShipTypes(1, "alliance")
	if err != nil {
		t.Error(err)
		return
	}
}
