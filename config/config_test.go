package config

import "testing"

func TestAppContext(t *testing.T) {
	_, err := ReadConfig("config.conf")
	if err != nil {
		t.Error(err)
		return
	}
}
