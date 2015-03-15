package evedata

import (
	"evedata/models"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var routes = Routes{
	Route{
		"marketRegions",
		"GET",
		"/J/marketRegions/",
		MarketRegions,
	}, Route{
		"marketItemLists",
		"GET",
		"/J/marketItemLists/",
		MarketItemLists,
	}, Route{
		"marketSellRegionItems",
		"GET",
		"/J/marketSellRegionItems/",
		MarketSellRegionItems,
	}, Route{
		"marketBuyRegionItems",
		"GET",
		"/J/marketBuyRegionItems/",
		MarketBuyRegionItems,
	}, Route{
		"agents",
		"GET",
		"/U/agents/",
		FindAgents,
	},
}

type appFunc func(*AppContext, http.ResponseWriter, *http.Request) (int, error)
type appHandler struct {
	*AppContext
	h appFunc
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc appFunc
}

type Routes []Route

func (a appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := a.h(a.AppContext, w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)

		loadUser(r, a.Db)

		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}

func NewRouter(ctx *AppContext) *mux.Router {

	router := mux.NewRouter().StrictSlash(false)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(appHandler{ctx, route.HandlerFunc})
	}

	return router
}

const ContextKey int = 0

func loadUser(r *http.Request, db *sqlx.DB) {
	uidC, err := r.Cookie("uid")
	if err != nil {
		return
	}

	passC, err := r.Cookie("pass")
	if err != nil {
		return
	}

	uid, err := strconv.Atoi(uidC.Value)
	if err != nil {
		return
	}

	models.SetUser(r, uid, passC.Value, db)
}
