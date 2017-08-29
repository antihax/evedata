package nail

import (
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
	nsq "github.com/nsqio/go-nsq"
)

type Nail struct {
	stop    chan bool
	wg      *sync.WaitGroup
	db      *sqlx.DB
	inQueue map[string]*nsq.Consumer
}

func NewNail(db *sqlx.DB, addresses []string) *Nail {
	n := &Nail{
		db:      db,
		wg:      &sync.WaitGroup{},
		stop:    make(chan bool),
		inQueue: make(map[string]*nsq.Consumer),
	}

	nsqcfg := nsq.NewConfig()

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
	for _, h := range s.inQueue {
		h.Stop()
	}
}

// Wrap handlers in a wait group we can properly account during shutdown.
func (s *Nail) wait(next nsq.Handler) nsq.Handler {
	return nsq.HandlerFunc(func(m *nsq.Message) error {
		s.wg.Add(1)
		defer s.wg.Done()
		return next.HandleMessage(m)
	})
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
