package eveConsumer

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/evedata/esi"
	"github.com/antihax/evedata/models"
	"github.com/garyburd/redigo/redis"
)

// Add market history items to the queue
func (c *EVEConsumer) marketMaintTrigger() error {

	// Skip if we are not ready
	cacheUntilTime, _, err := models.GetServiceState("marketMaint")
	if err != nil {
		return err
	}

	// Check if it is time to update the market history
	curTime := time.Now().UTC()
	if cacheUntilTime.Before(curTime) {
		// We wont repeat this for 24 hours just after it updates.
		newTime := curTime.Add(time.Hour * 3)

		err = models.SetServiceState("marketMaint", newTime, 1)
		if err != nil {
			return err
		}

		err = models.MaintMarket()
	}
	return err
}

// Add market history items to the queue
func (c *EVEConsumer) marketHistoryUpdateTrigger() error {

	// Skip if we are not ready
	cacheUntilTime, _, err := models.GetServiceState("marketHistory")
	if err != nil {
		return err
	}

	// Check if it is time to update the market history
	curTime := time.Now().UTC()
	if cacheUntilTime.Before(curTime) {
		// We wont repeat this for 24 hours just after it updates.
		curTime = curTime.Add(time.Hour * 24)
		newTime := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 30, 0, 0, time.UTC)

		err = models.SetServiceState("marketHistory", newTime, 1)
		if err != nil {
			return err
		}

		// Get lists to build our requests
		regions, err := models.GetMarketRegions()
		if err != nil {
			return err
		}
		types, err := models.GetMarketTypes()
		if err != nil {
			return err
		}

		// Get a redis connection from the pool
		red := c.ctx.Cache.Get()
		defer red.Close()

		// Load types into redis queue
		// Build a pipeline request to add the region IDs to redis
		for _, r := range regions {
			// Add regions into marketOrders just in case they disapear.
			// NX = Don't update score if element exists
			red.Send("ZADD", "EVEDATA_marketRegions", "NX", time.Now().UTC().Unix(), r.RegionID)
			for _, t := range types {
				red.Send("SADD", "EVEDATA_marketHistory", fmt.Sprintf("%d:%d", r.RegionID, t.TypeID))
			}
		}

		// Send the request to add
		red.Flush()
	}
	return err
}

func (c *EVEConsumer) marketOrderCheckQueue(r redis.Conn) error {
	ret, err := r.Do("SPOP", "EVEDATA_marketOrders")
	if err != nil {
		return err
	} else if ret == nil {
		return nil
	}
	v, err := redis.Int(ret, err)
	if err != nil {
		return err
	}

	var page int32 = 1
	c.marketRegionAddRegion(v, time.Now().UTC().Unix()+(60*60), r)
	for {
		b, res, err := c.ctx.ESI.MarketApi.GetMarketsRegionIdOrders((int32)(v), "all", map[string]interface{}{"page": page})
		if err != nil {
			return err
		} else if len(b) == 0 { // end of the pages
			break
		}
		var values []string
		for _, e := range b {
			var buy byte
			if e.IsBuyOrder == true {
				buy = 1
			} else {
				buy = 0
			}
			values = append(values, fmt.Sprintf("(%d,%f,%d,%d,%d,%d,%d,'%s',%d,%d,%d,UTC_TIMESTAMP())",
				e.OrderId, e.Price, e.VolumeRemain, e.TypeId, e.VolumeTotal, e.MinVolume,
				buy, e.Issued.UTC().Format("2006-01-02 15:04:05"), e.Duration, e.LocationId, (int32)(v)))
		}

		stmt := fmt.Sprintf(`INSERT IGNORE INTO evedata.market (orderID, price, remainingVolume, typeID, enteredVolume, minVolume, bid, issued, duration, stationID, regionID, reported)
				VALUES %s
				ON DUPLICATE KEY UPDATE price=VALUES(price),
					remainingVolume=VALUES(remainingVolume),
					issued=VALUES(issued),
					duration=VALUES(duration),
					reported=VALUES(reported),
					done=0;
					`, strings.Join(values, ",\n"))

		for {
			tx, err := c.ctx.Db.Begin()
			if err != nil {
				log.Printf("%s", err)
				continue
			}
			_, err = tx.Exec(stmt)
			if err != nil {
				log.Printf("%s", err)
				continue
			}

			err = tx.Commit()
			if err != nil {
				if strings.Contains(err.Error(), "1213") == false {
					log.Printf("Market: %v\n", err)
					break
				} else {
					continue
				}
			} else {
				break // success
			}
		}

		// Cache the greater of one hour, or the returned cache-control
		cacheUntil := max(time.Now().UTC().Add(time.Hour*1).Unix(), esi.CacheExpires(res).UTC().Unix())
		c.marketRegionAddRegion(v, cacheUntil, r)

		// Next page
		page++
	}
	return nil
}

func (c *EVEConsumer) marketRegionAddRegion(v int, t int64, r redis.Conn) {
	r.Do("ZADD", "EVEDATA_marketRegions", t, v)
}

func (c *EVEConsumer) marketHistoryCheckQueue(r redis.Conn) error {
	ret, err := r.Do("SPOP", "EVEDATA_marketHistory")
	if err != nil {
		return err
	} else if ret == nil {
		return nil
	}
	v, err := redis.String(ret, err)
	if err != nil {
		return err
	}

	data := strings.Split(v, ":")
	regionID, err := strconv.Atoi(data[0])
	typeID, err := strconv.Atoi(data[1])

	// Process Market History
	h, _, err := c.ctx.ESI.MarketApi.GetMarketsRegionIdHistory((int32)(regionID), (int32)(typeID), nil)
	if err != nil {
		// Something went wrong... let's try again..
		r.Do("SADD", "EVEDATA_marketHistory", v)
		return err
	}

	// There is nothing.
	if len(h) == 0 {
		return nil
	}

	var values []string

	for _, e := range h {
		values = append(values, fmt.Sprintf("('%s',%f,%f,%f,%d,%d,%d,%d)",
			e.Date.Format("2006-01-02"), e.Lowest, e.Highest, e.Average,
			e.Volume, e.OrderCount, typeID, regionID))
	}

	stmt := fmt.Sprintf("INSERT IGNORE INTO evedata.market_history (date, low, high, mean, quantity, orders, itemID, regionID) VALUES \n%s", strings.Join(values, ",\n"))

	for {
		tx, err := c.ctx.Db.Begin()
		if err != nil {
			log.Printf("%s", err)
			break
		}
		_, err = tx.Exec(stmt)
		if err != nil {
			log.Printf("%s", err)
			break
		}

		err = tx.Commit()
		if err != nil {
			if strings.Contains(err.Error(), "1213") == false {
				log.Printf("Market: %v\n", err)
				break
			} else {
				continue
			}
		} else {
			break // success
		}
	}

	return nil
}

// Check the regions for their cache time to expire
func (c *EVEConsumer) marketRegionCheckQueue(r redis.Conn) error {
	t := time.Now().UTC().Unix()

	// Get a list of regions by expired keys.
	if arr, err := redis.MultiBulk(r.Do("ZRANGEBYSCORE", "EVEDATA_marketRegions", 0, t)); err != nil {
		return err
	} else {
		// Add the region to the queue
		idList, _ := redis.Strings(arr, nil)
		for _, id := range idList {
			id, err := strconv.Atoi(id)
			if err != nil {
				return nil
			}
			if err := r.Send("SADD", "EVEDATA_marketOrders", id); err != nil {
				return err
			}
		}
	}

	// Removed the expired keys
	if err := r.Send("ZREMRANGEBYSCORE", "EVEDATA_marketRegions", 0, t); err != nil {
		return err
	}

	// Run the commands.
	err := r.Flush()

	return err
}
