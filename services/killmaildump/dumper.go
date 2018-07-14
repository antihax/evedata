// Package killmaildumper dumps killmails to json files for testing
package killmaildumper

import (
	"log"
	"os"

	nsq "github.com/nsqio/go-nsq"
)

// Dumper dumps killmails to json files for testing.
type Dumper struct {
	stop     chan bool
	consumer *nsq.Consumer
}

// NewDumper Service.
func NewDumper(consumerAddresses []string) *Dumper {
	// Setup a new artifice
	s := &Dumper{
		stop: make(chan bool),
	}

	path := "./json"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	nsqcfg := nsq.NewConfig()
	c, err := nsq.NewConsumer("killmail", "killmaildumper", nsqcfg)
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
func (s *Dumper) Close() {
	close(s.stop)
	s.consumer.Stop()
}
