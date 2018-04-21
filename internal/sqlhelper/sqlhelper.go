package sqlhelper

import (
	"time"

	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func Hash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	return string(hash), err
}

func CompareHash(pwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}

func NewTestDatabase() *sqlx.DB {
	database, err := setupDatabase("mysql", "root@tcp(127.0.0.1:3306)/eve?allowOldPasswords=1&parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	return database
}

func NewDatabase() *sqlx.DB {
	database, err := setupDatabase("mysql", os.Getenv("SQLAUTH")+"@tcp(sql.evedata:3306)/eve?allowOldPasswords=1&parseTime=true")
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
