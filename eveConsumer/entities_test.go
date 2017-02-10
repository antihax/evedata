package eveConsumer

import "testing"

func TestEntities(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	err := EntityAddToQueue(1, r)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestEntitiesTrigger(t *testing.T) {
	_, err := entitiesTrigger(eC)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestEntitiesConsumer(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	for {
		work, err := entitiesConsumer(eC, r)
		if err != nil {
			t.Error(err)
			return
		}
		if work == false {
			break
		}
	}
}
