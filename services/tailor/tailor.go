// Package tailor handles fitting attributes to the database
package tailor

import (
	"log"
	"os"

	nsq "github.com/nsqio/go-nsq"
)

// Tailor dumps killmails to json files for testing.
type Tailor struct {
	stop     chan bool
	consumer *nsq.Consumer
}

// NewTailor Service.
func NewTailor(consumerAddresses []string) *Tailor {
	// Setup a new artifice
	s := &Tailor{
		stop: make(chan bool),
	}

	nsqcfg := nsq.NewConfig()
	c, err := nsq.NewConsumer("killmail", "tailor", nsqcfg)
	if err != nil {
		log.Fatalln(err)
	}
	s.consumer = c
	c.AddHandler(nsq.HandlerFunc(s.killmailHandler))
	err = c.ConnectToNSQLookupds(consumerAddresses)
	if err != nil {
		log.Fatalln(err)
	}

	// Stop the logger being so verbose
	c.SetLogger(log.New(os.Stderr, "", log.Flags()), nsq.LogLevelError)
	return s
}

// Close the conservator service
func (s *Tailor) Close() {
	close(s.stop)
	s.consumer.Stop()
}
