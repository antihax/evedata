package conservator

import (
	"log"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type spawnFunc func(s *Conservator, consumer *nsq.Consumer)

// Structure for handling routes
type conservatorHandler struct {
	Topic     string
	SpawnFunc spawnFunc
}

var handlers []conservatorHandler

// AddHandler adds a nail handler
func addHandler(topic string, spawnFunc spawnFunc) {
	handlers = append(handlers, conservatorHandler{topic, spawnFunc})
}

func (s *Conservator) registerHandlers() error {
	nsqcfg := nsq.NewConfig()

	for _, h := range handlers {
		c, err := nsq.NewConsumer(h.Topic, "conservator", nsqcfg)
		if err != nil {
			log.Fatalln(err)
		}

		h.SpawnFunc(s, c)
		s.consumers[h.Topic] = c

		err = c.ConnectToNSQLookupds(s.consumerAddresses)
		if err != nil {
			log.Fatalln(err)
			return err
		}
	}
	return nil
}

// Wrap handlers in a wait group we can properly account during shutdown.
func (s *Conservator) wait(next nsq.Handler) nsq.Handler {
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
