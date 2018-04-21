package mailserver

import (
	"os"
	"testing"
	"time"

	"github.com/antihax/evedata/internal/redigohelper"
)

var mailserver *MailServer

func TestMain(m *testing.M) {
	redis := redigohelper.ConnectRedisTestPool()
	defer redis.Close()
	mailserver = NewMailServer(redis, "notreal", "reallynotreal")

	go func() {
		if err := mailserver.Run(); err != nil {
			panic(err)
		}
	}()
	retCode := m.Run()
	time.Sleep(time.Second * 5)
	mailserver.Close()

	time.Sleep(time.Second)

	os.Exit(retCode)
}
