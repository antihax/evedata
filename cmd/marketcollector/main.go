package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/antihax/evedata/internal/nsqhelper"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/services/marketcollector"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("evedata marketcollector: ")

	db := sqlhelper.NewDatabase()

	// Make a new service and send it into the background.
	marketcollector := marketcollector.NewMarketCollector(
		db,
		nsqhelper.Prod,
	)

	defer marketcollector.Close()

	// Run metrics
	http.Handle("/metrics", promhttp.Handler())

	go log.Fatalln(http.ListenAndServe(":3099", nil))

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
