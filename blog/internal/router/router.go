package router

import (
	"time"

	"blog/internal/config"
	"blog/internal/handler"
	"blog/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type Deps struct {
	Log    *zap.Logger
	Cfg    *config.Config
	Health *handler.HealthHandler
	User   *handler.UserHandler
	Post   *handler.PostHandler
}

func New(d Deps) *gin.Engine {
	if d.Cfg.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery(d.Log))
	r.Use(middleware.AccessLog(d.Log))
	r.Use(middleware.JWTAuth(d.Log))
	//r.Use(middleware.Timeout(8 * time.Second))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{middleware.RequestIDHeader},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	rl := middleware.NewIPRateLimiter(rate.Limit(10), 20, 10*time.Minute)
	r.Use(rl.Middleware())

	r.GET("/healthz", d.Health.Healthz)
	r.GET("/readyz", d.Health.Readyz)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", d.User.Create)
		v1.GET("/login/:username/:password", d.User.Login)
		v1.POST("/pst/write", d.Post.AddPost)
		v1.GET("/pst/:userid/:page/:size", d.Post.QueryPostByUser)

	}

	return r
}
