// Package artifice provides seqencing of timed triggers for pulling information.
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
	stop    chan bool
	esi     *goesi.APIClient
	esiSem  chan bool
	inQueue *redisqueue.RedisQueue
	wg      sync.WaitGroup
	db      *sqlx.DB
}

// NewSquirrel Service.
func NewSquirrel(redis *redis.Pool, ledis *redis.Pool, db *sqlx.DB) *Squirrel {
	// Get a caching http client
	cache := apicache.CreateHTTPClientCache(ledis)

	// Create our ESI API Client
	esiClient := goesi.NewAPIClient(cache, "EVEData-API-Squirrel")

	// Setup a new squirrel
	s := &Squirrel{
		wg:     sync.WaitGroup{},
		esiSem: make(chan bool, 100),
		db:     db,
		inQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),
		esi: esiClient,
	}

	return s
}

// Tick up the concurrency limit
func (s *Squirrel) esiSemStart() {
	s.esiSem <- true
	s.wg.Add(1)
}

// Completed
func (s *Squirrel) esiSemFinished() {
	<-s.esiSem
	s.wg.Done()
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
	s.wg.Add(1)
	defer s.wg.Done()
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
	return sqlhelper.DoSQL(s.db, stmt, args)
}
