package eveConsumer

import (
	"testing"

	"github.com/antihax/evedata/models"
	"github.com/garyburd/redigo/redis"
)

func TestContactSyncTrigger(t *testing.T) {
	err := contactSyncTrigger(eC)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestContactSyncConsumer(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()

	// Add a fake contact sync to the characters created above.
	err := models.AddContactSync(1001, 1001, 1002)
	if err != nil {
		t.Error(err)
		return
	}

	for {
		err := contactSyncConsumer(eC, r)
		if err != nil {
			t.Error(err)
			return
		}
		if i, _ := redis.Int(r.Do("SCARD", "EVEDATA_contactSyncQueue")); i == 0 {
			break
		}
	}
}
