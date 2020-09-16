package tokenserver

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/antihax/evedata/internal/msgpackcodec"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

var tokenserver *TokenServer

func TestMain(m *testing.M) {
	sql := sqlhelper.NewTestDatabase()
	redis := redigohelper.ConnectRedisTestPool()

	tokenserver = NewTokenServer(redis, sql, "notreal", "reallynotreal")

	go func() {
		if err := tokenserver.Run(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second * 1)
	retCode := m.Run()
	tokenserver.Close()

	redis.Close()
	sql.Close()

	time.Sleep(time.Second)

	os.Exit(retCode)
}

func TestTokens(t *testing.T) {
	tokenserver.tokenStore.SetToken(555, 555, &oauth2.Token{
		RefreshToken: "fake",
		AccessToken:  "really fake",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	})

	r, err := grpc.Dial("localhost:3002", grpc.WithInsecure(), grpc.WithCodec(&msgpackcodec.MsgPackCodec{}))
	assert.Nil(t, err)
	assert.NotNil(t, r)

	token := oauth2.Token{}
	err = r.Invoke(context.Background(), "/TokenStore/GetToken", &tokenstore.TokenRequest{CharacterID: 555, TokenCharacterID: 555}, &token)
	assert.Nil(t, err)

}
