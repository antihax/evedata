package discordauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

type Authenticator struct {
	httpClient *http.Client
	// Hide this...
	oauthConfig *oauth2.Config
}

// NewAuthenticator create a new EVE SSO Authenticator.
// Requires your application clientID, clientSecret, and redirectURL.
// RedirectURL must match exactly to what you registered with CCP.
func NewAuthenticator(client *http.Client, clientID string, clientSecret string, redirectURL string, scopes []string) *Authenticator {

	if client == nil {
		return nil
	}

	c := &Authenticator{}

	c.httpClient = client

	c.oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discordapp.com/api/oauth2/authorize",
			TokenURL: "https://discordapp.com/api/oauth2/token",
		},
		Scopes:      scopes,
		RedirectURL: redirectURL,
	}

	return c
}

// ChangeAuthURL changes the oauth2 configuration url for authentication
func (c *Authenticator) ChangeAuthURL(url string) {
	c.oauthConfig.Endpoint.AuthURL = url
}

// ChangeTokenURL changes the oauth2 configuration url for token
func (c *Authenticator) ChangeTokenURL(url string) {
	c.oauthConfig.Endpoint.TokenURL = url
}

// AuthorizeURL returns a url for an end user to authenticate with EVE SSO
// and return success to the redirectURL.
// It is important to create a significatly unique state for this request
// and verify the state matches when returned to the redirectURL.
func (c *Authenticator) AuthorizeURL(state string, onlineAccess bool, scopes []string) string {
	var url string

	// Generate the URL
	if onlineAccess == true {
		url = c.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline, oauth2.SetAuthURLParam("scope", strings.Join(scopes, " ")))
	} else {
		url = c.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("scope", strings.Join(scopes, " ")))
	}

	return url
}

// TokenExchange exchanges the code returned to the redirectURL with
// the CREST server to an access token. A caching client must be passed.
// This client MUST cache per CCP guidelines or face banning.
func (c *Authenticator) TokenExchange(code string) (*oauth2.Token, error) {
	tok, err := c.oauthConfig.Exchange(createContext(c.httpClient), code)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

// TokenSource creates a refreshable token that can be passed to ESI functions
func (c *Authenticator) TokenSource(token *oauth2.Token) oauth2.TokenSource {
	return c.oauthConfig.TokenSource(createContext(c.httpClient), token)
}

type MeResponse struct {
	ID            string
	UserName      string
	Discriminator string
	Avatar        string
	Verified      string
	Email         string
}

// Verify the client and collect user information.
func (c *Authenticator) Verify(auth oauth2.TokenSource) (*MeResponse, error) {
	v := &MeResponse{}
	_, err := c.doJSON("GET", "https://discordapp.com/api/users/@me", nil, v, "application/json;", auth)

	if err != nil {
		return nil, err
	}
	return v, nil
}

// Creates a new http.Request for a public resource.
func (c *Authenticator) newRequest(method, urlStr string, body interface{}, mediaType string) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, rel.String(), buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Calls a resource from the public CREST
func (c *Authenticator) doJSON(method, urlStr string, body interface{}, v interface{}, mediaType string, auth oauth2.TokenSource) (*http.Response, error) {

	req, err := c.newRequest(method, urlStr, body, mediaType)
	if err != nil {
		return nil, err
	}

	if auth != nil {
		// We were able to grab an oauth2 token from the context
		var latestToken *oauth2.Token
		if latestToken, err = auth.Token(); err != nil {
			return nil, err
		}
		latestToken.SetAuthHeader(req)
	}

	res, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(buf))
	}
	if err := json.Unmarshal([]byte(buf), v); err != nil {
		return nil, err
	}

	return res, nil
}

// Executes a request generated with newRequest
func (c *Authenticator) executeRequest(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusOK ||
		res.StatusCode == http.StatusCreated {
		return res, nil
	}
	return res, errors.New(res.Status)
}

// Add custom clients to the context.
func createContext(httpClient *http.Client) context.Context {
	parent := oauth2.NoContext
	ctx := context.WithValue(parent, oauth2.HTTPClient, httpClient)
	return ctx
}
