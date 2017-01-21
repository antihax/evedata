package models

import "testing"

func TestAddCRESTRef(t *testing.T) {
	err := AddCRESTRef(1000005, "https://crest-tq.eveonline.com/corporations/1000005/")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetCRESTRef(t *testing.T) {
	ref, err := GetCRESTRef(1000005)
	if err != nil {
		t.Error(err)
		return
	}

	if ref.CrestRef != "https://crest-tq.eveonline.com/corporations/1000005/" {
		t.Error("CrestRef does not match")
		return
	}
}
