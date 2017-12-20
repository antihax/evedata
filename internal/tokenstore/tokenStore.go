package tokenstore

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

// TokenStore provides storage and caching of OAuth2 Tokens
type TokenStore struct {
	redis *redis.Pool
	db    *sqlx.DB
	auth  *goesi.SSOAuthenticator
}

// NewTokenStore provides mechinism for caching and storing SSO Tokens
// If a refresh token changes, do remember to invalidate the cache
func NewTokenStore(redis *redis.Pool, db *sqlx.DB, auth *goesi.SSOAuthenticator) *TokenStore {
	t := &TokenStore{redis, db, auth}
	return t
}

// GetToken retreives a token from storage
func (c *TokenStore) GetToken(characterID int32, tokenCharacterID int32) (*oauth2.Token, error) {
	t, err := c.getTokenFromCache(characterID, tokenCharacterID)
	if err != nil || t == nil {
		t, err = c.getTokenFromDB(characterID, tokenCharacterID)
		if err != nil {
			return nil, err
		}
	}

	if t.Expiry.Before(time.Now()) {
		a, err := c.auth.TokenSource(t)
		if err != nil {
			c.tokenError(characterID, tokenCharacterID, 999, err.Error())
			return nil, err
		}
		token, err := a.Token()
		if err != nil {
			c.tokenError(characterID, tokenCharacterID, 999, err.Error())
			return nil, err
		}
		c.setTokenToCache(characterID, tokenCharacterID, token)
		c.updateTokenToDB(characterID, tokenCharacterID, token)

		tok := &oauth2.Token{
			Expiry:       token.Expiry,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
		}
		c.tokenSuccess(characterID, tokenCharacterID)
		return tok, nil
	}

	return t, nil
}

// SetToken a token to storage
func (c *TokenStore) SetToken(characterID int32, tokenCharacterID int32, token *oauth2.Token) error {
	err := c.setTokenToCache(characterID, tokenCharacterID, token)
	if err != nil {
		return err
	}

	return c.updateTokenToDB(characterID, tokenCharacterID, token)
}

// GetTokenSource retreives a token from storage and returns a token source
func (c *TokenStore) GetTokenSource(characterID int32, tokenCharacterID int32) (oauth2.TokenSource, error) {
	t, err := c.getTokenFromCache(characterID, tokenCharacterID)
	if err != nil || t == nil {
		t, err = c.getTokenFromDB(characterID, tokenCharacterID)
		if err != nil {
			return nil, err
		}
	}

	a, err := c.auth.TokenSource(t)
	if err != nil {
		c.tokenError(characterID, tokenCharacterID, 999, err.Error())
		return nil, err
	}

	if t.Expiry.Before(time.Now()) {
		token, err := a.Token()
		if err != nil {
			c.invalidateTokenCache(characterID, tokenCharacterID)
			c.tokenError(characterID, tokenCharacterID, 999, err.Error())
			return nil, err
		}
		c.setTokenToCache(characterID, tokenCharacterID, token)
		c.updateTokenToDB(characterID, tokenCharacterID, token)
		c.tokenSuccess(characterID, tokenCharacterID)
	}

	return a, err
}

func (c *TokenStore) getTokenFromDB(characterID int32, tokenCharacterID int32) (*oauth2.Token, error) {

	type CRESTToken struct {
		Expiry       time.Time `db:"expiry" json:"expiry,omitempty"`
		TokenType    string    `db:"tokenType" json:"tokenType,omitempty"`
		AccessToken  string    `db:"accessToken" json:"accessToken,omitempty"`
		RefreshToken string    `db:"refreshToken" json:"refreshToken,omitempty"`
	}
	token := CRESTToken{}

	if err := c.db.QueryRowx(
		`SELECT expiry, tokenType, accessToken, refreshToken
			FROM evedata.crestTokens
			WHERE characterID = ? AND tokenCharacterID = ?
			LIMIT 1`,
		characterID, tokenCharacterID).StructScan(&token); err != nil {

		return nil, err
	}

	tok := &oauth2.Token{
		Expiry:       token.Expiry,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
	}

	return tok, nil
}

func (c *TokenStore) getTokenFromCache(characterID int32, tokenCharacterID int32) (*oauth2.Token, error) {
	r := c.redis.Get()
	defer r.Close()
	tok := &oauth2.Token{}

	key := fmt.Sprintf("EVEDATA_TOKENSTORE_%d_%d", characterID, tokenCharacterID)

	v, err := redis.Bytes(r.Do("GET", key))
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, errors.New("Timed out waiting on token store")
	}

	err = gobcoder.GobDecoder(v, tok)
	if err != nil {
		return nil, err
	}

	return tok, nil
}

func (c *TokenStore) setTokenToCache(characterID int32, tokenCharacterID int32, token *oauth2.Token) error {
	r := c.redis.Get()
	defer r.Close()

	key := fmt.Sprintf("EVEDATA_TOKENSTORE_%d_%d", characterID, tokenCharacterID)

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	tok := &oauth2.Token{
		Expiry:       token.Expiry,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
	}

	err := enc.Encode(tok)
	if err != nil {
		return err
	}
	if err := r.Send("SETEX", key, 80000, b.Bytes()); err != nil {
		return err
	}
	return nil
}

func (c *TokenStore) invalidateTokenCache(characterID int32, tokenCharacterID int32) error {
	r := c.redis.Get()
	defer r.Close()

	key := fmt.Sprintf("EVEDATA_TOKENSTORE_%d_%d", characterID, tokenCharacterID)

	if err := r.Send("DEL", key); err != nil {
		return err
	}
	return nil
}

func (c *TokenStore) updateTokenToDB(characterID int32, tokenCharacterID int32, token *oauth2.Token) error {
	_, err := c.db.Exec(`
		UPDATE evedata.crestTokens 
		SET accessToken = ?,
			refreshToken = ?, 
			expiry = ?
		WHERE 
			characterID = ? AND
			tokenCharacterID = ?`,
		token.AccessToken,
		token.RefreshToken,
		token.Expiry,
		characterID,
		tokenCharacterID)
	return err
}

func (c *TokenStore) tokenError(characterID int32, tokenCharacterID int32, code int, status string) error {
	if _, err := c.db.Exec(`
		UPDATE evedata.crestTokens SET lastCode = ?, lastStatus = ?
		WHERE characterID = ? AND tokenCharacterID = ?`,
		code, status, characterID, tokenCharacterID); err != nil {
		return err
	}
	return nil
}

func (c *TokenStore) tokenSuccess(characterID int32, tokenCharacterID int32) error {
	if _, err := c.db.Exec(`
		UPDATE evedata.crestTokens SET lastCode = ?, lastStatus = ?
		WHERE characterID = ? AND tokenCharacterID = ?`,
		"200", "Ok", characterID, tokenCharacterID); err != nil {
		return err
	}
	return nil
}
