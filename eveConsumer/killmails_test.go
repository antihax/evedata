package eveConsumer

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestKillmailsPull(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	for {
		err := eC.killmailCheckQueue(r)
		if err != nil {
			t.Error(err)
			return
		}
		if i, _ := redis.Int(r.Do("SCARD", "EVEDATA_killQueue")); i == 0 {
			break
		}
	}
}
