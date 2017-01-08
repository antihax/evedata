package main

import (
	log "github.com/Sirupsen/logrus"
	evedata "github.com/antihax/evedata/server"
	_ "github.com/antihax/evedata/views"
)

// bootstrap interface
func main() {
	log.Println("Starting EVEData")
	evedata.GoServer()
	log.Println("Ended EVEData")
}
