package datapackages

import (
	"github.com/antihax/goesi/esi"
)

type ResolveItems struct {
	Item   esi.GetDogmaDynamicItemsTypeIdItemIdOk
	ItemID int64
	TypeID int32
}
