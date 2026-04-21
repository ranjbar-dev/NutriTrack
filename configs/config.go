package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	SMS      SMSConfig
	VAPID    VAPIDConfig
}

type AppConfig struct {
	Port               string `mapstructure:"PORT"`
	Env                string `mapstructure:"ENV"`
	TimeZone           string `mapstructure:"TIMEZONE"`
	CORSAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

type JWTConfig struct {
	AccessSecret  string `mapstructure:"JWT_ACCESS_SECRET"`
	RefreshSecret string `mapstructure:"JWT_REFRESH_SECRET"`
	AccessTTLMin  int    `mapstructure:"JWT_ACCESS_TTL_MIN"`
	RefreshTTLDay int    `mapstructure:"JWT_REFRESH_TTL_DAY"`
}

type SMSConfig struct {
	KavenegarAPIKey string `mapstructure:"KAVENEGAR_API_KEY"`
	OTPTemplate     string `mapstructure:"KAVENEGAR_OTP_TEMPLATE"`
}

type VAPIDConfig struct {
	PublicKey  string `mapstructure:"VAPID_PUBLIC_KEY"`
	PrivateKey string `mapstructure:"VAPID_PRIVATE_KEY"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig() // OK if .env doesn't exist; env vars take precedence

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.App.Port == "" {
		cfg.App.Port = "8080"
	}
	if cfg.App.TimeZone == "" {
		cfg.App.TimeZone = "Asia/Tehran"
	}
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}

	return &cfg, nil
}
