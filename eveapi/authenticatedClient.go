package eveapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

type AnonymousClient struct {
	httpClient *http.Client
	base       eveURI
}

type AuthenticatedClient struct {
	AnonymousClient
	tokenSource oauth2.TokenSource
}

func (c *AnonymousClient) UseTestServer(testServer bool) {
	if testServer == true {
		c.base = eveSisi
	} else {
		c.base = eveTQ
	}
}

// NewAuthenticatedClient assigns a token to a client.
func NewAuthenticatedClient(client *http.Client, tok CRESTToken) *AuthenticatedClient {
	c := &AuthenticatedClient{}
	c.base = eveTQ
	c.tokenSource = oauth2.StaticTokenSource(tok)
	c.httpClient = oauth2.NewClient(createContext(client), c.tokenSource)
	return c
}

type VerifyResponse struct {
	CharacterID        int64
	CharacterName      string
	ExpiresOn          string
	Scopes             string
	TokenType          string
	CharacterOwnerHash string
}

func decode(res *http.Response, ret interface{}) error {
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(buf), ret); err != nil {
		return err
	}
	return err
}

func (c *AuthenticatedClient) Verify() (*VerifyResponse, error) {
	v := &VerifyResponse{}
	res, err := c.httpClient.Get(c.base.Login + "/oauth/verify")
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	err = decode(res, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
