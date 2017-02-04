package models

import (
	"testing"
	"time"
)

func TestRetryTransaction(t *testing.T) {
	tx, err := Begin()
	if err != nil {
		t.Error(err)
		return
	}
	_, err = tx.Exec(`
		INSERT INTO evedata.states (nextCheck, value, state)VALUES(?,?,?) ON DUPLICATE KEY UPDATE nextCheck=VALUES(nextCheck), value=VALUES(value)
	`, time.Now().UTC(), 1, "testTransactionState")
	if err != nil {
		t.Error(err)
		return
	}

	err = RetryTransaction(tx)
	if err != nil {
		t.Error(err)
		return
	}
}
