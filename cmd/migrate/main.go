// cmd/migrate/main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Import your database driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Import the file source driver
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set.")
	}

	m, err := migrate.New(
		"file://migrations", // Path to your migration files
		databaseURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Example commands:
	args := os.Args[1:] // Get command-line arguments

	if len(args) == 0 {
		fmt.Println("Usage: go run cmd/migrate/main.go [up|down|goto VERSION|force VERSION|version|create NAME]")
		return
	}

	switch args[0] {
	case "up":
		fmt.Println("Applying all pending migrations...")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		fmt.Println("Migrations applied successfully (or no change).")
	case "down":
		fmt.Println("Rolling back one migration...")
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		fmt.Println("Migration rolled back successfully (or no change).")
	case "version":
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		fmt.Printf("Current database version: %d (dirty: %t)\n", version, dirty)
	case "goto":
		if len(args) < 2 {
			log.Fatal("Usage: goto VERSION")
		}
		versionStr := args[1]
		var version uint
		_, err := fmt.Sscanf(versionStr, "%d", &version)
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		fmt.Printf("Migrating to version %d...\n", version)
		if err := m.Migrate(version); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to migrate to version %d: %v", version, err)
		}
		fmt.Println("Migration completed.")
	case "force":
		if len(args) < 2 {
			log.Fatal("Usage: force VERSION")
		}
		versionStr := args[1]
		var version int
		_, err := fmt.Sscanf(versionStr, "%d", &version)
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		fmt.Printf("Forcing database version to %d (use with caution!)...\n", version)
		if err := m.Force(version); err != nil { // Force doesn't return ErrNoChange
			log.Fatalf("Failed to force version %d: %v", version, err)
		}
		fmt.Println("Force completed.")
	case "create":
		if len(args) < 2 {
			log.Fatal("Usage: create NAME")
		}
		name := args[1]
		// This part is tricky. 'migrate' programmatically doesn't offer 'create'.
		// You'd typically use the CLI 'migrate create' for this.
		fmt.Printf("To create new migrations, please use the 'migrate create' CLI command directly:\n")
		fmt.Printf("  migrate create -ext sql -dir migrations -seq %s\n", name)
	default:
		fmt.Println("Unknown command.")
	}
}
