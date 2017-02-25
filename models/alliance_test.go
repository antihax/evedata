package models

import (
	"testing"
	"time"
)

func TestAllianceAdd(t *testing.T) {
	err := UpdateAlliance(1, "Test Alliance Please Ignore", 10, "TEST", 4,
		time.Now().UTC(), time.Now().UTC())
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateCorporation(4, "Test Executor", "TEST2", 10,
		"Test Executor Corp", 1, 0, "somewhere", 50, time.Now().UTC())
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateCharacter(10, "Test Character", 1, 1, 4, 1, 1, "male", -10, time.Now().UTC())
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAllianceGet(t *testing.T) {
	alliance, err := GetAlliance(1)
	if err != nil {
		t.Error(err)
		return
	}
	if alliance.AllianceName != "Test Alliance Please Ignore" {
		t.Error("Could not find alliance 'Test Alliance Please Ignore'")
	}
}

func TestAllianceGetMembers(t *testing.T) {
	members, err := GetAllianceMembers(1)
	if err != nil {
		t.Error(err)
		return
	}
	if len(members) == 0 {
		t.Error("No members found")
	}
	if members[0].CorporationName != "Test Executor" {
		t.Error("Could not find member 'Test Executor'")
	}
}

func TestCorporationGet(t *testing.T) {
	corporation, err := GetCorporation(4)
	if err != nil {
		t.Error(err)
		return
	}
	if corporation.CorporationName != "Test Executor" {
		t.Error("Could not find corporation 'Test Executor'")
	}
}

func TestCharacterGet(t *testing.T) {
	character, err := GetCharacter(10)
	if err != nil {
		t.Error(err)
		return
	}
	if character.CharacterName != "Test Character" {
		t.Error("Could not find corporation 'Test Character'")
	}
}
