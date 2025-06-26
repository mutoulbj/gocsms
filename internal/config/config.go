package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         string
	OCPPPort           string
	PostgresDSN        string
	DBMaxOpenConns     int
	DBMaxIdleConns     int
	DBConnMaxLifeMin   int
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	JWTSecret          string
	JWTAccessTokenTTL  time.Duration
	JWTRefreshTokenTTL time.Duration
	JWTIssuer          string
}

func GocsmsConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found,using environment variables")
	}

	return &Config{
		ServerPort:         getEnv("SERVER_PORT", "8001"),
		OCPPPort:           getEnv("OCPP_PORT", "8003"),
		PostgresDSN:        getEnv("PG_DSN", ""),
		DBMaxOpenConns:     getEnvAsInt("DB_MAX_OPEN_CONNS", 50),
		DBMaxIdleConns:     getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifeMin:   getEnvAsInt("DB_CONN_MAX_LIFE_MIN", 60),
		RedisAddr:          getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            getEnvAsInt("REDIS_DB", 0),
		JWTSecret:          getEnv("JWT_SECRET", "1qaz2wsx3edc4rfv5tgb6yhn7ujm8ik9ol0p"),
		JWTAccessTokenTTL:  getEnvDuration("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
		JWTRefreshTokenTTL: getEnvDuration("JWT_REFRESH_TOKEN_TTL", 7*24*time.Hour),
		JWTIssuer:          getEnv("JWT_ISSUER", "gocsms"),
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
