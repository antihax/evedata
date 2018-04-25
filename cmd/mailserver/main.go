package main

import (
	"log"

	"os"
	"os/signal"
	"syscall"

	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/services/mailserver"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("evedata mailserver: ")
	redis := redigohelper.ConnectRedisLocalPool()
	// Make a new service and send it into the background.
	mailserver, err := mailserver.NewMailServer(redis, os.Getenv("ESI_CLIENTID_TOKENSTORE"), os.Getenv("ESI_SECRET_TOKENSTORE"))
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Starting MailServer\n")
	go func() { log.Println(mailserver.Run()) }()
	defer mailserver.Close()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
