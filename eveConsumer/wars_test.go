package eveConsumer

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestWarsUpdate(t *testing.T) {
	err := eC.updateWars()
	if err != nil {
		t.Error(err)
		return
	}
}

// Temp disable as we have no CREST Mock
/*func TestWarsCheckCREST(t *testing.T) {
	err := eC.collectWarsFromCREST()
	if err != nil {
		t.Error(err)
		return
	}
}*/

func TestWarsPull(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	for {
		err := eC.warCheckQueue(r)
		if err != nil {
			t.Error(err)
			return
		}
		if i, _ := redis.Int(r.Do("SCARD", "EVEDATA_warQueue")); i == 0 {
			break
		}
	}
}
