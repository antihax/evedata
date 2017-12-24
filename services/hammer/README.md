# Hammer Microservice

Hammer takes work from a HammerQueue and resolves the work from CCP's ESI Service and
feeding the data into NSQ topics for consumption from other services such as Nail.

# Usage

## Registering the Handler
```
func init() {
	registerConsumer("operation", consumerFunc)
}

func killmailConsumer(s *Hammer, parameter interface{}) {
	parameters := parameter.([]interface{})

	hash := parameters[0].(string)
	id := int32(parameters[1].(int))
    ... do stuff
```