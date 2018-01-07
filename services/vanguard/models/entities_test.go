package models

import "testing"

func TestGetEntityName(t *testing.T) {
	_, err := GetEntityName(2)
	if err != nil {
		t.Error(err)
		return
	}
}
func TestGetTypeName(t *testing.T) {
	_, err := GetTypeName(2)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetSystemName(t *testing.T) {
	_, err := GetSystemName(30000001)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetCelestialName(t *testing.T) {
	_, err := GetCelestialName(30000001)
	if err != nil {
		t.Error(err)
		return
	}
}
