package main

import (
	"jethur/config"
	authRoutes "jethur/grc/auth/routes"
	commonRoutes "jethur/grc/common/routes"
	complianceRoutes "jethur/grc/compliance/routes"
	incidentRoutes "jethur/grc/incident/routes"
	masterRoutes "jethur/grc/masters/routes"
	organizationRoutes "jethur/grc/organization/routes"
	riskRoutes "jethur/grc/risk/routes"
	settingRoutes "jethur/grc/settings/routes"
	"jethur/middlewares"
	"jethur/utility"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv()

	time.Local = time.UTC

	port := config.GetEnv("PORT", "8080")
	rootPath := config.GetEnv("ROOT_PATH", "")
	if rootPath == "" {
		utility.LogFatal("ROOT_PATH missing in env file", nil)
	}
	bundle := utility.LoadBundle(rootPath) //loading language bundles

	r := gin.Default()

	db := config.GetDatabaseInstance() // Obtain the database instance

	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.LocalizerMiddleware(bundle))
	r.Use(middlewares.RecoveryWithDatabaseLog(db))
	r.Use(middlewares.ThrottleMiddleware(1, 5))

	masterRoutes.RegisterRoutes(db, r)
	authRoutes.RegisterRoutes(db, r)
	complianceRoutes.RegisterRoutes(db, r)
	organizationRoutes.RegisterRoutes(db, r)
	settingRoutes.RegisterRoutes(db, r)
	commonRoutes.RegisterRoutes(db, r)
	riskRoutes.RegisterRoutes(db, r)
	incidentRoutes.RegisterRoutes(db, r)

	if err := r.Run(":" + port); err != nil {
		utility.LogFatal("Failed to run server", map[string]interface{}{
			"error": err.Error(),
			"port":  port,
		})
	}
}
