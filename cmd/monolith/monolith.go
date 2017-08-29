package main

import (
	"log"

	"github.com/antihax/evedata/evedata"
	_ "github.com/antihax/evedata/views"
)

// bootstrap interface
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting EVEData")
	evedata.GoServer()
	log.Println("Ended EVEData")
}
