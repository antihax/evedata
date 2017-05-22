package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/services/hammer"
)

func main() {
	redis := redigohelper.ConnectRedisPool(
		[]string{"sentinel1:26379", "sentinel2:26379", "sentinel3:26379"},
		os.Getenv("REDIS_PASSWORD"),
		"evedata",
		true,
	)

	// Make a new service and send it into the background.
	hammer := hammer.NewHammer(redis)
	go hammer.Run()
	defer hammer.Close()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
