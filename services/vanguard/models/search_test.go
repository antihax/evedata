package models

import (
	"testing"
)

func TestSearchMarketNames(t *testing.T) {
	res, err := SearchMarketNames("warp disruptor ii")
	if err != nil {
		t.Error(err)
		return
	}

	if res[0].TypeName != "Heavy Warp Disruptor II" {
		t.Errorf("Wrong item returned %s", res[0].TypeName)
		return
	}
}
