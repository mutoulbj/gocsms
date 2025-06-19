package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/mutoulbj/gocsms/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var BunDB *bun.DB

func Init() {
	db, err := sql.Open("pgx", config.GocsmsConfig().PostgresDSN)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	// set up pool options
	db.SetMaxOpenConns(config.GocsmsConfig().DBMaxOpenConns)
	db.SetMaxIdleConns(config.GocsmsConfig().DBMaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.GocsmsConfig().DBConnMaxLifeMin) * time.Minute)

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	BunDB = bun.NewDB(db, pgdialect.New())
}
