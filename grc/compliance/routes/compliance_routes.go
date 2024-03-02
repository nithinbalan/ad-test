package routes

import (
	"jethur/grc/compliance/controllers"
	"jethur/grc/compliance/repositories"
	"jethur/grc/compliance/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	//Compliance Category
	complianceCategoryRepository := repositories.NewComplianceCategoryRepository(db)
	complianceCategoryService := services.NewComplianceCategoryService(complianceCategoryRepository)
	complianceCategoryController := controllers.NewComplianceCategoryController(complianceCategoryService)

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
