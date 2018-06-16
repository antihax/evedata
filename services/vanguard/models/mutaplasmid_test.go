package models

import "testing"

func TestGetMutaplasmidData(t *testing.T) {
	_, err := GetMutaplasmidData("Warp Disruptor")
	if err != nil {
		if err.Error() != "no data found" {
			t.Error(err)
			return
		}
	}
}
