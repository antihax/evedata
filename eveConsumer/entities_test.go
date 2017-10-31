package eveConsumer

import "testing"

func TestEntities(t *testing.T) {
	err := EntityCorporationAddToQueue(1)
	if err != nil {
		t.Error(err)
		return
	}

	err = EntityCharacterAddToQueue(1)
	if err != nil {
		t.Error(err)
		return
	}

	err = EntityAllianceAddToQueue(1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestCharSearchConsumer(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()

	CharSearchAddToQueue([]interface{}{"croakroach", "some other dude"}, &r)

	for {
		work, err := charSearchConsumer(eC, &r)
		if err != nil {
			t.Error(err)
			return
		}
		if work == false {
			break
		}
	}
}
