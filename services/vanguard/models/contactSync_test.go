package models

import (
	"testing"
	"time"
)

func TestAddContactSync(t *testing.T) {
	err := AddContactSync(1, 1, 2)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetContactSyncs(t *testing.T) {
	syncs, err := GetContactSyncs(1)
	if err != nil {
		t.Error(err)
		return
	}
	err = syncs[0].Updated(time.Now().UTC())
	if err != nil {
		t.Error(err)
		return
	}
	err = syncs[0].Error("It broke")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteContactSync(t *testing.T) {
	err := DeleteContactSync(1, 1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetExpiredContactSyncs(t *testing.T) {
	_, err := GetExpiredContactSyncs()
	if err != nil {
		t.Error(err)
		return
	}
}
