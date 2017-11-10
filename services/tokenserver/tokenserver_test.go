package tokenserver

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/goesi"
)

func TestTokenServer(t *testing.T) {
	sql := sqlhelper.NewTestDatabase()

	redis := redigohelper.ConnectRedisTestPool()
	redConn := redis.Get()
	defer redConn.Close()
	redConn.Do("FLUSHALL")

	// Setup an authenticator
	auth := goesi.NewSSOAuthenticator(&http.Client{}, "FAAAK", "SOFAAAAAAAKE", "",
		[]string{"esi-universe.read_structures.v1",
			"esi-search.search_structures.v1",
			"esi-markets.structure_markets.v1"})
	auth.ChangeTokenURL("http://127.0.0.1:8080")
	auth.ChangeAuthURL("http://127.0.0.1:8080")

	ts := NewTokenServer(redis, sql, auth)
	go ts.Run()

	// sleep to allow server to start
	time.Sleep(time.Millisecond * 50)

	conn, err := grpc.Dial("localhost:4001", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := tokenstore.NewTokenStoreClient(conn)

	expire, _ := ptypes.TimestampProto(time.Now().Add(time.Hour))

	resp, err := client.SetToken(context.TODO(), &tokenstore.SetTokenRequest{
		CharacterID:      1,
		TokenCharacterID: 1,
		AccessToken:      "FAAAAAAAKE",
		RefreshToken:     "so fake",
		TokenType:        "Bearer",
		Expiry:           expire,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetOk() != true {
		t.Fatal("Not true")
	}

	token, err := client.GetToken(context.TODO(), &tokenstore.GetTokenRequest{
		CharacterID:      1,
		TokenCharacterID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if token.GetAccessToken() != "FAAAAAAAKE" {
		t.Fatal("Not FAAAAAAAKE")
	}

	token, err = client.GetToken(context.TODO(), &tokenstore.GetTokenRequest{
		CharacterID:      1,
		TokenCharacterID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if token.GetAccessToken() != "FAAAAAAAKE" {
		t.Fatal("Not FAAAAAAAKE")
	}
}
