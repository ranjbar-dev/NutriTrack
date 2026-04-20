package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/ranjbar-dev/nutritrack/backend/internal/config"
	"github.com/ranjbar-dev/nutritrack/backend/internal/handler"
	"github.com/ranjbar-dev/nutritrack/backend/internal/middleware"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository"
	"github.com/ranjbar-dev/nutritrack/backend/internal/service"
	customvalidator "github.com/ranjbar-dev/nutritrack/backend/internal/validator"
	"github.com/ranjbar-dev/nutritrack/backend/pkg/sms"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Setup zerolog
	logger := setupLogger(cfg.Environment)

	logger.Info().
		Str("environment", cfg.Environment).
		Str("port", cfg.Port).
		Msg("starting NutriTrack API")

	// Create root context with signal cancellation for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Connect to PostgreSQL with pool configuration
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse database URL")
	}
	poolConfig.MaxConns = 20
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create database pool")
	}
	defer pool.Close()

	// Verify database connection
	if err := pool.Ping(ctx); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping database")
	}
	logger.Info().Msg("database connection established")

	// Run migrations
	if err := runMigrations(cfg.DatabaseURL); err != nil {
		logger.Fatal().Err(err).Msg("failed to run migrations")
	}

	// Register custom validators (e.g. iranian_mobile)
	if err := customvalidator.RegisterCustomValidators(); err != nil {
		logger.Fatal().Err(err).Msg("failed to register custom validators")
	}

	// Initialize SMS sender — mock in development, Kavenegar in production (D-04)
	var smsSender sms.Sender
	if cfg.Environment == "development" {
		smsSender = sms.NewMockSender(logger)
	} else {
		smsSender = sms.NewKavenegarSender(cfg.SMSAPIKey, cfg.SMSTemplate)
	}

	// Initialize rate limiter: 3 requests per 10-minute window (D-27, AUTH-07, SEC-07)
	rateLimiter := middleware.NewRateLimiter(3, 10*time.Minute)

	// JWT secret as bytes
	jwtSecret := []byte(cfg.JWTSecret)

	// Initialize repositories
	userRepo := repository.NewUserRepository(pool)
	otpRepo := repository.NewOTPRepository(pool)
	tokenRepo := repository.NewTokenRepository(pool)
	foodRepo := repository.NewFoodRepository(pool)
	medRepo := repository.NewMedicationRepository(pool)
	planRepo := repository.NewDietPlanRepository(pool)
	trackingRepo := repository.NewTrackingRepository(pool)
	commRepo := repository.NewCommunicationRepository(pool)

	// Initialize services
	authService := service.NewAuthService(userRepo, otpRepo, tokenRepo, smsSender, jwtSecret, logger)
	userService := service.NewUserService(userRepo, logger)
	foodService := service.NewFoodService(foodRepo, logger)
	medService := service.NewMedicationService(medRepo, logger)
	planService := service.NewDietPlanService(planRepo, logger)
	trackingService := service.NewTrackingService(trackingRepo, cfg.UploadsDir, logger)
	commService := service.NewCommunicationService(commRepo, userRepo, cfg.UploadsDir, logger)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	adminHandler := handler.NewAdminHandler(userService)
	clientHandler := handler.NewClientHandler(userService)
	foodHandler := handler.NewFoodHandler(foodService)
	medHandler := handler.NewMedicationHandler(medService)
	planHandler := handler.NewDietPlanHandler(planService)
	trackingHandler := handler.NewTrackingHandler(trackingService, cfg.UploadsDir)
	commHandler := handler.NewCommunicationHandler(commService, cfg.UploadsDir)

	// Create Gin engine — use gin.New() not gin.Default() per D-07
	r := gin.New()

	// Global middleware chain: Recovery → SecurityHeaders → RequestID → Logger → CORS
	r.Use(middleware.Recovery())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS(cfg.FrontendURL))

	// Public routes — no auth required
	pub := r.Group("/api")
	{
		pub.GET("/health", handler.HealthCheck)
		pub.POST("/auth/login", authHandler.Login)
		pub.POST("/auth/otp/request", middleware.RateLimit(rateLimiter), authHandler.RequestOTP)
		pub.POST("/auth/otp/verify", middleware.RateLimit(rateLimiter), authHandler.VerifyOTP)
		pub.POST("/auth/refresh", authHandler.Refresh)
		pub.POST("/auth/logout", authHandler.Logout)
	}

	// Authenticated routes (any logged-in user, regardless of role)
	authed := r.Group("/api")
	authed.Use(middleware.Auth(jwtSecret))
	{
		authed.GET("/auth/me", authHandler.GetMe)
	}

	foods := r.Group("/api/foods")
	foods.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("nutritionist", "super_admin"))
	{
		foods.GET("", foodHandler.List)
		foods.POST("", foodHandler.Create)
		foods.GET("/:id", foodHandler.Get)
		foods.PUT("/:id", foodHandler.Update)
		foods.DELETE("/:id", foodHandler.Delete)
	}

	meds := r.Group("/api/medications")
	meds.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("nutritionist", "super_admin"))
	{
		meds.GET("", medHandler.List)
		meds.POST("", medHandler.Create)
		meds.GET("/:id", medHandler.Get)
		meds.PUT("/:id", medHandler.Update)
		meds.DELETE("/:id", medHandler.Delete)
	}

	// Admin routes — super_admin only
	admin := r.Group("/api/admin")
	admin.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("super_admin"))
	{
		admin.POST("/nutritionists", adminHandler.CreateNutritionist)
	}

	// Nutritionist routes
	nutri := r.Group("/api/nutritionist")
	nutri.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("nutritionist"))
	{
		nutri.POST("/clients", clientHandler.RegisterClient)
		nutri.GET("/clients", clientHandler.NutriListClients)
		nutri.GET("/clients/:clientId", clientHandler.NutriGetClientProfile)
		nutri.PATCH("/clients/:clientId/activate", clientHandler.NutriActivateClient)
		nutri.PATCH("/clients/:clientId/deactivate", clientHandler.NutriDeactivateClient)
		nutri.PATCH("/clients/:clientId/profile", clientHandler.NutriUpdateClientProfile)
		nutri.GET("/clients/:clientId/tracking/food", trackingHandler.NutriListFoodLogs)
		nutri.GET("/clients/:clientId/tracking/water", trackingHandler.NutriListWaterLogs)
		nutri.GET("/clients/:clientId/tracking/sleep", trackingHandler.NutriListSleepLogs)
		nutri.GET("/clients/:clientId/tracking/exercise", trackingHandler.NutriListExerciseLogs)
		nutri.GET("/clients/:clientId/tracking/medication", trackingHandler.NutriListMedicationLogs)
		nutri.GET("/clients/:clientId/tracking/body", trackingHandler.NutriListBodyMeasurements)
		nutri.GET("/clients/:clientId/tracking/body/weight-history", trackingHandler.NutriGetWeightHistory)
		nutri.POST("/clients/:clientId/body-measurements", trackingHandler.NutriUpsertBodyMeasurement)
		nutri.GET("/clients/:clientId/lab-results", trackingHandler.NutriListLabResults)
		nutri.GET("/clients/:clientId/lab-results/:labId/download", trackingHandler.NutriDownloadLabResult)
		nutri.GET("/food-requests", commHandler.NutriListFoodRequests)
		nutri.PATCH("/food-requests/:requestId/approve", commHandler.NutriApproveFoodRequest)
		nutri.PATCH("/food-requests/:requestId/reject", commHandler.NutriRejectFoodRequest)
	}

	// Client plan list: GET /api/clients/:clientId/plans (nutritionist or super_admin)
	nutriClientRoutes := r.Group("/api/clients")
	nutriClientRoutes.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("nutritionist", "super_admin"))
	{
		nutriClientRoutes.GET("/:clientId/plans", planHandler.ListClientPlans)
	}

	// Shared read route for full plan aggregate (nutritionist/super_admin/client)
	dietPlanRead := r.Group("/api/diet-plans")
	dietPlanRead.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("nutritionist", "super_admin", "client"))
	{
		dietPlanRead.GET("/:id", planHandler.GetPlanAggregate)
	}

	// Diet plan CRUD + sub-resources (nutritionist or super_admin)
	dietPlans := r.Group("/api/diet-plans")
	dietPlans.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("nutritionist", "super_admin"))
	{
		dietPlans.POST("", planHandler.CreatePlan)
		dietPlans.PATCH("/:id", planHandler.UpdatePlanHeader)
		dietPlans.PATCH("/:id/activate", planHandler.ActivatePlan)
		dietPlans.DELETE("/:id", planHandler.DeletePlan)

		dietPlans.POST("/:id/days", planHandler.AddDay)
		dietPlans.PUT("/:id/days/:dayId", planHandler.UpdateDay)
		dietPlans.DELETE("/:id/days/:dayId", planHandler.DeleteDay)

		dietPlans.POST("/:id/days/:dayId/meals", planHandler.AddMeal)
		dietPlans.PUT("/:id/days/:dayId/meals/:mealId", planHandler.UpdateMeal)
		dietPlans.DELETE("/:id/days/:dayId/meals/:mealId", planHandler.DeleteMeal)
		dietPlans.PATCH("/:id/days/:dayId/meals/:mealId/order", planHandler.ReorderMeal)

		dietPlans.POST("/:id/days/:dayId/meals/:mealId/options", planHandler.AddOption)
		dietPlans.DELETE("/:id/days/:dayId/meals/:mealId/options/:optId", planHandler.DeleteOption)

		dietPlans.POST("/:id/days/:dayId/meals/:mealId/options/:optId/items", planHandler.AddItem)
		dietPlans.PUT("/:id/days/:dayId/meals/:mealId/options/:optId/items/:itemId", planHandler.UpdateItem)
		dietPlans.DELETE("/:id/days/:dayId/meals/:mealId/options/:optId/items/:itemId", planHandler.DeleteItem)

		dietPlans.POST("/:id/days/:dayId/exercises", planHandler.AddExercise)
		dietPlans.PUT("/:id/days/:dayId/exercises/:exId", planHandler.UpdateExercise)
		dietPlans.DELETE("/:id/days/:dayId/exercises/:exId", planHandler.DeleteExercise)

		dietPlans.POST("/:id/medications", planHandler.AddMedication)
		dietPlans.PUT("/:id/medications/:medId", planHandler.UpdateMedication)
		dietPlans.DELETE("/:id/medications/:medId", planHandler.DeleteMedication)
	}

	// Client active-plan route (client role only)
	clientGroup := r.Group("/api/clients")
	clientGroup.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("client"))
	{
		clientGroup.GET("/me/active-plan", planHandler.GetActivePlan)
		clientGroup.GET("/me/plans", planHandler.ListMyPlans)
	}

	client := r.Group("/api/client")
	client.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("client"))
	{
		client.GET("/tracking/daily", trackingHandler.GetDailyDashboard)
		client.POST("/food-logs", trackingHandler.LogFood)
		client.GET("/food-logs", trackingHandler.ListFoodLogs)
		client.POST("/water-logs", trackingHandler.LogWater)
		client.GET("/water-logs", trackingHandler.ListWaterLogs)
		client.POST("/sleep-logs", trackingHandler.UpsertSleep)
		client.GET("/sleep-logs", trackingHandler.GetSleepLog)
		client.POST("/exercise-logs", trackingHandler.LogExercise)
		client.GET("/exercise-logs", trackingHandler.ListExerciseLogs)
		client.POST("/medication-logs", trackingHandler.LogMedication)
		client.GET("/medication-logs", trackingHandler.ListMedicationLogs)
		client.POST("/body-measurements", trackingHandler.UpsertBodyMeasurement)
		client.GET("/body-measurements", trackingHandler.GetBodyMeasurement)
		client.GET("/body-measurements/history", trackingHandler.GetMeasurementHistory)
		client.GET("/body-measurements/weight-history", trackingHandler.GetWeightHistory)
		client.POST("/lab-results", trackingHandler.UploadLabResult)
		client.GET("/lab-results", trackingHandler.ListLabResults)
		client.POST("/food-requests", commHandler.ClientCreateFoodRequest)
		client.GET("/food-requests", commHandler.ClientListFoodRequests)
	}

	// Messaging routes (shared — both client and nutritionist)
	// IMPORTANT: specific paths registered BEFORE wildcard /:partnerId
	msgs := r.Group("/api/messages")
	msgs.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("client", "nutritionist"))
	{
		msgs.GET("/unread-count", commHandler.GetUnreadCount)
		msgs.GET("/attachment/:messageId", commHandler.DownloadAttachment)
		msgs.POST("", commHandler.SendMessage)
		msgs.GET("/:partnerId", commHandler.ListMessages)
		msgs.GET("/:partnerId/poll", commHandler.PollNewMessages)
		msgs.PATCH("/:partnerId/read", commHandler.MarkRead)
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info().Str("addr", srv.Addr).Msg("HTTP server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	logger.Info().Msg("shutdown signal received")

	// Graceful shutdown with 5-second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("server forced shutdown")
	}

	logger.Info().Msg("server stopped")
}

// setupLogger configures zerolog output based on environment and returns the logger instance.
func setupLogger(env string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if env == "development" {
		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = logger
		return logger
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Logger = logger
	return logger
}

// runMigrations runs database migrations from the db/migrations directory.
func runMigrations(databaseURL string) error {
	// Resolve migrations path relative to the executable
	migrationsPath := findMigrationsDir()

	// The pgx/v5 driver for golang-migrate uses "pgx5://" scheme.
	// Convert standard "postgres://" or "postgresql://" URLs.
	migrateURL := databaseURL
	if strings.HasPrefix(migrateURL, "postgres://") {
		migrateURL = "pgx5://" + strings.TrimPrefix(migrateURL, "postgres://")
	} else if strings.HasPrefix(migrateURL, "postgresql://") {
		migrateURL = "pgx5://" + strings.TrimPrefix(migrateURL, "postgresql://")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		migrateURL,
	)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	version, dirty, _ := m.Version()
	log.Info().
		Uint("version", version).
		Bool("dirty", dirty).
		Msg("database migrations complete")

	return nil
}

// findMigrationsDir locates the migrations directory.
func findMigrationsDir() string {
	// Check common paths
	candidates := []string{
		"db/migrations",
		"backend/db/migrations",
		"../db/migrations",
	}

	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		if info, err := os.Stat(abs); err == nil && info.IsDir() {
			return abs
		}
	}

	// Default fallback
	return "db/migrations"
}
