package main

import (
	"jethur/config"
	authRoutes "jethur/grc/auth/routes"
	commonRoutes "jethur/grc/common/routes"
	complianceRoutes "jethur/grc/compliance/routes"
	libraryRoutes "jethur/grc/libraries/routes"
	masterRoutes "jethur/grc/masters/routes"
	organizationRoutes "jethur/grc/organization/routes"
	settingRoutes "jethur/grc/settings/routes"
	"jethur/middlewares"
	"jethur/utility"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv() // Load environment variables

	time.Local = time.UTC

	port := config.GetEnv("PORT", "8080")
	rootPath := config.GetEnv("ROOT_PATH", "")
	if rootPath == "" {
		log.Fatal("ROOT_PATH is not set")
	}
	bundle := utility.LoadBundle(rootPath) //loading language bundles

	r := gin.Default()

	db := config.GetDatabaseInstance() // Obtain the database instance

	//middlewares
	r.Use(middlewares.LocalizerMiddleware(bundle))
	r.Use(middlewares.RecoveryWithDatabaseLog(db))
	r.Use(middlewares.ThrottleMiddleware(1, 5))

	masterRoutes.RegisterRoutes(db, r)
	authRoutes.RegisterRoutes(db, r)
	complianceRoutes.RegisterRoutes(db, r)
	organizationRoutes.RegisterRoutes(db, r)
	libraryRoutes.RegisterRoutes(db, r)
	settingRoutes.RegisterRoutes(db, r)
	commonRoutes.RegisterRoutes(db, r)

	// Start the server on specified port
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
