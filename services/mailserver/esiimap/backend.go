package esiimap

import (
	"log"
	"strconv"
	"strings"

	"github.com/antihax/evedata/internal/redisqueue"

	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/antihax/goesi"
	"github.com/emersion/go-imap/backend"
)

type Backend struct {
	tokenAPI         *tokenstore.TokenServerAPI
	esi              *goesi.APIClient
	tokenAuth        *goesi.SSOAuthenticator
	cacheQueue       *redisqueue.RedisQueue
	cacheLookup      chan int32
	cacheMailingList chan int32
}

func New(tokenAPI *tokenstore.TokenServerAPI, esi *goesi.APIClient, tokenAuth *goesi.SSOAuthenticator, q *redisqueue.RedisQueue) *Backend {
	b := &Backend{
		tokenAPI,
		esi,
		tokenAuth,
		q,
		make(chan int32, 1000000),
		make(chan int32, 1000000),
	}

	// Start the cache lookup queue
	go b.precacheLookup()
	go b.precacheMailingLists()
	return b
}

func (s *Backend) Login(username, password string) (backend.User, error) {
	parts := strings.Split(username, "@")
	characterID, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	u, err := s.tokenAPI.GetMailUser(int32(characterID), password)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	ts := s.tokenAuth.TokenSource(u.Token)

	user := NewUser(username, ts, s, u.TokenCharacterID)
	return user, nil
}
