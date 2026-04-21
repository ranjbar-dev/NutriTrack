package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/ranjbar-dev/nutritrack/bootstrap"
	"github.com/ranjbar-dev/nutritrack/configs"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/handler"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

func New(db *pgxpool.Pool, rdb *redis.Client, cfg *configs.Config) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Wire all dependencies via DI container
	container := bootstrap.NewContainer(db, rdb, cfg)

	r := gin.New()

	// Global middleware chain (order matters)
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

	// API v1
	v1 := r.Group("/api/v1")

	// --- Auth routes (public) ---
	authHandler := handler.NewAuthHandler(container.AuthService)
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/login",      authHandler.Login)
		authGroup.POST("/otp/send",   authHandler.SendOTP)
		authGroup.POST("/otp/verify", authHandler.VerifyOTP)
		authGroup.POST("/refresh",    authHandler.Refresh)
	}

	// --- Protected routes (require JWT) ---
	protected := v1.Group("")
	protected.Use(middleware.RequireAuth(container.JWTService))
	{
		protected.POST("/auth/logout", authHandler.Logout)
		protected.GET("/auth/me",      authHandler.Me)
	}

	// 404 handler
	r.NoRoute(middleware.NotFound())

	return r
}
