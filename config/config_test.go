package config

import "testing"

func TestAppContext(t *testing.T) {
	_, err := ReadConfig("config-example.conf")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = ReadConfig("noconfig")
	if err == nil {
		t.Error(err)
		return
	}
}
