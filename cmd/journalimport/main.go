package main

import (
	"log"
	"strings"

	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/goesi"
)

// Add any new refTypes into the database
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("evedata journal import: ")

	db := sqlhelper.NewDatabase()

	for k, v := range goesi.JournalRefID {
		k = strings.Title(strings.Replace(k, "_", " ", -1))
		_, err := db.Exec("INSERT INTO evedata.walletJournalRefType (refTypeID, refTypeName) VALUES (?,?) ON DUPLICATE KEY UPDATE refTypeID = refTypeID", v, k)
		if err != nil {
			log.Println(err)
		}
	}
}
