package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/services/squirrel"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("evedata squirrel: ")

	db := sqlhelper.NewDatabase()
	// Run metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatalln(http.ListenAndServe(":3000", nil))
	}()

	// Make a new service and send it into the background.
	squirrel := squirrel.NewSquirrel(db)

	squirrel.Run()
	squirrel.Close()

	// Allow prometheus to collect final stats
	time.Sleep(20 * time.Second)
}
