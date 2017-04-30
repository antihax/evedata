package models

import "testing"

func TestGetSystemVertices(t *testing.T) {
	_, err := GetSystemVertices()
	if err != nil {
		t.Error(err)
		return
	}
}
