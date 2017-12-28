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
			t := triggers[chosen]
			start := time.Now()
			go func(start time.Time, t trigger, s *Artifice) {
				log.Printf("Running trigger %s\n", t.name)
				err := t.f(s)
				log.Printf("trigger %s ran for %s\n", t.name, time.Since(start).String())
				if err != nil {
					log.Println(err)
				}
			}(start, t, s)
		}
	}
}
