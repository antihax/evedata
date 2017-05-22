package nail

import (
	"sync"

	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
)

type Nail struct {
	stop    chan bool
	wg      *sync.WaitGroup
	db      *sqlx.DB
	inQueue *redisqueue.RedisQueue
}

func NewNail(db *sqlx.DB, redis *redis.Pool) *Nail {
	n := &Nail{
		db:   db,
		wg:   &sync.WaitGroup{},
		stop: make(chan bool),
		inQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-nail",
		),
	}

	return n
}

// Run the nail service
func (s *Nail) Run() {
	for {
		select {
		case <-s.stop:
			return
		}
	}
}

// Stop the nail service
func (s *Nail) Close() {
	close(s.stop)
	s.wg.Wait()
}
