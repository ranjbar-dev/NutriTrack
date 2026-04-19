// Package main provides a CLI tool for seeding the Super Admin account (D-03, AUTH-02).
// Usage: ADMIN_EMAIL=admin@example.com ADMIN_PASSWORD=secure123 DATABASE_URL=... go run cmd/seed/main.go
// This is run once during initial deployment, not during every startup.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

func main() {
	// Load required environment variables
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	dbURL := os.Getenv("DATABASE_URL")

	if email == "" || password == "" || dbURL == "" {
		fmt.Fprintln(os.Stderr, "ADMIN_EMAIL, ADMIN_PASSWORD, and DATABASE_URL must be set")
		os.Exit(1)
	}

	ctx := context.Background()

	// Connect to PostgreSQL
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping database: %v\n", err)
		os.Exit(1)
	}

	// Hash password with bcrypt cost 12 (AUTH-10)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to hash password: %v\n", err)
		os.Exit(1)
	}

	// Insert super_admin user via sqlc
	queries := sqlc.New(pool)
	_, err = queries.CreateUser(ctx, sqlc.CreateUserParams{
		Role:           sqlc.UserRoleSuperAdmin,
		FullName:       "مدیر سیستم",
		Email:          pgtype.Text{String: email, Valid: true},
		PasswordHash:   pgtype.Text{String: string(hash), Valid: true},
		Mobile:         pgtype.Text{Valid: false},
		DateOfBirth:    pgtype.Date{Valid: false},
		HeightCm:       pgtype.Float4{Valid: false},
		Gender:         sqlc.NullGenderType{Valid: false},
		NutritionistID: pgtype.UUID{Valid: false},
		Notes:          pgtype.Text{Valid: false},
	})
	if err != nil {
		// Handle duplicate email gracefully (idempotent)
		fmt.Fprintf(os.Stderr, "⚠️  Super Admin may already exist or error occurred: %v\n", err)
		fmt.Println("ℹ️  If the admin already exists, this is expected. No changes made.")
		os.Exit(0)
	}

	fmt.Println("✅ Super Admin seeded successfully")
	fmt.Printf("   Email: %s\n", email)
}
