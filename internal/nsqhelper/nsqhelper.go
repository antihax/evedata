package nsqhelper

import (
	nsq "github.com/nsqio/go-nsq"
)

var (
	Prod = []string{"nsqlookupd1.nsq.svc.cluster.local:4161", "nsqlookupd2.nsq.svc.cluster.local:4161"}
	Test = []string{"localhost:4161"}
)

func NewNSQConsumer(topicName, channelName string, maxInFlight int) (*nsq.Consumer, error) {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = maxInFlight
	return nsq.NewConsumer(topicName, channelName, cfg)
}

func NewNSQProducer() (*nsq.Producer, error) {
	cfg := nsq.NewConfig()
	return nsq.NewProducer("nsqd.nsq.svc.cluster.local:4150", cfg)
}

func NewTestNSQProducer() (*nsq.Producer, error) {
	cfg := nsq.NewConfig()
	return nsq.NewProducer("localhost:4150", cfg)
}
