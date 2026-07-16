package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	exePath, err := os.Executable()
	if err != nil {
		log.Println("Cannot find executable path, using existing environment variables")
		return
	}

	dir := filepath.Dir(exePath)
	envPath := filepath.Join(dir, ".env")

	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			log.Println("Failed to load .env from executable path:", err)
		} else {
			log.Println("Loaded .env from:", envPath)
		}
	} else {
		if err := godotenv.Load(); err != nil {
			log.Println(".env not found, using container/system environment variables")
		}
	}
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func GetBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}
