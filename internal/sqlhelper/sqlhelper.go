package sqlhelper

import (
	"time"

	"log"
	"os"

	"github.com/antihax/evedata/models"
	"github.com/jmoiron/sqlx"
)

func NewTestDatabase() *sqlx.DB {
	database, err := models.SetupDatabase("mysql", "root@tcp(127.0.0.1:3306)/eve?allowOldPasswords=1&parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	return database
}

func NewDatabase() *sqlx.DB {
	database, err := models.SetupDatabase("mysql", os.Getenv("SQLAUTH")+"@tcp(sql)/eve?allowOldPasswords=1&parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	return database
}

func setupDatabase(driver string, spec string) (*sqlx.DB, error) {
	var err error

	// Build Connection Pool
	database, err := sqlx.Connect(driver, spec)
	if err != nil {
		return nil, err
	}

	// Check we can connect
	if err = database.Ping(); err != nil {
		return nil, err
	}

	// Put some finite limits to prevent opening too many connections
	database.SetConnMaxLifetime(time.Minute * 2)
	database.SetMaxIdleConns(100)

	return database, nil
}
