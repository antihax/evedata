package models

import "testing"

func TestMaintKillMails(t *testing.T) {
	err := MaintKillMails()
	if err != nil {
		t.Error(err)
	}
}

func TestMaintMarket(t *testing.T) {
	err := MaintMarket()
	if err != nil {
		t.Error(err)
	}
}

func TestMaintContactSync(t *testing.T) {
	err := MaintContactSync()
	if err != nil {
		t.Error(err)
	}
}

func TestMaintOrphanCharacters(t *testing.T) {
	_, err := MaintOrphanCharacters()
	if err != nil {
		t.Error(err)
	}
}
