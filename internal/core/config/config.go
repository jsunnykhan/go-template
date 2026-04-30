package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Mode string
}

type ServerConfig struct {
	Host string
	Port int
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode)
}

type JWTConfig struct {
	SecretKey                    string
	ExpirationTimeInHours        int
	RefreshExpirationTimeInHours int
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Mode: getEnv("APP_MODE", "development"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnvInt("DATABASE_PORT", 3306),
			User:     getEnv("DATABASE_USER", "root"),
			Password: getEnv("DATABASE_PASSWORD", "root"),
			DBName:   getEnv("DATABASE_NAME", "myapp"),
			SSLMode:  getEnv("DATABASE_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey:                    getEnv("JWT_SECRET_KEY", "mysecretkey"),
			ExpirationTimeInHours:        getEnvInt("JWT_EXPIRATION_TIME_IN_HOURS", 4),
			RefreshExpirationTimeInHours: getEnvInt("JWT_REFRESH_EXPIRATION_TIME_IN_HOURS", 24),
		},
	}
	return cfg, nil
}

// ---------- helper -------------

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}
