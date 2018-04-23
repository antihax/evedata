package tokenstore

import (
	"context"
	"time"

	"google.golang.org/grpc/keepalive"

	"github.com/antihax/evedata/internal/msgpackcodec"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

type TokenServerAPI struct {
	grpc *grpc.ClientConn
}

func NewTokenServerAPI() (*TokenServerAPI, error) {
	r, err := grpc.Dial("tokenserver.evedata:32004",
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 5,
			Timeout:             time.Second * 10,
			PermitWithoutStream: true,
		}),
		grpc.WithInsecure(),
		grpc.WithCodec(&msgpackcodec.MsgPackCodec{}),
	)

	if err != nil {
		return nil, err

	}
	return &TokenServerAPI{r}, nil
}

func NewTokenServerInternalAPI() (*TokenServerAPI, error) {
	r, err := grpc.Dial("tokenserver.evedata:3002", grpc.WithInsecure(), grpc.WithCodec(&msgpackcodec.MsgPackCodec{}))
	if err != nil {
		r, err = grpc.Dial("tokenserver.evedata:32004", grpc.WithInsecure(), grpc.WithCodec(&msgpackcodec.MsgPackCodec{}))
		if err != nil {
			return nil, err
		}
	}
	return &TokenServerAPI{r}, nil
}

func (s *TokenServerAPI) GetToken(characterID, tokenCharacterID int32) (*oauth2.Token, error) {
	token := &oauth2.Token{}
	if err := s.grpc.Invoke(context.Background(), "/TokenStore/GetToken",
		&TokenRequest{
			CharacterID:      characterID,
			TokenCharacterID: tokenCharacterID,
		}, token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *TokenServerAPI) GetMailUser(username, password string) (*MailUser, error) {
	user := &MailUser{}
	if err := s.grpc.Invoke(context.Background(), "/TokenStore/GetMailUser",
		&MailUserRequest{
			Username: username,
			Password: password,
		}, user); err != nil {
		return nil, err
	}

	return user, nil
}

type TokenRequest struct {
	CharacterID      int32
	TokenCharacterID int32
}
type MailUserRequest struct {
	Username string
	Password string
}

type MailUser struct {
	CharacterID      int32
	TokenCharacterID int32
	Token            *oauth2.Token
}
