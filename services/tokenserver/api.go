package tokenserver

import (
	"context"
	"errors"
	"strings"

	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/internal/tokenstore"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

type TokenServerPlan interface {
	GetToken(ctx context.Context, t *tokenstore.TokenRequest) (*oauth2.Token, error)
}

var serviceDesc = grpc.ServiceDesc{
	ServiceName: "TokenStore",
	HandlerType: (*TokenServerPlan)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetToken",
			Handler:    GetTokenHandler,
		},
		{
			MethodName: "GetMailUser",
			Handler:    GetMailUserHandler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

// TokenRequest is the parameters needed for a request for a current oauth2 token set

func GetTokenHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(tokenstore.TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(*TokenServer).GetToken(ctx, in)
}

// GetToken gets a token from the store and returns it
func (s TokenServer) GetToken(ctx context.Context, t *tokenstore.TokenRequest) (*oauth2.Token, error) {
	return s.tokenStore.GetToken(t.CharacterID, t.TokenCharacterID)
}

func GetMailUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(tokenstore.MailUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(*TokenServer).GetMailUser(ctx, in)
}

// GetToken gets a token from the store and returns it
func (s TokenServer) GetMailUser(ctx context.Context, u *tokenstore.MailUserRequest) (*tokenstore.MailUser, error) {
	type MailUser struct {
		CharacterID      int32  `db:"characterID"`
		TokenCharacterID int32  `db:"tokenCharacterID"`
		Password         string `db:"mailPassword"`
	}

	t := MailUser{}
	err := s.db.QueryRowx(
		`	SELECT characterID, tokenCharacterID, mailPassword FROM evedata.crestTokens
			WHERE tokenCharacterID = ? AND mailPassword != "" AND scopes LIKE "%read_mail%"
			LIMIT 1;`, u.CharacterID).StructScan(&t)
	if err != nil {
		// Ignore this error.
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, errors.New("Access Denied")
		}
		return nil, err
	}

	if !sqlhelper.CompareHash(u.Password, t.Password) {
		return nil, errors.New("Access Denied")
	}

	token, err := s.tokenStore.GetToken(t.CharacterID, t.TokenCharacterID)
	if err != nil {
		return nil, err
	}

	return &tokenstore.MailUser{Token: token, CharacterID: t.CharacterID, TokenCharacterID: t.TokenCharacterID}, nil
}
