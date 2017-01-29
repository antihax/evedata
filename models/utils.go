package models

import (
	"log"
	"strings"

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
				log.Printf("Assets: %v\n", err)
				return err
			} else {
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
	var (
		err  error
		rows int64
	)
	for {
		x, err := database.Exec(sql, args...)
		if err == nil {
			rows, err = x.RowsAffected()
			return rows, err
		} else if strings.Contains(err.Error(), "1213") == false {
			return rows, err
		}
	}
	return rows, err
}
