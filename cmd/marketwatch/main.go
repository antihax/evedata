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
	"github.com/antihax/evedata/services/marketwatch"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("evedata-marketwatch: ")
	log.Println("starting evedata-marketwatch")

	db := sqlhelper.NewDatabase()

	mw := marketwatch.NewMarketWatch(
		os.Getenv("ESI_REFRESHKEY"),
		os.Getenv("ESI_CLIENTID_TOKENSTORE"),
		os.Getenv("ESI_SECRET_TOKENSTORE"),
		db,
		redigohelper.ConnectLedisProdPool(),
	)
	go mw.Run()

	// Run metrics
	http.Handle("/metrics", promhttp.Handler())

	log.Println("started evedata-marketwatch")
	go log.Fatalln(http.ListenAndServe(":3000", nil))

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
