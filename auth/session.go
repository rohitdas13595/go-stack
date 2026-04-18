package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

// SessionRow maps to _gostack_sessions (PRD §8).
type SessionRow struct {
	ID        string
	UserID    sql.NullInt64
	Payload   []byte
	ExpiresAt time.Time
}

// SessionStore persists server-side sessions in SQL.
type SessionStore struct {
	DB *sql.DB
}

// Get loads session payload JSON by id.
func (s *SessionStore) Get(ctx context.Context, id string) (map[string]any, error) {
	if s.DB == nil {
		return nil, sql.ErrNoRows
	}
	var raw []byte
	err := s.DB.QueryRowContext(ctx, `SELECT payload FROM _gostack_sessions WHERE id = ? AND expires_at > ?`, id, time.Now()).Scan(&raw)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// Save replaces a session row (SQLite-friendly).
func (s *SessionStore) Save(ctx context.Context, id string, userID int64, payload map[string]any, ttl time.Duration) error {
	if s.DB == nil {
		return sql.ErrConnDone
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	exp := time.Now().Add(ttl)
	_, _ = s.DB.ExecContext(ctx, `DELETE FROM _gostack_sessions WHERE id = ?`, id)
	_, err = s.DB.ExecContext(ctx,
		`INSERT INTO _gostack_sessions (id, user_id, payload, expires_at) VALUES (?,?,?,?)`,
		id, userID, b, exp)
	return err
}
