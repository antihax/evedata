package models

import (
	"fmt"
	"os"
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

func DumpDatabase(file string, db string) (err error) {
	f, err := os.Create(file)
	defer f.Close()

	f.WriteString(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;\n\n", db))

	f.WriteString(fmt.Sprintf("USE %s;\n\n", db))

	rows, err := database.Query(`SELECT table_name
			FROM information_schema.TABLES WHERE table_schema = ?;`, db)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var table, create string
		err = rows.Scan(&table)
		if err != nil {
			return err
		}
		row := database.QueryRow(fmt.Sprintf(`SHOW CREATE TABLE %s.%s;`, db, table))
		err = row.Scan(&table, &create)
		if err != nil {
			return err
		}
		f.WriteString(fmt.Sprintf("%s;\n\n", create))
	}
	return
}
