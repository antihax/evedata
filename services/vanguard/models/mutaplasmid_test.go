package models

import "testing"

func TestGetMutaplasmidData(t *testing.T) {
	_, err := GetMutaplasmidData("Warp Disruptor")
	if err != nil {
		t.Error(err)
		return
	}
}
