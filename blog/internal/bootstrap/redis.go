package bootstrap

import (
	"context"
	"errors"
	"time"

	"blog/internal/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})
}

func PingRedis(ctx context.Context, rdb *redis.Client) error {
	if rdb == nil {
		return errors.New("redis not initialized")
	}
	return rdb.Ping(ctx).Err()
}
