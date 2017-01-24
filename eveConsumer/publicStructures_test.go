package eveConsumer

import "testing"

func TestStructureTrigger(t *testing.T) {
	_, err := structuresTrigger(eC)
	if err != nil {
		t.Error(err)
		return
	}
}
