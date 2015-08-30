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

/*
{
   "uploadKeys":[
      {
         "name":"EVE Market Data Relay",
         "key":"0"
      },
      {
         "name":"EMDR",
         "key":"4132bdc57b0acc0e60abd3501eae58b8b13e4426"
      }
   ],
   "generator":{
      "version":"2.0.0.3758",
      "name":"EVEMon.MarketUnifiedUploader"
   },
   "currentTime":"2015-08-15T23:02:06+00:00",
   "resultType":"orders",
   "version":"0.1",
   "rowsets":[
      {
         "typeID":22546,
         "rows":[
            [
               235000000.0,
               5,
               32767,
               4211561384,
               5,
               1,
               false,
               "2015-08-05T05:03:12+00:00",
               90,
               61000912,
               30002162
            ],
            [
               232000000.0,
               5,
               32767,
               4211561733,
               5,
               1,
               false,
               "2015-08-05T05:03:42+00:00",
               90,
               61000912,
               30002162
            ],
            [
               231000000.0,
               6,
               32767,
               4211562813,
               6,
               1,
               false,
               "2015-08-05T05:05:00+00:00",
               90,
               61000912,
               30002162
            ],
            [
               243098686.0,
               1,
               32767,
               4153151948,
               6,
               1,
               false,
               "2015-06-17T21:31:13+00:00",
               90,
               61000752,
               30002184
            ],
            [
               23803000.8900000006,
               2,
               32767,
               4211988767,
               2,
               1,
               true,
               "2015-08-11T13:44:04+00:00",
               90,
               61000669,
               30002168
            ],
            [
               22992274.9899999984,
               1,
               32767,
               3820631074,
               1,
               1,
               true,
               "2015-05-30T13:57:36+00:00",
               90,
               60014915,
               30002175
            ],
            [
               23802002.9899999984,
               5,
               32767,
               4079095984,
               5,
               1,
               true,
               "2015-08-11T03:05:10+00:00",
               90,
               61000396,
               30002178
            ]
         ],
         "regionID":10000025,
         "generatedAt":"2015-08-13T06:07:38+00:00"
      }
   ],
   "columns":[
      "price",
      "volRemaining",
      "range",
      "orderID",
      "volEntered",
      "minVolume",
      "bid",
      "issueDate",
      "duration",
      "stationID",
      "solarSystemID"
   ]
}
*/

/*
{
   "uploadKeys":[
      {
         "name":"EVE Market Data Relay",
         "key":"0"
      },
      {
         "name":"EMDR",
         "key":"4132bdc57b0acc0e60abd3501eae58b8b13e4426"
      }
   ],
   "generator":{
      "version":"2.0.0.3758",
      "name":"EVEMon.MarketUnifiedUploader"
   },
   "currentTime":"2015-08-15T23:02:06+00:00",
   "resultType":"orders",
   "version":"0.1",
   "rowsets":[
{
   "uploadKeys":[
      {
         "name":"Eve-Central",
         "key":"0"
      },
      {
         "name":"EVE Market Data Relay",
         "key":"0"
      },
      {
         "name":"EMDR",
         "key":"8ffb8eba5a3cecc990a5e463b356650b05516d1e"
      }
   ],
   "generator":{
      "version":"1.0.0.0",
      "name":"EveHQ"
   },
   "currentTime":"2015-08-15T23:02:13+00:00",
   "resultType":"history",
   "version":"0.1",
   "rowsets":[
      {
         "typeID":9495,
         "rows":[
            [
               "2015-08-15T00:00:00+00:00",
               6,
               11,
               50000.02,
               465198.0,
               394126.91
            ]
         ],
         "regionID":10000032,
         "generatedAt":"2015-08-15T23:01:56+00:00"
      }
   ],
   "columns":[
      "date",
      "orders",
      "quantity",
      "low",
      "high",
      "average"
   ]
}*/

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
