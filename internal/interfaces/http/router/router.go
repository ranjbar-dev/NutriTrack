package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/ranjbar-dev/nutritrack/configs"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

func New(db *pgxpool.Pool, rdb *redis.Client, cfg *configs.Config) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global middleware (order matters)
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())

	// Health check (public, no auth)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "nutritrack",
			"version": "1.0.0",
		})
	})

	// API v1 — protected route groups registered by phases 2+
	v1 := r.Group("/api/v1")
	_ = v1

	// 404 handler
	r.NoRoute(middleware.NotFound())

	return r
}
