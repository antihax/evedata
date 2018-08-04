package nail

import (
	"github.com/antihax/evedata/internal/datapackages"

	"github.com/antihax/evedata/internal/gobcoder"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("structure", func(s *Nail, consumer *nsq.Consumer) {
		consumer.AddHandler(s.wait(nsq.HandlerFunc(s.structureHandler)))
	})
}

func (s *Nail) structureHandler(message *nsq.Message) error {
	b := datapackages.Structure{}
	err := gobcoder.GobDecoder(message.Body, &b)
	if err != nil {
		return err
	}

	// Push into the denormalized table. This table is volatile.
	err = s.doSQL(`INSERT INTO staStations
		(stationID, solarSystemID, corporationID, stationName, x, y, z, constellationID, regionID)
		VALUES(?,?,?,?,?,?,?,evedata.constellationIDBySolarSystem(solarSystemID),evedata.regionIDBySolarSystem(solarSystemID))
		ON DUPLICATE KEY UPDATE stationName=VALUES(stationName),solarSystemID=VALUES(solarSystemID),corporationID=VALUES(corporationID),
		x=VALUES(x),y=VALUES(y),z=VALUES(z),constellationID=evedata.constellationIDBySolarSystem(VALUES(solarSystemID)),regionID=evedata.regionIDBySolarSystem(VALUES(solarSystemID));`,
		b.StructureID, b.Structure.SolarSystemId, b.Structure.OwnerId, b.Structure.Name, b.Structure.Position.X, b.Structure.Position.Y, b.Structure.Position.Z)
	if err != nil {
		return err
	}

	// Insert into our table for tracking.
	err = s.doSQL(`INSERT INTO evedata.structures
		(stationID, solarSystemID, stationName, x, y, z, ownerID, typeID, updated)
		VALUES(?,?,?,?,?,?,?,?, UTC_TIMESTAMP())
		ON DUPLICATE KEY UPDATE stationName=VALUES(stationName),solarSystemID=VALUES(solarSystemID),ownerID=VALUES(ownerID),typeID=VALUES(typeID),
		x=VALUES(x),y=VALUES(y),z=VALUES(z);`,
		b.StructureID, b.Structure.SolarSystemId, b.Structure.Name,
		b.Structure.Position.X, b.Structure.Position.Y, b.Structure.Position.Z,
		b.Structure.OwnerId, b.Structure.TypeId)

	return err
}
