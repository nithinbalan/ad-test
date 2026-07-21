package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	awarenessRoutes := router.Group("/")

	awarenessRoutes.Use(middlewares.AuthMiddleware(app.DB))

	db := app.DB

	awarenessCampaignTypeController := app.AwarenessReportingModule.AwarenessCampaignTypeController
	awarenessChannelController := app.AwarenessReportingModule.AwarenessChannelController
	awarenessCampaignStatusController := app.AwarenessReportingModule.AwarenessCampaignStatusController
	awarenessPriorityController := app.AwarenessReportingModule.AwarenessPriorityController
	awarenessEvaluationRatingController := app.AwarenessReportingModule.AwarenessEvaluationRatingController
	awarenessCampaignController := app.AwarenessReportingModule.AwarenessCampaignController
	awarenessMaturityLevelController := app.AwarenessReportingModule.AwarenessMaturityLevelController
	awarenessAssessmentThemeController := app.AwarenessReportingModule.AwarenessAssessmentThemeController
	awarenessPostCampaignSurveyController := app.AwarenessReportingModule.AwarenessPostCampaignSurveyController
	awarenessHandlingStatusController := app.AwarenessReportingModule.AwarenessHandlingStatusController
	awarenessReportSubCategoryController := app.AwarenessReportingModule.AwarenessReportSubCategoryController
	awarenessReportController := app.AwarenessReportingModule.AwarenessReportController

	// awareness campaign types
	awarenessRoutes.GET("/awareness-reporting/awareness-campaign-types", awarenessCampaignTypeController.GetAllAwarenessCampaignTypes)
	awarenessRoutes.GET("/awareness-reporting/awareness-campaign-types/:id", awarenessCampaignTypeController.GetAwarenessCampaignTypeDetails)
	awarenessRoutes.POST("/awareness-reporting/awareness-campaign-types", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessCampaignTypeController.CreateAwarenessCampaignType)
	awarenessRoutes.PUT("/awareness-reporting/awareness-campaign-types/:id", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessCampaignTypeController.UpdateAwarenessCampaignType)
	awarenessRoutes.POST("/awareness-reporting/awareness-campaign-types/delete", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessCampaignTypeController.DeleteAwarenessCampaignTypes)
	awarenessRoutes.PUT("/awareness-reporting/awareness-campaign-types/activate", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessCampaignTypeController.ActivateAwarenessCampaignTypes)
	awarenessRoutes.PUT("/awareness-reporting/awareness-campaign-types/deactivate", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessCampaignTypeController.DeactivateAwarenessCampaignTypes)

	// awareness channels
	awarenessRoutes.GET("/awareness-reporting/awareness-channels", awarenessChannelController.GetAllAwarenessChannels)
	awarenessRoutes.GET("/awareness-reporting/awareness-channels/:id", awarenessChannelController.GetAwarenessChannelDetails)
	awarenessRoutes.POST("/awareness-reporting/awareness-channels", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessChannelController.CreateAwarenessChannel)
	awarenessRoutes.PUT("/awareness-reporting/awareness-channels/:id", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessChannelController.UpdateAwarenessChannel)
	awarenessRoutes.POST("/awareness-reporting/awareness-channels/delete", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessChannelController.DeleteAwarenessChannels)
	awarenessRoutes.PUT("/awareness-reporting/awareness-channels/activate", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessChannelController.ActivateAwarenessChannels)
	awarenessRoutes.PUT("/awareness-reporting/awareness-channels/deactivate", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessChannelController.DeactivateAwarenessChannels)

	// awareness campaign status
	awarenessRoutes.GET("/awareness-reporting/awareness-campaign-statuses", awarenessCampaignStatusController.GetAllAwarenessCampaignStatuses)
	awarenessRoutes.GET("/awareness-reporting/awareness-campaign-statuses/:id", awarenessCampaignStatusController.GetAwarenessCampaignStatusDetails)

	//awareness priorities
	awarenessRoutes.GET("/awareness-reporting/awareness-priorities", awarenessPriorityController.GetAllAwarenessPriorities)
	awarenessRoutes.GET("/awareness-reporting/awareness-priorities/:id", awarenessPriorityController.GetAwarenessPriorityDetails)

	//awareness evaluation
	awarenessRoutes.GET("/awareness-reporting/awareness-evaluation-ratings", awarenessEvaluationRatingController.GetAllAwarenessEvaluationRatings)
	awarenessRoutes.GET("/awareness-reporting/awareness-evaluation-ratings/:id", awarenessEvaluationRatingController.GetAwarenessEvaluationRatingDetails)

	// awareness campaign
	awarenessRoutes.GET("/awareness-reporting/awareness-campaigns", middlewares.RolePermissionMiddleware(db, "awareness_campaign", "view"), awarenessCampaignController.GetAllAwarenessCampaigns)
	awarenessRoutes.GET("/awareness-reporting/awareness-campaigns/:id", middlewares.RolePermissionMiddleware(db, "awareness_campaign", "view"), awarenessCampaignController.GetAwarenessCampaignDetails)
	awarenessRoutes.POST("/awareness-reporting/awareness-campaigns", middlewares.RolePermissionMiddleware(db, "awareness_campaign", "modify"), awarenessCampaignController.CreateAwarenessCampaign)
	awarenessRoutes.PUT("/awareness-reporting/awareness-campaigns/:id", middlewares.RolePermissionMiddleware(db, "awareness_campaign", "modify"), awarenessCampaignController.UpdateAwarenessCampaign)
	awarenessRoutes.POST("/awareness-reporting/awareness-campaigns/delete", middlewares.RolePermissionMiddleware(db, "awareness_campaign", "modify"), awarenessCampaignController.DeleteAwarenessCampaigns)
	awarenessRoutes.GET("/awareness-reporting/awareness-campaigns/export", middlewares.RolePermissionMiddleware(db, "awareness_campaign", "modify"), awarenessCampaignController.GetAwarenessCampaignExport)

	router.GET("/awareness-reporting/awareness-campaigns/:id/document-preview", awarenessCampaignController.GetAwarenessCampaignDocumentPreview)
	router.GET("/awareness-reporting/awareness-campaigns/:id/document-download", awarenessCampaignController.DownloadAwarenessCampaignDocument)

	// awareness maturity levels
	awarenessRoutes.GET("/awareness-reporting/awareness-maturity-levels", awarenessMaturityLevelController.GetAllAwarenessMaturityLevels)
	awarenessRoutes.GET("/awareness-reporting/awareness-maturity-levels/:id", awarenessMaturityLevelController.GetAwarenessMaturityLevelDetails)

	// awareness assessment themes
	awarenessRoutes.GET("/awareness-reporting/awareness-assessment-themes", awarenessAssessmentThemeController.GetAllAwarenessAssessmentThemes)
	awarenessRoutes.GET("/awareness-reporting/awareness-assessment-themes/:id", awarenessAssessmentThemeController.GetAwarenessAssessmentThemeDetails)
	awarenessRoutes.POST("/awareness-reporting/awareness-assessment-themes", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessAssessmentThemeController.CreateAwarenessAssessmentTheme)
	awarenessRoutes.PUT("/awareness-reporting/awareness-assessment-themes/:id", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessAssessmentThemeController.UpdateAwarenessAssessmentTheme)
	awarenessRoutes.POST("/awareness-reporting/awareness-assessment-themes/delete", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessAssessmentThemeController.DeleteAwarenessAssessmentThemes)
	awarenessRoutes.PUT("/awareness-reporting/awareness-assessment-themes/activate", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessAssessmentThemeController.ActivateAwarenessAssessmentThemes)
	awarenessRoutes.PUT("/awareness-reporting/awareness-assessment-themes/deactivate", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessAssessmentThemeController.DeactivateAwarenessAssessmentThemes)

	// awareness post campaign surveys
	awarenessRoutes.GET("/awareness-reporting/awareness-post-campaign-surveys", middlewares.RolePermissionMiddleware(db, "awareness_post_campaign_survey", "view"), awarenessPostCampaignSurveyController.GetAllAwarenessPostCampaignSurveys)
	awarenessRoutes.GET("/awareness-reporting/awareness-post-campaign-surveys/:id", middlewares.RolePermissionMiddleware(db, "awareness_post_campaign_survey", "view"), awarenessPostCampaignSurveyController.GetAwarenessPostCampaignSurveyDetails)
	awarenessRoutes.POST("/awareness-reporting/awareness-post-campaign-surveys", middlewares.RolePermissionMiddleware(db, "awareness_post_campaign_survey", "modify"), awarenessPostCampaignSurveyController.CreateAwarenessPostCampaignSurvey)
	awarenessRoutes.PUT("/awareness-reporting/awareness-post-campaign-surveys/:id", middlewares.RolePermissionMiddleware(db, "awareness_post_campaign_survey", "modify"), awarenessPostCampaignSurveyController.UpdateAwarenessPostCampaignSurvey)
	awarenessRoutes.POST("/awareness-reporting/awareness-post-campaign-surveys/delete", middlewares.RolePermissionMiddleware(db, "awareness_post_campaign_survey", "modify"), awarenessPostCampaignSurveyController.DeleteAwarenessPostCampaignSurveys)
	awarenessRoutes.GET("/awareness-reporting/awareness-post-campaign-surveys/export", middlewares.RolePermissionMiddleware(db, "awareness_post_campaign_survey", "view"), awarenessPostCampaignSurveyController.GetAwarenessPostCampaignSurveysExport)

	// awareness handling statuses
	awarenessRoutes.GET("/awareness-reporting/awareness-handling-statuses", awarenessHandlingStatusController.GetAllAwarenessHandlingStatuses)
	awarenessRoutes.GET("/awareness-reporting/awareness-handling-statuses/:id", awarenessHandlingStatusController.GetAwarenessHandlingStatusDetails)

	// awareness report sub categories
	awarenessRoutes.GET("/awareness-reporting/awareness-report-sub-categories", awarenessReportSubCategoryController.GetAllAwarenessReportSubCategories)
	awarenessRoutes.GET("/awareness-reporting/awareness-report-sub-categories/:id", awarenessReportSubCategoryController.GetAwarenessReportSubCategoryDetails)
	awarenessRoutes.POST("/awareness-reporting/awareness-report-sub-categories", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessReportSubCategoryController.CreateAwarenessReportSubCategory)
	awarenessRoutes.PUT("/awareness-reporting/awareness-report-sub-categories/:id", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessReportSubCategoryController.UpdateAwarenessReportSubCategory)
	awarenessRoutes.POST("/awareness-reporting/awareness-report-sub-categories/delete", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessReportSubCategoryController.DeleteAwarenessReportSubCategories)
	awarenessRoutes.PUT("/awareness-reporting/awareness-report-sub-categories/activate", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessReportSubCategoryController.ActivateAwarenessReportSubCategories)
	awarenessRoutes.PUT("/awareness-reporting/awareness-report-sub-categories/deactivate", middlewares.RolePermissionMiddleware(db, "awareness_reporting_setting", "modify"), awarenessReportSubCategoryController.DeactivateAwarenessReportSubCategories)

	// awareness reports
	awarenessRoutes.GET("/awareness-reporting/awareness-reports", middlewares.RolePermissionMiddleware(db, "awareness_report", "view"), awarenessReportController.GetAllAwarenessReports)
	awarenessRoutes.GET("/awareness-reporting/awareness-reports/:id", middlewares.RolePermissionMiddleware(db, "awareness_report", "view"), awarenessReportController.GetAwarenessReportDetails)
	awarenessRoutes.POST("/awareness-reporting/awareness-reports", middlewares.RolePermissionMiddleware(db, "awareness_report", "modify"), awarenessReportController.CreateAwarenessReport)
	awarenessRoutes.PUT("/awareness-reporting/awareness-reports/:id", middlewares.RolePermissionMiddleware(db, "awareness_report", "modify"), awarenessReportController.UpdateAwarenessReport)
	awarenessRoutes.POST("/awareness-reporting/awareness-reports/delete", middlewares.RolePermissionMiddleware(db, "awareness_report", "modify"), awarenessReportController.DeleteAwarenessReports)
	awarenessRoutes.GET("/awareness-reporting/awareness-reports/export", middlewares.RolePermissionMiddleware(db, "awareness_report", "view"), awarenessReportController.GetAwarenessReportExport)
}
