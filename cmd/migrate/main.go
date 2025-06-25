package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/mutoulbj/gocsms/internal/db"

	"github.com/uptrace/bun/migrate"
)

func main() {
	flag.Parse()
	args := flag.Args()

	db.Init()

	ctx := context.Background()
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(os.DirFS("migrations")); err != nil {
		log.Fatalf("discover migrations failed: %v", err)
	}
	mgr := migrate.NewMigrator(db.BunDB, migrations)

	if len(args) == 0 {
		log.Println("usage: go run cmd/migrate.go [init|up|down|reset]")
		os.Exit(1)
	}

	switch args[0] {
	case "init":
		if err := mgr.Init(ctx); err != nil {
			log.Fatalf("init failed: %v", err)
		}
	case "up":
		if _, err := mgr.Migrate(ctx); err != nil {
			log.Fatalf("up failed: %v", err)
		}
	case "down":
		if _, err := mgr.Rollback(ctx); err != nil {
			log.Fatalf("down failed: %v", err)
		}
	case "reset":
		if err := mgr.Reset(ctx); err != nil {
			log.Fatalf("reset failed: %v", err)
		}
	default:
		log.Println("invalid command")
		os.Exit(1)
	}
}
