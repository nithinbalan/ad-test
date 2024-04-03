package routes

import (
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/risk/controllers"
	riskRepositories "jethur/grc/risk/repositories"
	riskServices "jethur/grc/risk/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	riskRoutes := router.Group("/")

	riskRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	riskCategoryRepo := riskRepositories.NewRiskCategoryRepository(db)
	riskCategorySvc := riskServices.NewRiskCategoryService(riskCategoryRepo, activityLogSvc)
	riskCategoryController := controllers.NewRiskCategoryController(riskCategorySvc)

	riskRoutes.GET("/risk-categories", riskCategoryController.GetAllRiskCategories)
	riskRoutes.GET("/risk-categories/:id", riskCategoryController.GetRiskCategoryDetails)
	riskRoutes.POST("/risk-categories", riskCategoryController.CreateRiskCategory)
	riskRoutes.PUT("/risk-categories/:id", riskCategoryController.UpdateRiskCategory)
	riskRoutes.POST("/risk-categories/delete", riskCategoryController.DeleteRiskCategories)
	riskRoutes.PUT("/risk-categories/activate", riskCategoryController.ActivateRiskCategories)
	riskRoutes.PUT("/risk-categories/deactivate", riskCategoryController.DeactivateRiskCategories)

}
