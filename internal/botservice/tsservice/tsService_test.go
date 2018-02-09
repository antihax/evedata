package tsservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTS3(t *testing.T) {
	ts, err := NewTSService("localhost:10011", "serveradmin", "nothinguseful")
	assert.Nil(t, err)
	assert.NotNil(t, ts)
	serverList, err := ts.GetServerList()
	assert.Nil(t, err)
	assert.NotNil(t, serverList)
	assert.Equal(t, serverList["1"], "TeamSpeak ]I[ Server")

	err = ts.UseServer(1)
	assert.Nil(t, err)

	channelList, err := ts.GetChannelList()
	assert.Nil(t, err)
	assert.Equal(t, channelList["1"], "Default Channel")

	err = ts.SendMessageToChannel("1", "Test Channel Message")
	assert.Nil(t, err)

	// This should fail.
	err = ts.SendMessageToUser("1234567", "Test User Message")
	assert.NotNil(t, err)

	err = ts.SendMessageToServer("1", "Test Server Message")
	assert.Nil(t, err)

}
