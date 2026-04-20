package config

import (
	"errors"
	"os"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port        string // HTTP server port (default "8080")
	DatabaseURL string // PostgreSQL connection URL (required)
	JWTSecret   string // JWT signing secret, min 32 chars (required)
	Environment string // "development" or "production" (default "development")
	FrontendURL string // Frontend URL for CORS (default "http://localhost:3000")
	UploadsDir  string // Base directory for uploaded files (default "./uploads")
	SMSAPIKey   string // Kavenegar API key for OTP delivery
	SMSTemplate string // Kavenegar verification template name
}

// Load reads configuration from environment variables and validates required fields.
func Load() (*Config, error) {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Environment: getEnv("ENVIRONMENT", "development"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		UploadsDir:  getEnv("UPLOADS_DIR", "./uploads"),
		SMSAPIKey:   os.Getenv("SMS_API_KEY"),
		SMSTemplate: os.Getenv("SMS_TEMPLATE"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	if c.JWTSecret == "" {
		return errors.New("JWT_SECRET is required")
	}

	if len(c.JWTSecret) < 32 {
		return errors.New("JWT_SECRET must be at least 32 characters")
	}

	if c.Environment != "development" && c.Environment != "production" {
		return errors.New("ENVIRONMENT must be 'development' or 'production'")
	}

	return nil
}

// getEnv returns the value of the environment variable named by key,
// or defaultVal if the variable is not present.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
