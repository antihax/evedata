package config

import (
	"encoding/json"
	"os"
)

// Config stucture for the EVEData App
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
		Enabled bool
		Import  bool
		Upload  bool
		URL     string
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
	Domain string

	Redis struct {
		Address  string
		Password string
	}

	Discord struct {
		Enabled   bool
		ClientID  string
		SecretKey string
		Token     string
	}
}

// ReadConfig should be run at startup and output shared between microservices via context.
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

	if err != nil {
		return nil, err
	}
	return &c, nil
}
