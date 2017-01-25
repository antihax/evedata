package eveConsumer

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestKillmailsConsumer(t *testing.T) {
	err := killmailAddToQueue(1, "FAKE HASH")
	if err != nil {
		t.Error(err)
		return
	}
	err = killmailAddToQueue(2, "FAKE HASH")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestKillmailsConsumer(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	for {
		_, err := killmailsConsumer(eC, r)
		if err != nil {
			t.Error(err)
			return
		}
		if i, _ := redis.Int(r.Do("SCARD", "EVEDATA_killQueue")); i == 0 {
			break
		}
	}
}
