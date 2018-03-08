package models

import "testing"

func TestAddShare(t *testing.T) {
	err := AddShare(1, 1, 1, "kill")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetShares(t *testing.T) {
	_, err := GetShares(2)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetLocatorResponses(t *testing.T) {
	_, err := GetLocatorResponses(1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteShare(t *testing.T) {
	err := DeleteShare(1, 1, 1)
	if err != nil {
		t.Error(err)
		return
	}
}
