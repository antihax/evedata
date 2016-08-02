package eveapi

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// SSOAuthenticator provides interfacing to the CREST SSO. NewSSOAuthenticator is used to create
// this structure.
type SSOAuthenticator struct {
	// Hide this...
	oauthConfig *oauth2.Config
}

// Redirect type to hide oauth2 API
type CRESTToken oauth2.Token
type CRESTTokenP *oauth2.Token

// NewSSOAuthenticator create a new CREST SSO Authenticator.
// Requires your application clientID, clientSecret, and redirectURL.
// RedirectURL must match exactly to what you registered with CCP.
func NewSSOAuthenticator(clientID string, clientSecret string, redirectURL string, scopes []string) *SSOAuthenticator {
	client := &SSOAuthenticator{}
	client.oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.eveonline.com/oauth/authorize",
			TokenURL: "https://login.eveonline.com/oauth/token",
		},
		Scopes:      scopes,
		RedirectURL: redirectURL,
	}
	return client
}

// AuthorizeURL returns a url for an end user to authenticate with EVE SSO
// and return success to the redirectURL.
// It is important to create a significatly unique state for this request
// and verify the state matches when returned to the redirectURL.
func (c SSOAuthenticator) AuthorizeURL(state string, onlineAccess bool) string {
	if onlineAccess == true {
		return c.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	} else {
		return c.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	}
}

// TokenExchange exchanges the code returned to the redirectURL with
// the CREST server to an access token. A caching client must be passed.
// This client MUST cache per CCP guidelines or face banning.
func (c SSOAuthenticator) TokenExchange(client *http.Client, code string) (CRESTTokenP, error) {

	tok, err := c.oauthConfig.Exchange(createContext(client), code)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

// Add custom clients to the context.
func createContext(httpClient *http.Client) context.Context {
	parent := oauth2.NoContext
	ctx := context.WithValue(parent, oauth2.HTTPClient, httpClient)
	return ctx
}

// TokenToJSON helper function to convert a token to a storable format.
func TokenToJSON(token *CRESTToken) (string, error) {
	if d, err := json.Marshal(token); err != nil {
		return "", err
	} else {
		return string(d), nil
	}
}

// TokenFromJSON helper function to convert stored JSON to a token.
func TokenFromJSON(jsonStr string) (*CRESTToken, error) {
	var token CRESTToken
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		return nil, err
	}
	return &token, nil
}
