package models

import "testing"

func TestGetActiveWarsByID(t *testing.T) {
	_, err := GetActiveWarsByID(1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetPendingWarsByID(t *testing.T) {
	_, err := GetPendingWarsByID(1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetFinishedWarsByID(t *testing.T) {
	_, err := GetFinishedWarsByID(1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetActiveWarList(t *testing.T) {
	_, err := GetActiveWarList()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetWarsForEntityByID(t *testing.T) {
	_, err := GetWarsForEntityByID(1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetKnownAlliesByID(t *testing.T) {
	_, err := GetKnownAlliesByID(1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetFactionWarEntitiesForID(t *testing.T) {
	_, err := GetFactionWarEntitiesForID(FactionsByName["Caldari"])
	if err != nil {
		t.Error(err)
		return
	}
	_, err = GetFactionWarEntitiesForID(FactionsByName["Minmatar"])
	if err != nil {
		t.Error(err)
		return
	}
	_, err = GetFactionWarEntitiesForID(FactionsByName["Amarr"])
	if err != nil {
		t.Error(err)
		return
	}
	_, err = GetFactionWarEntitiesForID(FactionsByName["Gallente"])
	if err != nil {
		t.Error(err)
		return
	}
}
