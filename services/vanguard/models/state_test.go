package models

import (
	"testing"
	"time"
)

func TestSetServiceState(t *testing.T) {
	err := SetServiceState("testState", time.Now().UTC(), 1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetServiceStateByDays(t *testing.T) {
	err := SetServiceStateByDays("testStateDays", 0, 1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetServiceState(t *testing.T) {
	_, i, err := GetServiceState("testState")
	if err != nil {
		t.Error(err)
		return
	}
	if i != 1 {
		t.Error("testState incorrect")
		return
	}

	_, i, err = GetServiceState("testStateDays")
	if err != nil {
		t.Error(err)
		return
	}
	if i != 1 {
		t.Error("testStateDays incorrect")
		return
	}
}
