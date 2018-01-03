package models

import "time"

// SetServiceState sets state information: nextCheck, value (page, etc).
func SetServiceState(state string, cacheUntil time.Time, page int32) error {
	if _, err := database.Exec(`
		INSERT INTO evedata.states (nextCheck, value, state)VALUES(?,?,?) ON DUPLICATE KEY UPDATE nextCheck=VALUES(nextCheck), value=VALUES(value)
	`, cacheUntil.UTC(), page, state); err != nil {
		return err
	}
	return nil
}

// SetServiceStateByDays sets state information: nextCheck, value (page, etc).
func SetServiceStateByDays(state string, daysToCache int32, page int32) error {

	if _, err := database.Exec(`
		INSERT INTO evedata.states (nextCheck, value, state)VALUES(?,?,?) ON DUPLICATE KEY UPDATE nextCheck=VALUES(nextCheck), value=VALUES(value)
	`, daysToCache, page, state); err != nil {
		return err
	}
	return nil
}

// GetServiceState Get service state (page number, cache expired, etc).
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetServiceState(service string) (time.Time, int32, error) {
	type ServiceState struct {
		Value     int32     `db:"value"`
		NextCheck time.Time `db:"nextCheck"`
	}
	r := ServiceState{}
	if err := database.Get(&r, `
		SELECT value, nextCheck
			FROM evedata.states 
			WHERE state = ?
			LIMIT 1;
		`, service); err != nil {
		return time.Now().UTC(), 0, err
	}
	return r.NextCheck, r.Value, nil
}
