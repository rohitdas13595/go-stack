package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

// Manager holds named *sql.DB pools.
type Manager struct {
	mu sync.RWMutex
	pools map[string]*sql.DB
}

// NewManager creates an empty manager.
func NewManager() *Manager {
	return &Manager{pools: make(map[string]*sql.DB)}
}

// Register adds a named connection pool.
func (m *Manager) Register(name string, db *sql.DB) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pools[name] = db
}

// DB returns a pool by name; empty name = "default".
func (m *Manager) DB(name string) *sql.DB {
	if name == "" {
		name = "default"
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.pools[name]
}

// OpenSQLite opens a sqlite database with sensible pool settings.
func OpenSQLite(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(16)
	db.SetConnMaxLifetime(30 * time.Minute)
	return db, nil
}

// OpenPostgres opens postgres via pgx stdlib driver.
func OpenPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// OpenMySQL opens MySQL.
func OpenMySQL(dsn string) (*sql.DB, error) {
	return sql.Open("mysql", dsn)
}

// PingAll pings every registered pool.
func (m *Manager) PingAll(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for name, db := range m.pools {
		if err := db.PingContext(ctx); err != nil {
			return fmt.Errorf("ping %s: %w", name, err)
		}
	}
	return nil
}

// Stats returns pool stats for a name.
func (m *Manager) Stats(name string) sql.DBStats {
	db := m.DB(name)
	if db == nil {
		return sql.DBStats{}
	}
	return db.Stats()
}
