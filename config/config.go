// config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from the .env file.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .env file")
	}
}

// GetEnv retrieves the value of an environment variable or returns a default value if not set.
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
