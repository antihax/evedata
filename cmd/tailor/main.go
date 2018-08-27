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
	"github.com/antihax/evedata/services/tailor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	backblaze "gopkg.in/kothar/go-backblaze.v0"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("evedata tailor: ")

	db := sqlhelper.NewDatabase()

	b2, err := backblaze.NewB2(backblaze.Credentials{
		AccountID:      os.Getenv("B2_ACCOUNTID"),
		ApplicationKey: os.Getenv("B2_APPLICATION_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Make a new service and send it into the background.
	tailor := tailor.NewTailor(
		db,
		b2,
		nsqhelper.Prod,
	)

	defer tailor.Close()

	// Run metrics
	http.Handle("/metrics", promhttp.Handler())

	go log.Fatalln(http.ListenAndServe(":3000", nil))

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
