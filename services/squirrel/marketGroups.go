package squirrel

import (
	"context"
	"log"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/antihax/goesi/esi"
)

var marketGroupsChannel = make(chan esi.GetMarketsGroupsMarketGroupIdOk, 10000)

func init() {
	registerTrigger("marketGroups", func(s *Squirrel) error {
		groups, _, err := s.esi.ESI.MarketApi.GetMarketsGroups(context.Background(), nil)
		if err != nil {
			return err
		}

		wg := sync.WaitGroup{}
		for _, g := range groups {
			wg.Add(1)
			go func(g int32) {
				defer func() { wg.Done() }()
				group, _, err := s.esi.ESI.MarketApi.GetMarketsGroupsMarketGroupId(context.Background(), g, nil)
				if err != nil {
					return
				}
				marketGroupsChannel <- group
			}(g)
		}

		// Wait for everything to finish
		wg.Wait()

		// Close the group to wrap up any final items
		close(marketGroupsChannel)
		return nil
	})

	registerCollector("marketGroups", func(s *Squirrel) error {
		for {
			count := 0
			sql := sq.Insert("eve.invMarketGroups").Columns("marketGroupID", "parentGroupID", "marketGroupName", "description", "hasTypes")
			for g := range marketGroupsChannel {
				count++
				sql = sql.Values(g.MarketGroupId, g.ParentGroupId, g.Name, g.Description, len(g.Types) > 0)
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
			err = s.doSQL(sqlq+" ON DUPLICATE KEY UPDATE marketGroupID = marketGroupId", args...)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	})
}
