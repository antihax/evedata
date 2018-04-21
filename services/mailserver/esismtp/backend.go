package esismtp

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/antihax/goesi"
	smtp "github.com/emersion/go-smtp"
	"golang.org/x/oauth2"
)

func New(tokenAPI *tokenstore.TokenServerAPI, esi *goesi.APIClient, tokenAuth *goesi.SSOAuthenticator) *Backend {
	return &Backend{tokenAPI, esi, tokenAuth}
}

type Backend struct {
	tokenAPI  *tokenstore.TokenServerAPI
	esi       *goesi.APIClient
	tokenAuth *goesi.SSOAuthenticator
}

func (s *Backend) Login(username, password string) (smtp.User, error) {
	u, err := s.tokenAPI.GetMailUser(username, password)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	ts, err := s.tokenAuth.TokenSource(u.Token)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &User{
		username:    username,
		token:       ts,
		backend:     s,
		characterID: u.TokenCharacterID}, nil
}

// Require clients to authenticate using SMTP AUTH before sending emails
func (bkd *Backend) AnonymousLogin() (smtp.User, error) {
	return nil, smtp.ErrAuthRequired
}

type User struct {
	username    string
	token       oauth2.TokenSource
	backend     *Backend
	characterID int32
}

func (u *User) Send(from string, to []string, r io.Reader) error {
	log.Println("Sending message:", from, to)

	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
		log.Println("Data:", string(b))
	}
	return nil
}

func (u *User) Logout() error {
	return nil
}
