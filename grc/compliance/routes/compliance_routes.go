package routes

import (
	"jethur/grc/compliance/controllers"
	complianceRepositories "jethur/grc/compliance/repositories"
	complianceServices "jethur/grc/compliance/services"
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	complianceRoutes := router.Group("/")

	complianceRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	//Compliance Category
	complianceCategoryRepo := complianceRepositories.NewComplianceCategoryRepository(db)
	complianceCategorySvc := complianceServices.NewComplianceCategoryService(complianceCategoryRepo, activityLogSvc)
	complianceCategoryController := controllers.NewComplianceCategoryController(complianceCategorySvc)

	complianceRoutes.GET("/compliance-categories", complianceCategoryController.GetAllComplianceCategories)
	complianceRoutes.GET("/compliance-categories/:id", complianceCategoryController.GetComplianceCategoryDetails)
	complianceRoutes.POST("/compliance-categories", complianceCategoryController.CreateComplianceCategory)
	complianceRoutes.PUT("/compliance-categories/:id", complianceCategoryController.UpdateComplianceCategory)
	complianceRoutes.POST("/compliance-categories/delete", complianceCategoryController.DeleteComplianceCategories)
	complianceRoutes.PUT("/compliance-categories/activate", complianceCategoryController.ActivateComplianceCategories)
	complianceRoutes.PUT("/compliance-categories/deactivate", complianceCategoryController.DeactivateComplianceCategories)

	//Compliance Status
	complianceStatusRepo := complianceRepositories.NewComplianceStatusRepository(db)
	complianceStatusSvc := complianceServices.NewComplianceStatusService(complianceStatusRepo)
	complianceStatusController := controllers.NewComplianceStatusController(complianceStatusSvc)

	complianceRoutes.GET("/compliance-statuses", complianceStatusController.GetAllComplianceStatuses)
	complianceRoutes.GET("/compliance-statuses/:id", complianceStatusController.GetComplianceStatusDetails)

	//SLA Category
	SLACategoryRepo := complianceRepositories.NewSLACategoryRepository(db)
	SLACategorySvc := complianceServices.NewSLACategoryService(SLACategoryRepo, activityLogSvc)
	SLACategoryController := controllers.NewSLACategoryController(SLACategorySvc)

	complianceRoutes.GET("/sla-categories/:id", SLACategoryController.GetSLACategoryDetails)
	complianceRoutes.POST("/sla-categories", SLACategoryController.CreateCategory)
	complianceRoutes.GET("/sla-categories", SLACategoryController.GetAllSLACategories)
	complianceRoutes.PUT("/sla-categories/:id", SLACategoryController.UpdateCategory)
	complianceRoutes.POST("/sla-categories/delete", SLACategoryController.DeleteSLACategories)
	complianceRoutes.POST("/sla-categories/activate", SLACategoryController.ActivateSLACategories)
	complianceRoutes.POST("/sla-categories/deactivate", SLACategoryController.DeactivateSLACategories)

	//Compliance Register
	complianceRepositories := complianceRepositories.NewComplianceRepository(db)
	reviewFrequencyRepo := masterRepositories.NewReviewFrequencyRepository(db)
	severityLevelRepo := masterRepositories.NewSeverityLevelRepository(db)
	complianceSvc := complianceServices.NewComplianceService(complianceRepositories,
		activityLogSvc,
		complianceCategoryRepo, //@todo add org unit and user repo
		complianceStatusRepo,
		reviewFrequencyRepo,
		severityLevelRepo)
	complianceController := controllers.NewComplianceController(complianceSvc)

	complianceRoutes.GET("/compliance-register", complianceController.GetAllCompliances)
	complianceRoutes.GET("/compliance-register/:id", complianceController.GetComplianceDetails)
	complianceRoutes.POST("/compliance-register", complianceController.CreateCompliance)
	complianceRoutes.PUT("/compliance-register/:id", complianceController.UpdateCompliance)
	complianceRoutes.POST("/compliance-register/delete", complianceController.DeleteCompliances)
	complianceRoutes.PUT("/compliance-register/activate", complianceController.ActivateCompliances)
	complianceRoutes.PUT("/compliance-register/deactivate", complianceController.DeactivateCompliances)
}
