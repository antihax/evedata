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
func registerTrigger(name string, f triggerFunc, minutes int) {
	triggers = append(triggers, trigger{name, f, time.NewTicker(time.Duration(minutes) * time.Minute)})

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
			log.Printf("Running trigger %s\n", trigger.name)
			err := trigger.f(s)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
