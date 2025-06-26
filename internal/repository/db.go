package repository

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/mutoulbj/gocsms/internal/config"
)

// GocsmsBunDB initializes a Bun ORM database connection for PostgreSQL
func NewBunDB(cfg *config.Config, log *logrus.Logger) *bun.DB {
	sqldb, err := sql.Open("pgx", config.GocsmsConfig().PostgresDSN)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	// 设置连接池参数，来自配置文件
	sqldb.SetMaxOpenConns(config.GocsmsConfig().DBMaxOpenConns)
	sqldb.SetMaxIdleConns(config.GocsmsConfig().DBMaxIdleConns)
	sqldb.SetConnMaxLifetime(time.Duration(config.GocsmsConfig().DBConnMaxLifeMin) * time.Minute)

	ctx := context.Background()
	if err := sqldb.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	log.Info("successfully connected to PostgreSQL")
	return db
}

// GocsmsRedisClient initializes a Redis client
func NewRedisClient(cfg *config.Config, log *logrus.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}

	log.Info("Successfully connected to Redis")
	return client
}
