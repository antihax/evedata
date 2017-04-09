package config

import (
	"time"

	"github.com/BurntSushi/toml"
)

// Config stucture for the EVEData App
type Config struct {
	UserAgent     string
	GenerateStats bool
	Database      struct {
		Import string
		Driver string
		Spec   string
	}

	Store struct {
		Key    string
		Domain string
	}

	EVEConsumer struct {
		Enabled      bool
		ZKillEnabled bool
		ZKillID      string
		Consumers    int
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
		ESIAccessToken struct {
			ClientID    string
			SecretKey   string
			RedirectURL string

			AccessToken  string
			RefreshToken string
			Expiry       time.Time
			TokenType    string
		}
	}

	Redis struct {
		Address    string
		Password   string
		Sentinel   bool
		Addresses  []string
		MasterName string
	}

	Discord struct {
		Enabled   bool
		ClientID  string
		SecretKey string
		Token     string
	}
}

// ReadConfig should be run at startup and output shared between services via context.
func ReadConfig() (*Config, error) {
	c := Config{}

	if _, err := toml.DecodeFile("config/config.conf", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
