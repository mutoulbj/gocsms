package db

import (
	"context"
	"log"
	"time"
	"database/sql"

	"github.com/mutoulbj/gocsms/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var BunDB *bun.DB

func Init() {
	db, err := sql.Open("pgx", config.Cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	// 设置连接池参数，来自配置文件
	db.SetMaxOpenConns(config.Cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(config.Cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.Cfg.DBConnMaxLifeMin) * time.Minute)

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	BunDB = bun.NewDB(db, pgdialect.New())
}
