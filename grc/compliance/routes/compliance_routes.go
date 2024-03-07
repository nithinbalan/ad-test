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

	activityLogRepo := activityLogRepository.NewActivityLogRepository(db)
	activityLogSvc := activityLogServie.NewActivityLogService(activityLogRepo)

	complianceCategoryRepo := complianceRepo.NewComplianceCategoryRepository(db)
	complianceCategorySvc := complianceService.NewComplianceCategoryService(complianceCategoryRepo, activityLogSvc)
	complianceCategoryController := controllers.NewComplianceCategoryController(complianceCategorySvc)

	complianceRoutes := router.Group("/")

	complianceRoutes.Use(middlewares.AuthMiddleware())

	complianceRoutes.GET("/compliance-categories", complianceCategoryController.GetAllComplianceCategories)
	complianceRoutes.GET("/compliance-categories/:id", complianceCategoryController.GetComplianceCategoryDetails)
	complianceRoutes.POST("/compliance-categories", complianceCategoryController.CreateCategory)
	complianceRoutes.PUT("/compliance-categories/:id", complianceCategoryController.UpdateCategory)
	complianceRoutes.POST("/compliance-categories/delete", complianceCategoryController.DeleteComplianceCategories)
	complianceRoutes.POST("/compliance-categories/activate", complianceCategoryController.ActivateComplianceCategories)
	complianceRoutes.POST("/compliance-categories/deactivate", complianceCategoryController.DeactivateComplianceCategories)
}
