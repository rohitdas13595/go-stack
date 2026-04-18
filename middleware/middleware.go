package middleware

import (
	"compress/gzip"
	"context"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Logger writes structured request logs via slog.
func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrap := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(wrap, r)
			slog.Info("http",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrap.status,
				"dur_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Recover catches panics and returns 500.
func Recover() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					slog.Error("panic", "recover", rec)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// CORS adds Access-Control headers.
func CORS(origins []string, methods []string) func(http.Handler) http.Handler {
	allowOrigin := "*"
	if len(origins) > 0 {
		allowOrigin = strings.Join(origins, ",")
	}
	allowMethods := "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	if len(methods) > 0 {
		allowMethods = strings.Join(methods, ",")
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", allowMethods)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequestID ensures X-Request-ID.
func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = time.Now().UTC().Format("20060102150405.000000000")
			}
			w.Header().Set("X-Request-ID", id)
			next.ServeHTTP(w, r)
		})
	}
}

// Timeout cancels the request context after d.
func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type gzipRW struct {
	http.ResponseWriter
	w *gzip.Writer
}

func (g *gzipRW) Write(b []byte) (int, error) {
	return g.w.Write(b)
}

// Compress gzip-compresses responses when Accept-Encoding allows.
func Compress() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			gr := &gzipRW{ResponseWriter: w, w: gz}
			next.ServeHTTP(gr, r)
		})
	}
}

// RateLimit is a simple in-memory token bucket per client IP.
func RateLimit(rps int) func(http.Handler) http.Handler {
	if rps <= 0 {
		rps = 100
	}
	var mu sync.Mutex
	buckets := map[string]*bucket{}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			mu.Lock()
			b, ok := buckets[ip]
			if !ok {
				b = &bucket{tokens: rps, last: time.Now()}
				buckets[ip] = b
			}
			now := time.Now()
			elapsed := now.Sub(b.last).Seconds()
			b.tokens += int(elapsed * float64(rps))
			if b.tokens > rps {
				b.tokens = rps
			}
			b.last = now
			if b.tokens <= 0 {
				mu.Unlock()
				http.Error(w, "rate limit", http.StatusTooManyRequests)
				return
			}
			b.tokens--
			mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}
}

type bucket struct {
	tokens int
	last   time.Time
}
