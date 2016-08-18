package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	Database struct {
		Import string
		Driver string
		Spec   string
	}
	Store struct {
		Key string
	}
	EMDRCrestBridge struct {
		Enabled       bool
		Import        bool
		Upload        bool
		URL           string
		MaxGoRoutines int64
	}
	CREST struct {
		SSO struct {
			ClientID    string
			SecretKey   string
			RedirectURL string
		}
		Token struct {
			ClientID    string
			SecretKey   string
			RedirectURL string
		}
	}
	ServerIP         string
	Domain           string
	MemcachedAddress string
}

func ReadConfig() (*Config, error) {
	c := Config{}

	// read configuration
	file, err := os.Open("config/config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}
	c.ServerIP, err = checkIP()
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func checkIP() (string, error) {
	rsp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(buf)), nil
}
