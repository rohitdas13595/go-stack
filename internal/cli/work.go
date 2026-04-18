package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rohitdas13595/go-stack/jobs"
	"github.com/spf13/cobra"
)

var workQueues string

var workCmd = &cobra.Command{
	Use:   "work",
	Short: "Run Redis job worker (dequeue default queue)",
	Run: func(cmd *cobra.Command, args []string) {
		addr := os.Getenv("REDIS_URL")
		if addr == "" {
			addr = "localhost:6379"
		}
		rdb := redis.NewClient(&redis.Options{Addr: addr})
		q := &jobs.RedisQueue{Client: rdb}
		queue := "default"
		if workQueues != "" {
			queue = workQueues
		}
		ctx := context.Background()
		for {
			j, err := q.Dequeue(ctx, queue, 5*time.Second)
			if err == redis.Nil {
				continue
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			fmt.Println("job", j.Name, string(j.Payload))
		}
	},
}

func init() {
	workCmd.Flags().StringVar(&workQueues, "queues", "default", "queue name")
}
