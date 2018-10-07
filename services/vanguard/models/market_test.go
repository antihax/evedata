package models

import (
	"testing"
)

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
	_, err := database.Exec(`
			INSERT INTO evedata.market VALUES
			 	(4792612441, 10000002, 60003760, 37306, 1, 168828720.00, 1, 2, 2, '2017-03-11 03:38:58', 90, UTC_TIMESTAMP,1),
				(4827094797, 10000002, 60003760, 37306, 0, 289904000.00, 1, 1, 10, '2017-04-16 08:19:54', 90, UTC_TIMESTAMP,1),
				(4830712070, 10000002, 60003760, 37306, 1, 186000000.00, 1, 1, 1, '2017-04-16 07:11:10', 90, UTC_TIMESTAMP,1),
				(4836973592, 10000002, 60003760, 37306, 1, 186000016.00, 1, 1, 1, '2017-04-16 14:38:51', 90, UTC_TIMESTAMP,1),
				(4837213223, 10000002, 60003760, 37306, 1, 180688672.00, 1, 1, 1, '2017-04-15 00:36:14', 90, UTC_TIMESTAMP,1),
				(4837341606, 10000002, 60003760, 37306, 1, 186000016.00, 1, 1, 1, '2017-04-16 15:05:48', 90, UTC_TIMESTAMP,1),
				(4837352850, 10000002, 60003760, 37306, 1, 186000016.00, 1, 3, 3, '2017-04-16 15:06:39', 90, UTC_TIMESTAMP,1),
				(4837660888, 10000002, 60003760, 37306, 1, 180688688.00, 1, 2, 2, '2017-04-15 17:29:36', 90, UTC_TIMESTAMP,1),
				(4838092384, 10000002, 60003760, 37306, 1, 186000016.00, 1, 2, 2, '2017-04-16 13:42:42', 3, UTC_TIMESTAMP,1),
				(4838181372, 10000002, 60003760, 37306, 1, 186000016.00, 1, 4, 5, '2017-04-16 15:13:52', 3, UTC_TIMESTAMP,1),
				(4838183905, 10000002, 60003760, 37306, 1, 186000000.00, 1, 1, 1, '2017-04-16 07:46:03', 3, UTC_TIMESTAMP,1),
				(4838184927, 10000002, 60003760, 37306, 1, 186000016.00, 1, 3, 3, '2017-04-16 15:07:41', 3, UTC_TIMESTAMP,1),
				(4838186007, 10000002, 60003760, 37306, 1, 186000016.00, 1, 1, 1, '2017-04-16 14:36:24', 3, UTC_TIMESTAMP,1),
				(4838438363, 10000002, 60003760, 37306, 1, 186000016.00, 1, 1, 1, '2017-04-16 14:46:45', 3, UTC_TIMESTAMP,1),
				(4838575470, 10000002, 60003760, 37306, 1, 186000016.00, 1, 1, 1, '2017-04-16 15:07:26', 1, UTC_TIMESTAMP,1),
				(4839407265, 10000002, 60003760, 37306, 1, 186000000.00, 1, 1, 1, '2017-04-16 07:57:13', 90, UTC_TIMESTAMP,1),
				(4839690690, 10000002, 60003760, 37306, 1, 186000016.00, 1, 2, 2, '2017-04-16 13:02:41', 90, UTC_TIMESTAMP,1),
				(4840452883, 10000002, 60003760, 37306, 0, 289904992.00, 1, 12, 12, '2017-04-15 15:29:58', 90, UTC_TIMESTAMP,1),
				(4840918903, 10000002, 60003760, 37306, 0, 289904000.00, 1, 5, 12, '2017-04-15 23:29:16', 90, UTC_TIMESTAMP,1),
				(4841277767, 10000002, 60003760, 37306, 1, 186000016.00, 1, 3, 3, '2017-04-16 15:22:49', 90, UTC_TIMESTAMP,1),
				(4841279374, 10000002, 60003760, 37306, 1, 186000016.00, 1, 4, 4, '2017-04-16 14:56:17', 3, UTC_TIMESTAMP,1),
				(4841291003, 10000002, 60003760, 37306, 0, 289904000.00, 1, 2, 10, '2017-04-16 10:12:37', 90, UTC_TIMESTAMP,1),
				(4841308374, 10000002, 60003760, 37306, 1, 186000000.00, 1, 1, 1, '2017-04-16 12:05:55', 90, UTC_TIMESTAMP,1),
				(4841341678, 10000002, 60003760, 37306, 1, 186000016.00, 1, 1, 1, '2017-04-16 15:13:57', 90, UTC_TIMESTAMP,1)
				ON DUPLICATE KEY UPDATE orderID=orderID;
		`)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = database.Exec(`
			INSERT INTO evedata.market_vol VALUES( 54.0000, 10000002, 37306) ON DUPLICATE KEY UPDATE number=number;
		`)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = GetArbitrageCalculator(60003760, 1, 19999999999, 0.01, 0.01, "delta")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = GetArbitrageCalculator(60003760, 1, 19999999999, 0.01, 0.01, "percentage")
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

func TestMarketUnderValued(t *testing.T) {
	_, err := MarketUnderValued(10000002, 10000002, 10000044, 0.2)
	if err != nil {
		t.Error(err)
		return
	}
}
