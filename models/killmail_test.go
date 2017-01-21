package models

import (
	"testing"
	"time"
)

func TestAddKillmail(t *testing.T) {
	err := AddKillmail(1, 1, time.Now(), 1, 1, 1, "FAKEHASH", 1, 2323, 1.0, 1.0, 1.0, 10, 1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAddKillmailAttacker(t *testing.T) {
	err := AddKillmailAttacker(1, 2, 1, 1, 1, true, 2323, 1, -1.34)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAddKillmailItems(t *testing.T) {
	err := AddKillmailItems(1, 1, 1, 22, 23, 0)
	if err != nil {
		t.Error(err)
		return
	}
}

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
