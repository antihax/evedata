package nail

import (
	"testing"

	"github.com/antihax/evedata/internal/nsqhelper"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/services/hammer"
	"github.com/stretchr/testify/assert"
)

var (
	testWork []redisqueue.Work = []redisqueue.Work{
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(56733821)}},
	}
)

func TestNail(t *testing.T) {
	redis := redigohelper.ConnectRedisTestPool()
	defer redis.Close()

	producer, err := nsqhelper.NewTestNSQProducer()
	assert.Nil(t, err)

	hammer := hammer.NewHammer(redis, producer)
	hammer.ChangeBasePath("http://127.0.0.1:8080/latest")
	go hammer.Run()
	defer hammer.Close()

	sql := sqlhelper.NewTestDatabase()
	defer sql.Close()

	nail := NewNail(sql, nsqhelper.Test)
	go nail.Run()

	err = hammer.QueueWork(testWork)
	assert.Nil(t, err)
	defer nail.Close()
}
