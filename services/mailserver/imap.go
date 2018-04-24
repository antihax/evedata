package mailserver

import (
	"crypto/tls"
	_ "fmt"
	"log"
	"os"

	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/antihax/evedata/services/mailserver/esiimap"
	"github.com/antihax/evedata/services/mailserver/esismtp"
	"github.com/antihax/goesi"
	imap "github.com/emersion/go-imap/server"
	smtp "github.com/emersion/go-smtp"
)

// IMAP S 993
// SMTP S 465

type IMAPServer struct {
	imap *imap.Server
	smtp *smtp.Server
}

func NewIMAPServer(tls *tls.Config, tokenAPI *tokenstore.TokenServerAPI, esi *goesi.APIClient, tokenAuth *goesi.SSOAuthenticator, q *redisqueue.RedisQueue) (*IMAPServer, error) {
	imap := imap.New(esiimap.New(tokenAPI, esi, tokenAuth, q))
	smtp := smtp.NewServer(esismtp.New(tokenAPI, esi, tokenAuth, q))
	imap.Addr = ":1993"
	smtp.Addr = ":1587"
	imap.Debug = os.Stdout
	imap.ErrorLog = log.New(os.Stdout, "INFO: ", log.Lshortfile)

	smtp.TLSConfig = tls
	imap.TLSConfig = tls

	smtp.Domain = "localhost"
	smtp.Debug = os.Stdout

	return &IMAPServer{
		imap: imap,
		smtp: smtp,
	}, nil
}

func (s *IMAPServer) Run() {
	log.Printf("Starting SMTP\n")
	go func() { log.Fatal(s.smtp.ListenAndServeTLS()) }()
	log.Printf("Starting IMAP\n")
	go func() { log.Fatal(s.imap.ListenAndServeTLS()) }()
}
