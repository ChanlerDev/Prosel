package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Cors     CorsConfig
	Auth     AuthConfig
}

type AppConfig struct {
	Environment string
	Version     string
}

type HTTPConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host          string
	Port          string
	User          string
	Password      string
	Name          string
	SSLMode       string
	MigrationsDir string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type CorsConfig struct {
	AllowedOrigins []string
}

type AuthConfig struct {
	JWTSecret          string
	JWTIssuer          string
	AccessTokenMinutes int
	RefreshTokenHours  int
	PasswordBcryptCost int
}

func Load() Config {
	return Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			Version:     getEnv("APP_VERSION", "dev"),
		},
		HTTP: HTTPConfig{Port: getEnv("PORT", "8080")},
		Database: DatabaseConfig{
			Host:          getEnv("DB_HOST", "localhost"),
			Port:          getEnv("DB_PORT", "5432"),
			User:          getEnv("DB_USER", "prosel"),
			Password:      getEnv("DB_PASSWORD", "prosel"),
			Name:          getEnv("DB_NAME", "prosel"),
			SSLMode:       getEnv("DB_SSLMODE", "disable"),
			MigrationsDir: getEnv("MIGRATIONS_DIR", "migrations"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Cors: CorsConfig{AllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ",")},
		Auth: AuthConfig{
			JWTSecret:          getEnv("JWT_SECRET", "dev-secret-change-me"),
			JWTIssuer:          getEnv("JWT_ISSUER", "prosel"),
			AccessTokenMinutes: getEnvInt("ACCESS_TOKEN_MINUTES", 15),
			RefreshTokenHours:  getEnvInt("REFRESH_TOKEN_HOURS", 168),
			PasswordBcryptCost: getEnvInt("PASSWORD_BCRYPT_COST", 12),
		},
	}
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC", c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

func (c HTTPConfig) Address() string {
	return ":" + c.Port
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func StartupTimeout() time.Duration {
	return 10 * time.Second
}
