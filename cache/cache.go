package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store is a minimal cache interface.
type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// Memory is an in-process cache.
type Memory struct {
	data map[string]memItem
}

type memItem struct {
	val []byte
	exp time.Time
}

// NewMemory creates memory store.
func NewMemory() *Memory {
	return &Memory{data: make(map[string]memItem)}
}

func (m *Memory) Get(ctx context.Context, key string) ([]byte, error) {
	it, ok := m.data[key]
	if !ok {
		return nil, redis.Nil
	}
	if !it.exp.IsZero() && time.Now().After(it.exp) {
		delete(m.data, key)
		return nil, redis.Nil
	}
	return it.val, nil
}

func (m *Memory) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	m.data[key] = memItem{val: val, exp: exp}
	return nil
}

func (m *Memory) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

// RedisStore uses Redis strings.
type RedisStore struct {
	Client *redis.Client
	Prefix string
}

func (r *RedisStore) k(key string) string {
	if r.Prefix == "" {
		r.Prefix = "gostack:cache:"
	}
	return r.Prefix + key
}

func (r *RedisStore) Get(ctx context.Context, key string) ([]byte, error) {
	if r.Client == nil {
		return nil, fmt.Errorf("cache: nil client")
	}
	s, err := r.Client.Get(ctx, r.k(key)).Bytes()
	return s, err
}

func (r *RedisStore) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	return r.Client.Set(ctx, r.k(key), val, ttl).Err()
}

func (r *RedisStore) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, r.k(key)).Err()
}

// GetJSON decodes JSON from store.
func GetJSON[T any](ctx context.Context, s Store, key string) (T, bool, error) {
	var zero T
	b, err := s.Get(ctx, key)
	if err == redis.Nil {
		return zero, false, nil
	}
	if err != nil {
		return zero, false, err
	}
	if err := json.Unmarshal(b, &zero); err != nil {
		return zero, false, err
	}
	return zero, true, nil
}

// SetJSON encodes value.
func SetJSON(ctx context.Context, s Store, key string, v any, ttl time.Duration) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.Set(ctx, key, b, ttl)
}
