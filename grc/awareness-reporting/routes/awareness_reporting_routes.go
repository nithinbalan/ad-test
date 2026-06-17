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
	awarenessChannelController := app.AwarenessReportingModule.AwarenessChannelController
	awarenessCampaignStatusController := app.AwarenessReportingModule.AwarenessCampaignStatusController

	// awareness campaign types
	awarenessRoutes.GET("/awareness-reporting/awareness-campaign-types", awarenessCampaignTypeController.GetAllAwarenessCampaignTypes)
	awarenessRoutes.GET("/awareness-reporting/awareness-campaign-types/:id", awarenessCampaignTypeController.GetAwarenessCampaignTypeDetails)
	awarenessRoutes.POST("/awareness-reporting/awareness-campaign-types",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessCampaignTypeController.CreateAwarenessCampaignType,)
	awarenessRoutes.PUT("/awareness-reporting/awareness-campaign-types/:id",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessCampaignTypeController.UpdateAwarenessCampaignType,)
	awarenessRoutes.POST("/awareness-reporting/awareness-campaign-types/delete",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessCampaignTypeController.DeleteAwarenessCampaignTypes,)
	awarenessRoutes.PUT("/awareness-reporting/awareness-campaign-types/activate",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessCampaignTypeController.ActivateAwarenessCampaignTypes,)
	awarenessRoutes.PUT("/awareness-reporting/awareness-campaign-types/deactivate",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessCampaignTypeController.DeactivateAwarenessCampaignTypes,)

	// awareness channels
	awarenessRoutes.GET("/awareness-reporting/awareness-channels", awarenessChannelController.GetAllAwarenessChannels)
	awarenessRoutes.GET("/awareness-reporting/awareness-channels/:id", awarenessChannelController.GetAwarenessChannelDetails)
	awarenessRoutes.POST("/awareness-reporting/awareness-channels",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessChannelController.CreateAwarenessChannel,)
	awarenessRoutes.PUT("/awareness-reporting/awareness-channels/:id",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessChannelController.UpdateAwarenessChannel,)
	awarenessRoutes.POST("/awareness-reporting/awareness-channels/delete",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessChannelController.DeleteAwarenessChannels,)
	awarenessRoutes.PUT("/awareness-reporting/awareness-channels/activate",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessChannelController.ActivateAwarenessChannels,)
	awarenessRoutes.PUT("/awareness-reporting/awareness-channels/deactivate",middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"),awarenessChannelController.DeactivateAwarenessChannels,)

	// awareness campaign status
	awarenessRoutes.GET("/awareness-reporting/awareness-campaign-statuses",awarenessCampaignStatusController.GetAllAwarenessCampaignStatuses,)
	awarenessRoutes.GET("/awareness-reporting/awareness-campaign-statuses/:id",awarenessCampaignStatusController.GetAwarenessCampaignStatusDetails,)
}