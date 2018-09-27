package marketwatch

import (
	"sync"
	"time"

	"github.com/antihax/goesi/esi"
)

// Order wrapper to find last touch time.
// Cast structure to market
type Order struct {
	Touched time.Time
	Order   esi.GetMarketsRegionIdOrders200Ok
}

// OrderChange Details of what changed on an order
type OrderChange struct {
	OrderID      int64     `json:"order_id"`
	LocationId   int64     `json:"location_id"`
	TypeID       int32     `json:"type_id"`
	VolumeChange int32     `json:"volume_change,omitempty"`
	VolumeRemain int32     `json:"volume_remain,omitempty"`
	Price        float64   `json:"price"`
	Duration     int32     `json:"duration,omitempty"`
	IsBuyOrder   bool      `json:"is_buy_order,omitempty"`
	Issued       time.Time `json:"issued,omitempty"`
	Changed      bool      `json:"-"`
	TimeChanged  time.Time `json:"time_changed"`
}

// storeData returns changes or true if the item is new
func (s *MarketWatch) storeData(locationID int64, order Order) (OrderChange, bool) {
	change := OrderChange{
		OrderID:     order.Order.OrderId,
		LocationId:  order.Order.LocationId,
		TypeID:      order.Order.TypeId,
		Issued:      order.Order.Issued,
		IsBuyOrder:  order.Order.IsBuyOrder,
		TimeChanged: time.Now().UTC(), // We know this was within 5 minutes of this time
	}
	sMap := s.getMarketStore(locationID)
	v, loaded := sMap.LoadOrStore(order.Order.OrderId, order)
	if loaded {
		cOrder := v.(Order)
		if order.Order.VolumeRemain != cOrder.Order.VolumeRemain ||
			order.Order.Price != cOrder.Order.Price ||
			order.Order.Duration != cOrder.Order.Duration {
			change.Changed = true
			change.VolumeChange = cOrder.Order.VolumeRemain - order.Order.VolumeRemain
			change.VolumeRemain = order.Order.VolumeRemain
			change.Price = order.Order.Price
			change.Duration = order.Order.Duration
		}
		sMap.Store(order.Order.OrderId, order)
		return change, false
	}
	return change, true
}

func (s *MarketWatch) expireOrders(locationID int64, t time.Time) []OrderChange {
	sMap := s.getMarketStore(locationID)
	changes := []OrderChange{}

	// Find any expired orders
	sMap.Range(
		func(k, v interface{}) bool {
			o := v.(Order)
			if t.After(o.Touched) {
				changes = append(changes, OrderChange{
					OrderID:      o.Order.OrderId,
					LocationId:   o.Order.LocationId,
					TypeID:       o.Order.TypeId,
					Issued:       o.Order.Issued,
					IsBuyOrder:   o.Order.IsBuyOrder,
					Changed:      true,
					VolumeChange: o.Order.VolumeRemain,
					VolumeRemain: 0,
					Price:        o.Order.Price,
					Duration:     o.Order.Duration,
					TimeChanged:  time.Now().UTC(), // We know this was within 5 minutes of this time
				})
			}
			return true
		})

	// Delete them out of the map
	for _, c := range changes {
		sMap.Delete(c.OrderID)
	}

	return changes
}

// getMarketStore for a location
func (s *MarketWatch) getMarketStore(locationID int64) *sync.Map {
	s.mmutex.RLock()
	defer s.mmutex.RUnlock()
	return s.market[locationID]
}

// createMarketStore for a location
func (s *MarketWatch) createMarketStore(locationID int64) {
	s.mmutex.Lock()
	defer s.mmutex.Unlock()
	s.market[locationID] = &sync.Map{}
}

// getStructureState for a location
func (s *MarketWatch) getStructureState(locationID int64) *Structure {
	s.smutex.RLock()
	defer s.smutex.RUnlock()
	return s.structures[locationID]
}

// createStructureState for a location
func (s *MarketWatch) createStructureState(locationID int64) *Structure {
	state := &Structure{}
	s.smutex.Lock()
	defer s.smutex.Unlock()
	s.structures[locationID] = state
	return state
}
