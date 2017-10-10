package nail

import (
	"testing"

	"github.com/antihax/evedata/internal/nsqhelper"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/services/hammer"
)

var (
	testWork []redisqueue.Work = []redisqueue.Work{
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(56733821)}},
	}
)

func TestNail(t *testing.T) {
	sql := sqlhelper.NewTestDatabase()

	redis := redigohelper.ConnectRedisTestPool()

	producer, err := nsqhelper.NewTestNSQProducer()
	if err != nil {
		t.Fatal(err)
	}

	hammer := hammer.NewHammer(redis, producer, "123400", "faaaaaaake", "sofake")
	hammer.ChangeBasePath("http://127.0.0.1:8080")
	hammer.ChangeTokenPath("http://127.0.0.1:8080")

	go hammer.Run()

	nail := NewNail(sql, nsqhelper.Test)
	go nail.Run()

	err = hammer.QueueWork(testWork)
	if err != nil {
		t.Fatal(err)
	}

	nail.Close()
	hammer.Close()
	redis.Close()
	sql.Close()
}
