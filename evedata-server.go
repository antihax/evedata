package main

import (
	evedata "github.com/antihax/evedata/server"
	_ "github.com/antihax/evedata/views"
)

// bootstrap interface
func main() {
	evedata.GoServer()
}
