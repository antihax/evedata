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
	name    string
	f       triggerFunc
	ticker  *time.Ticker
	running bool
}

// Register a trigger to a queue operation.
func registerTrigger(name string, f triggerFunc, ticker *time.Ticker) {
	triggers = append(triggers, trigger{name, f, ticker, false})
}

func (s *Artifice) runTriggers() {
	for {
		cases := make([]reflect.SelectCase, len(triggers))
		for i, ch := range triggers {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.ticker.C)}
		}
		chosen, _, ok := reflect.Select(cases)
		if ok {
			go func(t *trigger, s *Artifice) {
				// prevent running duplicate tasks
				if !t.running {
					// set running state
					t.running = true

					// unset running state when the function ends
					defer func(b *bool) {
						*b = false
					}(&t.running)

					// run the task
					if err := t.f(s); err != nil {
						log.Println(err)
					}
				} else {
					log.Printf("already running %s\n", t.name)
				}
			}(&triggers[chosen], s)
		}
	}
}
