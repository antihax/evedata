package tokenserver

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

type TokenServerPlan interface {
	GetToken(ctx context.Context, t *TokenRequest) (*oauth2.Token, error)
}

var serviceDesc = grpc.ServiceDesc{
	ServiceName: "TokenStore",
	HandlerType: (*TokenServerPlan)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetToken",
			Handler:    GetTokenHandler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

// TokenRequest is the parameters needed for a request for a current oauth2 token set
type TokenRequest struct {
	CharacterID      int32
	TokenCharacterID int32
}

func GetTokenHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(*TokenServer).GetToken(ctx, in)
}

// GetToken gets a token from the store and returns it
func (s TokenServer) GetToken(ctx context.Context, t *TokenRequest) (*oauth2.Token, error) {
	return s.tokenStore.GetToken(t.CharacterID, t.TokenCharacterID)
}

type MailUserRequest struct {
	Username string
	Password string
}

func GetMailUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MailUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(*TokenServer).GetMailUser(ctx, in)
}

// GetToken gets a token from the store and returns it
func (s TokenServer) GetMailUser(ctx context.Context, u *MailUserRequest) (*oauth2.Token, error) {
	type mailUser struct {
		CharacterID      int32 `db:"characterID"`
		TokenCharacterID int32 `db:"tokenCharacterID"`
	}
	t := mailUser{}

	if err := s.db.QueryRowx(
		`	SELECT characterID, tokenCharacterID FROM evedata.crestTokens
			WHERE characterName = ? AND mailPassword = ? AND mailPassword != ""
			LIMIT 1;`).StructScan(&t); err != nil {
		// Ignore this error.
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, errors.New("Access Denied")
		}
		return nil, err
	}

	return s.tokenStore.GetToken(t.CharacterID, t.TokenCharacterID)
}
