package eveConsumer

import (
	"evedata/appContext"
	"log"
	"time"
)

type EveConsumer struct {
	ctx *appContext.AppContext

	stopChannel chan bool
}

func NewEVEConsumer(ctx *appContext.AppContext) *EveConsumer {
	e := &EveConsumer{ctx, make(chan bool)}

	return e
}

func (c *EveConsumer) goConsumer() {
	log.Printf("EVEConsumer: Running\n")
	rate := time.Second * 60
	throttle := time.Tick(rate)
	for {

		select {
		case <-c.stopChannel:
			return
		default:
			c.checkWars()
		}
		<-throttle
	}
	log.Printf("EVEConsumer: Shutting Down\n")
}

func (c *EveConsumer) goTriggers() {
	log.Printf("EVEConsumer: Running\n")
	rate := time.Second * 60
	throttle := time.Tick(rate)
	for {

		select {
		case <-c.stopChannel:
			return
		default:
			c.contactSync()
		}
		<-throttle
	}
	log.Printf("EVEConsumer: Shutting Down\n")
}

func (c *EveConsumer) RunConsumer() {
	go c.goConsumer()
	go c.goTriggers()
	log.Printf("EVEConsumer: Started\n")
}
func (c *EveConsumer) StopConsumer() {
	log.Printf("EVEConsumer: Stopped\n")
}
