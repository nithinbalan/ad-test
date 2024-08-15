package main

import (
	"jethur/config"
	auditRoutes "jethur/grc/audit/routes"
	authRoutes "jethur/grc/auth/routes"
	commonRoutes "jethur/grc/common/routes"
	complianceRoutes "jethur/grc/compliance/routes"
	documentRoutes "jethur/grc/documents/routes"
	incidentRoutes "jethur/grc/incident/routes"
	libraryRoutes "jethur/grc/library/routes"
	masterRoutes "jethur/grc/masters/routes"
	meetingRoutes "jethur/grc/meeting/routes"
	organizationRoutes "jethur/grc/organization/routes"
	riskRoutes "jethur/grc/risk/routes"
	settingRoutes "jethur/grc/settings/routes"
	taskRoutes "jethur/grc/task/routes"
	"jethur/initializers"
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

	app := initializers.NewApp(db)

	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.LocalizerMiddleware(bundle))
	r.Use(middlewares.RecoveryWithDatabaseLog(db))
	r.Use(middlewares.ThrottleMiddleware(10, 20))

	masterRoutes.RegisterRoutes(app, r)
	meetingRoutes.RegisterRoutes(app, r)
	authRoutes.RegisterRoutes(app, r)
	complianceRoutes.RegisterRoutes(app, r)
	organizationRoutes.RegisterRoutes(app, r)
	settingRoutes.RegisterRoutes(app, r)
	libraryRoutes.RegisterRoutes(app, r)
	commonRoutes.RegisterRoutes(db, r)
	riskRoutes.RegisterRoutes(app, r)
	incidentRoutes.RegisterRoutes(app, r)
	taskRoutes.RegisterRoutes(app, r)
	documentRoutes.RegisterRoutes(app, r)
	auditRoutes.RegisterRoutes(app, r)

	if err := r.Run(":" + port); err != nil {
		utility.LogFatal("Failed to run server", map[string]interface{}{
			"error": err.Error(),
			"port":  port,
		})
	}
}
