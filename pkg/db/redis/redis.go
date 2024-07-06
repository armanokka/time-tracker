package redis

import (
	"context"
	"github.com/armanokka/test_task_Effective_mobile/config"
	"github.com/redis/go-redis/v9"
	"time"
)

// Returns new redis client
func NewRedisClient(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		MinIdleConns: cfg.Redis.MinIdleConns,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout) * time.Second,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
