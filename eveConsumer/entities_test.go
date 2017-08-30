package eveConsumer

import "testing"

func TestEntities(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	err := EntityCorporationAddToQueue(1, &r)
	if err != nil {
		t.Error(err)
		return
	}

	err = EntityCharacterAddToQueue(1, &r)
	if err != nil {
		t.Error(err)
		return
	}

	err = EntityAllianceAddToQueue(1, &r)
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
		work, err := corporationConsumer(eC, &r)
		if err != nil {
			t.Error(err)
			return
		}
		work, err = characterConsumer(eC, &r)
		if err != nil {
			t.Error(err)
			return
		}
		work, err = allianceConsumer(eC, &r)
		if err != nil {
			t.Error(err)
			return
		}
		if work == false {
			break
		}
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
