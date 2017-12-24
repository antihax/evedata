package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/antihax/evedata/internal/nsqhelper"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/services/discordbottemp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("evedata discord: ")
	db := sqlhelper.NewDatabase()
	redis := redigohelper.ConnectRedisProdPool()
	// Make a new service and send it into the background.
	discord := discordbottemp.NewDiscordBot(redis, db, nsqhelper.Prod, os.Getenv("DISCORD_TOKEN"))
	go discord.Run()

	// Run metrics
	http.Handle("/metrics", promhttp.Handler())
	go log.Fatalln(http.ListenAndServe(":3000", nil))

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	// Stop the service gracefully.
	discord.Close()
}
