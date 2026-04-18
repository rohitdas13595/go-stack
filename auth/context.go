package auth

import "context"

type ctxKey int

const userKey ctxKey = iota

// WithUser attaches a user value to the request context.
func WithUser(ctx context.Context, user any) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// UserFromContext returns the user from context.
func UserFromContext(ctx context.Context) any {
	if ctx == nil {
		return nil
	}
	return ctx.Value(userKey)
}
