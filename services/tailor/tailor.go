// Package tailor handles fitting attributes to the database
package tailor

import (
	"log"
	"os"

	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/jmoiron/sqlx"
	nsq "github.com/nsqio/go-nsq"
)

// Tailor dumps killmails to json files for testing.
type Tailor struct {
	stop     chan bool
	consumer *nsq.Consumer
	db       *sqlx.DB
}

// NewTailor Service.
func NewTailor(db *sqlx.DB, consumerAddresses []string) *Tailor {
	// Setup a new artifice
	s := &Tailor{
		stop: make(chan bool),
		db:   db,
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

	go s.killmailConsumer()
	return s
}

// Close the conservator service
func (s *Tailor) Close() {
	close(s.stop)
	s.consumer.Stop()
}

func (s *Tailor) doSQL(stmt string, args ...interface{}) error {
	return sqlhelper.DoSQL(s.db, stmt, args...)
}
