package eveConsumer

import (
	"os"
	"testing"
	"time"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/models"
	"github.com/garyburd/redigo/redis"
)

var (
	ctx appContext.AppContext
	eC  *EVEConsumer
)

func TestMain(m *testing.M) {
	ctx = appContext.NewTestAppContext()
	eC = NewEVEConsumer(&ctx)

	// Create service states
	models.SetServiceState("wars", time.Now().UTC(), 1)
	models.SetServiceState("alliances", time.Now().UTC(), 1)
	models.SetServiceState("marketHistory", time.Now().UTC(), 1)
	models.SetServiceState("marketMaint", time.Now().UTC(), 1)
	models.SetServiceState("npcCorps", time.Now().UTC(), 1)
	models.SetServiceState("structures", time.Now().UTC(), 1)

	// Run the tests
	retCode := m.Run()
	os.Exit(retCode)
}

func TestWarsUpdate(t *testing.T) {
	err := eC.updateWars()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestWarsCheckCREST(t *testing.T) {
	err := eC.collectWarsFromCREST()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestWarsWarsPull(t *testing.T) {
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

func TestEntities(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	err := EntityAddToQueue(1, r)
	if err != nil {
		t.Error(err)
		return
	}
}

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

func TestEntitiesPull(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	for {
		err := eC.entityCheckQueue(r)
		if err != nil {
			t.Error(err)
			return
		}
		if i, _ := redis.Int(r.Do("SCARD", "EVEDATA_entityQueue")); i == 0 {
			break
		}
	}
}

func TestUpdateEntities(t *testing.T) {
	err := eC.updateEntities()
	if err != nil {
		t.Error(err)
		return
	}
}
