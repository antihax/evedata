package appContext

import (
	"database/sql"
	"evedata/config"
	"evedata/eveapi"
	"net/http"

	"github.com/bradleypeabody/gorilla-sessions-memcache"
	"github.com/jmoiron/sqlx"
)

// AppContext provides access to handles throughout the app.
type AppContext struct {
	Conf  *config.Config
	Db    *sqlx.DB
	Store *gsm.MemcacheStore
	EVE   *eveapi.AnonymousClient

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
