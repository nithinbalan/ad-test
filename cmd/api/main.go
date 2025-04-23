package main

import (
	"flag"
	"jethur/config"
	assetRoutes "jethur/grc/asset/routes"
	auditInsightsRoutes "jethur/grc/audit-insights/routes"
	auditRoutes "jethur/grc/audit/routes"
	authRoutes "jethur/grc/auth/routes"
	businessContinuityRoutes "jethur/grc/business-continuity/routes"
	commonRoutes "jethur/grc/common/routes"
	complianceRoutes "jethur/grc/compliance/routes"
	documentRoutes "jethur/grc/documents/routes"
	incidentRoutes "jethur/grc/incident/routes"
	kpiRoutes "jethur/grc/kpi/routes"
	libraryRoutes "jethur/grc/library/routes"
	logRoutes "jethur/grc/logs/routes"
	masterRoutes "jethur/grc/masters/routes"
	meetingRoutes "jethur/grc/meeting/routes"
	organizationRoutes "jethur/grc/organization/routes"
	organizationalContextRoutes "jethur/grc/organizational-context/routes"
	reportRoutes "jethur/grc/reports/routes"
	riskRoutes "jethur/grc/risk/routes"
	settingRoutes "jethur/grc/settings/routes"
	strategyRoutes "jethur/grc/strategy/routes"
	taskRoutes "jethur/grc/task/routes"
	"jethur/initializers"
	"jethur/middlewares"
	"jethur/utility"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv() // Load environment variables first

	// Define command-line flags
	checkPDF := flag.Bool("check-pdf", false, "Check if PDF conversion is working")
	flag.Parse()

	// Handle the `--check-pdf` flag
	if *checkPDF {
		utility.RunPDFConversionCheck()
		return
	}

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
	r.Use(middlewares.LicenseMiddleware(db))

	masterRoutes.RegisterRoutes(app, r)
	logRoutes.RegisterRoutes(app, r)
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
	kpiRoutes.RegisterRoutes(app, r)
	auditInsightsRoutes.RegisterRoutes(app, r)
	reportRoutes.RegisterRoutes(app, r)
	assetRoutes.RegisterRoutes(app, r)
	businessContinuityRoutes.RegisterRoutes(app, r)
	organizationalContextRoutes.RegisterRoutes(app, r)
	strategyRoutes.RegisterRoutes(app, r)

	if err := r.Run(":" + port); err != nil {
		utility.LogFatal("Failed to run server", map[string]interface{}{
			"error": err.Error(),
			"port":  port,
		})
	}
}
