package tsBotService

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHQ(t *testing.T) {
	ts, err := NewTSService("localhost:10011", "serveradmin", "nothinguseful")
	assert.Nil(t, err)
	assert.NotNil(t, ts)
}
