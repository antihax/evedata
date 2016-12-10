package models

import "time"

// SetServiceState sets state information: nextCheck, value (page, etc).
func SetServiceState(state string, cacheUntil time.Time, page int32) error {
	if _, err := database.Exec(`
		UPDATE states SET nextCheck = ? WHERE state = ? LIMIT 1
	`, cacheUntil, page); err != nil {
		return err
	}
	return nil
}

// SetServiceStateByDays sets state information: nextCheck, value (page, etc).
func SetServiceStateByDays(state string, daysToCache int32, page int32) error {

	if _, err := database.Exec(`
		UPDATE states SET nextCheck = DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY) WHERE state = ? LIMIT 1
	`, daysToCache, page); err != nil {
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
			FROM states 
			WHERE state = ?
			LIMIT 1;
		`, service); err != nil {
		return time.Now(), 0, err
	}
	return r.NextCheck, r.Value, nil
}
