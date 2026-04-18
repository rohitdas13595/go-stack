package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Job is a unit of work identified by name + payload.
type Job struct {
	Name    string
	Payload json.RawMessage
}

// Option configures dispatch.
type Option func(*dispatchOpts)

type dispatchOpts struct {
	delay    time.Duration
	queue    string
	maxRetry int
}

// Delay schedules after duration.
func Delay(d time.Duration) Option {
	return func(o *dispatchOpts) { o.delay = d }
}

// Queue sets queue name.
func Queue(name string) Option {
	return func(o *dispatchOpts) { o.queue = name }
}

// MaxRetries sets max attempts.
func MaxRetries(n int) Option {
	return func(o *dispatchOpts) { o.maxRetry = n }
}

// RedisQueue is a minimal Redis-backed queue (LIST).
type RedisQueue struct {
	Client *redis.Client
	Prefix string
}

func (q *RedisQueue) key(name string) string {
	if q.Prefix == "" {
		q.Prefix = "gostack:jobs"
	}
	return q.Prefix + ":" + name
}

// Enqueue pushes JSON job to Redis list.
func (q *RedisQueue) Enqueue(ctx context.Context, queue string, j Job) error {
	if q.Client == nil {
		return fmt.Errorf("jobs: nil redis client")
	}
	b, err := json.Marshal(j)
	if err != nil {
		return err
	}
	return q.Client.LPush(ctx, q.key(queue), b).Err()
}

// Dequeue blocks up to timeout waiting for a job.
func (q *RedisQueue) Dequeue(ctx context.Context, queue string, timeout time.Duration) (Job, error) {
	res, err := q.Client.BRPop(ctx, timeout, q.key(queue)).Result()
	if err != nil {
		return Job{}, err
	}
	if len(res) < 2 {
		return Job{}, fmt.Errorf("jobs: unexpected redis reply")
	}
	var j Job
	if err := json.Unmarshal([]byte(res[1]), &j); err != nil {
		return Job{}, err
	}
	return j, nil
}
