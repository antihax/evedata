package models

import "testing"

func TestAddCRESTRef(t *testing.T) {
	err := AddCRESTRef(1000005, "https://crest-tq.eveonline.com/corporations/1000005/")
	if err != nil {
		t.Error(err)
		return
	}
	err = AddCRESTRef(90000051, "https://crest-tq.eveonline.com/characters/90000051/")
	if err != nil {
		t.Error(err)
		return
	}
	err = AddCRESTRef(99000001, "https://crest-tq.eveonline.com/alliances/99000001/")
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
