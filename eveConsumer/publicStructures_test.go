package eveConsumer

import "testing"

func TestStructureCheck(t *testing.T) {
	err := eC.collectStructuresFromESI()
	if err != nil {
		t.Error(err)
		return
	}
}
