package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Database struct {
		Import string
		Driver string
		Spec   string
	}
	Store struct {
		Key []byte
	}
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
	return &c, nil
}
