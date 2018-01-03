package models

import "testing"

func TestWalletSetup(t *testing.T) {
	database.Exec(``)
}

func TestGetWalletSummary(t *testing.T) {
	_, err := GetWalletSummary(1, 0)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = GetWalletSummary(1, 1)
	if err != nil {
		t.Error(err)
		return
	}
}
