package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

////////////////////////////////////////////////////////////////////////////////

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR,default=localhost:6379"`
	Password string `env:"REDIS_PASSWORD,default="`
	DB       int    `env:"REDIS_DB,default=0"`
}

////////////////////////////////////////////////////////////////////////////////

func GetRedis(ctx context.Context, cfg RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	////////////////////////////////////////////////////////////////////////////

	if err := retry(3, 3*time.Second, func(i int) error {
		statusCmd := client.Ping(ctx)
		err := statusCmd.Err()
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to create connection pool (%s): %w", cfg.Addr, err)
	}

	return client, nil
}
