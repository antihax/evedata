package eveConsumer

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/eveapi"
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

func TestMarketAddRegion(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	eC.marketRegionAddRegion(1, time.Now().UTC().Unix(), r)
}

// Setup some dummy scopes to test authenticated ESI stuff.
// Expiration must be set well past now to prevent
// accessing a nil authentiticator
func TestScopeSetup(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()

	scopes := []string{
		eveapi.ScopeCharacterContractsRead,
		eveapi.ScopeCharacterMarketOrdersRead,
		eveapi.ScopeCharacterResearchRead,
		eveapi.ScopeCharacterWalletRead,
		"esi-assets.read_assets.v1",
		"esi-characters.read_contacts.v1",
		"esi-characters.write_contacts.v1",
		"esi-corporations.read_corporation_membership.v1",
		"esi-location.read_location.v1",
		"esi-location.read_ship_type.v1",
		"esi-planets.manage_planets.v1",
		"esi-search.search_structures.v1",
		"esi-skills.read_skills.v1",
		"esi-ui.open_window.v1",
		"esi-ui.write_waypoint.v1",
		"esi-universe.read_structures.v1",
		"esi-wallet.read_character_wallet.v1",
	}
	tok := eveapi.CRESTToken{
		AccessToken:  "FAKE",
		RefreshToken: "So Fake",
		Expiry:       time.Now().Add(time.Hour * 100000),
		TokenType:    "Bearer"}

	err := models.AddCRESTToken(1001, 1001, "dude", &tok, strings.Join(scopes, " "))
	if err != nil {
		t.Error(err)
		return
	}
	tok2 := eveapi.CRESTToken{
		AccessToken:  "FAKE",
		RefreshToken: "So Fake",
		Expiry:       time.Now().Add(time.Hour * 100000),
		TokenType:    "Bearer"}

	err = models.AddCRESTToken(1001, 1002, "dude 2", &tok2, strings.Join(scopes, " "))
	if err != nil {
		t.Error(err)
		return
	}

	// 147035273 has some wars returned from the mock... lets throw these in.
	err = models.UpdateCharacter(1001, "dude", 1, 1, 147035273, 0, "Gallente", -10, time.Now())
	if err != nil {
		t.Error(err)
		return
	}
	err = models.UpdateCharacter(1002, "dude 2", 1, 1, 147035273, 0, "Gallente", -10, time.Now())
	if err != nil {
		t.Error(err)
		return
	}

	err = models.UpdateCorporation(147035273, "Dude Corp", "TEST2", 10, 60000004,
		"Test Executor Corp", 0, 0, "somewhere", 50, 1000, time.Now())
	if err != nil {
		t.Error(err)
		return
	}

	err = models.UpdateCorporation(145904674, "Assaulting", "BADDUDES", 10, 60000004,
		"Test Executor Corp", 0, 0, "somewhere", 50, 1000, time.Now())
	if err != nil {
		t.Error(err)
		return
	}

	eC.assetsShouldUpdate()
}

func TestEntitiesFromCrest(t *testing.T) {
	err := eC.collectEntitiesFromCREST()
	if err != nil {
		t.Error(err)
		return
	}
}

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

func TestContactSyncCheck(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()

	// Add a fake contact sync to the characters created above.
	err := models.AddContactSync(1001, 1001, 1002)
	if err != nil {
		t.Error(err)
		return
	}
	eC.contactSync()

	for {
		err := eC.contactSyncCheckQueue(r)
		if err != nil {
			t.Error(err)
			return
		}
		if i, _ := redis.Int(r.Do("SCARD", "EVEDATA_contactSyncQueue")); i == 0 {
			break
		}
	}
}

func TestStructureCheck(t *testing.T) {
	err := eC.collectStructuresFromESI()
	if err != nil {
		t.Error(err)
		return
	}
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

func TestMarketMaintTrigger(t *testing.T) {
	err := eC.marketMaintTrigger()
	if err != nil {
		t.Error(err)
		return
	}
}
