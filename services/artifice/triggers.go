package artifice

import (
	"log"
	"reflect"
	"time"
)

var (
	triggers []trigger
)

type triggerFunc func(*Artifice) error

type trigger struct {
	name   string
	f      triggerFunc
	ticker *time.Ticker
}

// Register a trigger to a queue operation.
func registerTrigger(name string, f triggerFunc, ticker *time.Ticker) {
	triggers = append(triggers, trigger{name, f, ticker})
}

func (s *Artifice) runTriggers() {
	for {
		cases := make([]reflect.SelectCase, len(triggers))
		for i, ch := range triggers {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.ticker.C)}
		}
		chosen, _, ok := reflect.Select(cases)
		if ok {
			trigger := triggers[chosen]
			start := time.Now()
			log.Printf("Running trigger %s\n", trigger.name)
			err := trigger.f(s)
			log.Printf("trigger %s ran for %s\n", trigger.name, time.Since(start).String())
			if err != nil {
				log.Println(err)
			}
		}
	}
}
