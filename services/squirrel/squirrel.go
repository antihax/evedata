// Package squirrel collects static data from ESI and dumps it into the db
package squirrel

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
)

// Squirrel Collects and dumps all static data from ESI
type Squirrel struct {
	esi     *goesi.APIClient
	wg      sync.WaitGroup
	inQueue *redisqueue.RedisQueue
	db      *sqlx.DB
}

// NewSquirrel Service.
func NewSquirrel(redis *redis.Pool, db *sqlx.DB) *Squirrel {
	// Get a caching http client
	cache := apicache.CreateHTTPClientCache(redis)

	// Limit concurrency of outbound calls
	cache.Transport = LimiterTransport{next: cache.Transport}

	// Create our ESI API Client
	esiClient := goesi.NewAPIClient(cache, "EVEData-API-Squirrel")

	// Setup a new squirrel
	s := &Squirrel{
		db: db,
		inQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),
		esi: esiClient,
	}

	return s
}

// Close the service
func (s *Squirrel) Close() {
	s.wg.Wait()
}

// ChangeBasePath for ESI (sisi/mock/tranquility)
func (s *Squirrel) ChangeBasePath(path string) {
	s.esi.ChangeBasePath(path)
}

// Run the service
func (s *Squirrel) Run() {
	s.runTriggers()
}

func getPages(r *http.Response) (int32, error) {
	// Decode the page into int32. Return if this fails as there were no extra pages.
	pagesInt, err := strconv.Atoi(r.Header.Get("x-pages"))
	if err != nil {
		return 0, err
	}
	pages := int32(pagesInt)
	return pages, err
}

// QueueWork directly
func (s *Squirrel) QueueWork(work []redisqueue.Work, priority int) error {
	return s.inQueue.QueueWork(work, priority)
}

func (s *Squirrel) doSQL(stmt string, args ...interface{}) error {
	return sqlhelper.DoSQL(s.db, stmt, args...)
}
