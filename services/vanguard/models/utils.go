package models

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func Begin() (*sqlx.Tx, error) {
	return database.Beginx()
}

// Retry on deadlocks
func RetryTransaction(tx *sqlx.Tx) error {
	for {
		err := tx.Commit()
		if err != nil {
			if strings.Contains(err.Error(), "1213") == false {
				return err
			} else {
				time.Sleep(500 * time.Millisecond)
				continue
			}
		} else {
			return err
		}
	}
}

// Escape MySQL string
func Escape(value string) string {
	replace := map[string]string{"'": `\'`}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}
