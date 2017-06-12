package nsqhelper

import (
	nsq "github.com/nsqio/go-nsq"
)

var (
	Prod = []string{"nsqlookupd1.nsq", "nsqlookupd2.nsq"}
	Test = []string{"localhost:4161"}
)

func NewNSQConsumer(topicName, channelName string, maxInFlight int) (*nsq.Consumer, error) {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = maxInFlight
	return nsq.NewConsumer(topicName, channelName, cfg)
}

func NewNSQProducer() (*nsq.Producer, error) {
	cfg := nsq.NewConfig()
	return nsq.NewProducer("nsqd.nsq:4150", cfg)
}

func NewTestNSQProducer() (*nsq.Producer, error) {
	cfg := nsq.NewConfig()
	return nsq.NewProducer("localhost:4150", cfg)
}
