package squirrel

import (
	"context"
	"log"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/goesi/esi"
)

type npcCorpOffer struct {
	esi.GetLoyaltyStoresCorporationIdOffers200Ok
	CorporationID int32
}

var npcCorpChannel = make(chan npcCorpOffer, 10000)

func init() {
	registerTrigger("npcCorporations", func(s *Squirrel) error {
		s.esiSemStart()
		corporations, _, err := s.esi.ESI.CorporationApi.GetCorporationsNpccorps(context.Background(), nil)
		s.esiSemFinished()
		if err != nil {
			return err
		}

		wg := sync.WaitGroup{}
		work := []redisqueue.Work{}
		for _, corp := range corporations {
			// Queue for corp data collection
			work = append(work, redisqueue.Work{Operation: "corporation", Parameter: corp})

			// Get corporation loyalty store
			s.esiSemStart()
			wg.Add(1)
			go func(corp int32) {
				defer func() { s.esiSemFinished(); wg.Done() }()
				offers, _, err := s.esi.ESI.LoyaltyApi.GetLoyaltyStoresCorporationIdOffers(context.Background(), corp, nil)
				if err != nil {
					return
				}
				for _, o := range offers {
					npcCorpChannel <- npcCorpOffer{o, corp}
				}
			}(corp)
		}

		// Queue the corporations for hammer to find
		if err := s.QueueWork(work, redisqueue.Priority_Lowest); err != nil {
			log.Println(err)
		}

		// Wait for everything to finish
		wg.Wait()

		// Close the group to wrap up any final items
		close(npcCorpChannel)
		return nil
	})

	registerCollector("npcCorporations", func(s *Squirrel) error {
		for {
			count := 0
			offerSql := sq.Insert("evedata.lpOffers").Columns("offerID", "corporationID", "typeID", "quantity", "lpCost", "akCost", "iskCost")
			reqSql := sq.Insert("evedata.lpOfferRequirements").Columns("offerID", "typeID", "quantity")
			for g := range npcCorpChannel {
				count++
				offerSql = offerSql.Values(g.OfferId, g.CorporationID, g.TypeId, g.Quantity, g.LpCost, g.AkCost, g.IskCost)
				if count > 80 {
					break
				}
				for _, r := range g.RequiredItems {
					reqSql = reqSql.Values(g.OfferId, r.TypeId, r.Quantity)
				}
			}
			if count == 0 {
				break
			}
			sqlq, args, err := offerSql.ToSql()
			if err != nil {
				log.Println(err)
			} else {
				err = s.doSQL(sqlq+" ON DUPLICATE KEY UPDATE offerID = offerID", args...)
				if err != nil {
					log.Println(err)
				}
			}

			sqlq, args, err = reqSql.ToSql()
			if err != nil {
				log.Println(err)

			} else {
				err = s.doSQL(sqlq+" ON DUPLICATE KEY UPDATE offerID = offerID", args...)
				if err != nil {
					log.Println(err)
				}
			}
		}
		return nil
	})
}
