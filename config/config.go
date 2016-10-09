package config

import "github.com/BurntSushi/toml"

// Config stucture for the EVEData App
type Config struct {
	Database struct {
		Import string
		Driver string
		Spec   string
	}
	Store struct {
		Key    string
		Domain string
	}
	EMDRCrestBridge struct {
		Enabled bool
		Import  bool
		Upload  bool
		URL     string
	}
	EVEConsumer struct {
		Enabled bool
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

	if _, err := toml.DecodeFile("config/config.conf", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
