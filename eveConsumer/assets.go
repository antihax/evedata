package eveConsumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/antihax/evedata/esi"
	"github.com/garyburd/redigo/redis"
)

// Perform contact sync for wardecs
func (c *EVEConsumer) assetsShouldUpdate() {
	r := c.ctx.Cache.Get()
	defer r.Close()

	// Gather characters for update. Group for optimized updating.
	rows, err := c.ctx.Db.Query(
		`SELECT characterID, tokenCharacterID FROM crestTokens WHERE 
		assetCacheUntil < UTC_TIMESTAMP() AND lastStatus NOT LIKE "%Invalid refresh token%";`)
	if err != nil {
		log.Printf("Assets: Failed query: %v", err)
		return
	}

	// Loop updatable characters
	for rows.Next() {
		var (
			char      int64 // Source char
			tokenChar int64 // Token Char
		)

		err = rows.Scan(&char, &tokenChar)
		if err != nil {
			log.Printf("Assets: Failed scan: %v", err)
			continue
		}
		_, err = r.Do("SADD", "EVEDATA_assetQueue", fmt.Sprintf("%d:%d", char, tokenChar))
		if err != nil {
			log.Printf("Assets: Failed scan: %v", err)
			continue
		}
	}
	rows.Close()
}

func (c *EVEConsumer) assetsCheckQueue(r redis.Conn) error {
	ret, err := r.Do("SPOP", "EVEDATA_assetQueue")
	if err != nil {
		return err
	} else if ret == nil {
		return nil
	}

	v, err := redis.String(ret, err)
	if err != nil {
		return err
	}

	dest := strings.Split(v, ":")

	if len(dest) != 2 {
		return errors.New("Invalid asset string")
	}

	char, err := strconv.ParseInt(dest[0], 10, 64)
	if err != nil {
		return err
	}
	tokenChar, err := strconv.ParseInt(dest[1], 10, 64)
	if err != nil {
		return err
	}

	token, err := c.getToken(char, tokenChar)

	// authentication token context for destination char
	auth := context.WithValue(context.TODO(), esi.ContextOAuth2, token)

	assets, res, err := c.ctx.ESI.AssetsApi.GetCharactersCharacterIdAssets(auth, (int32)(tokenChar), nil)
	if err != nil {
		syncError(char, tokenChar, res, err)
	} else {
		syncSuccess(char, tokenChar, 200, "OK")

		for {
			tx, err := c.ctx.Db.Beginx()
			if err != nil {
				return err
			}

			tx.Exec("DELETE FROM evedata.assets WHERE characterID = ?", tokenChar)
			for _, asset := range assets {
				tx.Exec(`INSERT INTO evedata.assets
							(locationID, typeID, quantity, characterID, 
							locationFlag, itemID, locationType, isSingleton)
							VALUES (?,?,?,?,?,?,?,?);`,
					asset.LocationId, asset.TypeId, asset.Quantity, tokenChar,
					asset.LocationFlag, asset.ItemId, asset.LocationType, asset.IsSingleton)
			}

			tx.Exec(`UPDATE crestTokens SET assetCacheUntil = ? 
						WHERE characterID = ? AND tokenCharacterID = ?`,
				esi.CacheExpires(res), char, tokenChar)

			err = tx.Commit()
			if err != nil {
				fmt.Printf("Assets: %v\n", err)
			} else {
				break
			}
		}
	}

	return err
}
