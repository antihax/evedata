package eveConsumer

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestEntitiesFromCrest(t *testing.T) {
	err := eC.collectEntitiesFromCREST()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestEntities(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	err := EntityAddToQueue(1, r)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestEntitiesPull(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	for {
		err := eC.entityCheckQueue(r)
		if err != nil {
			t.Error(err)
			return
		}
		if i, _ := redis.Int(r.Do("SCARD", "EVEDATA_entityQueue")); i == 0 {
			break
		}
	}
}

func TestUpdateEntities(t *testing.T) {
	err := eC.updateEntities()
	if err != nil {
		t.Error(err)
		return
	}
}
