package kafkahelper

import "github.com/Shopify/sarama"

var (
	productionBrokers []string = []string{"front:32181", "front:32182", "front:32183"}
	testBrokers       []string = []string{"127.0.0.1:9092"}
)

func NewKafkaClient() (sarama.Client, error) {
	conf := sarama.NewConfig()
	return sarama.NewClient(productionBrokers, conf)
}

func NewKafkaTestClient() (sarama.Client, error) {
	conf := sarama.NewConfig()
	return sarama.NewClient(testBrokers, conf)
}
