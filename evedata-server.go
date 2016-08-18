package main

import (
	"evedata/server"
	_ "evedata/views"
)

// bootstrap interface

func main() {
	evedata.GoServer()
}
