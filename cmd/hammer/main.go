package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/antihax/evedata/internal/nsqhelper"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/services/hammer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	redis := redigohelper.ConnectRedisProdPool()

	producer, err := nsqhelper.NewNSQProducer()
	if err != nil {
		log.Panicln(err)
	}

	// Make a new service and send it into the background.
	hammer := hammer.NewHammer(redis, producer)
	go hammer.Run()
	defer hammer.Close()

	// Run metrics
	http.Handle("/metrics", promhttp.Handler())
	go log.Fatalln(http.ListenAndServe(":3000", nil))

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
