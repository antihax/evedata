package nail

import (
	"sync"

	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	nsq "github.com/nsqio/go-nsq"
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

type spawnFunc func(s *Nail, consumer *nsq.Consumer)

// Structure for handling routes
type nailHandler struct {
	Topic     string
	SpawnFunc spawnFunc
}

var handlers []nailHandler

func AddHandler(topic string, spawnFunc spawnFunc) {
	handlers = append(handlers, nailHandler{topic, spawnFunc})
}
