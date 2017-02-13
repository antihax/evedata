package eveConsumer

import "testing"

func TestAssetTrigger(t *testing.T) {
	_, err := assetsTrigger(eC)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAssetConsumer(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	for {
		work, err := assetsConsumer(eC, &r)
		if err != nil {
			t.Error(err)
			return
		}
		if work == false {
			break
		}
	}
}
