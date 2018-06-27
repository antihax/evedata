// Package artifice provides seqencing of timed triggers for pulling information.
package squirrel

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redisqueue"
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

// DoSQL executes a sql statement
func (s *Squirrel) doSQL(stmt string, args ...interface{}) error {
	for {
		_, err := s.RetryExec(stmt, args...)
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

// RetryExecTillNoRows retries the exec until we get no error (deadlocks) and no results are returned
func (s *Squirrel) RetryExecTillNoRows(sql string, args ...interface{}) error {
	for {
		rows, err := s.RetryExec(sql, args...)
		if err != nil {
			return err
		}
		if rows == 0 {
			break
		}
	}
	return nil
}

// RetryExec retries the exec until we get no error (deadlocks)
func (s *Squirrel) RetryExec(sql string, args ...interface{}) (int64, error) {
	var rows int64
	for {
		res, err := s.db.Exec(sql, args...)
		if err == nil {
			rows, err = res.RowsAffected()
			return rows, err
		} else if strings.Contains(err.Error(), "1213") == false {
			return rows, err
		}
	}
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
