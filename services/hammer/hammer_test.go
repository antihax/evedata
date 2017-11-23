package hammer

import (
	"testing"
	"time"

	"github.com/antihax/evedata/internal/nsqhelper"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
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
		{Operation: "characterWalletTransactions", Parameter: []int32{int32(1), int32(1)}},
		//{Operation: "characterWalletJournal", Parameter: []int32{int32(1), int32(1)}},
		{Operation: "characterAssets", Parameter: []int32{int32(1), int32(1)}},
		{Operation: "characterNotifications", Parameter: []int32{int32(1), int32(1)}},
		{Operation: "loyaltyStore", Parameter: int32(1000001)},
		{Operation: "wheeeeeeeeee", Parameter: int32(1000001)},
	}
)

func TestHammerService(t *testing.T) {
	sql := sqlhelper.NewTestDatabase()
	// Setup a hammer service
	redis := redigohelper.ConnectRedisTestPool()
	defer redis.Close()

	producer, err := nsqhelper.NewTestNSQProducer()
	assert.Nil(t, err)

	hammer := NewHammer(redis, sql, producer, "123400", "faaaaaaake", "sofake", "1232423423", "now with 200% more fake!")
	hammer.ChangeBasePath("http://127.0.0.1:8080")
	hammer.ChangeTokenPath("http://127.0.0.1:8080")
	hammer.tokenStore.SetToken(1, 1, &oauth2.Token{
		RefreshToken: "fake",
		AccessToken:  "really fake",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	})

	// Run Hammer
	go hammer.Run()

	err = hammer.AddAlliance(99002200)
	assert.Nil(t, err)

	err = hammer.AddCorporation(99002200)
	assert.Nil(t, err)

	err = hammer.AddCharacter(99002200)
	assert.Nil(t, err)

	// Load the work into the queue
	err = hammer.QueueWork(testWork)
	assert.Nil(t, err)

	time.Sleep(time.Second)
	// Wait for the consumers to finish

	hammer.Close()
	time.Sleep(time.Second)
}
