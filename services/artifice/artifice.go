// Package artifice provides seqencing of timed triggers for pulling information.
package artifice

import (
	"log"
	"strconv"
	"time"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

// Artifice handles the scheduling of routine tasks.
type Artifice struct {
	stop    chan bool
	inQueue *redisqueue.RedisQueue
	esi     *goesi.APIClient
	redis   *redis.Pool
	db      *sqlx.DB
	mail    chan esi.PostCharactersCharacterIdMailMail

	// authentication
	token       *oauth2.TokenSource
	tokenCharID int32
	auth        *goesi.SSOAuthenticator
}

// NewArtifice Service.
func NewArtifice(redis *redis.Pool, db *sqlx.DB, clientID string, secret string, refresh string, refreshCharID string) *Artifice {

	if clientID == "" {
		log.Fatalln("Missing clientID")
	}
	if secret == "" {
		log.Fatalln("Missing secret")
	}
	if refresh == "" {
		log.Fatalln("Missing refresh token")
	}
	if refreshCharID == "" {
		log.Fatalln("Missing refresh CharID")
	}
	// Get a caching http client
	cache := apicache.CreateHTTPClientCache(redis)

	// Create our ESI API Client
	esiClient := goesi.NewAPIClient(cache, "EVEData-API-Artifice")

	// Setup an authenticator
	auth := goesi.NewSSOAuthenticator(cache, clientID, secret, "", []string{})

	tok := &oauth2.Token{
		Expiry:       time.Now(),
		AccessToken:  "",
		RefreshToken: refresh,
		TokenType:    "Bearer",
	}

	charID, err := strconv.Atoi(refreshCharID)
	if err != nil {
		log.Fatalln(err)
	}

	// Build our token
	token := auth.TokenSource(tok)

	// Setup a new artifice
	s := &Artifice{
		stop: make(chan bool),
		inQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),
		db:    db,
		auth:  auth,
		mail:  make(chan esi.PostCharactersCharacterIdMailMail),
		esi:   esiClient,
		redis: redis,

		tokenCharID: int32(charID),
		token:       &token,
	}

	return s
}

// Close the service
func (s *Artifice) Close() {
	close(s.stop)
}

// ChangeBasePath for ESI (sisi/mock/tranquility)
func (s *Artifice) ChangeBasePath(path string) {
	s.esi.ChangeBasePath(path)
}

// ChangeTokenPath for ESI (sisi/mock/tranquility)
func (s *Artifice) ChangeTokenPath(path string) {
	s.auth.ChangeTokenURL(path)
	s.auth.ChangeAuthURL(path)
}

// QueueWork directly
func (s *Artifice) QueueWork(work []redisqueue.Work, priority int) error {
	return s.inQueue.QueueWork(work, priority)
}

// QueueSize returns the size of the queue
func (s *Artifice) QueueSize() (int, error) {
	return s.inQueue.Size()
}

// Run the service
func (s *Artifice) Run() {
	go s.startup()
	go s.zkillboardPost()
	go s.warKillmails()
	go s.runMetrics()
	go s.mailRunner()
	s.runTriggers()
}

// DoSQL executes a sql statement
func (s *Artifice) doSQL(stmt string, args ...interface{}) error {
	return sqlhelper.DoSQL(s.db, stmt, args...)
}

type CharacterPairs struct {
	CharacterID      int32 `db:"characterID"`
	TokenCharacterID int32 `db:"tokenCharacterID"`
	AllianceID       int32 `db:"allianceID"`
	CorporationID    int32 `db:"corporationID"`
}

func (s *Artifice) GetCharactersForScope(scope string) ([]CharacterPairs, error) {
	pairs := []CharacterPairs{}
	err := s.db.Select(&pairs,
		`SELECT characterID, tokenCharacterID FROM evedata.crestTokens T
			WHERE lastStatus != "invalid_token" AND scopes LIKE ?`, "%"+scope+"%")
	return pairs, err
}

func (s *Artifice) GetAllianceForScope(scope string) ([]CharacterPairs, error) {
	pairs := []CharacterPairs{}
	err := s.db.Select(&pairs,
		`SELECT characterID, tokenCharacterID, allianceID, corporationID FROM evedata.crestTokens T
			WHERE lastStatus != "invalid_token" AND scopes LIKE ? AND allianceID > 0
			GROUP BY allianceID
			`, "%"+scope+"%")
	return pairs, err
}

func (s *Artifice) GetCorporationForScope(scope string) ([]CharacterPairs, error) {
	pairs := []CharacterPairs{}
	err := s.db.Select(&pairs,
		`SELECT characterID, tokenCharacterID, allianceID, corporationID FROM evedata.crestTokens T
			WHERE lastStatus != "invalid_token" AND scopes LIKE ? 
			GROUP BY corporationID
			`, "%"+scope+"%")
	return pairs, err
}
