package models

import "testing"

func TestGetLocalIntel(t *testing.T) {

	var names []interface{}
	names = append(names, "dude")
	names = append(names, "Test Character")

	_, err := GetLocalIntel(names)
	if err != nil {
		t.Error(err)
	}

}
