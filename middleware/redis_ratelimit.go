package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisRateLimit limits requests per second per IP using Redis INCR + EXPIRE.
func RedisRateLimit(rdb *redis.Client, perSecond int) func(http.Handler) http.Handler {
	if perSecond <= 0 {
		perSecond = 10
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rdb == nil {
				next.ServeHTTP(w, r)
				return
			}
			key := "gostack:rl:" + r.RemoteAddr
			ctx := context.Background()
			n, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			if n == 1 {
				_ = rdb.Expire(ctx, key, time.Second).Err()
			}
			if int(n) > perSecond {
				http.Error(w, "rate limit", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
