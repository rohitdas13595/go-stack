package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const migrationsTable = `_gostack_migrations`

// Migration is one versioned migration.
type Migration struct {
	Version string
	Name    string
	UpSQL   string
	DownSQL string
}

// EnsureTable creates the migrations tracking table.
func EnsureTable(ctx context.Context, db *sql.DB, dialect string) error {
	var q string
	switch dialect {
	case "postgres", "pgx":
		q = `CREATE TABLE IF NOT EXISTS ` + migrationsTable + ` (
			id SERIAL PRIMARY KEY,
			version VARCHAR(255) NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			batch INT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`
	default:
		q = `CREATE TABLE IF NOT EXISTS ` + migrationsTable + ` (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version VARCHAR(255) NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			batch INTEGER NOT NULL,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`
	}
	_, err := db.ExecContext(ctx, q)
	return err
}

// LoadDir reads all .sql files from dir sorted by filename.
func LoadDir(dir string) ([]Migration, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".sql" {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)
	var out []Migration
	for _, f := range files {
		b, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			return nil, err
		}
		m, err := parseSQLFile(f, string(b))
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

var upMark = regexp.MustCompile(`(?m)^--\s*\+gostack:up\s*$`)
var downMark = regexp.MustCompile(`(?m)^--\s*\+gostack:down\s*$`)

func parseSQLFile(filename, body string) (Migration, error) {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	parts := strings.SplitN(base, "_", 2)
	version := parts[0]
	name := base
	if len(parts) == 2 {
		name = parts[1]
	}
	idx := upMark.FindStringIndex(body)
	if idx == nil {
		return Migration{}, fmt.Errorf("migrate: missing -- +gostack:up in %s", filename)
	}
	rest := body[idx[1]:]
	didx := downMark.FindStringIndex(rest)
	var upSQL, downSQL string
	if didx == nil {
		upSQL = strings.TrimSpace(rest)
	} else {
		upSQL = strings.TrimSpace(rest[:didx[0]])
		downSQL = strings.TrimSpace(rest[didx[1]:])
	}
	return Migration{Version: version, Name: name, UpSQL: upSQL, DownSQL: downSQL}, nil
}

// Status returns applied migration versions in order.
func Status(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, `SELECT version FROM `+migrationsTable+` ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vs []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return vs, rows.Err()
}

func insertApplied(ctx context.Context, tx *sql.Tx, dialect, version, name string, batch int) error {
	switch dialect {
	case "postgres", "pgx":
		_, err := tx.ExecContext(ctx, `INSERT INTO `+migrationsTable+` (version, name, batch) VALUES ($1,$2,$3)`, version, name, batch)
		return err
	default:
		_, err := tx.ExecContext(ctx, `INSERT INTO `+migrationsTable+` (version, name, batch) VALUES (?,?,?)`, version, name, batch)
		return err
	}
}

func deleteVersion(ctx context.Context, tx *sql.Tx, dialect, version string) error {
	switch dialect {
	case "postgres", "pgx":
		_, err := tx.ExecContext(ctx, `DELETE FROM `+migrationsTable+` WHERE version = $1`, version)
		return err
	default:
		_, err := tx.ExecContext(ctx, `DELETE FROM `+migrationsTable+` WHERE version = ?`, version)
		return err
	}
}

// Migrate applies all pending migrations in one batch.
func Migrate(ctx context.Context, db *sql.DB, dialect string, ms []Migration) error {
	if err := EnsureTable(ctx, db, dialect); err != nil {
		return err
	}
	if dialect == "postgres" || dialect == "pgx" {
		if _, err := db.ExecContext(ctx, `SELECT pg_advisory_lock(872014)`); err != nil {
			return err
		}
		defer db.ExecContext(context.Background(), `SELECT pg_advisory_unlock(872014)`)
	}
	applied, err := Status(ctx, db)
	if err != nil {
		return err
	}
	have := map[string]struct{}{}
	for _, v := range applied {
		have[v] = struct{}{}
	}
	var batch int
	_ = db.QueryRowContext(ctx, `SELECT COALESCE(MAX(batch),0) FROM `+migrationsTable).Scan(&batch)
	batch++

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, m := range ms {
		if _, ok := have[m.Version]; ok {
			continue
		}
		if _, err := tx.ExecContext(ctx, m.UpSQL); err != nil {
			return fmt.Errorf("migrate up %s: %w", m.Version, err)
		}
		if err := insertApplied(ctx, tx, dialect, m.Version, m.Name, batch); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Rollback rolls back the latest batch.
func Rollback(ctx context.Context, db *sql.DB, dialect string, ms []Migration) error {
	if err := EnsureTable(ctx, db, dialect); err != nil {
		return err
	}
	var batch int
	if err := db.QueryRowContext(ctx, `SELECT COALESCE(MAX(batch),0) FROM `+migrationsTable).Scan(&batch); err != nil {
		return err
	}
	if batch == 0 {
		return nil
	}
	var rows *sql.Rows
	var err error
	switch dialect {
	case "postgres", "pgx":
		rows, err = db.QueryContext(ctx, `SELECT version FROM `+migrationsTable+` WHERE batch = $1 ORDER BY id DESC`, batch)
	default:
		rows, err = db.QueryContext(ctx, `SELECT version FROM `+migrationsTable+` WHERE batch = ? ORDER BY id DESC`, batch)
	}
	if err != nil {
		return err
	}
	var versions []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			rows.Close()
			return err
		}
		versions = append(versions, v)
	}
	rows.Close()

	byVer := map[string]Migration{}
	for _, m := range ms {
		byVer[m.Version] = m
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, v := range versions {
		m, ok := byVer[v]
		if !ok {
			return fmt.Errorf("migrate: unknown version %s for rollback", v)
		}
		if strings.TrimSpace(m.DownSQL) != "" {
			if _, err := tx.ExecContext(ctx, m.DownSQL); err != nil {
				return fmt.Errorf("migrate down %s: %w", v, err)
			}
		}
		if err := deleteVersion(ctx, tx, dialect, v); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Dialect normalizes driver name.
func Dialect(driver string) string {
	switch driver {
	case "pgx", "postgres":
		return "pgx"
	case "mysql":
		return "mysql"
	default:
		return "sqlite"
	}
}
