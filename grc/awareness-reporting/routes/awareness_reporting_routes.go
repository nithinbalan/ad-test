package routes

import (
	"fmt"
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	fmt.Println("Awareness routes loaded")
	awarenessRoutes := router.Group("/")

	awarenessRoutes.Use(middlewares.AuthMiddleware(app.DB))

	db := app.DB

	awarenessCampaignTypeController := app.AwarenessReportingModule.AwarenessCampaignTypeController

	awarenessRoutes.GET("/awareness-reporting/awareness-campaign-types", awarenessCampaignTypeController.GetAllAwarenessCampaignTypes)
	awarenessRoutes.GET("/awareness-reporting/awareness-campaign-types/:id", awarenessCampaignTypeController.GetAwarenessCampaignTypeDetails)
	awarenessRoutes.POST("/awareness-reporting/awareness-campaign-types",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessCampaignTypeController.CreateAwarenessCampaignType,)
	awarenessRoutes.PUT("/awareness-reporting/awareness-campaign-types/:id",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessCampaignTypeController.UpdateAwarenessCampaignType,)
	awarenessRoutes.POST("/awareness-reporting/awareness-campaign-types/delete",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessCampaignTypeController.DeleteAwarenessCampaignTypes,)
	awarenessRoutes.PUT("/awareness-reporting/awareness-campaign-types/activate",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessCampaignTypeController.ActivateAwarenessCampaignTypes,)
	awarenessRoutes.PUT("/awareness-reporting/awareness-campaign-types/deactivate",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessCampaignTypeController.DeactivateAwarenessCampaignTypes,)
}