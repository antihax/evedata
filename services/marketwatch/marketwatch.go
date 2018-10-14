package marketwatch

import (
	"log"
	"sync"
	"time"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/goesi/esi"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"

	"github.com/antihax/goesi"
	"golang.org/x/oauth2"
)

// MarketWatch provides CCP Market Data
type MarketWatch struct {
	// goesi client
	esi *goesi.APIClient

	// authentication
	doAuth    bool
	token     *oauth2.TokenSource
	tokenAuth *goesi.SSOAuthenticator

	// SQL Queues
	orderChan          chan []esi.GetMarketsRegionIdOrders200Ok
	contractChan       chan []FullContract
	orderChangeChan    chan []OrderChange
	contractChangeChan chan []ContractChange
	orderDeleteChan    chan []OrderChange
	contractDeleteChan chan []ContractChange

	// Database pool
	db *sqlx.DB

	// data store
	market     map[int64]*sync.Map
	structures map[int64]*Structure
	contracts  map[int64]*sync.Map
	mmutex     sync.RWMutex // Market mutex for the main map
	cmutex     sync.RWMutex // Contract mutex for the main map
	smutex     sync.RWMutex // Structure mutex for the whole map
}

// NewMarketWatch creates a new MarketWatch microservice
func NewMarketWatch(refresh, tokenClientID, tokenSecret string, db *sqlx.DB, ledis *redis.Pool) *MarketWatch {
	// Get a caching http client
	cache := apicache.CreateLimitedHTTPClientCache(ledis)

	// Create our ESI API Client
	esiClient := goesi.NewAPIClient(cache, "EVEData-API-MarketWatch")

	// Setup an authenticator for our user tokens
	doAuth := false
	if tokenClientID == "" || tokenSecret == "" || refresh == "" {
		log.Println("Warning: Missing authentication parameters so only regional market will be polled")
	} else {
		doAuth = true
	}
	auth := goesi.NewSSOAuthenticator(cache, tokenClientID, tokenSecret, "", []string{})

	// Setup our token for structures
	tok := &oauth2.Token{
		Expiry:       time.Now(),
		AccessToken:  "",
		RefreshToken: refresh,
		TokenType:    "Bearer",
	}
	token := auth.TokenSource(tok)

	return &MarketWatch{
		// ESI Client
		esi: esiClient,

		// database pool
		db: db,

		// ESI SSO Handler
		doAuth:    doAuth,
		token:     &token,
		tokenAuth: auth,

		// SQL channels
		orderChan:          make(chan []esi.GetMarketsRegionIdOrders200Ok, 10000),
		contractChan:       make(chan []FullContract, 10000),
		orderChangeChan:    make(chan []OrderChange, 10000),
		contractChangeChan: make(chan []ContractChange, 10000),
		orderDeleteChan:    make(chan []OrderChange, 10000),
		contractDeleteChan: make(chan []ContractChange, 10000),

		// Market Data Map
		market:     make(map[int64]*sync.Map),
		structures: make(map[int64]*Structure),
		contracts:  make(map[int64]*sync.Map),
	}
}

// Run starts the market watch service
func (s *MarketWatch) Run() error {
	go s.startUpMarketWorkers()
	go s.sqlPumps()
	return nil
}

func (s *MarketWatch) doSQL(stmt string, args ...interface{}) error {
	return sqlhelper.DoSQL(s.db, stmt, args...)
}
