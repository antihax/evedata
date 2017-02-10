package eveConsumer

import (
	"time"

	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"

	"golang.org/x/net/context"
)

func init() {
	addTrigger("structures", structuresTrigger)
	addConsumer("structures", structureConsumer, "EVEDATA_structureQueue")
}

func structureConsumer(c *EVEConsumer, r redis.Conn) (bool, error) {
	// POP some work of the queue
	ret, err := r.Do("SPOP", "EVEDATA_structureQueue")
	if err != nil {
		return false, err
	} else if ret == nil {
		return false, nil
	}

	v, err := redis.Int64(ret, err)
	if err != nil {
		return false, err
	}

	err = c.updateStructure(v)
	return true, err
}

func structuresTrigger(c *EVEConsumer) (bool, error) {
	nextCheck, _, err := models.GetServiceState("structures")
	if err != nil {
		return false, err
	} else if nextCheck.After(time.Now().UTC()) {
		return false, nil
	}

	w, r, err := c.ctx.ESI.V1.UniverseApi.GetUniverseStructures(nil)
	cache := goesi.CacheExpires(r)
	if err != nil {
		return false, err
	}

	// Update state so we dont have two polling at once.
	err = models.SetServiceState("structures", cache, 1)
	if err != nil {
		return false, err
	}

	redis := c.ctx.Cache.Get()
	for _, s := range w {
		// Build a pipeline request to add the structure IDs to redis
		redis.Send("SADD", "EVEDATA_structureQueue", s)

		// Send the request to add

	}
	redis.Flush()
	redis.Close()

	stations, err := c.ctx.EVE.ConquerableStationsListXML()
	if err != nil {
		return false, err
	}

	for _, s := range stations.Stations {
		_, err = c.ctx.Db.Exec(`
		INSERT INTO staStations
			(	stationID, 
				solarSystemID, 
				stationName, 
				corporationID, 
				constellationID, 
				regionID)
			VALUES(?,?,?,?, evedata.constellationIDBySolarSystem(solarSystemID), evedata.regionIDBySolarSystem(solarSystemID))

			ON DUPLICATE KEY UPDATE 	stationName=VALUES(stationName),
										corporationID=VALUES(corporationID);`,
			s.StationID,
			s.SolarSystemID,
			s.StationName,
			s.CorporationID)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (c *EVEConsumer) updateStructure(s int64) error {
	ctx := context.WithValue(context.TODO(), goesi.ContextOAuth2, c.ctx.ESIPublicToken)
	struc, _, err := c.ctx.ESI.V1.UniverseApi.GetUniverseStructuresStructureId(ctx, s, nil)
	if err != nil {
		return err
	}

	// Push into the denormalized table. This table is volitile.
	_, err = c.ctx.Db.Exec(`INSERT INTO staStations
					(stationID, solarSystemID, stationName, x, y, z, constellationID, regionID)
					VALUES(?,?,?,?,?,?,evedata.constellationIDBySolarSystem(solarSystemID),evedata.regionIDBySolarSystem(solarSystemID))
					ON DUPLICATE KEY UPDATE stationName=VALUES(stationName),solarSystemID=VALUES(solarSystemID),
					x=VALUES(x),y=VALUES(y),z=VALUES(z),constellationID=evedata.constellationIDBySolarSystem(VALUES(solarSystemID)),regionID=evedata.regionIDBySolarSystem(VALUES(solarSystemID));`,
		s, struc.SolarSystemId, struc.Name, struc.Position.X, struc.Position.Y, struc.Position.Z)
	if err != nil {
		return err
	}

	// Insert into our table for tracking.
	_, err = c.ctx.Db.Exec(`INSERT INTO evedata.structures
					(stationID, solarSystemID, stationName, x, y, z, updated)
					VALUES(?,?,?,?,?,?, UTC_TIMESTAMP())
					ON DUPLICATE KEY UPDATE stationName=VALUES(stationName),solarSystemID=VALUES(solarSystemID),
					x=VALUES(x),y=VALUES(y),z=VALUES(z);`,
		s, struc.SolarSystemId, struc.Name, struc.Position.X, struc.Position.Y, struc.Position.Z)

	if err != nil {
		return err
	}

	return nil
}

func (c *EVEConsumer) updateStation(s int64) error {
	ctx := context.WithValue(context.TODO(), goesi.ContextOAuth2, c.ctx.ESIPublicToken)
	struc, _, err := c.ctx.ESI.V1.UniverseApi.GetUniverseStructuresStructureId(ctx, s, nil)
	if err != nil {
		return err
	}

	_, err = c.ctx.Db.Exec(`INSERT INTO staStations
					(stationID, solarSystemID, stationName, x, y, z, constellationID, regionID)
					VALUES(?,?,?,?,?,?,evedata.constellationIDBySolarSystem(solarSystemID),evedata.regionIDBySolarSystem(solarSystemID))
					ON DUPLICATE KEY UPDATE stationName=VALUES(stationName),solarSystemID=VALUES(solarSystemID),
					x=VALUES(x),y=VALUES(y),z=VALUES(z),constellationID=evedata.constellationIDBySolarSystem(VALUES(solarSystemID)),regionID=evedata.regionIDBySolarSystem(VALUES(solarSystemID));`,
		s, struc.SolarSystemId, struc.Name, struc.Position.X, struc.Position.Y, struc.Position.Z)
	if err != nil {
		return err
	}

	return nil
}
