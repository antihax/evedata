package eveConsumer

import (
	"testing"

	"github.com/antihax/evedata/models"
)

func TestContactSyncTrigger(t *testing.T) {
	_, err := contactSyncTrigger(eC)
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
		work, err := contactSyncConsumer(eC, r)
		if err != nil {
			t.Error(err)
			return
		}
		if work == false {
			break
		}
	}
}
