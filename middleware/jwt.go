package middleware

import (
	"net/http"
	"strings"

	"github.com/rohitdas13595/go-stack/auth"
)

// JWT validates Authorization: Bearer and attaches claims UserID as string user placeholder.
func JWT(secret []byte, loadUser func(r *http.Request, userID string) (any, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
				next.ServeHTTP(w, r)
				return
			}
			raw := strings.TrimSpace(h[7:])
			c, err := auth.Parse(secret, raw)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			var u any = c.UserID
			if loadUser != nil {
				if loaded, err := loadUser(r, c.UserID); err == nil && loaded != nil {
					u = loaded
				}
			}
			next.ServeHTTP(w, r.WithContext(auth.WithUser(r.Context(), u)))
		})
	}
}
