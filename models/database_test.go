package models

import (
	"os"
	"testing"

	"log"

	_ "github.com/go-sql-driver/mysql"
)

func TestMain(m *testing.M) {
	var err error

	// Connect to test database on 127.0.0.1
	// This will need ./sql/eve.sql and ./sql/evedata.sql imported.

	// [GLOBAL] database is global in models package: see database.go
	database, err = SetupDatabase("mysql", "root@tcp(127.0.0.1:3306)/eve?allowOldPasswords=1&parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}

	retCode := m.Run()
	database.Close()
	os.Exit(retCode)
}
