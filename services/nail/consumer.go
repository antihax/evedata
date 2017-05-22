package nail

type consumer struct {
	name string
	f    consumerFunc
}

type consumerFunc func(*Nail, []byte) error

func addConsumer(name string, f consumerFunc) {
	consumers = append(consumers, consumer{name, f})
	consumerMap[name] = f
}

func wait(s *Nail, f consumerFunc, p []byte) {
	s.wg.Add(1)
	defer s.wg.Done()
	f(s, p)
}

// For handling Consumers
var (
	consumers   []consumer
	consumerMap map[string]consumerFunc
)

func init() {
	consumerMap = make(map[string]consumerFunc)
}
