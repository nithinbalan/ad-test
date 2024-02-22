package routes

import (
	"jethur/grc/compliance/controllers"
	"jethur/grc/compliance/repositories/implementations"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	//Compliance Category
	complianceCategoryRepository := implementations.NewComplianceCategoryRepository()
	complianceCategoryController := controllers.NewComplianceCategoryController(complianceCategoryRepository)
	router.GET("/compliance-categories", complianceCategoryController.GetAllComplianceCategories)
}
