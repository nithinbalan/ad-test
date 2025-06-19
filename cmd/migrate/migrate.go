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
		log.Fatal("Usage: migrate.exe [up|down]")
	}
	dir := extractMigrationsToTemp()
	defer os.RemoveAll(dir) // Clean up after run

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN not set")
	}

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	goose.SetDialect("sqlserver")

	var gooseErr error
	switch os.Args[1] {
	case "up":
		gooseErr = goose.Up(db, dir)
	case "down":
		gooseErr = goose.DownTo(db, dir, 0)
	case "seed":
		runSeeders(db)
	default:
		log.Fatal("Invalid command. Use up, down, or seed")
	}

	if gooseErr != nil {
		log.Fatalf("Migration error: %v", gooseErr)
	}

	fmt.Printf("Migration '%s' completed successfully.\n", os.Args[1])
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

func runSeeders(db *sql.DB) {
	err := fs.WalkDir(assets.SeederFS, "seeders", func(path string, d fs.DirEntry, err error) error {
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
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("error executing seeder %s: %w", path, err)
		}
		fmt.Printf("Executed seeder: %s\n", path)
		return nil
	})
	if err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}
}
