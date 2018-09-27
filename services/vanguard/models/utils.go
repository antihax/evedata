package models

import (
	"fmt"
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

func FormatValue(v float64) string {
	switch {
	case v > 1000000000000:
		return fmt.Sprintf("%.1ft", v/1000000000000)
	case v > 1000000000:
		return fmt.Sprintf("%.1fb", v/1000000000)
	case v > 1000000:
		return fmt.Sprintf("%.1fm", v/1000000)
	case v > 1000:
		return fmt.Sprintf("%.1fk", v/1000)
	default:
		return fmt.Sprintf("%.1f", v)
	}
}
