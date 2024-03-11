package routes

import (
	"jethur/grc/compliance/controllers"
	complianceRepo "jethur/grc/compliance/repositories"
	complianceService "jethur/grc/compliance/services"
	activityLogRepository "jethur/grc/logs/repositories"
	activityLogServie "jethur/grc/logs/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	complianceRoutes := router.Group("/")

	complianceRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := activityLogRepository.NewActivityLogRepository(db)
	activityLogSvc := activityLogServie.NewActivityLogService(activityLogRepo)

	//Compliance Category
	complianceCategoryRepo := complianceRepo.NewComplianceCategoryRepository(db)
	complianceCategorySvc := complianceService.NewComplianceCategoryService(complianceCategoryRepo, activityLogSvc)
	complianceCategoryController := controllers.NewComplianceCategoryController(complianceCategorySvc)

	complianceRoutes.GET("/compliance-categories", complianceCategoryController.GetAllComplianceCategories)
	complianceRoutes.GET("/compliance-categories/:id", complianceCategoryController.GetComplianceCategoryDetails)
	complianceRoutes.POST("/compliance-categories", complianceCategoryController.CreateComplianceCategory)
	complianceRoutes.PUT("/compliance-categories/:id", complianceCategoryController.UpdateComplianceCategory)
	complianceRoutes.POST("/compliance-categories/delete", complianceCategoryController.DeleteComplianceCategories)
	complianceRoutes.POST("/compliance-categories/activate", complianceCategoryController.ActivateComplianceCategories)
	complianceRoutes.POST("/compliance-categories/deactivate", complianceCategoryController.DeactivateComplianceCategories)

	//Compliance Status
	complianceStatusRepo := complianceRepo.NewComplianceStatusRepository(db)
	complianceStatusSvc := complianceService.NewComplianceStatusService(complianceStatusRepo)
	complianceStatusController := controllers.NewComplianceStatusController(complianceStatusSvc)

	complianceRoutes.GET("/compliance-statuses", complianceStatusController.GetAllComplianceStatuses)

	//Compliance Register
	complianceRepo := complianceRepo.NewComplianceRepository(db)
	complianceSvc := complianceService.NewComplianceService(complianceRepo, activityLogSvc)
	complianceController := controllers.NewComplianceController(complianceSvc)

	complianceRoutes.GET("/compliance-register", complianceController.GetAllCompliances)
	complianceRoutes.GET("/compliance-register/:id", complianceController.GetComplianceDetails)
	complianceRoutes.POST("/compliance-register", complianceController.CreateCompliance)
	complianceRoutes.PUT("/compliance-register/:id", complianceController.UpdateCompliance)
	complianceRoutes.POST("/compliance-register/delete", complianceController.DeleteCompliances)
	complianceRoutes.PUT("/compliance-register/activate", complianceController.ActivateCompliances)
	complianceRoutes.PUT("/compliance-register/deactivate", complianceController.DeactivateCompliances)
}
