package eveConsumer

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"
	"golang.org/x/oauth2"
)

var (
	ctx appContext.AppContext
	eC  *EVEConsumer
)

func TestMain(m *testing.M) {
	ctx = appContext.NewTestAppContext()
	eC = NewEVEConsumer(&ctx)
	r := ctx.Cache.Get()

	// Create service states
	models.SetServiceState("wars", time.Now().UTC(), 1)
	models.SetServiceState("alliances", time.Now().UTC(), 1)
	models.SetServiceState("marketHistory", time.Now().UTC(), 1)
	models.SetServiceState("marketMaint", time.Now().UTC(), 1)
	models.SetServiceState("npcCorps", time.Now().UTC(), 1)
	models.SetServiceState("structures", time.Now().UTC(), 1)

	scopes := []string{
		goesi.ScopeCharacterContractsRead,
		goesi.ScopeCharacterMarketOrdersRead,
		goesi.ScopeCharacterResearchRead,
		goesi.ScopeCharacterWalletRead,
		goesi.ScopeCharacterContactsRead,
		goesi.ScopeCharacterContactsWrite,
		"esi-assets.read_assets.v1",
		"esi-characters.read_contacts.v1",
		"esi-characters.write_contacts.v1",
		"esi-characters.read_notifications.v1",
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
	tok := oauth2.Token{
		AccessToken:  "FAKE",
		RefreshToken: "So Fake",
		Expiry:       time.Now().Add(time.Hour * 100000),
		TokenType:    "Bearer"}

	err := models.AddCRESTToken(1001, 1001, "dude", &tok, strings.Join(scopes, " "))
	if err != nil {
		log.Fatal(err)
		return
	}
	tok2 := oauth2.Token{
		AccessToken:  "FAKE",
		RefreshToken: "So Fake",
		Expiry:       time.Now().Add(time.Hour * 100000),
		TokenType:    "Bearer"}

	err = models.AddCRESTToken(1001, 1002, "dude 2", &tok2, strings.Join(scopes, " "))
	if err != nil {
		log.Fatal(err)
		return
	}

	// 147035273 has some wars returned from the mock... lets throw these in.
	err = models.UpdateCharacter(1001, "dude", 1, 1, 147035273, 0, 1, "Female", -10, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}
	err = models.UpdateCharacter(1002, "dude 2", 1, 1, 147035273, 0, 2, "Female", -10, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	err = models.UpdateCorporation(147035273, "Dude Corp", "TEST2", 10,
		0, 0, 50, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	err = models.UpdateCorporation(145904674, "Assaulting", "BADDUDES", 10,
		0, 0, 50, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	r.Close()

	// Run the tests
	retCode := m.Run()
	os.Exit(retCode)
}
