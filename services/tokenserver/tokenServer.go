// Package artifice provides seqencing of timed triggers for pulling information.
package tokenserver

import (
	"net"
	"sync"
	"time"

	"google.golang.org/grpc/keepalive"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/msgpackcodec"
	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

// TokenServer provides token information.
type TokenServer struct {
	stop       chan bool
	wg         *sync.WaitGroup
	redis      *redis.Pool
	db         *sqlx.DB
	tokenAPI   *grpc.Server
	tokenStore *tokenstore.TokenStore

	// authentication
	auth *goesi.SSOAuthenticator
}

// NewTokenServer Service.
func NewTokenServer(redis *redis.Pool, db *sqlx.DB, clientID, secret string) *TokenServer {

	// Get a caching http client
	cache := apicache.CreateHTTPClient()

	// Setup a token authenticator
	auth := goesi.NewSSOAuthenticator(cache, clientID, secret, "", []string{})

	// Setup RPC server
	server := grpc.NewServer(grpc.CustomCodec(&msgpackcodec.MsgPackCodec{}),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    time.Second * 5,
				Timeout: time.Second * 10,
			}),
	)

	// create Token Store
	tokenStore := tokenstore.NewTokenStore(redis, db, auth)

	// Setup a new TokenServer
	s := &TokenServer{
		stop:       make(chan bool),
		wg:         &sync.WaitGroup{},
		db:         db,
		auth:       auth,
		redis:      redis,
		tokenAPI:   server,
		tokenStore: tokenStore,
	}

	return s
}

// Run the service
func (s *TokenServer) Run() error {
	lis, err := net.Listen("tcp", ":3002")
	if err != nil {
		return err
	}

	s.tokenAPI.RegisterService(&serviceDesc, s)
	err = s.tokenAPI.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

// Close the service
func (s *TokenServer) Close() {
	s.tokenAPI.GracefulStop()
	close(s.stop)
	s.wg.Wait()
}
