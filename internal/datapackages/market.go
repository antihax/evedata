package datapackages

import "github.com/antihax/goesi/esi"

type MarketOrders struct {
	Orders   []esi.GetMarketsRegionIdOrders200Ok
	RegionID int32
}
