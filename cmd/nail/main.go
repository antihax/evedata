package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/antihax/evedata/internal/nsqhelper"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/services/nail"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db := sqlhelper.NewDatabase()

	// Make a new service and send it into the background.
	nail := nail.NewNail(db, nsqhelper.Prod)
	go nail.Run()

	// Run metrics
	http.Handle("/metrics", promhttp.Handler())
	go log.Fatalln(http.ListenAndServe(":3000", nil))

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	// Stop the service gracefully.
	nail.Close()
}
