package tokenStore

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

type TokenStore struct {
	redis *redis.Pool
	db    *sqlx.DB
	auth  *goesi.SSOAuthenticator
}

func NewTokenStore(redis *redis.Pool, db *sqlx.DB, auth *goesi.SSOAuthenticator) *TokenStore {
	t := &TokenStore{redis, db, auth}
	return t
}

func (c *TokenStore) GetTokenSource(characterID int64, tokenCharacterID int64) (goesi.CRESTTokenSource, error) {
	t, err := c.getTokenFromCache(characterID, tokenCharacterID)
	if err != nil || t == nil {
		t, err = c.getTokenFromDB(characterID, tokenCharacterID)
		if err != nil {
			return nil, err
		}
	}

	a, err := c.auth.TokenSource(t)
	if err != nil {
		return nil, err
	}

	if t.Expiry.Before(time.Now()) {
		token, err := a.Token()
		if err != nil {
			return nil, err
		}
		c.setTokenToCache(characterID, tokenCharacterID, token)
		c.updateTokenToDB(characterID, tokenCharacterID, token)
	}

	return a, err
}

func (c *TokenStore) getTokenFromDB(characterID int64, tokenCharacterID int64) (*goesi.CRESTToken, error) {

	type CRESTToken struct {
		Expiry       time.Time `db:"expiry" json:"expiry,omitempty"`
		TokenType    string    `db:"tokenType" json:"tokenType,omitempty"`
		AccessToken  string    `db:"accessToken" json:"accessToken,omitempty"`
		RefreshToken string    `db:"refreshToken" json:"refreshToken,omitempty"`
	}
	token := &CRESTToken{}

	if err := c.db.QueryRowx(
		`SELECT expiry, tokenType, accessToken, refreshToken
			FROM evedata.crestTokens
			WHERE characterID = ? AND tokenCharacterID = ?
			LIMIT 1`,
		characterID, tokenCharacterID).StructScan(token); err != nil {

		return nil, err
	}

	tok := &goesi.CRESTToken{
		Expiry:       token.Expiry,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
	}

	return tok, nil
}

func (c *TokenStore) getTokenFromCache(characterID int64, tokenCharacterID int64) (*goesi.CRESTToken, error) {
	r := c.redis.Get()
	defer r.Close()
	tok := &goesi.CRESTToken{}

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

func (c *TokenStore) setTokenToCache(characterID int64, tokenCharacterID int64, token *oauth2.Token) error {
	r := c.redis.Get()
	defer r.Close()

	key := fmt.Sprintf("EVEDATA_TOKENSTORE_%d_%d", characterID, tokenCharacterID)

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	tok := &goesi.CRESTToken{
		Expiry:       token.Expiry,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
	}

	err := enc.Encode(tok)
	if err != nil {
		return err
	}
	if err := r.Send("SET", key, b.Bytes()); err != nil {
		return err
	}
	return nil
}

func (c *TokenStore) updateTokenToDB(characterID int64, tokenCharacterID int64, token *oauth2.Token) error {
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
