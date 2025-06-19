package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort       string
	OCPPPort         string
	PostgresDSN      string
	DBMaxOpenConns   int
	DBMaxIdleConns   int
	DBConnMaxLifeMin int
	RedisAddr        string
	RedisPassword    string
	RedisDB          int
}

func GocsmsConfig() *Config {
    if err := godotenv.Load(); err != nil {
		log.Println(".env file not found,using environment variables")
	}

	return &Config{
        ServerPort:       getEnv("PORT", "8001"),
		OCPPPort:         getEnv("OCPP_PORT", "8003"),
		PostgresDSN:      getEnv("PG_DSN", ""),
		DBMaxOpenConns:   getEnvAsInt("DB_MAX_OPEN_CONNS", 50),
		DBMaxIdleConns:   getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifeMin: getEnvAsInt("DB_CONN_MAX_LIFE_MIN", 60),
		RedisAddr:        getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          getEnvAsInt("REDIS_DB", 0),
	}
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
