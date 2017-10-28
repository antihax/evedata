package eveConsumer

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/evedata/models"
	"github.com/garyburd/redigo/redis"
)

func init() {

	//addConsumer("market", marketHistoryConsumer, "EVEDATA_marketHistory")
	//addConsumer("market", marketRegionConsumer, "")

	addTrigger("market", marketMaintTrigger)
	//addTrigger("market", marketHistoryTrigger)

}

// Add market history items to the queue
func marketMaintTrigger(c *EVEConsumer) (bool, error) {

	// Skip if we are not ready
	cacheUntilTime, _, err := models.GetServiceState("marketMaint")
	if err != nil {
		return false, err
	}

	// Check if it is time to update the market history
	curTime := time.Now().UTC()
	if cacheUntilTime.Before(curTime) {
		newTime := curTime.Add(time.Hour * 1)

		err = models.SetServiceState("marketMaint", newTime, 1)
		if err != nil {
			return false, err
		}

		err = models.MaintMarket()
	}
	return true, err
}

// Add market history items to the queue
func marketHistoryTrigger(c *EVEConsumer) (bool, error) {

	// Skip if we are not ready
	cacheUntilTime, _, err := models.GetServiceState("marketHistory")
	if err != nil {
		return false, err
	}

	// Check if it is time to update the market history
	curTime := time.Now().UTC()
	if cacheUntilTime.Before(curTime) {
		// We wont repeat this for 24 hours just after it updates.
		curTime = curTime.Add(time.Hour * 24)
		newTime := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 30, 0, 0, time.UTC)

		err = models.SetServiceState("marketHistory", newTime, 1)
		if err != nil {
			return false, err
		}

		// Get lists to build our requests
		regions, err := models.GetMarketRegions()
		if err != nil {
			return false, err
		}
		types, err := models.GetMarketTypes()
		if err != nil {
			return false, err
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
	return true, err
}

func (c *EVEConsumer) marketRegionAddRegion(v int, t int64, redisPtr *redis.Conn) {
	r := *redisPtr
	r.Do("ZADD", "EVEDATA_marketRegions", t, v)
}

func marketHistoryConsumer(c *EVEConsumer, redisPtr *redis.Conn) (bool, error) {
	r := *redisPtr
	ret, err := r.Do("SPOP", "EVEDATA_marketHistory")
	if err != nil {
		return false, err
	} else if ret == nil {
		return false, nil
	}
	v, err := redis.String(ret, err)
	if err != nil {
		return false, err
	}

	data := strings.Split(v, ":")
	regionID, err := strconv.Atoi(data[0])
	typeID, err := strconv.Atoi(data[1])

	// Process Market History
	h, res, err := c.ctx.ESI.ESI.MarketApi.GetMarketsRegionIdHistory(nil, (int32)(regionID), (int32)(typeID), nil)
	if err != nil {
		if res.StatusCode >= 500 {
			// Something went wrong... let's try again..
			r.Do("SADD", "EVEDATA_marketHistory", v)
		}
		return false, err
	}

	// There is nothing.
	if len(h) == 0 {
		return false, nil
	}

	var values []string

	ignoreBefore := time.Now().UTC().Add(time.Hour * 24 * -2)

	for _, e := range h {
		orderDate, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			return false, err
		}

		if orderDate.After(ignoreBefore) {
			values = append(values, fmt.Sprintf("(%q,%f,%f,%f,%d,%d,%d,%d)",
				e.Date, e.Lowest, e.Highest, e.Average,
				e.Volume, e.OrderCount, typeID, regionID))
		}
	}

	if len(values) == 0 {
		return false, nil
	}

	stmt := fmt.Sprintf("INSERT INTO evedata.market_history (date, low, high, mean, quantity, orders, itemID, regionID) VALUES \n%s ON DUPLICATE KEY UPDATE date=date", strings.Join(values, ",\n"))

	tx, err := models.Begin()
	if err != nil {
		return false, err
	}
	_, err = tx.Exec(stmt)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	err = models.RetryTransaction(tx)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Check the regions for their cache time to expire
func marketRegionConsumer(c *EVEConsumer, redisPtr *redis.Conn) (bool, error) {
	r := *redisPtr
	t := time.Now().UTC().Unix()

	// Get a list of regions by expired keys.
	if arr, err := redis.MultiBulk(r.Do("ZRANGEBYSCORE", "EVEDATA_marketRegions", 0, t)); err != nil {
		return false, err
	} else {
		// Add the region to the queue
		idList, _ := redis.Strings(arr, nil)
		for _, id := range idList {
			id, err := strconv.Atoi(id)
			if err != nil {
				return false, nil
			}
			if err := r.Send("SADD", "EVEDATA_marketOrders", id); err != nil {
				return false, err
			}
		}
	}

	// Removed the expired keys
	if err := r.Send("ZREMRANGEBYSCORE", "EVEDATA_marketRegions", 0, t); err != nil {
		return false, err
	}

	// Run the commands.
	err := r.Flush()

	return true, err
}
