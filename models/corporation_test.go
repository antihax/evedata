package models

import (
	"log"
	"testing"
	"time"
)

func TestUpdateCorporation(t *testing.T) {
	err := UpdateCorporation(147035273, "Dude Corp", "TEST2", 10,
		0, 0, "somewhere", 50, time.Now().UTC())
	if err != nil {
		log.Fatal(err)
		return
	}
}

func TestGetCorporation(t *testing.T) {
	corp, err := GetCorporation(147035273)
	if err != nil {
		t.Error(err)
		return
	}
	if corp.MemberCount != 50 {
		t.Error("corporation memberCount does not match")
		return
	}
}
