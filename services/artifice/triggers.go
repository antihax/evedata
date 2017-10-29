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
	daily  bool
	hour   int
}

// Register a trigger to a queue operation.
func registerTrigger(name string, f triggerFunc, ticker *time.Ticker) {
	triggers = append(triggers, trigger{name, f, ticker, false, 0})
}

// Register a daily trigger to a queue operation.
func registerDailyTrigger(name string, f triggerFunc, hour int) {
	ticker := time.NewTicker(getNextTickDuration(hour))
	triggers = append(triggers, trigger{name, f, ticker, true, hour})
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
			if trigger.daily {
				trigger.ticker.Stop()
				trigger.ticker = time.NewTicker(getNextTickDuration(trigger.hour))
			}
			log.Printf("Running trigger %s\n", trigger.name)
			err := trigger.f(s)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func getNextTickDuration(hour int) time.Duration {
	now := time.Now()
	nextTick := time.Date(now.UTC().Year(), now.UTC().Month(), now.UTC().Day(), hour, 1, 1, 1, time.UTC)
	if nextTick.Before(now) {
		nextTick = nextTick.Add(24 * time.Hour)
	}
	return nextTick.Sub(time.Now())
}
