package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAllianceAdd(t *testing.T) {
	err := UpdateAlliance(2, "Test Alliance Please Ignore", 10, "TEST", 4,
		time.Now().UTC(), time.Now().UTC())
	assert.Nil(t, err)
	err = UpdateCorporation(4, "Test Executor", "TEST2", 10,
		1, 0, 50, time.Now().UTC())
	assert.Nil(t, err)
	err = UpdateCharacter(10, "Test Character", 1, 1, 4, 1, 1, "male", -10, time.Now().UTC())
	assert.Nil(t, err)
}

func TestAllianceGet(t *testing.T) {
	alliance, err := GetAlliance(2)
	assert.Nil(t, err)
	if alliance.AllianceName != "Test Alliance Please Ignore" {
		t.Error("Could not find alliance 'Test Alliance Please Ignore'")
	}
}

func TestAllianceGetMembers(t *testing.T) {
	members, err := GetAllianceMembers(1)
	assert.Nil(t, err)
	if len(members) == 0 {
		t.Error("No members found")
	}
	if members[0].CorporationName != "Test Executor" {
		t.Error("Could not find member 'Test Executor'")
	}
}

func TestCorporationGet(t *testing.T) {
	corporation, err := GetCorporation(4)
	assert.Nil(t, err)
	if corporation.CorporationName != "Test Executor" {
		t.Error("Could not find corporation 'Test Executor'")
	}
}

func TestCharacterGet(t *testing.T) {
	character, err := GetCharacter(10)
	assert.Nil(t, err)
	if character.CharacterName != "Test Character" {
		t.Error("Could not find corporation 'Test Character'")
	}
}

func TestGetAllianceHistory(t *testing.T) {
	_, err := GetAllianceHistory(2)
	assert.Nil(t, err)
}

func TestGetAllianceJoinHistory(t *testing.T) {
	_, err := GetAllianceJoinHistory(2)
	assert.Nil(t, err)
}
