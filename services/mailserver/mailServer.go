// Package artifice provides seqencing of timed triggers for pulling information.
package mailserver

import (
	"crypto/tls"
	"log"
	"sync"

	"github.com/antihax/evedata/internal/redisqueue"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
)

// MailServer provides token information.
type MailServer struct {
	stop chan bool
	wg   *sync.WaitGroup

	tokenAPI  *tokenstore.TokenServerAPI
	esi       *goesi.APIClient
	tokenAuth *goesi.SSOAuthenticator
	redis     *redis.Pool
}

// NewMailServer Service.
func NewMailServer(redis *redis.Pool, clientID, secret string) *MailServer {

	// Get a caching http client
	httpClient := apicache.CreateHTTPClientCache(redis)

	// Setup a token authenticator
	auth := goesi.NewSSOAuthenticator(httpClient, clientID, secret, "", []string{})
	esiClient := goesi.NewAPIClient(httpClient, "EVEData-Mail-Server")
	// Setup a new MailServer
	s := &MailServer{
		stop:      make(chan bool),
		wg:        &sync.WaitGroup{},
		esi:       esiClient,
		tokenAuth: auth,
		redis:     redis,
	}

	return s
}

// Run the hammer service
func (s *MailServer) Run() error {
	tokenServer, err := tokenstore.NewTokenServerAPI()
	if err != nil {
		log.Println(err)
		return err
	}
	s.tokenAPI = tokenServer

	cert, err := tls.LoadX509KeyPair(
		"/etc/letsencrypt/live/mail.evedata.org/fullchain1.pem",
		"/etc/letsencrypt/live/mail.evedata.org/privkey1.pem",
	)
	if err != nil {
		log.Println(err)
		return err
	}

	q := redisqueue.NewRedisQueue(s.redis, "mailserver_queue")

	imap, err := NewIMAPServer(&tls.Config{Certificates: []tls.Certificate{cert}}, s.tokenAPI, s.esi, s.tokenAuth, q)
	if err != nil {
		log.Println(err)
		return err
	}
	imap.Run()
	return nil
}

// Close the hammer service
func (s *MailServer) Close() {

	close(s.stop)
	s.wg.Wait()
}
