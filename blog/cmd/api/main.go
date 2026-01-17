package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"blog/internal/bootstrap"
	"blog/internal/config"
	"blog/internal/handler"
	"blog/internal/logger"
	"blog/internal/model"
	"blog/internal/mq"
	"blog/internal/repo"
	"blog/internal/router"
	"blog/internal/service"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	cfg   *config.Config
	log   *zap.Logger
	db    *gorm.DB
	sqlDB *sql.DB
	redis *redis.Client

	rmq  *bootstrap.Rabbit
	prod *mq.Producer
	cons *mq.Consumer

	httpSrv *http.Server
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log, err := logger.New(cfg.Log)
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	app := &App{cfg: cfg, log: log}
	if err := app.Start(); err != nil {
		log.Fatal("start_failed", zap.Error(err))
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	app.Shutdown()
}

func (a *App) Start() error {
	db, sqlDB, err := bootstrap.InitMySQL(a.cfg.MySQL, a.log)
	if err != nil {
		return err
	}
	a.db = db
	a.sqlDB = sqlDB

	if err := a.db.AutoMigrate(&model.User{}); err != nil {
		a.log.Warn("auto_migrate_failed", zap.Error(err))
	}

	a.redis = bootstrap.InitRedis(a.cfg.Redis)

	rmq, err := bootstrap.InitRabbitMQ(a.cfg.RabbitMQ)
	if err != nil {
		return err
	}
	a.rmq = rmq
	a.prod = mq.NewProducer(rmq.Ch, a.cfg.RabbitMQ.Exchange, a.cfg.RabbitMQ.RoutingKey)

	a.cons = mq.NewConsumer(
		rmq.Ch,
		a.cfg.RabbitMQ.Queue,
		a.cfg.RabbitMQ.ConsumerTag,
		a.cfg.RabbitMQ.Prefetch,
		a.cfg.RabbitMQ.Concurrency,
		a.log.Named("consumer"),
		mq.DefaultJSONHandler(a.log.Named("mq")),
	)
	if err := a.cons.Start(context.Background()); err != nil {
		return err
	}

	readyFn := bootstrap.ReadyFn(bootstrap.ReadyDeps{
		SQLDB:  a.sqlDB,
		Redis:  a.redis,
		RMQCfg: a.cfg.RabbitMQ,
	})
	health := &handler.HealthHandler{ReadyFn: readyFn}

	userRepo := repo.NewUserRepo(a.db)
	userSvc := service.NewUserService(userRepo, a.prod, a.log)
	userHandler := handler.NewUserHandler(userSvc, a.log)

	postRepo := repo.NewPostRepo(a.db)
	postSvc := service.NewPostService(postRepo, a.prod, a.log)
	postHandler := handler.NewPostHandler(postSvc, a.log)

	engine := router.New(router.Deps{
		Log:    a.log,
		Cfg:    a.cfg,
		Health: health,
		User:   userHandler,
		Post:   postHandler,
	})

	a.httpSrv = &http.Server{
		Addr:              a.cfg.App.HTTPAddr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		a.log.Info("http_listen", zap.String("addr", a.cfg.App.HTTPAddr))
		if err := a.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Fatal("http_failed", zap.Error(err))
		}
	}()
	return nil
}

func (a *App) Shutdown() {
	timeout := a.cfg.App.ShutdownTimeout()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	a.log.Info("shutdown_begin", zap.Duration("timeout", timeout))

	if a.httpSrv != nil {
		_ = a.httpSrv.Shutdown(ctx)
	}
	if a.cons != nil {
		a.cons.Stop()
	}

	if a.rmq != nil {
		if a.rmq.Ch != nil {
			_ = a.rmq.Ch.Close()
		}
		if a.rmq.Conn != nil {
			_ = a.rmq.Conn.Close()
		}
	}

	if a.redis != nil {
		_ = a.redis.Close()
	}
	if a.sqlDB != nil {
		_ = a.sqlDB.Close()
	}
	a.log.Info("shutdown_done")
}
