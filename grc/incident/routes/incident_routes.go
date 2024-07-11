package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	incidentRoutes := router.Group("/")

	incidentRoutes.Use(middlewares.AuthMiddleware())

	incidentCategoryController := app.IncidentModule.IncidentCategoryController
	incidentSubCategoryController := app.IncidentModule.IncidentSubCategoryController
	incidentTypeController := app.IncidentModule.IncidentTypeController
	incidentStatusController := app.IncidentModule.IncidentStatusController
	incidentRegisterController := app.IncidentModule.IncidentRegisterController
	incidentInvestigationStatusController := app.IncidentModule.IncidentInvestigationStatusController
	incidentInvestigationController := app.IncidentModule.IncidentInvestigationController
	incidentActionPlanController := app.IncidentModule.IncidentActionPlanController
	incidentActionPlanStatusUpdateController := app.IncidentModule.IncidentActionPlanStatusUpdateController
	incidentRootCauseAnalysisController := app.IncidentModule.IncidentRootCauseAnalysisController
	incidentCommunicationLogController := app.IncidentModule.IncidentCommunicationLogController
	incidentClosureController := app.IncidentModule.IncidentClosureController
	incidentPostReviewEffectivenessController := app.IncidentModule.IncidentPostReviewEffectivenessController
	incidentPostReviewController := app.IncidentModule.IncidentPostReviewController
	incidentAssociatedItemController := app.IncidentModule.IncidentAssociatedItemController

	//Incident Category
	incidentRoutes.GET("/incident-categories", incidentCategoryController.GetAllIncidentCategories)
	incidentRoutes.GET("/incident-categories/:id", incidentCategoryController.GetIncidentCategoryDetails)
	incidentRoutes.POST("/incident-categories", incidentCategoryController.CreateIncidentCategory)
	incidentRoutes.PUT("/incident-categories/:id", incidentCategoryController.UpdateIncidentCategory)
	incidentRoutes.POST("/incident-categories/delete", incidentCategoryController.DeleteIncidentCategories)
	incidentRoutes.PUT("/incident-categories/activate", incidentCategoryController.ActivateIncidentCategories)
	incidentRoutes.PUT("/incident-categories/deactivate", incidentCategoryController.DeactivateIncidentCategories)

	//Incident Sub Category
	incidentRoutes.GET("/incident-sub-categories", incidentSubCategoryController.GetAllIncidentSubCategories)
	incidentRoutes.GET("/incident-sub-categories/:id", incidentSubCategoryController.GetIncidentSubCategoryDetails)
	incidentRoutes.POST("/incident-sub-categories", incidentSubCategoryController.CreateIncidentSubCategory)
	incidentRoutes.PUT("/incident-sub-categories/:id", incidentSubCategoryController.UpdateIncidentSubCategory)
	incidentRoutes.POST("/incident-sub-categories/delete", incidentSubCategoryController.DeleteIncidentSubCategories)
	incidentRoutes.PUT("/incident-sub-categories/activate", incidentSubCategoryController.ActivateIncidentSubCategories)
	incidentRoutes.PUT("/incident-sub-categories/deactivate", incidentSubCategoryController.DeactivateIncidentSubCategories)

	//Incident Type
	incidentRoutes.GET("/incident-types", incidentTypeController.GetAllIncidentTypes)
	incidentRoutes.GET("/incident-types/:id", incidentTypeController.GetIncidentTypeDetails)
	incidentRoutes.POST("/incident-types", incidentTypeController.CreateIncidentType)
	incidentRoutes.PUT("/incident-types/:id", incidentTypeController.UpdateIncidentType)
	incidentRoutes.POST("/incident-types/delete", incidentTypeController.DeleteIncidentTypes)
	incidentRoutes.PUT("/incident-types/activate", incidentTypeController.ActivateIncidentTypes)
	incidentRoutes.PUT("/incident-types/deactivate", incidentTypeController.DeactivateIncidentTypes)

	//Incident Status
	incidentRoutes.GET("/incident-statuses", incidentStatusController.GetAllIncidentStatuses)
	incidentRoutes.GET("/incident-statuses/:id", incidentStatusController.GetIncidentStatusDetails)

	//Incident Register
	incidentRoutes.GET("/incident-register", incidentRegisterController.GetAllIncidentRegisters)
	incidentRoutes.GET("/incident-register/:id", incidentRegisterController.GetIncidentRegisterDetails)
	incidentRoutes.POST("/incident-register", incidentRegisterController.CreateIncidentRegister)
	incidentRoutes.PUT("/incident-register/:id", incidentRegisterController.UpdateIncidentRegister)
	incidentRoutes.POST("/incident-register/delete", incidentRegisterController.DeleteIncidentRegisters)

	//Incident Investigation Status
	incidentRoutes.GET("/incident-investigation-statuses", incidentInvestigationStatusController.GetAllIncidentInvestigationStatuses)
	incidentRoutes.GET("/incident-investigation-statuses/:id", incidentInvestigationStatusController.GetIncidentInvestigationStatusDetails)

	//Incident Investigation
	incidentRoutes.GET("/incidents/:incidentID/investigations", incidentInvestigationController.GetAllIncidentInvestigations)
	incidentRoutes.GET("/incidents/:incidentID/investigations/:id", incidentInvestigationController.GetIncidentInvestigationDetails)
	incidentRoutes.POST("/incidents/:incidentID/investigations", incidentInvestigationController.CreateIncidentInvestigation)
	incidentRoutes.PUT("/incidents/:incidentID/investigations/:id", incidentInvestigationController.UpdateIncidentInvestigation)
	incidentRoutes.POST("/incidents/:incidentID/investigations/delete", incidentInvestigationController.DeleteIncidentInvestigations)

	//Incident Action Plan
	incidentRoutes.GET("/incident-action-plans", incidentActionPlanController.GetAllIncidentActionPlans)
	incidentRoutes.GET("/incidents/:incidentID/action-plans", incidentActionPlanController.GetAllIncidentActionPlans)
	incidentRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId", incidentActionPlanController.GetIncidentActionPlanDetails)
	incidentRoutes.POST("/incidents/:incidentID/action-plans", incidentActionPlanController.CreateIncidentActionPlan)
	incidentRoutes.PUT("/incidents/:incidentID/action-plans/:actionPlanId", incidentActionPlanController.UpdateIncidentActionPlan)
	incidentRoutes.POST("/incidents/:incidentID/action-plans/delete", incidentActionPlanController.DeleteIncidentActionPlans)

	//Incident Action Plan Status Updates
	incidentRoutes.GET("/incidents/:incidentID/action-plans/:actionPlanId/status-updates", incidentActionPlanStatusUpdateController.GetAllIncidentActionPlanStatusUpdates)
	incidentRoutes.POST("/incidents/:incidentID/action-plans/:actionPlanId/status-updates", incidentActionPlanStatusUpdateController.CreateIncidentActionPlanStatusUpdate)
	incidentRoutes.PUT("/incidents/:incidentID/action-plans/:actionPlanId/status-updates/:id", incidentActionPlanStatusUpdateController.UpdateIncidentActionPlanStatusUpdate)

	//Incident Root Cause Analyses
	incidentRoutes.GET("/incidents/:incidentID/root-cause-analysis", incidentRootCauseAnalysisController.GetAllIncidentRootCauseAnalyses)
	incidentRoutes.POST("/incidents/:incidentID/root-cause-analysis", incidentRootCauseAnalysisController.CreateIncidentRootCauseAnalysis)
	incidentRoutes.PUT("/incidents/:incidentID/root-cause-analysis/:id", incidentRootCauseAnalysisController.UpdateIncidentRootCauseAnalysis)
	incidentRoutes.DELETE("/incidents/:incidentID/root-cause-analysis/:id", incidentRootCauseAnalysisController.DeleteIncidentRootCauseAnalyses)

	//Incident Communication logs
	incidentRoutes.GET("/incidents/:incidentID/communication-logs", incidentCommunicationLogController.GetAllIncidentCommunicationLogs)
	incidentRoutes.GET("/incidents/:incidentID/communication-logs/:id", incidentCommunicationLogController.GetIncidentCommunicationLogDetails)
	incidentRoutes.POST("/incidents/:incidentID/communication-logs", incidentCommunicationLogController.CreateIncidentCommunicationLog)
	incidentRoutes.PUT("/incidents/:incidentID/communication-logs/:id", incidentCommunicationLogController.UpdateIncidentCommunicationLog)
	incidentRoutes.POST("/incidents/:incidentID/communication-logs/delete", incidentCommunicationLogController.DeleteIncidentCommunicationLogs)

	//Incident Closure
	incidentRoutes.GET("/incidents/:incidentID/closures", incidentClosureController.GetIncidentClosureDetails)
	incidentRoutes.POST("/incidents/:incidentID/closures", incidentClosureController.CreateIncidentClosure)
	incidentRoutes.PUT("/incidents/:incidentID/closures", incidentClosureController.UpdateIncidentClosure)
	incidentRoutes.DELETE("/incidents/:incidentID/closures", incidentClosureController.DeleteIncidentClosures)

	//Incident Post Review Effectiveness
	incidentRoutes.GET("/incident-post-review-effectivenesses", incidentPostReviewEffectivenessController.GetAllIncidentPostReviewEffectivenesses)
	incidentRoutes.GET("/incident-post-review-effectivenesses/:id", incidentPostReviewEffectivenessController.GetIncidentPostReviewEffectivenessDetails)

	//Incident PostReview
	incidentRoutes.GET("/incidents/:incidentID/post-reviews", incidentPostReviewController.GetAllIncidentPostReviews)
	incidentRoutes.GET("/incidents/:incidentID/post-reviews/:id", incidentPostReviewController.GetIncidentPostReviewDetails)
	incidentRoutes.POST("/incidents/:incidentID/post-reviews", incidentPostReviewController.CreateIncidentPostReview)
	incidentRoutes.PUT("/incidents/:incidentID/post-reviews/:id", incidentPostReviewController.UpdateIncidentPostReview)
	incidentRoutes.POST("/incidents/:incidentID/post-reviews/delete", incidentPostReviewController.DeleteIncidentPostReviews)

	//Incident Associated Items
	incidentRoutes.GET("/incidents/:incidentID/controls", incidentAssociatedItemController.GetAllIncidentControls)
	incidentRoutes.GET("/incidents/:incidentID/compliance-requirements", incidentAssociatedItemController.GetAllIncidentComplianceRequirements)
	incidentRoutes.GET("/incidents/:incidentID/risks", incidentAssociatedItemController.GetAllIncidentRisks)
	incidentRoutes.GET("/incidents/:incidentID/processes", incidentAssociatedItemController.GetAllIncidentProcesses)
	incidentRoutes.GET("/incidents/:incidentID/slas", incidentAssociatedItemController.GetAllIncidentSLAs)

}
