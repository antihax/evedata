package emdrConsumer

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"evedata/appContext"
	"io"
	"io/ioutil"
	"log"

	zmq "github.com/pebbe/zmq4"

	_ "github.com/go-sql-driver/mysql"
)

// [TODO] Update and complete
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

func goEMDRConsumer(c *appContext.AppContext) {

	client, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		log.Fatal(err)
	}

	// Connect
	log.Print("Subscribing to EMDR")
	if err = client.Connect("tcp://relay-us-central-1.eve-emdr.com:8050"); err != nil {
		log.Fatal(err)
	}
	client.SetSubscribe("")

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
