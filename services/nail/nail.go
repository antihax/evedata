package nail

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	nsq "github.com/nsqio/go-nsq"
)

// Nail handles storage of data to SQL
type Nail struct {
	stop           chan bool
	wg             *sync.WaitGroup
	db             *sqlx.DB
	inQueue        map[string]*nsq.Consumer
	redis          *redis.Pool
	outQueue       *redisqueue.RedisQueue
	characterRaces map[int32]string
}

// NewNail creates a new storage engine
func NewNail(redis *redis.Pool, db *sqlx.DB, addresses []string) *Nail {
	n := &Nail{
		db:             db,
		wg:             &sync.WaitGroup{},
		stop:           make(chan bool),
		inQueue:        make(map[string]*nsq.Consumer),
		characterRaces: make(map[int32]string),
		outQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),
	}
	go n.loadStaticData()
	nsqcfg := nsq.NewConfig()
	nsqcfg.MaxInFlight = 50
	nsqcfg.MsgTimeout = time.Minute * 5

	for _, h := range handlers {
		c, err := nsq.NewConsumer(h.Topic, "nail", nsqcfg)
		if err != nil {
			log.Fatalln(err)
		}

		h.SpawnFunc(n, c)
		n.inQueue[h.Topic] = c

		err = c.ConnectToNSQLookupds(addresses)
		if err != nil {
			log.Fatalln(err)
		}

		// Stop the logger being so verbose
		c.SetLogger(log.New(os.Stderr, "", log.Flags()), nsq.LogLevelError)
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

// Close stop the nail service
func (s *Nail) Close() {
	close(s.stop)
	for _, h := range s.inQueue {
		h.Stop()
	}
	s.wg.Wait()
}

// Wrap handlers in a wait group we can properly account during shutdown.
func (s *Nail) wait(next nsq.Handler) nsq.Handler {
	return nsq.HandlerFunc(func(m *nsq.Message) error {
		s.wg.Add(1)
		defer s.wg.Done()
		err := next.HandleMessage(m)
		if err != nil {
			log.Printf("%s\n", err)
			m.Requeue(time.Second)
		} else {
			m.Finish()
		}
		return err
	})
}

type spawnFunc func(s *Nail, consumer *nsq.Consumer)

// Structure for handling routes
type nailHandler struct {
	Topic     string
	SpawnFunc spawnFunc
}

var handlers []nailHandler

// AddHandler adds a nail handler
func AddHandler(topic string, spawnFunc spawnFunc) {
	handlers = append(handlers, nailHandler{topic, spawnFunc})
}

// DoSQL executes a sql statement
func (s *Nail) doSQL(stmt string, args ...interface{}) error {
	return sqlhelper.DoSQL(s.db, stmt, args...)
}

func (s *Nail) loadStaticData() error {
	rows, err := s.db.Query(`SELECT raceID, raceName FROM chrRaces;`)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id   int32
			race string
		)
		err := rows.Scan(&id, &race)
		if err != nil {
			log.Fatalln(err)
		}
		s.characterRaces[id] = race
	}
	return nil
}
