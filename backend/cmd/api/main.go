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
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Setup zerolog
	setupLogger(cfg.Environment)

	log.Info().
		Str("environment", cfg.Environment).
		Str("port", cfg.Port).
		Msg("starting NutriTrack API")

	// Create root context with signal cancellation for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Connect to PostgreSQL
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database pool")
	}
	defer pool.Close()

	// Verify database connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database")
	}
	log.Info().Msg("database connection established")

	// Run migrations
	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	// Create Gin engine — use gin.New() not gin.Default() per D-07
	// (we add our own middleware in Plan 03)
	r := gin.New()

	// Register routes
	r.GET("/api/health", handler.HealthCheck)

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
		log.Info().Str("addr", srv.Addr).Msg("HTTP server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Info().Msg("shutdown signal received")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("server forced shutdown")
	}

	log.Info().Msg("server stopped")
}

// setupLogger configures zerolog output based on environment.
func setupLogger(env string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
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
