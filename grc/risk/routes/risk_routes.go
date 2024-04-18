package routes

import (
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/risk/controllers"
	riskRepositories "jethur/grc/risk/repositories"
	riskServices "jethur/grc/risk/services"

	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"

	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	riskRoutes := router.Group("/")

	riskRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	severityLevelRepo := masterRepositories.NewSeverityLevelRepository(db)
	severityLevelSvc := masterServices.NewSeverityLevelService(severityLevelRepo)

	riskCategoryRepo := riskRepositories.NewRiskCategoryRepository(db)
	riskCategorySvc := riskServices.NewRiskCategoryService(riskCategoryRepo, activityLogSvc)
	riskCategoryController := controllers.NewRiskCategoryController(riskCategorySvc)

	potentialRiskStatusRepo := riskRepositories.NewPotentialRiskStatusRepository(db)
	potentialRiskStatusSvc := riskServices.NewPotentialRiskStatusService(potentialRiskStatusRepo)
	potentialRiskStatusController := controllers.NewPotentialRiskStatusController(potentialRiskStatusSvc)

	riskLevelRepo := riskRepositories.NewRiskLevelRepository(db)
	riskLevelSvc := riskServices.NewRiskLevelService(riskLevelRepo)
	riskLevelController := controllers.NewRiskLevelController(riskLevelSvc)

	potentialRiskRepo := riskRepositories.NewPotentialRiskRepository(db)
	potentialRiskSvc := riskServices.NewPotentialRiskService(
		potentialRiskRepo,
		activityLogSvc,
		severityLevelSvc,
		riskCategorySvc,
		potentialRiskStatusSvc,
	)
	potentialRiskController := controllers.NewPotentialRiskController(potentialRiskSvc)

	potentialRiskStatusUpdateRepo := riskRepositories.NewPotentialRiskStatusUpdateRepository(db)
	potentialRiskStatusUpdateSvc := riskServices.NewPotentialRiskStatusUpdateService(potentialRiskStatusUpdateRepo, activityLogSvc, potentialRiskSvc, potentialRiskStatusSvc)
	potentialRiskStatusUpdateController := controllers.NewPotentialRiskStatusUpdateController(potentialRiskStatusUpdateSvc)

	riskAppetiteStatementRepo := riskRepositories.NewRiskAppetiteStatementRepository(db)
	riskAppetiteStatementSvc := riskServices.NewRiskAppetiteStatementService(
		riskAppetiteStatementRepo,
		activityLogSvc,
		riskCategorySvc,
	)
	riskAppetiteStatementController := controllers.NewRiskAppetiteStatementController(riskAppetiteStatementSvc)

	riskRoutes.GET("/risk-categories", riskCategoryController.GetAllRiskCategories)
	riskRoutes.GET("/risk-categories/:id", riskCategoryController.GetRiskCategoryDetails)
	riskRoutes.POST("/risk-categories", riskCategoryController.CreateRiskCategory)
	riskRoutes.PUT("/risk-categories/:id", riskCategoryController.UpdateRiskCategory)
	riskRoutes.POST("/risk-categories/delete", riskCategoryController.DeleteRiskCategories)
	riskRoutes.PUT("/risk-categories/activate", riskCategoryController.ActivateRiskCategories)
	riskRoutes.PUT("/risk-categories/deactivate", riskCategoryController.DeactivateRiskCategories)

	riskRoutes.GET("/potential-risk-statuses", potentialRiskStatusController.GetAllPotentialRiskStatuses)
	riskRoutes.GET("/potential-risk-statuses/:id", potentialRiskStatusController.GetPotentialRiskStatusDetails)

	riskRoutes.GET("/potential-risks", potentialRiskController.GetAllPotentialRisks)
	riskRoutes.GET("/potential-risks/:id", potentialRiskController.GetPotentialRiskDetails)
	riskRoutes.POST("/potential-risks", potentialRiskController.CreatePotentialRisk)
	riskRoutes.PUT("/potential-risks/:id", potentialRiskController.UpdatePotentialRisk)
	riskRoutes.POST("/potential-risks/delete", potentialRiskController.DeletePotentialRisks)

	riskRoutes.GET("/risk-levels", riskLevelController.GetAllRiskLevels)
	riskRoutes.GET("/risk-levels/:id", riskLevelController.GetRiskLevelDetails)

	riskRoutes.GET("/potential-risks/:id/status-updates", potentialRiskStatusUpdateController.GetAllPotentialRiskStatusUpdates)
	riskRoutes.POST("/potential-risks/:id/status-updates", potentialRiskStatusUpdateController.CreatePotentialRiskStatusUpdate)

	riskRoutes.GET("/risk-appetite-statements", riskAppetiteStatementController.GetAllRiskAppetiteStatements)
	riskRoutes.GET("/risk-appetite-statements/:id", riskAppetiteStatementController.GetRiskAppetiteStatementDetails)
	riskRoutes.POST("/risk-appetite-statements", riskAppetiteStatementController.CreateRiskAppetiteStatement)
	riskRoutes.PUT("/risk-appetite-statements/:id", riskAppetiteStatementController.UpdateRiskAppetiteStatement)
	riskRoutes.POST("/risk-appetite-statements/delete", riskAppetiteStatementController.DeleteRiskAppetiteStatements)
	riskRoutes.PUT("/risk-appetite-statements/activate", riskAppetiteStatementController.ActivateRiskAppetiteStatements)
	riskRoutes.PUT("/risk-appetite-statements/deactivate", riskAppetiteStatementController.DeactivateRiskAppetiteStatements)
}
