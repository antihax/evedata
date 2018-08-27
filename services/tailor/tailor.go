// Package tailor handles fitting attributes to the database
package tailor

import (
	"log"
	"os"
	"time"

	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/jmoiron/sqlx"
	nsq "github.com/nsqio/go-nsq"
	backblaze "gopkg.in/kothar/go-backblaze.v0"
)

// Tailor dumps killmails to json files for testing.
type Tailor struct {
	stop     chan bool
	consumer *nsq.Consumer
	db       *sqlx.DB
	b2       *backblaze.B2
}

// NewTailor Service.
func NewTailor(db *sqlx.DB, b2 *backblaze.B2, consumerAddresses []string) *Tailor {
	// Setup a new artifice
	s := &Tailor{
		stop: make(chan bool),
		db:   db,
		b2:   b2,
	}

	nsqcfg := nsq.NewConfig()
	nsqcfg.MaxInFlight = 50
	nsqcfg.MsgTimeout = time.Second * 30
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
