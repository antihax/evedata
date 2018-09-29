package datapackages

import "github.com/antihax/goesi/esi"

type MarketOrders struct {
	Orders   []esi.GetMarketsRegionIdOrders200Ok
	RegionID int32
}

type MarketHistory struct {
	History  []esi.GetMarketsRegionIdHistory200Ok
	RegionID int32
	TypeID   int32
}

type Structure struct {
	Structure   esi.GetUniverseStructuresStructureIdOk
	StructureID int64
}

type StructureOrders struct {
	Orders      []esi.GetMarketsStructuresStructureId200Ok
	StructureID int64
}

type CharacterStructure struct {
	Structure   esi.GetUniverseStructuresStructureIdOk
	StructureID int64
	CharacterID int32
}

func WintoUnixTimestamp(t int64) int64 {
	return int64(t-116444736000000000) / 1e7
}
