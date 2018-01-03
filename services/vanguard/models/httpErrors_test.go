package models

import (
	"log"
	"net/http"
	"testing"
)

func TestAddHTTPError(t *testing.T) {
	req, err := http.NewRequest("GET", "/hi", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	res := &http.Response{}
	err = AddHTTPError(req, res)
	if err != nil {
		log.Fatal(err)
		return
	}
}
