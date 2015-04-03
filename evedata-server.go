package main

import (
	"evedata/server"
	"evedata/views"
)

func main() {
	views.Init()
	evedata.GoServer()
}
