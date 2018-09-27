// Package marketcollector handles fitting attributes to the database
package marketcollector

import (
	"encoding/json"
	"log"
	"time"

	"github.com/antihax/eve-marketwatch/marketwatch"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

var wsDialer = websocket.Dialer{
	Subprotocols:     []string{"p1", "p2"},
	ReadBufferSize:   1024 * 1024 * 500,
	WriteBufferSize:  1024,
	HandshakeTimeout: 30 * time.Second,
}

// MarketCollector gathers changes from eve-marketwatch and posts to SQL
type MarketCollector struct {
	stop             chan bool
	ws               *websocket.Conn
	db               *sqlx.DB
	messageChan      chan *Message
	orderHistoryChan chan []marketwatch.OrderChange
}

// NewMarketCollector Service.
func NewMarketCollector(db *sqlx.DB) *MarketCollector {
	// Setup a new artifice
	s := &MarketCollector{
		stop:             make(chan bool),
		db:               db,
		messageChan:      make(chan *Message, 10000),
		orderHistoryChan: make(chan []marketwatch.OrderChange, 10000),
	}
	c, _, err := wsDialer.Dial("ws://marketwatch.evedata:3005/?market=1&contract=1", nil)
	if err != nil {
		log.Fatalln(err)
	}
	s.ws = c
	go s.readPump()
	go s.sqlPump()

	// Restart once an hour to get full market tick
	restart := time.NewTimer(1 * time.Hour)
	go func() {
		<-restart.C
		s.Close()
	}()

	return s
}

// Message wraps different payloads for the websocket interface
type Message struct {
	Action  string           `json:"action"`
	Payload *json.RawMessage `json:"payload"`
}

func (s *MarketCollector) readPump() {
	defer s.Close()
	for {
		message := &Message{}
		err := s.ws.ReadJSON(message)
		if err != nil {
			log.Fatalln("read:", err)
		}
		s.messageChan <- message
	}
}

// Close the conservator service
func (s *MarketCollector) Close() {
	close(s.stop)
	s.ws.Close()
}

func (s *MarketCollector) doSQL(stmt string, args ...interface{}) error {
	return sqlhelper.DoSQL(s.db, stmt, args...)
}
