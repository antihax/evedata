package models

import "testing"

func TestGetMarketHistory(t *testing.T) {
	_, err := GetMarketHistory(1, 10000002)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetArbitrageCalculatorStations(t *testing.T) {
	_, err := GetArbitrageCalculatorStations()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetGetArbitrageCalculator(t *testing.T) {
	_, err := GetArbitrageCalculator(999, 1, 1999999999, 1, 1, "delta")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = GetArbitrageCalculator(999, 1, 1999999999, 1, 1, "percentage")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetMarketRegions(t *testing.T) {
	_, err := GetMarketRegions()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetMarketTypes(t *testing.T) {
	_, err := GetMarketTypes()
	if err != nil {
		t.Error(err)
		return
	}
}
func TestMarketRegionItems(t *testing.T) {
	_, err := MarketRegionItems(10000002, 41, highSec, true)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = MarketRegionItems(10000002, 41, highSec, false)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = MarketRegionItems(0, 41, highSec, false)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = MarketRegionItems(0, 41, highSec, true)
	if err != nil {
		t.Error(err)
		return
	}
}
