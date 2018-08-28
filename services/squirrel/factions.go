package squirrel

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/antihax/goesi/esi"
)

var factionsChannel = make(chan esi.GetUniverseFactions200Ok, 10000)

func init() {
	registerTrigger("factions", func(s *Squirrel) error {
		factions, _, err := s.esi.ESI.UniverseApi.GetUniverseFactions(context.Background(), nil)
		if err != nil {
			return err
		}
		for _, f := range factions {
			factionsChannel <- f
		}
		// Close the group to wrap up any final items
		close(factionsChannel)
		return nil
	})

	registerCollector("factions", func(s *Squirrel) error {
		for {
			count := 0
			sql := sq.Insert("eve.eveNames").Columns(
				"itemID", "itemName", "categoryID", "groupID", "typeID",
			)
			for g := range factionsChannel {
				count++
				sql = sql.Values(g.FactionId, g.Name, 1, 19, 30)
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
			err = s.doSQL(sqlq+" ON DUPLICATE KEY UPDATE itemID = itemID", args...)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	})
}
