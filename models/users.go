package models

import (
	"crypto/md5"
	"database/sql"
	"net/http"

	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"

	"fmt"
	"time"
)

// Users is a structure containing authenticated user data.
type Users struct {
	UID       int64          `db:"uid"`
	Password  string         `db:"password"`
	UserName  string         `db:"userName"`
	Created   time.Time      `db:"created"`
	LastSeen  time.Time      `db:"lastseen"`
	EMail     string         `db:"email"`
	Name      sql.NullString `db:"name"`
	EMailMe   int64          `db:"emailme"`
	EveMailMe int64          `db:"evemailme"`
}

// GetUser takes a username and password hash and returns a User struct
func GetUser(r *http.Request) *Users {
	U := context.Get(r, "user")
	if U == nil {
		return nil
	}
	return U.(*Users)
}

// GetUser takes a username and password hash and returns a User struct
func SetUser(r *http.Request, user int, pass string, db *sqlx.DB) {
	U := Users{}
	var err error
	err = db.Get(&U, "SELECT * FROM users WHERE uid =? AND password =?", user, pass)
	if err != nil {
		context.Set(r, "user", nil)
	} else {
		context.Set(r, "user", &U)
	}

	return
}

// AuthenticateUser takes a username and password and returns a User struct
func AuthenticateUser(user string, pass string, db *sqlx.DB) *Users {
	U := Users{}
	var err error
	passB := []byte(pass)
	passMD5 := md5.Sum(passB)
	passHash := fmt.Sprintf("%x", passMD5)
	err = db.Get(&U, "SELECT * FROM users WHERE userName =? AND password =?", user, passHash)
	if err != nil {
		return nil
	}

	return &U
}
