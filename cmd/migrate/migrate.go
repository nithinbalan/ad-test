package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"jethur/assets"
	"log"
	"os"
	"path/filepath"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	_ = godotenv.Load()

	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate.exe [up|down|seed|unseed]")
	}
	command := os.Args[1]

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN not set")
	}

	switch command {
	case "up", "down":
		handleMigration(command, dsn)
	case "seed":
		if err := runSeederUp(dsn); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
	case "unseed":
		if err := runSeederDown(dsn); err != nil {
			log.Fatalf("Unseeding failed: %v", err)
		}
	default:
		log.Fatal("Invalid command. Use up, down, seed, or unseed")
	}

	fmt.Printf("Migration '%s' completed successfully.\n", command)
}

func handleMigration(direction, dsn string) {
	dir := extractMigrationsToTemp()
	defer os.RemoveAll(dir)

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	goose.SetDialect("sqlserver")

	switch direction {
	case "up":
		if err := goose.Up(db, dir); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
	case "down":
		if err := goose.DownTo(db, dir, 0); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
	}
}

func extractMigrationsToTemp() string {
	tmpDir, err := os.MkdirTemp("", "migrations")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}

	err = fs.WalkDir(assets.MigrationFS, "migrations", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		content, err := assets.MigrationFS.ReadFile(path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(tmpDir, filepath.Base(path))
		return os.WriteFile(destPath, content, 0644)
	})
	if err != nil {
		log.Fatalf("Failed to extract embedded migrations: %v", err)
	}

	return tmpDir
}

func runSeederUp(dsn string) error {
	tmpDir, err := extractSeedersToTemp()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return fmt.Errorf("DB connection error: %w", err)
	}
	defer db.Close()

	goose.SetDialect("sqlserver")
	goose.SetTableName("goose_seeder_versions")

	return goose.Up(db, tmpDir)
}

func runSeederDown(dsn string) error {
	tmpDir, err := extractSeedersToTemp()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return fmt.Errorf("DB connection error: %w", err)
	}
	defer db.Close()

	goose.SetDialect("sqlserver")
	goose.SetTableName("goose_seeder_versions")

	return goose.DownTo(db, tmpDir, 0)
}

func extractSeedersToTemp() (string, error) {
	tmpDir, err := os.MkdirTemp("", "seeders")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	err = fs.WalkDir(assets.SeederFS, "seeders", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		content, err := assets.SeederFS.ReadFile(path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(tmpDir, filepath.Base(path))
		return os.WriteFile(destPath, content, 0644)
	})
	if err != nil {
		return "", fmt.Errorf("failed to extract embedded seeders: %w", err)
	}

	return tmpDir, nil
}
