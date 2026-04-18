package gostack

import (
	"context"

	"github.com/rohitdas13595/go-stack/orm"
)

// Find loads a row by primary key using the default DB.
func Find[T any](ctx context.Context, id any) (*T, error) {
	return orm.Find[T](ctx, DB(), id)
}

// Query starts a typed query on the default DB.
func Query[T any]() *orm.QueryBuilder[T] {
	return orm.Query[T](DB())
}
