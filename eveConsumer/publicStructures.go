package eveConsumer

import (
	"evedata/esi"
	"evedata/models"
	"log"
	"time"

	"golang.org/x/net/context"
)

func (c *EVEConsumer) checkPublicStructures() {

	err := c.collectStructuresFromESI()
	if err != nil {
		log.Printf("EVEConsumer: collecting structures: %v", err)
	}
}

func (c *EVEConsumer) collectStructuresFromESI() error {
	nextCheck, _, err := models.GetServiceState("structures")
	if err != nil {
		return err
	} else if nextCheck.After(time.Now()) {
		return nil
	}

	log.Printf("EVEConsumer: collecting structures")
	w, cache, err := c.ctx.ESI.UniverseApi.GetUniverseStructures(nil)
	if err != nil {
		return err
	}

	// Update state so we dont have two polling at once.
	err = models.SetServiceState("structures", cache, 1)
	if err != nil {
		return err
	}

	for _, s := range w {
		c.updateStructure(s)
	}
	return nil
}

func (c *EVEConsumer) updateStructure(s int64) error {
	ctx := context.WithValue(context.TODO(), esi.ContextAuth, c.ctx.ESIPublicToken)
	struc, _, err := c.ctx.ESI.UniverseApi.GetUniverseStructuresStructureId(ctx, s, nil)
	if err != nil {
		return err
	}

	_, err = c.ctx.Db.Exec(`INSERT INTO staStations
					(stationID, solarSystemID, stationName, x, y, z, constellationID, regionID)
					VALUES(?,?,?,?,?,?,constellationIDBySolarSystem(solarSystemID),regionIDBySolarSystem(solarSystemID))
					ON DUPLICATE KEY UPDATE stationName=VALUES(stationName),solarSystemID=VALUES(solarSystemID),
					x=VALUES(x),y=VALUES(y),z=VALUES(z),constellationID=constellationIDBySolarSystem(VALUES(solarSystemID)),regionID=regionIDBySolarSystem(VALUES(solarSystemID));`,
		s, struc.SolarSystemId, struc.Name, struc.Position.X, struc.Position.Y, struc.Position.Z)
	if err != nil {
		return err
	}

	return nil
}
