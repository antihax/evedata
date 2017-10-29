package artifice

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/stretchr/testify/assert"
)

var artificeInstance *Artifice

func TestMain(m *testing.M) {
	sql := sqlhelper.NewTestDatabase()

	redis := redigohelper.ConnectRedisTestPool()
	redConn := redis.Get()
	defer redConn.Close()
	redConn.Do("FLUSHALL")

	artificeInstance = NewArtifice(redis, sql, "123400", "faaaaaaake", "sofake")

	go artificeInstance.Run()
	retCode := m.Run()
	time.Sleep(time.Second * 5)
	artificeInstance.Close()

	redis.Close()
	sql.Close()

	os.Exit(retCode)
}

func TestTriggers(t *testing.T) {
	for _, trigger := range triggers {
		err := trigger.f(artificeInstance)
		if err != nil {
			log.Println(err)
		}
		assert.Nil(t, err)
	}
}
