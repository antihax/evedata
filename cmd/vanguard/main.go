package main

import (
	"log"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
	_ "github.com/antihax/evedata/services/vanguard/views"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/context"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("evedata vanguard: ")

	r := redigohelper.ConnectRedisProdPool()
	db := sqlhelper.NewDatabase()

	// Make a new service and send it into the background.
	vanguard := vanguard.NewVanguard(r, db,
		os.Getenv("ESI_REFRESHKEY"),
		os.Getenv("ESI_CLIENTID_TOKENSTORE"),
		os.Getenv("ESI_SECRET_TOKENSTORE"),
		os.Getenv("ESI_CLIENTID_SSO"),
		os.Getenv("ESI_SECRET_SSO"),
		os.Getenv("COOKIE_SECRET"),
		os.Getenv("DOMAIN"),
	)
	log.Printf("Setup Router\n")
	rtr := vanguard.NewRouter()
	defer vanguard.Close()

	// Handle command line arguments
	if len(os.Args) > 1 {
		if os.Args[1] == "dumpdb" {
			// Dump the database to sql file.
			log.Printf("Dumping Database to evedata.sql\n")
			err := models.DumpDatabase("./sql/evedata.sql", "evedata")
			if err != nil {
				log.Fatalln(err)
			}
		} else if os.Args[1] == "flushcache" {
			// Erase http cache in redis
			log.Printf("Flushing Redis\n")
			conn := r.Get()
			defer conn.Close()
			keys, err := redis.Strings(conn.Do("KEYS", "*rediscache*"))
			if err != nil {
				log.Println(err)
			} else {
				for _, key := range keys {
					conn.Do("DEL", key)
					log.Printf("Deleting %s\n", key)
				}
			}
		} else if os.Args[1] == "flushstructures" {
			// Erase http cache in redis
			log.Printf("Flushing structures\n")
			conn := r.Get()
			defer conn.Close()
			keys, err := redis.Strings(conn.Do("KEYS", "*structure_failure*"))
			if err != nil {
				log.Println(err)
			} else {
				for _, key := range keys {
					conn.Do("DEL", key)
					log.Printf("Deleting %s\n", key)
				}
			}

		} else if os.Args[1] == "flushredis" {
			// Erase everything in redis for modified deployments
			log.Printf("Flushing Redis\n")
			conn := r.Get()
			defer conn.Close()
			conn.Do("FLUSHALL")
		}
	}

	log.Printf("Start Listening\n")
	go log.Fatalln(http.ListenAndServe(":3000", context.ClearHandler(rtr)))

	log.Printf("In production\n")
	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}
