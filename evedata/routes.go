package evedata

import (
	"context"
	"log"
	"mime"
	"net/http"
	"net/http/pprof"

	"github.com/antihax/evedata/appContext"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Context keys
type key int

const globalsKey key = 1 // Application Context (redis, config, ESI, etc)
const sessionKey key = 2 // User session data

// Structure for handling routes
type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Globals for routes so we can add them in `func init()`
var routes []route
var authRoutes []route
var notFoundHandler *route

func init() {
	// report correct
	mime.AddExtensionType(".svg", "image/svg+xml")
}

// AddRoute adds a non-authenticated web handler to the route list
// this should be called by func init() within the views package
func AddRoute(name string, method string, pattern string, handlerFunc http.HandlerFunc) {
	routes = append(routes, route{name, method, pattern, handlerFunc})
}

// AddAuthRoute adds an authenticated web handler to the route list
// this should be called by func init() within the views package
func AddAuthRoute(name string, method string, pattern string, handlerFunc http.HandlerFunc) {
	authRoutes = append(authRoutes, route{name, method, pattern, handlerFunc})
}

// AddNotFoundHandler provides a 404 handler
func AddNotFoundHandler(handlerFunc http.HandlerFunc) {
	notFoundHandler = &route{"404", "GET", "", handlerFunc}
}

// Middleware to add global AppContext to a request.Context
func contextWithGlobals(ctx context.Context, a *appContext.AppContext) context.Context {
	return context.WithValue(ctx, globalsKey, a)
}

// GlobalsFromContext returns attached AppContext from a request.Context
func GlobalsFromContext(ctx context.Context) *appContext.AppContext {
	return ctx.Value(globalsKey).(*appContext.AppContext)
}

// Middleware to add user session data to a request.Context
func contextWithSession(ctx context.Context, r *http.Request) context.Context {
	a := GlobalsFromContext(ctx)
	s, err := a.Store.Get(r, "session")
	if err != nil {
		log.Printf("%q", err)
		return ctx
	}
	return context.WithValue(ctx, sessionKey, s)
}

// SessionFromContext returns user session data from a request.Context
func SessionFromContext(ctx context.Context) *sessions.Session {
	return ctx.Value(sessionKey).(*sessions.Session)
}

// Handle authenticated requests
func authedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := contextWithGlobals(req.Context(), GetContext())
		ctx = contextWithSession(ctx, req)
		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}

// Handle normal requests
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := contextWithGlobals(req.Context(), GetContext())
		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}

// Add /debug information to the router. Make sure this is not exposed publicly.
func attachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.Handle("/debug/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/heap", pprof.Handler("heap"))
	router.Handle("/debug/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/block", pprof.Handler("block"))
	router.Handle("/debug/mutex", pprof.Handler("mutex"))
}

// Serve favicon.ico
func ServeFavIconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon/favicon.ico")
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
			Handler(middleware(route.HandlerFunc))
	}

	// Add authenticted routes
	for _, route := range authRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(authedMiddleware(route.HandlerFunc))
	}

	// Serve FavIcon
	router.Methods("GET").Path("/favicon.ico").HandlerFunc(ServeFavIconHandler)

	// Serve CSS
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/",
		http.FileServer(http.Dir("static/css"))))

	// Serve local images
	router.PathPrefix("/i/").Handler(http.StripPrefix("/i/",
		http.FileServer(http.Dir("static/i"))))

	// Serve third party images
	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/",
		http.FileServer(http.Dir("static/images"))))

	// Serve java script
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/",
		http.FileServer(http.Dir("static/js"))))

	// Server web fonts
	router.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/",
		http.FileServer(http.Dir("static/fonts"))))

	// Deal with 404s
	if notFoundHandler != nil {
		router.NotFoundHandler = middleware(notFoundHandler.HandlerFunc)
	}

	/*******************************************************************

	!! WARNING !!

	Make sure to deny access to the following handlers via reverse proxy.

	Both /metrics and /debug should be private and return 403 publicly.

	*******************************************************************/
	attachProfiler(router)

	// prometheus handler
	router.Methods("GET").Path("/metrics").Handler(promhttp.Handler())

	return router
}
