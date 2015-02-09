package evedata

import (
	"evedata/config"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

// appContext provides access to handles throughout the app.
type AppContext struct {
	conf  *config.Config
	db    *sqlx.DB
	store *sessions.Store
	/*
		templates map[string]*template.Template
		decoder   *schema.Decoder*/
}

func GoServer() {

	var err error

	// Make a new app context.
	ctx := &AppContext{}

	// Read configuation.
	ctx.conf, err = config.ReadConfig()
	if err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	// Connect to the database
	ctx.db, err = sqlx.Connect(ctx.conf.Database.Driver, ctx.conf.Database.Spec)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	// Check the connection is successful.
	err = ctx.db.Ping()
	if err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}

	ctx.Store = sessions.NewCookieStore(ctx.conf.Store.Key)

	// Allocate the router
	rtr := NewRouter(ctx)

	log.Println("Listening port 3000...")
	http.ListenAndServe(":3000", rtr)
}
