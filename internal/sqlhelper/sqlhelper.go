package sqlhelper

import (
	"strings"
	"time"

	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // SQL Driver
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Hash a password with bcrypt
func Hash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	return string(hash), err
}

// CompareHash of a bcrypt hash
func CompareHash(pwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}

// NewTestDatabase creates a connection to a test database on 127.0.0.1:3306 with no password
func NewTestDatabase() *sqlx.DB {
	database, err := setupDatabase("mysql", "root@tcp(127.0.0.1:3306)/eve?allowOldPasswords=1&parseTime=true&tls=skip-verify")
	if err != nil {
		log.Fatalln(err)
	}
	return database
}

// NewDatabase creates a new connection for sql.storage
func NewDatabase() *sqlx.DB {
	database, err := setupDatabase("mysql", os.Getenv("SQLAUTH")+"@tcp(sql.storage.svc.cluster.local:3306)/eve?allowOldPasswords=1&parseTime=true&tls=skip-verify")
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
	database.SetMaxIdleConns(10)
	database.SetMaxOpenConns(50)

	return database, nil
}

// DoSQL executes a sql statement
func DoSQL(db *sqlx.DB, stmt string, args ...interface{}) error {
	for {
		_, err := RetryExec(db, stmt, args...)
		if err != nil {
			if !strings.Contains(err.Error(), "1213") && !strings.Contains(err.Error(), "1205") {
				return err
			}
			time.Sleep(50 * time.Millisecond)
			continue
		} else {
			return err
		}
	}
}

// RetryExecTillNoRows retries the exec until we get no error (deadlocks) and no results are returned
func RetryExecTillNoRows(db *sqlx.DB, sql string, args ...interface{}) error {
	for {
		rows, err := RetryExec(db, sql, args...)
		if err != nil {
			return err
		}
		if rows == 0 {
			break
		}
	}
	return nil
}

// RetryExec retries the exec until we get no error (deadlocks)
func RetryExec(db *sqlx.DB, sql string, args ...interface{}) (int64, error) {
	var rows int64
	for {
		res, err := db.Exec(sql, args...)
		if err == nil {
			rows, err = res.RowsAffected()
			return rows, err
		} else if !strings.Contains(err.Error(), "1213") && !strings.Contains(err.Error(), "1205") {
			return rows, err
		}
	}
}

// RetryTransaction on deadlocks
func RetryTransaction(tx *sqlx.Tx) error {
	for {
		err := tx.Commit()
		if err != nil {
			if !strings.Contains(err.Error(), "1213") && !strings.Contains(err.Error(), "1205") {
				return err
			}
			time.Sleep(50 * time.Millisecond)
			continue
		} else {
			return err
		}
	}
}

// IToB converts a boolean to integer
func IToB(b bool) int {
	if b {
		return 1
	}
	return 0
}
