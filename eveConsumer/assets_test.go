package eveConsumer

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestAssetPull(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	for {
		err := eC.assetsCheckQueue(r)
		if err != nil {
			t.Error(err)
			return
		}
		if i, _ := redis.Int(r.Do("SCARD", "EVEDATA_assetQueue")); i == 0 {
			break
		}
	}
}
