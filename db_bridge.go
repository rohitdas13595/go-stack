package gostack

import (
	"database/sql"

	"github.com/rohitdas13595/go-stack/db"
)

var dbManager *db.Manager

// SetDBManager registers the database manager for package-level DB().
func SetDBManager(m *db.Manager) {
	dbManager = m
}

// DB returns a named connection pool (default when name omitted).
func DB(names ...string) *sql.DB {
	n := ""
	if len(names) > 0 {
		n = names[0]
	}
	if dbManager == nil {
		return nil
	}
	return dbManager.DB(n)
}
