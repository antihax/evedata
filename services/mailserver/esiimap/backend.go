package esiimap

import (
	"context"
	"log"

	"github.com/antihax/evedata/internal/redisqueue"

	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/antihax/go-imap/backend"
	"github.com/antihax/goesi"
)

func New(tokenAPI *tokenstore.TokenServerAPI, esi *goesi.APIClient, tokenAuth *goesi.SSOAuthenticator, q *redisqueue.RedisQueue) *Backend {
	return &Backend{tokenAPI, esi, tokenAuth, q}
}

type Backend struct {
	tokenAPI   *tokenstore.TokenServerAPI
	esi        *goesi.APIClient
	tokenAuth  *goesi.SSOAuthenticator
	cacheQueue *redisqueue.RedisQueue
}

func (s *Backend) lookupAddresses(ids []int32) ([]string, []string, error) {
	names, err := s.cacheQueue.GetCacheInBulk("addressName", ids)
	if err != nil {
		return nil, nil, err
	}
	types, err := s.cacheQueue.GetCacheInBulk("addressType", ids)
	if err != nil {
		return nil, nil, err
	}

	missing := []int32{}
	missingIdx := []int{}

	for i := range ids {
		if names[i] == "" || types[i] == "" {
			missing = append(missing, ids[i])
			missingIdx = append(missingIdx, i)
		}
	}
	if len(missing) > 0 {
		lookup, _, err := s.esi.ESI.UniverseApi.PostUniverseNames(context.Background(), missing, nil)
		if err != nil {
			log.Printf("%v : %s", missing, err)
			return nil, nil, err
		}
		for i, e := range lookup {
			names[missingIdx[i]] = e.Name
			types[missingIdx[i]] = e.Category
		}

		err = s.cacheQueue.SetCacheInBulk("addressName", ids, names)
		if err != nil {
			return nil, nil, err
		}
		err = s.cacheQueue.SetCacheInBulk("addressType", ids, types)
		if err != nil {
			return nil, nil, err
		}
	}

	return names, types, nil
}

func (s *Backend) Login(username, password string) (backend.User, error) {
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
		characterID: u.TokenCharacterID,
		mailboxes:   make(map[string]*Mailbox),
	}, nil
}
