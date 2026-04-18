package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// CSRF sets a cookie and validates X-CSRF-Token on mutating requests (minimal).
func CSRF() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("gostack_csrf")
			token := ""
			if err != nil || c.Value == "" {
				b := make([]byte, 16)
				_, _ = rand.Read(b)
				token = hex.EncodeToString(b)
				http.SetCookie(w, &http.Cookie{Name: "gostack_csrf", Value: token, Path: "/", HttpOnly: false})
			} else {
				token = c.Value
			}
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				if r.Header.Get("X-CSRF-Token") != token && r.FormValue("_csrf") != token {
					http.Error(w, "csrf", http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
