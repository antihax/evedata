package eveConsumer

import "testing"

func TestStructureTrigger(t *testing.T) {
	_, err := structuresTrigger(eC)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestStructureConsumer(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	for {
		work, err := structureConsumer(eC, &r)
		if err != nil {
			t.Error(err)
			return
		}
		if work == false {
			break
		}
	}
}
