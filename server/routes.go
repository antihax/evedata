package evedata

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var routes = Routes{
	Route{
		"marketRegions",
		"GET",
		"/J/marketRegions/",
		MarketRegions,
	},
	Route{
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
