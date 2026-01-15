package bootstrap

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"blog/internal/config"
	"github.com/redis/go-redis/v9"
)

type ReadyDeps struct {
	SQLDB  *sql.DB
	Redis  *redis.Client
	RMQCfg config.RabbitMQConfig
}

func ReadyFn(deps ReadyDeps) func() error {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if deps.SQLDB == nil {
			return errors.New("mysql not initialized")
		}
		if err := deps.SQLDB.PingContext(ctx); err != nil {
			return err
		}
		if err := PingRedis(ctx, deps.Redis); err != nil {
			return err
		}
		if err := PingRabbitMQ(deps.RMQCfg); err != nil {
			return err
		}
		return nil
	}
}
