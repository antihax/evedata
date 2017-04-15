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

	if res[0].TypeName != "Warp Disruptor II" {
		t.Error("Wrong item returned")
		return
	}
}
