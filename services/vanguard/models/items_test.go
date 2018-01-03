package models

import "testing"

func TestGetItem(t *testing.T) {
	item, err := GetItem(22)
	if err != nil {
		t.Error(err)
		return
	}
	if item.TypeName != "Arkonor" {
		t.Error("Could not get the correct item")
		return
	}
}

func TestGetItemAttributes(t *testing.T) {
	item, err := GetItemAttributes(22)
	if err != nil {
		t.Error(err)
		return
	}
	if len(*item) == 0 {
		t.Error("Could not get attributes")
		return
	}
}
