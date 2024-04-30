package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	exePath, err := os.Executable()
	if err != nil {
		panic("Cannot find executable path")
	}
	dir := filepath.Dir(exePath)
	envPath := filepath.Join(dir, ".env")

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		if err := godotenv.Load(); err != nil {
			panic("Error loading default .env file")
		}
	} else {
		if err := godotenv.Load(envPath); err != nil {
			panic("Error loading .env file from " + envPath)
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
