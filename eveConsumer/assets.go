package eveConsumer

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
)

func init() {
	addConsumer("assets", assetsConsumer, "EVEDATA_assetQueue")
	addTrigger("assets", assetsTrigger)
}

// assetsConsumer handles gathering assets from API
func assetsConsumer(c *EVEConsumer, redisPtr *redis.Conn) (bool, error) {
	// Dereference the redis pointer.
	r := *redisPtr

	// POP some work of the queue
	ret, err := r.Do("SPOP", "EVEDATA_assetQueue")
	if err != nil {
		return false, err
	} else if ret == nil {
		return false, nil
	}

	// Get the asset string
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

	char64, err := strconv.ParseInt(dest[0], 10, 64)
	if err != nil {
		return false, err
	}
	tokenChar64, err := strconv.ParseInt(dest[1], 10, 64)
	if err != nil {
		return false, err
	}
	char := int32(char64)
	tokenChar := int32(tokenChar64)

	// Get the OAuth2 Token from the database.
	token, err := c.ctx.TokenStore.GetTokenSource(char, tokenChar)
	if err != nil {
		return false, err
	}

	// Put the token into a context for the API client
	auth := context.WithValue(context.TODO(), goesi.ContextOAuth2, token)

	fmt.Printf("%+v\n", auth)
	assets, res, err := c.ctx.ESI.ESI.AssetsApi.GetCharactersCharacterIdAssets(auth, tokenChar, nil)
	if err != nil {
		tokenError(char, tokenChar, res, err)
		return false, err
	}

	tokenSuccess(char, tokenChar, 200, "OK")

	tx, err := models.Begin()
	if err != nil {
		return false, err
	}

	// Delete all the current assets. Reinsert everything.
	tx.Exec("DELETE FROM evedata.assets WHERE characterID = ?", tokenChar)

	var values []string

	// Skip if we have no assests
	if len(assets) != 0 {

		// Dump all assets into the DB.
		for _, asset := range assets {
			values = append(values, fmt.Sprintf("(%d,%d,%d,%d,%q,%d,%q,%v)",
				asset.LocationId, asset.TypeId, asset.Quantity, tokenChar,
				asset.LocationFlag, asset.ItemId, asset.LocationType, asset.IsSingleton))
		}
		stmt := fmt.Sprintf(`INSERT INTO evedata.assets
							(locationID, typeID, quantity, characterID, 
							locationFlag, itemID, locationType, isSingleton)
		VALUES %s ON DUPLICATE KEY UPDATE locationID = locationID;`, strings.Join(values, ",\n"))

		_, err = tx.Exec(stmt)
		if err != nil {
			tx.Rollback()
			return false, err
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

	// Gather characters for update. Group for optimized updating.
	rows, err := c.ctx.Db.Query(
		`SELECT characterID, tokenCharacterID FROM evedata.crestTokens WHERE 
		assetCacheUntil < UTC_TIMESTAMP() AND lastStatus NOT LIKE "%400 Bad Request%" AND 
		scopes LIKE "%esi-assets.read_assets.v1%";`)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	r := c.ctx.Cache.Get()
	defer r.Close()

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
	return true, err
}
