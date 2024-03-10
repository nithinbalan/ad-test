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
	complianceRoutes.POST("/compliance-categories", complianceCategoryController.CreateCategory)
	complianceRoutes.PUT("/compliance-categories/:id", complianceCategoryController.UpdateCategory)
	complianceRoutes.POST("/compliance-categories/delete", complianceCategoryController.DeleteComplianceCategories)
	complianceRoutes.POST("/compliance-categories/activate", complianceCategoryController.ActivateComplianceCategories)
	complianceRoutes.POST("/compliance-categories/deactivate", complianceCategoryController.DeactivateComplianceCategories)

	//Compliance Status
	complianceStatusRepo := complianceRepo.NewComplianceStatusRepository(db)
	complianceStatusSvc := complianceService.NewComplianceStatusService(complianceStatusRepo)
	complianceStatusController := controllers.NewComplianceStatusController(complianceStatusSvc)

	complianceRoutes.GET("/compliance-statuses", complianceStatusController.GetAllComplianceStatuses)

	//Sla Category
	slaCategoryRepo := complianceRepo.NewSlaCategoryRepository(db)
	slaCategorySvc := complianceService.NewSlaCategoryService(slaCategoryRepo, activityLogSvc)
	slaCategoryController := controllers.NewSlaCategoryController(slaCategorySvc)

	complianceRoutes.GET("/sla-categories", slaCategoryController.GetAllSlaCategories)
	complianceRoutes.GET("/sla-categories/:id", slaCategoryController.GetSlaCategoryDetails)
	complianceRoutes.POST("/sla-categories", slaCategoryController.CreateCategory)
	complianceRoutes.PUT("/sla-categories/:id", slaCategoryController.UpdateCategory)
	complianceRoutes.POST("/sla-categories/delete", slaCategoryController.DeleteSlaCategories)
	complianceRoutes.POST("/sla-categories/activate", slaCategoryController.ActivateSlaCategories)
	complianceRoutes.POST("/sla-categories/deactivate", slaCategoryController.DeactivateSlaCategories)
}
