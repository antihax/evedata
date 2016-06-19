package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"evedata/config"
	"io"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
	zmq "github.com/pebbe/zmq4"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *sqlx.DB
	conf *config.Config
)

func main() {
	goConsumer()
}

type emdrHeader struct {
	UploadKeys []struct {
		Name string
		Key  string
	}
	Generator struct {
		Version string
		Name    string
	}
	CurrentTime string
	ResultType  string
	Version     string
}

func goConsumer() {
	var err error
	// Read configuation.
	log.Print("Reading Configuration")
	conf, err = config.ReadConfig()
	if err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	// Build Connection Pool
	log.Print("Building Database Pool")
	db, err = sqlx.Connect(conf.Database.Driver, conf.Database.Spec)
	if err != nil {
		log.Fatalf("Cannot build database pool: %v", err)
	}

	// Check we can connect
	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	client, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		log.Fatal(err)
	}

	// Connect
	log.Print("Subscribing to EMDR")
	err = client.Connect("tcp://relay-us-central-1.eve-emdr.com:8050")
	client.SetSubscribe("")
	if err != nil {
		log.Fatal(err)
	}

	// Endless loop.
	log.Print("Waiting for content")
	for {
		// Receive message from ZeroMQ.
		msg, err := client.Recv(0)
		if err != nil {
			log.Fatal(err)
		}

		// Prepare to decode.
		decoded, err := zLibDecode(msg)
		if err != nil {
			log.Fatal(err)
		}

		decoder := json.NewDecoder(bytes.NewReader(decoded))
		for {
			var m emdrHeader
			if err := decoder.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

		}
	}
}

func zLibDecode(encoded string) (decoded []byte, err error) {
	b := bytes.NewBufferString(encoded)
	pipeline, err := zlib.NewReader(b)

	if err == nil {
		defer pipeline.Close()
		decoded, err = ioutil.ReadAll(pipeline)
	}

	return
}
