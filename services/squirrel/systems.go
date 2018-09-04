package squirrel

import (
	"context"
	"log"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/antihax/goesi/esi"
)

// Wrap
type constellation struct {
	C        esi.GetUniverseConstellationsConstellationIdOk
	RegionID int32
}

type system struct {
	S               esi.GetUniverseSystemsSystemIdOk
	RegionID        int32
	ConstellationID int32
}

var regionsChannel = make(chan esi.GetUniverseRegionsRegionIdOk, 10000)
var constellationsChannel = make(chan constellation, 10000)
var systemsChannel = make(chan system, 10000)

func (s *Squirrel) getSystem(regionID, constellationID, systemID int32) error {
	sys, _, err := s.esi.ESI.UniverseApi.GetUniverseSystemsSystemId(context.Background(), systemID, nil)
	if err != nil {
		return err
	}
	systemsChannel <- system{S: sys, RegionID: regionID, ConstellationID: constellationID}

	return nil
}

func (s *Squirrel) getConstellation(regionID, constellationID int32) error {
	cons, _, err := s.esi.ESI.UniverseApi.GetUniverseConstellationsConstellationId(context.Background(), constellationID, nil)
	if err != nil {
		return err
	}
	constellationsChannel <- constellation{C: cons, RegionID: regionID}
	for _, c := range cons.Systems {
		err := s.getSystem(regionID, constellationID, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	registerTrigger("regions", func(s *Squirrel) error {
		regions, _, err := s.esi.ESI.UniverseApi.GetUniverseRegions(context.Background(), nil)
		if err != nil {
			return err
		}
		wg := sync.WaitGroup{}
		for _, g := range regions {
			wg.Add(1)
			go func(g int32) {
				defer func() { wg.Done() }()
				group, _, err := s.esi.ESI.UniverseApi.GetUniverseRegionsRegionId(context.Background(), g, nil)
				if err != nil {
					return
				}
				regionsChannel <- group

				for _, c := range group.Constellations {
					err := s.getConstellation(g, c)
					if err != nil {
						return
					}
				}
			}(g)
		}

		// Wait for everything to finish
		wg.Wait()

		// Close the group to wrap up any final items
		close(regionsChannel)
		return nil
	})

	registerCollector("regions", func(s *Squirrel) error {
		for {
			count := 0
			sql := sq.Insert("eve.mapRegions").Columns(
				"regionID", "regionName",
			)
			for g := range regionsChannel {
				count++
				sql = sql.Values(g.RegionId, g.Name)
				if count > 80 {
					break
				}
			}
			if count == 0 {
				break
			}
			sqlq, args, err := sql.ToSql()
			if err != nil {
				log.Println(err)
				return err
			}
			err = s.doSQL(sqlq+" ON DUPLICATE KEY UPDATE regionID = regionID", args...)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		close(constellationsChannel)
		return nil
	})

	registerCollector("constellations", func(s *Squirrel) error {
		for {
			count := 0
			sql := sq.Insert("eve.mapConstellations").Columns(
				"regionID", "constellationID", "constellationName",
			)
			for g := range constellationsChannel {
				count++
				sql = sql.Values(g.RegionID, g.C.ConstellationId, g.C.Name)
				if count > 80 {
					break
				}
			}
			if count == 0 {
				break
			}
			sqlq, args, err := sql.ToSql()
			if err != nil {
				log.Println(err)
				return err
			}
			err = s.doSQL(sqlq+" ON DUPLICATE KEY UPDATE constellationID = constellationID", args...)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		close(systemsChannel)
		return nil
	})

	registerCollector("systems", func(s *Squirrel) error {
		for {
			count := 0
			sql := sq.Insert("eve.mapSolarSystems").Columns(
				"regionID", "constellationID", "solarSystemID", "solarSystemName",
				"security", "securityClass", "x", "y", "z",
			)
			for g := range systemsChannel {
				count++
				sql = sql.Values(g.RegionID, g.ConstellationID, g.S.SystemId, g.S.Name,
					g.S.SecurityStatus, g.S.SecurityClass, g.S.Position.X, g.S.Position.Y, g.S.Position.Z,
				)
				if count > 80 {
					break
				}
			}
			if count == 0 {
				break
			}
			sqlq, args, err := sql.ToSql()
			if err != nil {
				log.Println(err)
				return err
			}
			err = s.doSQL(sqlq+" ON DUPLICATE KEY UPDATE solarSystemID = solarSystemID", args...)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	})

}
