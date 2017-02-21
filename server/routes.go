package evedata

import (
	"log"
	"mime"
	"net/http"

	"github.com/antihax/evedata/appContext"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type appFunc func(*appContext.AppContext, http.ResponseWriter, *http.Request) (int, error)
type appHandler struct {
	*appContext.AppContext
	h appFunc
}

type appAuthFunc func(*appContext.AppContext, http.ResponseWriter, *http.Request, *sessions.Session) (int, error)
type appAuthHandler struct {
	*appContext.AppContext
	h appAuthFunc
}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc appFunc
}

type authRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc appAuthFunc
}

var routes []route
var authRoutes []authRoute

var notFoundHandler *route

func init() {
	// report correct
	mime.AddExtensionType(".svg", "image/svg+xml")
}

// AddRoute adds a non-authenticated web handler to the route list
// this should be called by func init() within the views package
func AddRoute(name string, method string, pattern string, handlerFunc appFunc) {
	routes = append(routes, route{name, method, pattern, handlerFunc})
}

// AddAuthRoute adds an authenticated web handler to the route list
// this should be called by func init() within the views package
func AddAuthRoute(name string, method string, pattern string, handlerFunc appAuthFunc) {
	authRoutes = append(authRoutes, authRoute{name, method, pattern, handlerFunc})
}

// AddNotFoundHandler provides a 404 handler
func AddNotFoundHandler(handlerFunc appFunc) {
	notFoundHandler = &route{"404", "GET", "", handlerFunc}
}

func (a appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := a.h(a.AppContext, w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		http.Error(w, err.Error(), status)
		return
	}
}

func (a appAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s, err := a.AppContext.Store.Get(r, "session")
	if err != nil {
		log.Printf("Session Store  %v", err)
	}
	status, err := a.h(a.AppContext, w, r, s)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		http.Error(w, err.Error(), status)
		return
	}
}

// NewRouter sets up the routes that were added.
func NewRouter(ctx *appContext.AppContext) *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	// Add public routes
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(appHandler{ctx, route.HandlerFunc})
	}

	// Add authenticted routes
	for _, route := range authRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(appAuthHandler{ctx, route.HandlerFunc})
	}

	// prometheus handler
	router.Methods("GET").Path("/metrics").Handler(promhttp.Handler())

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

	if notFoundHandler != nil {
		router.NotFoundHandler = appHandler{ctx, notFoundHandler.HandlerFunc}
	}

	return router
}
