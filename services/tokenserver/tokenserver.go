package tokenserver

import (
	"context"
	"net"

	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/ptypes"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":4001"
)

// TokenServer impliments gRPC for our tokenStore
type TokenServer struct {
	tokenStore *tokenstore.TokenStore
}

// GetToken returns a token from the store
func (s *TokenServer) GetToken(ctx context.Context, in *tokenstore.TokenRequest) (*tokenstore.Token, error) {

	tok, err := s.tokenStore.GetToken(in.GetCharacterID(), in.GetTokenCharacterID())
	if err != nil {
		return nil, err
	}

	expire, err := ptypes.TimestampProto(tok.Expiry)
	if err != nil {
		return nil, err
	}

	return &tokenstore.Token{
		RefreshToken: tok.RefreshToken,
		AccessToken:  tok.AccessToken,
		Expiry:       expire,
		TokenType:    tok.TokenType,
	}, err
}

// NewTokenServer creates a new token server
func NewTokenServer(redis *redis.Pool, db *sqlx.DB, auth *goesi.SSOAuthenticator) *TokenServer {
	ts := tokenstore.NewTokenStore(redis, db, auth)
	return &TokenServer{ts}
}

// Run starts the service
func (s *TokenServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	tokenstore.RegisterTokenStoreServer(server, s)

	// Register reflection service on gRPC server.
	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		return err
	}
	return nil
}
