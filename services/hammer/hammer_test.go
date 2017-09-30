package hammer

import (
	"log"
	"testing"
	"time"

	"sync"

	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/evedata/internal/nsqhelper"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/goesi/esi"
	nsq "github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
)

var (
	testWork = []redisqueue.Work{
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(1)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(2)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(3)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(4)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(5)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(6)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(7)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(8)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(9)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(10)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(11)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(12)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(13)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(14)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(15)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(16)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(17)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(18)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(19)}},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(20)}},
	}
)

func TestHammerService(t *testing.T) {

	// Setup a hammer service
	redis := redigohelper.ConnectRedisTestPool()
	defer redis.Close()

	producer, err := nsqhelper.NewTestNSQProducer()
	assert.Nil(t, err)

	hammer := NewHammer(redis, producer)
	hammer.ChangeBasePath("http://127.0.0.1:8080")
	defer hammer.Close()

	// Create a counter to ensure we get results for all work
	wg := &sync.WaitGroup{}
	wg.Add(len(testWork) + 1)

	// Run Hammer
	go hammer.Run()

	// Setup 20 consumers to test multi consumers.
	for i := 0; i < 20; i++ {
		go func() {
			consumer, err := nsqhelper.NewNSQConsumer("killmail", "hammer-test", 1)
			assert.Nil(t, err)
			consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
				k := esi.GetKillmailsKillmailIdKillmailHashOk{}
				err := gobcoder.GobDecoder(message.Body, &k)
				if err != nil {
					log.Println(err)
					return nil
				}

				assert.Equal(t, int32(56733821), k.KillmailId)

				wg.Done()

				// Hang the consumer so others get a chance.
				time.Sleep(time.Second)
				return nil
			}))
			consumer.ConnectToNSQLookupds(nsqhelper.Test)
		}()
	}

	// Load the work into the queue
	err = hammer.QueueWork(testWork)
	assert.Nil(t, err)

	// Wait for the consumers to finish\
	wg.Done()
	wg.Wait()
}
