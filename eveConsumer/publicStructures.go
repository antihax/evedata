package eveConsumer

import "log"

func (c *EVEConsumer) checkPublicStructures() {
	log.Printf("EVEConsumer: collecting structures")
	err := c.collectStructuresFromESI()
	if err != nil {
		log.Printf("EVEConsumer: collecting structures: %v", err)
	}
}

func (c *EVEConsumer) collectStructuresFromESI() error {
	r := struct {
		Value int
		Wait  int
	}{0, 0}

	if err := c.ctx.Db.Get(&r, `
		SELECT value, TIME_TO_SEC(TIMEDIFF(nextCheck, UTC_TIMESTAMP())) AS wait
			FROM states 
			WHERE state = 'structures'
			LIMIT 1;
		`); err != nil {
		return err
	}

	if r.Wait >= 0 {
		return nil
	}

	w, err := c.ctx.ESI.UniverseApi.GetUniverseStructures(nil)
	if err != nil {
		return err
	}

	// Update state so we dont have two polling at once.
	_, err = c.ctx.Db.Exec("UPDATE states SET value = 1, nextCheck =? WHERE state = 'structures' LIMIT 1")

	for _, s := range w {
		c.updateStructure(s)
	}
	return nil
}

func (c *EVEConsumer) updateStructure(s int64) error {
	struc, err := c.ctx.ESI.UniverseApi.GetUniverseStructuresStructureId(c.ctx.ESIPublicToken, s, nil)
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
