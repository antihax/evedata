package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/antihax/evedata/internal/nsqhelper"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/services/artifice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	redis := redigohelper.ConnectRedisProdPool()

	producer, err := nsqhelper.NewNSQProducer()
	if err != nil {
		log.Panicln(err)
	}

	// Make a new service and send it into the background.
	artifice := artifice.NewArtifice(redis, producer, os.Getenv("ESI_CLIENTID"), os.Getenv("ESI_SECRET"), os.Getenv("ESI_REFRESHKEY"))
	go artifice.Run()
	defer artifice.Close()

	// Run metrics
	http.Handle("/metrics", promhttp.Handler())

	go log.Fatalln(http.ListenAndServe(":3000", nil))

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
