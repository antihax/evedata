package marketwatch

import (
	"sync"
	"time"

	"github.com/antihax/goesi/esi"
)

// Contract wrapper to find last touch time.
// Cast structure to market
type Contract struct {
	Touched  time.Time
	Contract FullContract
}

// FullContract adds all three esi returns together
type FullContract struct {
	Contract esi.GetContractsPublicRegionId200Ok          `json:"contract"`
	Items    []esi.GetContractsPublicItemsContractId200Ok `json:"items,omitempty"`
	Bids     []esi.GetContractsPublicBidsContractId200Ok  `json:"bids,omitempty"`
}

// ContractChange Details of what changed on an contract
// Really only price and bids can change
type ContractChange struct {
	ContractId  int32                                       `json:"contract_id"`
	LocationId  int64                                       `json:"location_id"`
	Expired     bool                                        `json:"expired,omitempty"`
	DateExpired time.Time                                   `json:"date_expired,omitempty"`
	Changed     bool                                        `json:"-"`
	Bids        []esi.GetContractsPublicBidsContractId200Ok `json:"bids,omitempty"`
	Price       float64                                     `json:"price,omitempty"`
	Type_       string                                      `json:"type,omitempty"`
	TimeChanged time.Time                                   `json:"time_changed,omitempty"`
}

// storeContract returns changes or true if the item is new
func (s *MarketWatch) storeContract(locationID int64, c Contract) (ContractChange, bool) {
	sMap := s.getContractStore(locationID)
	change := ContractChange{
		ContractId:  c.Contract.Contract.ContractId,
		LocationId:  c.Contract.Contract.StartLocationId,
		TimeChanged: time.Now().UTC(), // We know this was within 30 minutes of this time
	}

	v, loaded := sMap.LoadOrStore(c.Contract.Contract.ContractId, c)
	if loaded {
		contract := v.(Contract)
		if len(contract.Contract.Bids) != len(c.Contract.Bids) {
			change.Price = contract.Contract.Contract.Price
			change.Bids = contract.Contract.Bids
			change.Type_ = contract.Contract.Contract.Type_
			change.DateExpired = contract.Contract.Contract.DateExpired
			change.Changed = true
		}
		sMap.Store(contract.Contract.Contract.ContractId, contract)
		return change, false
	}
	return change, true
}

func (s *MarketWatch) expireContracts(locationID int64, t time.Time) []ContractChange {
	sMap := s.getContractStore(locationID)
	changes := []ContractChange{}

	// Find any expired contracts
	sMap.Range(
		func(k, v interface{}) bool {
			o := v.(Contract)
			if t.After(o.Touched) {

				expired := false
				if o.Contract.Contract.DateExpired.Before(time.Now()) {
					expired = true
				}
				changes = append(changes, ContractChange{
					ContractId:  o.Contract.Contract.ContractId,
					LocationId:  o.Contract.Contract.StartLocationId,
					Price:       o.Contract.Contract.Price,
					Bids:        o.Contract.Bids,
					Type_:       o.Contract.Contract.Type_,
					DateExpired: o.Contract.Contract.DateExpired,
					Changed:     true,
					Expired:     expired,
					TimeChanged: time.Now().UTC(), // We know this was within 30 minutes of this time
				})
			}
			return true
		})

	// Delete them out of the map
	for _, c := range changes {
		sMap.Delete(c.ContractId)
	}

	return changes
}

// getContractStore for a location
func (s *MarketWatch) getContractStore(locationID int64) *sync.Map {
	s.cmutex.RLock()
	defer s.cmutex.RUnlock()
	return s.contracts[locationID]
}

// createContractStore for a location
func (s *MarketWatch) createContractStore(locationID int64) {
	s.cmutex.Lock()
	defer s.cmutex.Unlock()
	s.contracts[locationID] = &sync.Map{}
}
