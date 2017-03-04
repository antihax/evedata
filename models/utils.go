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

// Retry the exec until we get no error (deadlocks) and no results are returned
func RetryExecTillNoRows(sql string, args ...interface{}) error {
	for {
		rows, err := RetryExec(sql, args...)
		if err != nil {
			return err
		}
		if rows == 0 {
			break
		}
	}
	return nil
}

// Retry the exec until we get no error (deadlocks)
func RetryExec(sql string, args ...interface{}) (int64, error) {
	var rows int64
	for {
		res, err := database.Exec(sql, args...)
		if err == nil {
			rows, err = res.RowsAffected()
			return rows, err
		} else if strings.Contains(err.Error(), "1213") == false {
			return rows, err
		}
	}
}
