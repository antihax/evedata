package models

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

func Begin() (*sqlx.Tx, error) {
	return database.Beginx()
}

// Escape MySQL string
func Escape(value string) string {
	replace := map[string]string{"'": `\'`}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}
