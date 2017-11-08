package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/antihax/goesi"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/services/tokenserver"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("evedata tokenserver: ")

	db := sqlhelper.NewDatabase()
	redis := redigohelper.ConnectRedisProdPool()
	cache := apicache.CreateHTTPClientCache(redis)
	auth := goesi.NewSSOAuthenticator(cache, os.Getenv("ESI_CLIENTID"), os.Getenv("ESI_SECRET"), "", []string{})

	// Make a new service and send it into the background.
	tokenServer := tokenserver.NewTokenServer(redis, db, auth)
	go tokenServer.Run()

	// Run metrics
	http.Handle("/metrics", promhttp.Handler())
	go log.Fatalln(http.ListenAndServe(":3000", nil))

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
