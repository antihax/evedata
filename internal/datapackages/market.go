package datapackages

import "github.com/antihax/goesi/esi"

type MarketOrders struct {
	Orders   []esi.GetMarketsRegionIdOrders200Ok
	RegionID int32
}

type Structure struct {
	Structure   esi.GetUniverseStructuresStructureIdOk
	StructureID int64
}
type StructureOrders struct {
	Orders      []esi.GetMarketsStructuresStructureId200Ok
	StructureID int64
}
