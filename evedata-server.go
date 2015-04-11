package main

import (
	"evedata/server"
	_ "evedata/views" // bootstrap the interface
)

func main() {
	evedata.GoServer()
}
