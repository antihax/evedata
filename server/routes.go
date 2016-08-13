package evedata

import (
	"evedata/appContext"
	"log"
	"mime"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var routes Routes

func init() {
	// report correct
	mime.AddExtensionType(".svg", "image/svg+xml")
}

func AddRoute(r Route) {
	routes = append(routes, r)
}

type appFunc func(*appContext.AppContext, http.ResponseWriter, *http.Request, *sessions.Session) (int, error)
type appHandler struct {
	*appContext.AppContext
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

	s, _ := a.AppContext.Store.Get(r, "session")

	status, err := a.h(a.AppContext, w, r, s)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		http.Error(w, err.Error(), status)
	}
}

func NewRouter(ctx *appContext.AppContext) *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(appHandler{ctx, route.HandlerFunc})
	}

	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/",
		http.FileServer(http.Dir("static/css"))))

	router.PathPrefix("/i/").Handler(http.StripPrefix("/i/",
		http.FileServer(http.Dir("static/i"))))

	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/",
		http.FileServer(http.Dir("static/images"))))

	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/",
		http.FileServer(http.Dir("static/js"))))
	router.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/",
		http.FileServer(http.Dir("static/fonts"))))
	return router
}

const ContextKey int = 0
