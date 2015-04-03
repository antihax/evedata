package evedata

import (
	"evedata/config"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

// appContext provides access to handles throughout the app.
type AppContext struct {
	Conf  *config.Config
	Db    *sqlx.DB
	Store *sessions.Store
	/*
		templates map[string]*template.Template
		decoder   *schema.Decoder*/
}

func GoServer() {

	var err error

	// Make a new app context.8
	ctx := &AppContext{}

	// Read configuation.
	ctx.Conf, err = config.ReadConfig()
	if err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	// Connect to the database
	ctx.Db, err = sqlx.Connect(ctx.Conf.Database.Driver, ctx.Conf.Database.Spec)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	// Check the connection is successful.
	err = ctx.Db.Ping()
	if err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}

	// Allocate the routes
	rtr := NewRouter(ctx)

	log.Println("Listening port 3000...")
	http.ListenAndServe(":3000", context.ClearHandler(rtr))
}
