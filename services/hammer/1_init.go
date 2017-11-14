package hammer

// For handling Consumers
var (
	consumers   []consumer
	consumerMap map[string]consumerFunc
)

func init() {
	consumerMap = make(map[string]consumerFunc)
}
