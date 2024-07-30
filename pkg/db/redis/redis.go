package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Addr         string
	DB           int
	Password     string
	MinIdleConns int
	PoolTimeout  int
	PoolSize     int
}

func NewRedisClient(ctx context.Context, cfg *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
		Password:     cfg.Password,
		DB:           cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
