package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	database      *sqlx.DB
	SQLTimeFormat = "2006-01-02 15:04:05"
)

// Set Database handle
func SetDatabase(DB *sqlx.DB) {
	database = DB
}

func SetupDatabase(driver string, spec string) (*sqlx.DB, error) {
	var err error

	// Build Connection Pool
	if database, err = sqlx.Connect(driver, spec); err != nil {
		return nil, err
	}

	// Check we can connect
	if err = database.Ping(); err != nil {
		return nil, err
	}

	// Put some finite limits to prevent opening too many connections
	database.SetConnMaxLifetime(time.Minute * 20)
	database.SetMaxIdleConns(10)

	SetDatabase(database)
	return database, nil
}
