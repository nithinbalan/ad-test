package routes

import (
	"jethur/grc/compliance/controllers"
	"jethur/grc/compliance/repositories/implementations"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	//Compliance Category
	complianceCategoryRepository := implementations.NewComplianceCategoryRepository(db)
	complianceCategoryController := controllers.NewComplianceCategoryController(complianceCategoryRepository, db)

	complianceRoutes := router.Group("/")

	complianceRoutes.Use(middlewares.AuthMiddleware())

	complianceRoutes.GET("/compliance-categories", complianceCategoryController.GetAllComplianceCategories)
}
