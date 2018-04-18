package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/services/tokenserver"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("evedata tokenserver: ")
	redis := redigohelper.ConnectRedisProdPool()

	db := sqlhelper.NewDatabase()

	// Make a new service and send it into the background.
	tokenserver := tokenserver.NewTokenServer(redis, db, os.Getenv("ESI_CLIENTID_TOKENSTORE"), os.Getenv("ESI_SECRET_TOKENSTORE"))
	go tokenserver.Run()
	defer tokenserver.Close()

	// Run metrics
	http.Handle("/metrics", promhttp.Handler())

	go log.Fatalln(http.ListenAndServe(":3000", nil))

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
