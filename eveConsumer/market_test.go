package eveConsumer

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
)

func TestMarketAddRegion(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	eC.marketRegionAddRegion(1, time.Now().UTC().Unix(), r)
}

func TestMarketRegionCheck(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	err := eC.marketRegionCheckQueue(r)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMarketHistoryUpdateTrigger(t *testing.T) {
	err := eC.marketHistoryUpdateTrigger()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMarketOrderPull(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	for {
		err := eC.marketOrderCheckQueue(r)
		if err != nil {
			t.Error(err)
			return
		}
		if i, _ := redis.Int(r.Do("SCARD", "EVEDATA_marketOrders")); i == 0 {
			break
		}
	}
}

// This is bugged due to the ESI Spec.
func TestMarketHistoryPull(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	j := 0
	for {
		j++
		err := eC.marketHistoryCheckQueue(r)
		if err != nil {
			t.Error(err)
			return
		}
		if i, _ := redis.Int(r.Do("SCARD", "EVEDATA_marketHistory")); i == 0 || j > 5 {
			break
		}
	}
}

func TestMarketMaintTrigger(t *testing.T) {
	err := eC.marketMaintTrigger()
	if err != nil {
		t.Error(err)
		return
	}
}
