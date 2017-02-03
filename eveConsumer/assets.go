package eveConsumer

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"
	"github.com/antihax/goesi/v1"
	"github.com/garyburd/redigo/redis"
)

func init() {
	addConsumer("assets", assetsConsumer, "EVEDATA_assetQueue")
	addTrigger("assets", assetsTrigger)
}

func assetsConsumer(c *EVEConsumer, r redis.Conn) (bool, error) {
	// POP some work of the queue
	ret, err := r.Do("SPOP", "EVEDATA_assetQueue")
	if err != nil {
		return false, err
	} else if ret == nil {
		return false, nil
	}

	v, err := redis.String(ret, err)
	if err != nil {
		return false, err
	}

	// Split our char:tokenChar string
	dest := strings.Split(v, ":")

	// Quick sanity check
	if len(dest) != 2 {
		return false, errors.New("Invalid asset string")
	}

	char, err := strconv.ParseInt(dest[0], 10, 64)
	if err != nil {
		return false, err
	}
	tokenChar, err := strconv.ParseInt(dest[1], 10, 64)
	if err != nil {
		return false, err
	}

	// Get the OAuth2 Token from the database.
	token, err := c.getToken(char, tokenChar)

	// Put the token into a context for the API client
	auth := context.WithValue(context.TODO(), goesiv1.ContextOAuth2, token)

	assets, res, err := c.ctx.ESI.V1.AssetsApi.GetCharactersCharacterIdAssets(auth, (int32)(tokenChar), nil)
	if err != nil {
		tokenError(char, tokenChar, res, err)
		return false, err
	} else {
		tokenSuccess(char, tokenChar, 200, "OK")

		tx, err := models.Begin()
		if err != nil {
			return false, err
		}

		// Delete all the current assets. Reinsert everything.
		tx.Exec("DELETE FROM evedata.assets WHERE characterID = ?", tokenChar)

		// Dump all assets into the DB.
		for _, asset := range assets {
			tx.Exec(`INSERT INTO evedata.assets
							(locationID, typeID, quantity, characterID, 
							locationFlag, itemID, locationType, isSingleton)
							VALUES (?,?,?,?,?,?,?,?);`,
				asset.LocationId, asset.TypeId, asset.Quantity, tokenChar,
				asset.LocationFlag, asset.ItemId, asset.LocationType, asset.IsSingleton)
		}

		// Update our cacheUntil flag
		tx.Exec(`UPDATE evedata.crestTokens SET assetCacheUntil = ? 
						WHERE characterID = ? AND tokenCharacterID = ?`,
			goesi.CacheExpires(res), char, tokenChar)

		// Retry the transaction if we get deadlocks
		err = models.RetryTransaction(tx)
		if err != nil {
			return false, err
		}
	}

	return true, err
}

// Update character assets
func assetsTrigger(c *EVEConsumer) (bool, error) {
	r := c.ctx.Cache.Get()
	defer r.Close()

	// Gather characters for update. Group for optimized updating.
	rows, err := c.ctx.Db.Query(
		`SELECT characterID, tokenCharacterID FROM evedata.crestTokens WHERE 
		assetCacheUntil < UTC_TIMESTAMP() AND lastStatus NOT LIKE "%400 Bad Request%" AND 
		scopes LIKE "%esi-assets.read_assets.v1%";`)
	if err != nil {
		return false, err
	}

	// Loop updatable characters
	for rows.Next() {
		var (
			char      int64 // Source char
			tokenChar int64 // Token Char
		)

		err = rows.Scan(&char, &tokenChar)
		if err != nil {
			return false, err
		}

		// Add the job to the queue
		_, err = r.Do("SADD", "EVEDATA_assetQueue", fmt.Sprintf("%d:%d", char, tokenChar))
		if err != nil {
			return false, err
		}
	}
	err = rows.Close()
	return true, err
}
