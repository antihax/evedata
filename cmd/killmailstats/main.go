package main

import (
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/services/killmailstats"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("evedata killmailstats: ")
	log.Printf("Starting killmailstats Microservice Go: %s\n", runtime.Version())

	// Make a new service and send it into the background.
	kms := killmailstats.NewKillmailStats(sqlhelper.NewDatabase())
	go kms.Run()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	// Stop the service gracefully.
	kms.Close()
}
