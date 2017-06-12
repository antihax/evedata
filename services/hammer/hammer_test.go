package hammer

import (
	"log"
	"testing"

	"sync"

	"github.com/antihax/goesi/v1"

	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/evedata/internal/nsqhelper"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/redisqueue"
	nsq "github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
)

var (
	testWork []redisqueue.Work = []redisqueue.Work{
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(1)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(2)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(3)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(4)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(5)}},
	}
)

func TestHammerService(t *testing.T) {
	redis := redigohelper.ConnectRedisTestPool()
	defer redis.Close()

	producer, err := nsqhelper.NewTestNSQProducer()
	assert.Nil(t, err)

	hammer := NewHammer(redis, producer)
	hammer.ChangeBasePath("http://127.0.0.1:8080/latest")
	defer hammer.Close()

	// Create a counter to ensure we get results for all work
	wg := &sync.WaitGroup{}
	wg.Add(len(testWork))

	go hammer.Run()

	// Consume the queued data
	{
		consumer, err := nsqhelper.NewNSQConsumer("killmail", "hammer-test", 5)
		assert.Nil(t, err)
		consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
			k := goesiv1.GetKillmailsKillmailIdKillmailHashOk{}

			err := gobcoder.GobDecoder(message.Body, &k)
			if err != nil {
				log.Println(err)
				return nil
			}

			assert.Equal(t, int32(56733821), k.KillmailId)
			wg.Done()
			return nil
		}))
		consumer.ConnectToNSQLookupds(nsqhelper.Test)
	}

	err = hammer.QueueWork(testWork)
	assert.Nil(t, err)
	wg.Wait()
}
