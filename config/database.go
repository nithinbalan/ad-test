package config

import (
	"log"
	"sync"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDatabaseInstance() *gorm.DB {
	once.Do(func() {
		dsn := GetEnv("DATABASE_DSN", "")
		var err error // Declare error variable to capture the error from gorm.Open
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		log.Println("Database connection is successful")
	})
	return db
}
