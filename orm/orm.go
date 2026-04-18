package orm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Model base with ID and timestamps (embed in user models).
type Model struct {
	ID        int64        `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

// Tabler optionally provides table name.
type Tabler interface {
	TableName() string
}

func tableName[T any]() string {
	var zero T
	zt := reflect.TypeOf(zero)
	if zt.Kind() == reflect.Ptr {
		zt = zt.Elem()
	}
	var inst any = reflect.New(zt).Interface()
	if t, ok := inst.(Tabler); ok {
		return t.TableName()
	}
	return snakePlural(zt.Name())
}

func snakePlural(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		if r >= 'A' && r <= 'Z' {
			b.WriteRune(r - 'A' + 'a')
		} else {
			b.WriteRune(r)
		}
	}
	name := b.String()
	if name != "" && !strings.HasSuffix(name, "s") {
		name += "s"
	}
	return name
}

func collectFields(rt reflect.Type, rv reflect.Value) ([]string, []reflect.Value) {
	var cols []string
	var vals []reflect.Value
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		fv := rv.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			c, v := collectFields(f.Type, fv)
			cols = append(cols, c...)
			vals = append(vals, v...)
			continue
		}
		tag := f.Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}
		cols = append(cols, tag)
		vals = append(vals, fv.Addr())
	}
	return cols, vals
}

func addrInterfaces(vals []reflect.Value) []any {
	out := make([]any, len(vals))
	for i := range vals {
		out[i] = vals[i].Interface()
	}
	return out
}

// Find loads by primary key (SQLite/MySQL style `?` placeholder).
func Find[T any](ctx context.Context, db *sql.DB, id any) (*T, error) {
	if db == nil {
		return nil, fmt.Errorf("orm: nil db")
	}
	var zero T
	rt := reflect.TypeOf(zero)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	rv := reflect.New(rt).Elem()
	cols, ptrs := collectFields(rt, rv)
	if len(cols) == 0 {
		return nil, fmt.Errorf("orm: no db-tagged fields")
	}
	tn := tableName[T]()
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE id = ? LIMIT 1`, strings.Join(cols, ", "), tn)
	row := db.QueryRowContext(ctx, q, id)
	if err := row.Scan(addrInterfaces(ptrs)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	out := rv.Addr().Interface().(*T)
	return out, nil
}

// QueryBuilder is a minimal typed query.
type QueryBuilder[T any] struct {
	db     *sql.DB
	table  string
	wheres []string
	args   []any
	order  string
	limit  int
	offset int
}

// Query starts a query for T.
func Query[T any](db *sql.DB) *QueryBuilder[T] {
	if db == nil {
		return &QueryBuilder[T]{db: nil, table: tableName[T]()}
	}
	return &QueryBuilder[T]{db: db, table: tableName[T]()}
}

// Where ANDs a predicate (use SQL fragments with ? placeholders).
func (q *QueryBuilder[T]) Where(expr string, args ...any) *QueryBuilder[T] {
	q.wheres = append(q.wheres, expr)
	q.args = append(q.args, args...)
	return q
}

// OrderBy sets ORDER BY clause (raw SQL).
func (q *QueryBuilder[T]) OrderBy(o string) *QueryBuilder[T] {
	q.order = o
	return q
}

// Limit sets LIMIT.
func (q *QueryBuilder[T]) Limit(n int) *QueryBuilder[T] {
	q.limit = n
	return q
}

// Offset sets OFFSET.
func (q *QueryBuilder[T]) Offset(n int) *QueryBuilder[T] {
	q.offset = n
	return q
}

func (q *QueryBuilder[T]) selectSQL() (string, []any, error) {
	if q.db == nil {
		return "", nil, fmt.Errorf("orm: nil db")
	}
	var zero T
	rt := reflect.TypeOf(zero)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	rv := reflect.New(rt).Elem()
	cols, _ := collectFields(rt, rv)
	if len(cols) == 0 {
		return "", nil, fmt.Errorf("orm: no db-tagged fields")
	}
	where := strings.Join(q.wheres, " AND ")
	if where == "" {
		where = "1=1"
	}
	sqlStr := fmt.Sprintf(`SELECT %s FROM %s WHERE %s`, strings.Join(cols, ", "), q.table, where)
	if q.order != "" {
		sqlStr += " ORDER BY " + q.order
	}
	if q.limit > 0 {
		sqlStr += fmt.Sprintf(" LIMIT %d", q.limit)
	}
	if q.offset > 0 {
		sqlStr += fmt.Sprintf(" OFFSET %d", q.offset)
	}
	return sqlStr, q.args, nil
}

// First returns first row or nil.
func (q *QueryBuilder[T]) First(ctx context.Context) (*T, error) {
	orig := q.limit
	q.limit = 1
	sqlStr, args, err := q.selectSQL()
	q.limit = orig
	if err != nil {
		return nil, err
	}
	row := q.db.QueryRowContext(ctx, sqlStr, args...)
	var zero T
	rt := reflect.TypeOf(zero)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	rv := reflect.New(rt).Elem()
	_, ptrs := collectFields(rt, rv)
	if err := row.Scan(addrInterfaces(ptrs)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	out := rv.Addr().Interface().(*T)
	return out, nil
}

// All returns matching rows.
func (q *QueryBuilder[T]) All(ctx context.Context) ([]T, error) {
	sqlStr, args, err := q.selectSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zero T
	rt := reflect.TypeOf(zero)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	var out []T
	for rows.Next() {
		rv := reflect.New(rt).Elem()
		_, ptrs := collectFields(rt, rv)
		if err := rows.Scan(addrInterfaces(ptrs)...); err != nil {
			return nil, err
		}
		out = append(out, rv.Addr().Interface().(T))
	}
	return out, rows.Err()
}

// Count returns row count for current WHERE clause.
func (q *QueryBuilder[T]) Count(ctx context.Context) (int64, error) {
	if q.db == nil {
		return 0, fmt.Errorf("orm: nil db")
	}
	where := strings.Join(q.wheres, " AND ")
	if where == "" {
		where = "1=1"
	}
	sqlStr := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, q.table, where)
	var n int64
	err := q.db.QueryRowContext(ctx, sqlStr, q.args...).Scan(&n)
	return n, err
}
