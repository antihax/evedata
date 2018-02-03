// Package artifice provides seqencing of timed triggers for pulling information.
package artifice

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

// Artifice handles the scheduling of routine tasks.
type Artifice struct {
	stop     chan bool
	hammerWG *sync.WaitGroup
	inQueue  *redisqueue.RedisQueue
	esi      *goesi.APIClient
	redis    *redis.Pool
	db       *sqlx.DB

	// authentication
	token       *oauth2.TokenSource
	tokenCharID int32
	auth        *goesi.SSOAuthenticator
}

// NewArtifice Service.
func NewArtifice(redis *redis.Pool, db *sqlx.DB, clientID string, secret string, refresh string, refreshCharID string) *Artifice {

	// Get a caching http client
	cache := apicache.CreateHTTPClientCache(redis)

	// Create our ESI API Client
	esi := goesi.NewAPIClient(cache, "EVEData-API-Artifice")

	// Setup an authenticator
	auth := goesi.NewSSOAuthenticator(cache, clientID, secret, "", []string{})

	tok := &oauth2.Token{
		Expiry:       time.Now(),
		AccessToken:  "",
		RefreshToken: refresh,
		TokenType:    "Bearer",
	}

	charID, err := strconv.Atoi(refreshCharID)

	// Build our token
	token, err := auth.TokenSource(tok)
	if err != nil {
		log.Fatalln(err)
	}

	// Setup a new artifice
	s := &Artifice{
		stop:     make(chan bool),
		hammerWG: &sync.WaitGroup{},
		inQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),
		db:          db,
		auth:        auth,
		esi:         esi,
		redis:       redis,
		tokenCharID: int32(charID),
		token:       &token,
	}

	return s
}

// Close the hammer service
func (s *Artifice) Close() {
	close(s.stop)
	s.hammerWG.Wait()
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

// Run the hammer service
func (s *Artifice) Run() {
	go s.startup()
	go s.zkillboardPost()
	go s.warKillmails()
	go s.runMetrics()
	s.runTriggers()
}

// RetryTransaction on deadlocks
func retryTransaction(tx *sqlx.Tx) error {
	for {
		err := tx.Commit()
		if err != nil {
			if !strings.Contains(err.Error(), "1213") {
				return err
			}
			time.Sleep(250 * time.Millisecond)
			continue
		} else {
			return err
		}
	}
}

// DoSQL executes a sql statement
func (s *Artifice) doSQL(stmt string, args ...interface{}) error {
	for {
		err := s.doSQLTranq(stmt, args...)
		if err != nil {
			if !strings.Contains(err.Error(), "1213") {

				return err
			}
			time.Sleep(250 * time.Millisecond)
			continue
		} else {
			return err
		}
	}
}

// DoSQL executes a sql statement
func (s *Artifice) doSQLTranq(stmt string, args ...interface{}) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(stmt, args...)
	if err != nil {
		return err
	}

	err = retryTransaction(tx)
	if err != nil {
		return err
	}
	return nil
}
