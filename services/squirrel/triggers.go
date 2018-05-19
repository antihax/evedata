package squirrel

import (
	"log"
)

var (
	triggers   []trigger
	collectors []collector
)

type triggerFunc func(*Squirrel) error

type trigger struct {
	name    string
	f       triggerFunc
	running bool
}

type collectorFunc func(*Squirrel) error

type collector struct {
	name    string
	f       collectorFunc
	running bool
}

// Register a trigger to a queue operations.
func registerTrigger(name string, f triggerFunc) {
	triggers = append(triggers, trigger{name, f, false})
}

// Register a collector to consume data
func registerCollector(name string, f collectorFunc) {
	collectors = append(collectors, collector{name, f, false})
}

func (s *Squirrel) runTriggers() {
	for chosen := range collectors {
		s.wg.Add(1)
		go func(t *collector, s *Squirrel) {
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
				s.wg.Done()
			} else {
				log.Printf("already running collector %s\n", t.name)
			}
		}(&collectors[chosen], s)
	}

	for chosen := range triggers {
		s.wg.Add(1)
		go func(t *trigger, s *Squirrel) {
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
				s.wg.Done()
			} else {
				log.Printf("already running trigger %s\n", t.name)
			}
		}(&triggers[chosen], s)
	}
}
