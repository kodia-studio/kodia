package main

import (
	"fmt"
	"os"
	"strconv"

	infradb "github.com/kodia-studio/kodia/internal/infrastructure/database"
	migrationspkg "github.com/kodia-studio/kodia/internal/infrastructure/database/migrations/go"
	"github.com/kodia-studio/kodia/internal/infrastructure/logger"
	"github.com/kodia-studio/kodia/pkg/config"
	dbpkg "github.com/kodia-studio/kodia/pkg/database"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(cfg.IsDevelopment())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Connect to database
	db, err := infradb.New(cfg, log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Create migrator
	migrator := dbpkg.NewMigrator(db, log)

	// Ensure migrations table exists
	if err := migrator.EnsureTable(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to ensure migrations table: %v\n", err)
		os.Exit(1)
	}

	// Get registered migrations and convert to the format expected by migrator
	registryEntries := migrationspkg.All()
	migrations := make([]struct {
		Name      string
		Migration dbpkg.Migration
	}, len(registryEntries))

	for i, entry := range registryEntries {
		migrations[i].Name = entry.Name
		migration, ok := entry.Migration.(dbpkg.Migration)
		if !ok {
			fmt.Fprintf(os.Stderr, "Invalid migration type for %s\n", entry.Name)
			os.Exit(1)
		}
		migrations[i].Migration = migration
	}

	// Parse command
	command := os.Args[1]

	switch command {
	case "up":
		if err := handleUp(migrator, migrations); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "down":
		steps := 1
		if len(os.Args) > 2 {
			s, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid step count: %v\n", err)
				os.Exit(1)
			}
			steps = s
		}
		if err := handleDown(migrator, migrations, steps); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "rollback":
		if err := handleDown(migrator, migrations, 1); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "status":
		if err := handleStatus(migrator, migrations); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "fresh":
		if err := handleFresh(migrator, migrations); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "reset":
		if err := handleReset(migrator, migrations); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleUp(migrator *dbpkg.Migrator, migrations []struct {
	Name      string
	Migration dbpkg.Migration
}) error {
	fmt.Println("Running pending migrations...")
	if err := migrator.RunAll(migrations); err != nil {
		return err
	}
	fmt.Println("✅ All migrations completed successfully")
	return nil
}

func handleDown(migrator *dbpkg.Migrator, migrations []struct {
	Name      string
	Migration dbpkg.Migration
}, steps int) error {
	fmt.Printf("Rolling back %d batch(es)...\n", steps)
	if err := migrator.Rollback(migrations, steps); err != nil {
		return err
	}
	fmt.Printf("✅ Successfully rolled back %d batch(es)\n", steps)
	return nil
}

func handleStatus(migrator *dbpkg.Migrator, migrations []struct {
	Name      string
	Migration dbpkg.Migration
}) error {
	statuses, err := migrator.Status(migrations)
	if err != nil {
		return err
	}

	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Migration Status                         ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")

	for _, status := range statuses {
		statusIcon := "⏳"
		if !status.Pending {
			statusIcon = "✅"
		}

		batchStr := "—"
		if status.Batch != nil {
			batchStr = fmt.Sprintf("%d", *status.Batch)
		}

		runAtStr := "—"
		if status.RunAt != nil {
			runAtStr = status.RunAt.Format("2006-01-02 15:04:05")
		}

		fmt.Printf("║ %s %-30s │ Batch: %-3s │ %s\n", statusIcon, status.Name, batchStr, runAtStr)
	}

	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	return nil
}

func handleFresh(migrator *dbpkg.Migrator, migrations []struct {
	Name      string
	Migration dbpkg.Migration
}) error {
	fmt.Println("⚠️  This will DROP all tables and re-run all migrations!")
	fmt.Print("Are you sure? (type 'yes' to confirm): ")

	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "yes" {
		fmt.Println("❌ Operation cancelled")
		return nil
	}

	fmt.Println("Running fresh migration...")
	if err := migrator.Fresh(migrations); err != nil {
		return err
	}
	fmt.Println("✅ Fresh migration completed successfully")
	return nil
}

func handleReset(migrator *dbpkg.Migrator, migrations []struct {
	Name      string
	Migration dbpkg.Migration
}) error {
	fmt.Println("⚠️  This will rollback all migrations and re-run them!")
	fmt.Print("Are you sure? (type 'yes' to confirm): ")

	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "yes" {
		fmt.Println("❌ Operation cancelled")
		return nil
	}

	fmt.Println("Running reset...")
	if err := migrator.Reset(migrations); err != nil {
		return err
	}
	fmt.Println("✅ Reset completed successfully")
	return nil
}

func printUsage() {
	usage := `
Kodia Database Migration CLI

Usage:
  migrate <command> [arguments]

Commands:
  up                 Run all pending migrations
  down [steps]       Rollback last N batches (default: 1)
  rollback           Alias for down 1
  status             Show migration status
  fresh              Drop all tables and re-run all migrations
  reset              Rollback all migrations and re-run them

Examples:
  migrate up
  migrate down
  migrate down 3
  migrate status
  migrate fresh
  migrate reset
`
	fmt.Println(usage)
}
