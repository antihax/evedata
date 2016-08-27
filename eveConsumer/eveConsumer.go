package eveConsumer

import (
	"evedata/appContext"
	"log"
	"time"
)

// EveConsumer provides the microservice which conducts backend
// polling of EVE Crest and XML resources as needed.
type EVEConsumer struct {
	ctx                 *appContext.AppContext
	consumerStopChannel chan bool
	triggersStopChannel chan bool
}

// NewEVEConsumer creates a new EveConsumer
func NewEVEConsumer(ctx *appContext.AppContext) *EVEConsumer {
	e := &EVEConsumer{ctx, make(chan bool), make(chan bool)}

	return e
}

func (c *EVEConsumer) goConsumer() {
	log.Printf("EVEConsumer: Running Consumer\n")
	rate := time.Second * 60 * 15
	throttle := time.Tick(rate)
	for {

		select {
		case <-c.consumerStopChannel:
			log.Printf("EVEConsumer: Shutting Down\n")
			return
		default:
			c.checkAlliances()
			c.checkWars()
		}
		<-throttle
	}
}

func (c *EVEConsumer) goTriggers() {
	log.Printf("EVEConsumer: Running Triggers\n")
	rate := time.Second * 60 * 15
	throttle := time.Tick(rate)
	for {

		select {
		case <-c.triggersStopChannel:
			log.Printf("EVEConsumer: Shutting Down\n")
			return
		default:
			c.contactSync()
			c.updateDatabase()
		}
		<-throttle
	}
}

// RunConsumer starts the consumer and returns.
func (c *EVEConsumer) RunConsumer() {
	go c.goConsumer()
	go c.goTriggers()
	log.Printf("EVEConsumer: Started\n")
}

// StopConsumer shuts down any running go routines and returns.
func (c *EVEConsumer) StopConsumer() {
	c.consumerStopChannel <- true
	c.triggersStopChannel <- true
	log.Printf("EVEConsumer: Stopped\n")
}
