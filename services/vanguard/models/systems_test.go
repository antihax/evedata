package models

import "testing"

func TestGetSystemVertices(t *testing.T) {
	_, err := GetSystemVertices()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetSystemCelestials(t *testing.T) {
	_, err := GetSystemCelestials(30003280)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetNullSystems(t *testing.T) {
	_, err := GetNullSystems()
	if err != nil {
		t.Error(err)
		return
	}
}
