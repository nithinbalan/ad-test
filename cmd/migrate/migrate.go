package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	// Check if an argument was provided (either "up" or "down")
	if len(os.Args) < 2 {
		log.Fatal("You must pass 'up' or 'down' as an argument")
	}
	migrationDirection := os.Args[1]

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the MSSQL database connection string from an environment variable
	dbURL := os.Getenv("DATABASE_DSN")
	if dbURL == "" {
		log.Fatal("DATABASE_DSN environment variable is not set.")
	}

	// Specify the directory where your migration files are located
	migrationsDir := os.Getenv("MIGRATION_DIRECTORY")

	// Construct the Goose command, using the argument passed to determine the migration direction
	cmd := exec.Command("goose", "-dir", migrationsDir, "mssql", dbURL, migrationDirection)

	// Set the standard output and standard error to os.Stdout and os.Stderr respectively
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the Goose command
	if err := cmd.Run(); err != nil {
		log.Fatalf("Goose migration '%s' failed: %v", migrationDirection, err)
	}

	log.Printf("Migration '%s' completed successfully.", migrationDirection)
}
