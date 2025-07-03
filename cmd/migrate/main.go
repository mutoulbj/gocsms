package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mutoulbj/gocsms/internal/config"
	"github.com/mutoulbj/gocsms/pkg/db"

	"github.com/uptrace/bun/migrate"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 初始化配置
	appConfig := config.NewConfig()

	// 初始化数据库连接
	database, err := db.NewDB(&appConfig.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	// 解析命令行参数
	flag.Parse()
	args := flag.Args()

	ctx := context.Background()
	migrations := migrate.NewMigrations()

	// 注意：确保migrations目录存在且路径正确
	if err := migrations.Discover(os.DirFS("migrations")); err != nil {
		log.Fatalf("discover migrations failed: %v", err)
	}

	// 使用初始化后的数据库连接
	mgr := migrate.NewMigrator(database.DB, migrations)

	if len(args) == 0 {
		// 修正Usage提示，使其与实际命令一致
		log.Println("usage: go run cmd/migrate/main.go [init|up|down|reset]")
		os.Exit(1)
	}

	switch args[0] {
	case "init":
		if err := mgr.Init(ctx); err != nil {
			log.Fatalf("init failed: %v", err)
		}
		log.Println("Migration initialized successfully")
	case "up":
		group, err := mgr.Migrate(ctx)
		if err != nil {
			log.Fatalf("up failed: %v", err)
		}
		if group.ID == 0 {
			log.Println("No migrations to run")
		} else {
			log.Printf("Migrated to %s\n", group)
		}
	case "down":
		group, err := mgr.Rollback(ctx)
		if err != nil {
			log.Fatalf("down failed: %v", err)
		}
		if group.ID == 0 {
			log.Println("No migrations to rollback")
		} else {
			log.Printf("Rolled back %s\n", group)
		}
	case "reset":
		if err := mgr.Reset(ctx); err != nil {
			log.Fatalf("reset failed: %v", err)
		}
		log.Println("Migration reset successfully")
	default:
		log.Println("invalid command")
		os.Exit(1)
	}
}
