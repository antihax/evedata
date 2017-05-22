package hammer

import (
	"testing"

	"sync"

	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/goesi/v1"
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

	hammer := NewHammer(redis)
	hammer.ChangeBasePath("http://127.0.0.1:8080/latest")
	defer hammer.Close()

	// Create a counter to ensure we get results for all work
	wg := &sync.WaitGroup{}
	wg.Add(len(testWork))

	go hammer.Run()

	// Consume the queued data
	go func() {
		for {
			w, err := hammer.outQueue.GetWork()
			assert.Nil(t, err)

			k, ok := w.Parameter.(goesiv1.GetKillmailsKillmailIdKillmailHashOk)
			assert.True(t, ok)
			assert.Equal(t, int32(56733821), k.KillmailId)
			wg.Done()
		}
	}()

	err := hammer.QueueWork(testWork)
	assert.Nil(t, err)
	wg.Wait()
}
