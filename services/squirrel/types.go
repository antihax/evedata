package squirrel

import (
	"context"
	"log"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
)

var typesChannel = make(chan esi.GetUniverseTypesTypeIdOk, 10000)

func init() {
	registerTrigger("types", func(s *Squirrel) error {
		var page int32 = 1
		types := []int32{}
		for {
			t, _, err := s.esi.ESI.UniverseApi.GetUniverseTypes(context.Background(),
				&esi.GetUniverseTypesOpts{Page: optional.NewInt32(page)})
			if err != nil {
				return err
			} else if len(t) == 0 { // end of the pages
				break
			}
			types = append(types, t...)
			page++
		}

		wg := sync.WaitGroup{}
		for _, g := range types {
			wg.Add(1)
			go func(g int32) {
				defer func() { wg.Done() }()
				t, _, err := s.esi.ESI.UniverseApi.GetUniverseTypesTypeId(context.Background(), g, nil)
				if err != nil {
					return
				}
				typesChannel <- t
			}(g)
		}

		// Wait for everything to finish
		wg.Wait()

		// Close the group to wrap up any final items
		close(typesChannel)
		return nil
	})

	registerCollector("types", func(s *Squirrel) error {
		for {
			count := 0
			sql := sq.Insert("eve.invTypes").Columns(
				"typeID", "groupID", "typeName", "description", "mass", "volume", "capacity",
				"portionSize", "published", "marketGroupID", "iconID", "graphicID",
			)
			for g := range typesChannel {
				count++
				sql = sql.Values(
					g.TypeId, g.GroupId, g.Name, g.Description, g.Mass, g.Volume, g.Capacity,
					g.PortionSize, g.Published, g.MarketGroupId, g.IconId, g.GraphicId,
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
			err = s.doSQL(sqlq+" ON DUPLICATE KEY UPDATE typeName = VALUES(typeName), description = VALUES(description)", args...)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	})
}
