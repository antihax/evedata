// Package artifice provides seqencing of timed triggers for pulling information.
package mailserver

import (
	"log"
	"os"
	"sync"

	"github.com/antihax/evedata/services/mailserver/esiimap"
	"github.com/antihax/evedata/services/mailserver/esismtp"
	imap "github.com/emersion/go-imap/server"
	smtp "github.com/emersion/go-smtp"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
)

// MailServer provides token information.
type MailServer struct {
	stop      chan bool
	wg        *sync.WaitGroup
	imap      *imap.Server
	smtp      *smtp.Server
	tokenAPI  *tokenstore.TokenServerAPI
	esi       *goesi.APIClient
	tokenAuth *goesi.SSOAuthenticator
	redis     *redis.Pool
}

// NewMailServer Service.
func NewMailServer(redis *redis.Pool, clientID, secret string) (*MailServer, error) {

	// Get a caching http client
	httpClient := apicache.CreateHTTPClientCache()

	// Setup a token authenticator
	auth := goesi.NewSSOAuthenticator(httpClient, clientID, secret, "", []string{})
	esiClient := goesi.NewAPIClient(httpClient, "EVEData-Mail-Server")

	tokenServer, err := tokenstore.NewTokenServerAPI()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	q := redisqueue.NewRedisQueue(redis, "mailserver_queue")

	imap := imap.New(esiimap.New(tokenServer, esiClient, auth, q))
	smtp := smtp.NewServer(esismtp.New(tokenServer, esiClient, auth, q))

	// haproxy handles encryption
	imap.AllowInsecureAuth = true
	smtp.AllowInsecureAuth = true

	if len(os.Args) > 1 && os.Args[1] == "debug" {
		imap.Addr = ":1993"
		smtp.Addr = ":1465"
		imap.Debug = os.Stdout
		smtp.Debug = os.Stdout

	} else {
		imap.Addr = ":993"
		smtp.Addr = ":465"
	}

	imap.ErrorLog = log.New(os.Stdout, "INFO: ", log.Lshortfile)

	smtp.Domain = "localhost"

	// Setup a new MailServer
	s := &MailServer{
		stop:      make(chan bool),
		wg:        &sync.WaitGroup{},
		esi:       esiClient,
		tokenAuth: auth,
		tokenAPI:  tokenServer,
		redis:     redis,
		imap:      imap,
		smtp:      smtp,
	}

	return s, nil
}

// Run the service
func (s *MailServer) Run() error {
	log.Printf("Starting SMTP\n")
	go func() { log.Fatal(s.smtp.ListenAndServe()) }()
	log.Printf("Starting IMAP\n")
	go func() { log.Fatal(s.imap.ListenAndServe()) }()
	return nil
}

// Close the service
func (s *MailServer) Close() {

	close(s.stop)
	s.wg.Wait()
}
