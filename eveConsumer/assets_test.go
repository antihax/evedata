package eveConsumer

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestAssetTrigger(t *testing.T) {
	err := assetsTrigger(eC)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAssetConsumer(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	for {
		err := assetsConsumer(eC, r)
		if err != nil {
			t.Error(err)
			return
		}
		if i, _ := redis.Int(r.Do("SCARD", "EVEDATA_assetQueue")); i == 0 {
			break
		}
	}
}
