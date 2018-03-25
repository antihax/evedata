package conservator

import (
	"os"
	"testing"

	"github.com/antihax/evedata/internal/nsqhelper"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/prometheus/common/log"
	"github.com/stretchr/testify/assert"
)

var conserv *Conservator

func TestMain(m *testing.M) {
	sql := sqlhelper.NewTestDatabase()
	// Setup a hammer service
	redis := redigohelper.ConnectRedisTestPool()
	defer redis.Close()

	conserv = NewConservator(redis, sql, nsqhelper.Test, "TEST")

	inserts := []string{
		`INSERT INTO evedata.integrations
			(integrationID, entityID, address, authentication, type, services, options) 
			VALUES
			(1, 234, "127.0.0.1:10011", "serveradmin:nothinguseful", "ts3", "auth", ""),
			(2, 567, "127.0.0.2:10011", "serveradmin:nothinguseful", "ts3", "", "")
			ON DUPLICATE KEY UPDATE integrationID=integrationID`,
		`INSERT INTO evedata.integrationChannels
			(integrationID, channelID, services, options) 
			VALUES
			(1, 12345, "kill", ""),
			(1, 12346, "locator", ""),
			(2, 12347, "kill,locator,structure", "")
			ON DUPLICATE KEY UPDATE integrationID=integrationID`,
		`INSERT INTO evedata.sharing
			(characterID, tokenCharacterID, entityID, types) 
			VALUES
			(1123123, 24234234, 235, "kill,locator,structure"),
			(1123123, 24234234, 234, "kill"),
			(1123125, 24234235, 234, "locator"),
			(1123123, 24234234, 567, "war,locator,structure")
			ON DUPLICATE KEY UPDATE characterID=characterID`,
	}

	for _, insert := range inserts {
		if _, err := conserv.db.Exec(insert); err != nil {
			log.Fatal(err)
		}
	}

	// Run tests
	r := m.Run()
	conserv.Close()
	os.Exit(r)
}

func TestConservator(t *testing.T) {
	err := conserv.loadServices()
	assert.Nil(t, err)

	si, ok := conserv.services.Load(int32(1))
	assert.True(t, ok)
	assert.NotNil(t, si)

	service := si.(Service)
	assert.Equal(t, "ts3", service.Type)

	err = conserv.loadChannels()
	assert.Nil(t, err)
	ci, ok := conserv.channels.Load("12345")
	assert.True(t, ok)
	assert.NotNil(t, ci)

	channel := ci.(Channel)
	assert.Equal(t, "kill", channel.Services)

	err = conserv.loadShares()
	assert.Nil(t, err)

	assert.Equal(t, int32(234), conserv.notifications["kill"][int32(24234234)][0].EntityID)
	assert.Equal(t, int32(234), conserv.notifications["locator"][int32(24234235)][0].EntityID)
	assert.Equal(t, int32(567), conserv.notifications["structure"][int32(24234234)][0].EntityID)

	// Test we can properly delete entries
	conserv.db.Exec("DELETE FROM evedata.sharing")
	err = conserv.loadShares()
	assert.Nil(t, err)
	assert.Zero(t, len(conserv.notifications["structure"]))
	assert.Zero(t, len(conserv.notifications["kill"]))
	assert.Zero(t, len(conserv.notifications["locator"]))
}
