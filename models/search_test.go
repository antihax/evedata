package models

import (
	"fmt"
	"testing"
)

func TestSearchMarketNames(t *testing.T) {
	res, err := SearchMarketNames("warp disruptor ii")
	if err != nil {
		t.Error(err)
		return
	}

	if res[0].TypeName != "Heavy Warp Disruptor II" {
		fmt.Printf("%+v\n", res[0])
		t.Error("Wrong item returned")
		return
	}
}
