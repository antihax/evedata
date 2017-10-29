package hammer

import (
	"testing"
	"time"

	"github.com/antihax/evedata/internal/nsqhelper"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/redisqueue"
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
		{Operation: "marketOrders", Parameter: int32(1)},
		{Operation: "war", Parameter: int32(1)},
		{Operation: "alliance", Parameter: int32(1)},
		{Operation: "corporation", Parameter: int32(1)},
		{Operation: "character", Parameter: int32(1)},
		{Operation: "marketHistoryTrigger", Parameter: int32(1)},
		{Operation: "structure", Parameter: int64(1)},
		{Operation: "marketHistory", Parameter: []int32{1, 1}},
	}
)

func TestHammerService(t *testing.T) {

	// Setup a hammer service
	redis := redigohelper.ConnectRedisTestPool()
	defer redis.Close()

	producer, err := nsqhelper.NewTestNSQProducer()
	assert.Nil(t, err)

	hammer := NewHammer(redis, producer, "123400", "faaaaaaake", "sofake")
	hammer.ChangeBasePath("http://127.0.0.1:8080")
	hammer.ChangeTokenPath("http://127.0.0.1:8080")

	// Run Hammer
	go hammer.Run()

	// Load the work into the queue
	err = hammer.QueueWork(testWork)
	assert.Nil(t, err)

	time.Sleep(time.Second)
	// Wait for the consumers to finish

	hammer.Close()
}
