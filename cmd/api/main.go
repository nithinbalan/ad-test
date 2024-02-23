package main

import (
	"jethur/config"
	authRoutes "jethur/grc/auth/routes"
	masterRoutes "jethur/grc/masters/routes"
	"jethur/middleware"
	"jethur/utility"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv() // Load environment variables

	port := config.GetEnv("PORT", "8080")
	rootPath := config.GetEnv("ROOT_PATH", "")
	if rootPath == "" {
		log.Fatal("ROOT_PATH is not set")
	}
	bundle := utility.LoadBundle(rootPath) //loading language bundles

	r := gin.Default()

	db := config.GetDatabaseInstance() // Obtain the database instance

	//middlewares
	r.Use(middleware.LocalizerMiddleware(bundle))
	r.Use(middleware.RecoveryWithDatabaseLog(db))
	r.Use(middleware.ThrottleMiddleware(1, 5))

	masterRoutes.RegisterRoutes(db, r)
	authRoutes.RegisterRoutes(db, r)

	// Start the server on specified port
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
