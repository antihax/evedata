package marketwatch

import (
	"log"
	"sync"
	"time"

	"github.com/antihax/goesi/esi"

	"github.com/Masterminds/squirrel"
)

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (s *MarketWatch) saveContractItems(id int32, c []esi.GetContractsPublicItemsContractId200Ok) error {
	if len(c) > 0 {
		items := squirrel.Insert("evedata.contractItems").Columns(
			"contractID", "recordID", "isBPC", "itemID", "typeID",
			"isIncluded", "ME", "TE", "runs", "quantity",
		)

		for _, g := range c {
			items = items.Values(
				id, g.RecordId, boolToInt(g.IsBlueprintCopy), g.ItemId, g.TypeId,
				boolToInt(g.IsIncluded), g.MaterialEfficiency, g.TimeEfficiency, g.Runs, g.Quantity,
			)
		}

		sqlq, args, err := items.ToSql()
		if err != nil {
			log.Println(err)
			return err
		}

		err = s.doSQL(sqlq+`ON DUPLICATE KEY UPDATE itemID=itemID;`, args...)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (s *MarketWatch) saveContractBids(id int32, c []esi.GetContractsPublicBidsContractId200Ok) error {
	if len(c) > 0 {
		items := squirrel.Insert("evedata.contractBids").Columns(
			"contractID", "bidID", "dateBid", "amount",
		)

		for _, g := range c {
			items = items.Values(
				id, g.BidId, g.DateBid, g.Amount,
			)
		}

		sqlq, args, err := items.ToSql()
		if err != nil {
			log.Println(err)
			return err
		}

		err = s.doSQL(sqlq+`ON DUPLICATE KEY UPDATE bidID=bidID;`, args...)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (s *MarketWatch) saveContractAdditions(c []FullContract) {
	contract := squirrel.Insert("evedata.contracts").Columns(
		"contractID", "buyout", "collateral", "dateExpired", "dateIssued", "daysToComplete",
		"endLocationId", "forCorporation", "issuerCorporationID", "issuerID", "price",
		"reward", "locationID", "title", "type", "volume",
	)
	start := time.Now()
	count := 0
	length := len(c)

	wg := sync.WaitGroup{}
	sem := make(chan bool, 8)

	for _, g := range c {
		c := &g.Contract
		count++
		err := s.saveContractItems(c.ContractId, g.Items)
		if err != nil {
			continue
		}
		if c.Type_ == "auction" {
			err := s.saveContractBids(c.ContractId, g.Bids)
			if err != nil {
				continue
			}
		}

		contract = contract.Values(
			c.ContractId, c.Buyout, c.Collateral, c.DateExpired, c.DateIssued, c.DaysToComplete,
			c.EndLocationId, boolToInt(c.ForCorporation), c.IssuerCorporationId, c.IssuerId, c.Price,
			c.Reward, c.StartLocationId, c.Title, c.Type_, c.Volume,
		)

		// If we have 500 or reached the end, dump what we have
		if count%500 == 0 || count == length {
			sqlq, args, err := contract.ToSql()
			if err != nil {
				log.Println(err)
			} else {
				wg.Add(1)
				sem <- true
				go func(sqlq string, args []interface{}, start time.Time, count int) {
					defer func() { <-sem; wg.Done() }()
					err = s.doSQL(sqlq+`ON DUPLICATE KEY UPDATE price=VALUES(price), updated=UTC_TIMESTAMP();`, args...)
					if err != nil {
						log.Println(err)
					}
				}(sqlq, args, start, count)
			}
			// restart the query
			start = time.Now()
			contract = squirrel.Insert("evedata.contracts").Columns(
				"contractID", "buyout", "collateral", "dateExpired", "dateIssued", "daysToComplete",
				"endLocationId", "forCorporation", "issuerCorporationID", "issuerID", "price",
				"reward", "locationID", "title", "type", "volume",
			)
		}
	}
	wg.Wait()
}

func (s *MarketWatch) saveContractChanges(c []ContractChange) {
	for _, g := range c {
		if g.Type_ == "auction" {
			err := s.saveContractBids(g.ContractId, g.Bids)
			if err != nil {
				continue
			}
		}
		sqlq, args, err := squirrel.Update("evedata.contracts").
			Set("price", g.Price).
			Set("dateExpired", g.DateExpired).
			Where(squirrel.Eq{"contractID": g.ContractId}).ToSql()
		if err != nil {
			log.Println(err)
		} else {
			err = s.doSQL(sqlq, args...)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
