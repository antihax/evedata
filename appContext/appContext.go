package appContext

import (
	"database/sql"
	"evedata/config"
	"evedata/eveapi"
	"net/http"

	redistore "gopkg.in/boj/redistore.v1"

	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
)

// AppContext provides access to handles throughout the app.
type AppContext struct {
	Conf  *config.Config
	Db    *sqlx.DB
	Store *redistore.RediStore
	EVE   *eveapi.AnonymousClient

	Cache redis.Pool

	SSOAuthenticator   *eveapi.SSOAuthenticator
	TokenAuthenticator *eveapi.SSOAuthenticator
	HTTPClient         *http.Client

	Bridge struct {
		HistoryUpdate *sql.Stmt
		OrderMark     *sql.Stmt
		OrderUpdate   *sql.Stmt
		KillInsert    *sql.Stmt
	}
}
