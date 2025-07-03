package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	Database DatabaseConfig
	Redis RedisConfig
	JWT JWTConfig
}

type ServerConfig struct {
	ServerPort 	   string
	OCPPPort       string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

type DatabaseConfig struct {
	Host         string
    Port         string
    User         string
    Password     string
    DBName       string
    SSLMode      string
    MaxOpenConns int
    MaxIdleConns int
    MaxLifetime  time.Duration
}

type RedisConfig struct {
	Host		 string
	Port		 string
	Password	 string
	DB			 int
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found,using environment variables")
	}
	return &Config{
		Server: ServerConfig{
			ServerPort:         getEnv("SERVER_PORT", "8001"),
			OCPPPort:           getEnv("OCPP_PORT", "8003"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			User:         getEnv("DB_USER", "user"),
			Password:     getEnv("DB_PASSWORD", "password"),
			DBName:       getEnv("DB_NAME", "dbname"),
			SSLMode:      getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 50),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxLifetime:  getEnvDuration("DB_MAX_LIFETIME", 60*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:      getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "1qaz2wsx3edc4rfv5tgb6yhn7ujm8ik9ol0p"),
			AccessTokenTTL:  getEnvDuration("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: getEnvDuration("JWT_REFRESH_TOKEN_TTL", 7*24*time.Hour),
			Issuer:          getEnv("JWT_ISSUER", "gocsms"),
		},
	}
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if valueStr := os.Getenv(key); valueStr != "" {
		if d, err := time.ParseDuration(valueStr); err == nil {
			return d
		}
	}
	return defaultValue
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvAsInt(key string, fallback int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}
