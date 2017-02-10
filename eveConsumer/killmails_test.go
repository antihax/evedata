package eveConsumer

import "testing"

func TestKillmailsAddToQueue(t *testing.T) {
	err := eC.killmailAddToQueue(1, "FAKE HASH")
	if err != nil {
		t.Error(err)
		return
	}
	err = eC.killmailAddToQueue(2, "FAKE HASH")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestKillmailsConsumer(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	for {
		work, err := killmailsConsumer(eC, r)
		if err != nil {
			t.Error(err)
			return
		}
		if work == false {
			break
		}
	}
}
