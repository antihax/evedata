// Package killmailstats processes killmail report data
package killmailstats

import (
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/jmoiron/sqlx"
)

// KillmailStats processes killmail report data
type KillmailStats struct {
	db *sqlx.DB
}

// NewKillmailStats Service.
func NewKillmailStats(db *sqlx.DB) *KillmailStats {
	// Setup a new squirrel
	s := &KillmailStats{
		db: db,
	}

	return s
}

// Close the service
func (s *KillmailStats) Close() {

}

// Run the service
func (s *KillmailStats) Run() {
	s.runStats()
}

func (s *KillmailStats) doSQL(stmt string, args ...interface{}) error {
	return sqlhelper.DoSQL(s.db, stmt, args...)
}
