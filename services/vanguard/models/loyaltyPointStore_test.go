package models

import "testing"

func TestAddLPOffer(t *testing.T) {
	err := AddLPOffer(1, 1, 1, 1, 1, 1, 1)
	if err != nil {
		t.Error(err)
		return
	}
}
func TestAddLPOfferRequirements(t *testing.T) {
	err := AddLPOfferRequirements(1, 1, 1)
	if err != nil {
		t.Error(err)
		return
	}
}
func TestGetISKPerLP(t *testing.T) {
	_, err := GetISKPerLP("Some Corp")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetISKPerLPCorporations(t *testing.T) {
	_, err := GetISKPerLPCorporations()
	if err != nil {
		t.Error(err)
		return
	}
}
