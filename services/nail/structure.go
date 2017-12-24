package nail

import (
	"fmt"
	"strings"

	"github.com/antihax/evedata/internal/datapackages"

	"github.com/antihax/evedata/internal/gobcoder"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("structureOrders", spawnStructureMarketConsumer)
	AddHandler("structure", spawnStructureConsumer)
}

func spawnStructureMarketConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.structureMarketHandler)))
}

func (s *Nail) structureMarketHandler(message *nsq.Message) error {
	b := datapackages.StructureOrders{}
	err := gobcoder.GobDecoder(message.Body, &b)
	if err != nil {
		return err
	}

	if len(b.Orders) == 0 {
		return nil
	}

	var values []string
	for _, e := range b.Orders {
		var buy byte
		if e.IsBuyOrder {
			buy = 1
		} else {
			buy = 0
		}
		values = append(values, fmt.Sprintf("(%d,%f,%d,%d,%d,%d,%d,%q,%d,%d,evedata.regionIDByStructureID(%d),UTC_TIMESTAMP())",
			e.OrderId, e.Price, e.VolumeRemain, e.TypeId, e.VolumeTotal, e.MinVolume,
			buy, e.Issued.UTC().Format("2006-01-02 15:04:05"), e.Duration, e.LocationId, e.LocationId))
	}

	stmt := fmt.Sprintf(`INSERT INTO evedata.market (orderID, price, remainingVolume, typeID, enteredVolume, minVolume, bid, issued, duration, stationID, regionID, reported)
			VALUES %s
			ON DUPLICATE KEY UPDATE price=VALUES(price),
				remainingVolume=VALUES(remainingVolume),
				issued=VALUES(issued),
				duration=VALUES(duration),
				reported=VALUES(reported);
				`, strings.Join(values, ",\n"))
	return s.doSQL(stmt)
}

func spawnStructureConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.structureHandler)))
}

func (s *Nail) structureHandler(message *nsq.Message) error {
	b := datapackages.Structure{}
	err := gobcoder.GobDecoder(message.Body, &b)
	if err != nil {
		return err
	}

	// Push into the denormalized table. This table is volitile.
	err = s.doSQL(`INSERT INTO staStations
		(stationID, solarSystemID, stationName, x, y, z, constellationID, regionID)
		VALUES(?,?,?,?,?,?,evedata.constellationIDBySolarSystem(solarSystemID),evedata.regionIDBySolarSystem(solarSystemID))
		ON DUPLICATE KEY UPDATE stationName=VALUES(stationName),solarSystemID=VALUES(solarSystemID),
		x=VALUES(x),y=VALUES(y),z=VALUES(z),constellationID=evedata.constellationIDBySolarSystem(VALUES(solarSystemID)),regionID=evedata.regionIDBySolarSystem(VALUES(solarSystemID));`,
		b.StructureID, b.Structure.SolarSystemId, b.Structure.Name, b.Structure.Position.X, b.Structure.Position.Y, b.Structure.Position.Z)
	if err != nil {
		return err
	}

	// Insert into our table for tracking.
	err = s.doSQL(`INSERT INTO evedata.structures
		(stationID, solarSystemID, stationName, x, y, z, updated)
		VALUES(?,?,?,?,?,?, UTC_TIMESTAMP())
		ON DUPLICATE KEY UPDATE stationName=VALUES(stationName),solarSystemID=VALUES(solarSystemID),
		x=VALUES(x),y=VALUES(y),z=VALUES(z);`,
		b.StructureID, b.Structure.SolarSystemId, b.Structure.Name, b.Structure.Position.X, b.Structure.Position.Y, b.Structure.Position.Z)

	return err
}
