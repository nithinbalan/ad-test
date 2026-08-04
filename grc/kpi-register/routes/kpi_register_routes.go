package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	keyPerformanceIndicatorRoutes := router.Group("/")

	keyPerformanceIndicatorRoutes.Use(middlewares.AuthMiddleware(app.DB))
	unauthenticatedRoutes := router.Group("/")

	// unauthenticatedRoutes := router.Group("/")

	db := app.DB

	kpiCategoryController := app.KPIRegisterModule.KpiCategoryController
	kpiStatusController := app.KPIRegisterModule.KpiStatusController
	kpiReviewStatusController := app.KPIRegisterModule.KpiReviewStatusController
	entryTypeController := app.KPIRegisterModule.EntryTypeController
	dataSourceController := app.KPIRegisterModule.DataSourceController
	kpiRegisterController := app.KPIRegisterModule.KPIRegisterController
	kpiDashboardController := app.KPIRegisterModule.KPIDashboardController
	kpiValueUpdateController := app.KPIRegisterModule.KPIValueUpdateController
	kpiAchievementLogicController := app.KPIRegisterModule.KPIAchievementLogicController
	kpiAssociatedItemController := app.KPIRegisterModule.KPIAssociatedItemController

	// kpi category
	keyPerformanceIndicatorRoutes.GET("/kpi-categories", kpiCategoryController.GetAllKPICategories)
	keyPerformanceIndicatorRoutes.POST("/kpi-categories", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiCategoryController.CreateKPICategory)
	keyPerformanceIndicatorRoutes.GET("/kpi-categories/:id", kpiCategoryController.GetKPICategoryDetails)
	keyPerformanceIndicatorRoutes.PUT("/kpi-categories/:id", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiCategoryController.UpdateKPICategory)
	keyPerformanceIndicatorRoutes.POST("/kpi-categories/delete", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiCategoryController.DeleteKPICategories)
	keyPerformanceIndicatorRoutes.PUT("/kpi-categories/activate", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiCategoryController.ActivateKPICategories)
	keyPerformanceIndicatorRoutes.PUT("/kpi-categories/deactivate", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiCategoryController.DeactivateKPICategories)

	// kpi status
	keyPerformanceIndicatorRoutes.GET("/kpi-statuses", kpiStatusController.GetAllKPIStatuses)
	keyPerformanceIndicatorRoutes.GET("/kpi-statuses/:id", kpiStatusController.GetKPIStatusDetails)

	keyPerformanceIndicatorRoutes.GET("/kpi-review-statuses", kpiReviewStatusController.GetAllKPIReviewStatuses)
	keyPerformanceIndicatorRoutes.GET("/kpi-review-statuses/:id", kpiReviewStatusController.GetKPIReviewStatusDetails)

	keyPerformanceIndicatorRoutes.GET("/kpi-achievement-logics", kpiAchievementLogicController.GetAllKPIAchievementLogics)
	keyPerformanceIndicatorRoutes.GET("/kpi-achievement-logics/:id", kpiAchievementLogicController.GetKPIAchievementLogicDetails)

	// entry type
	keyPerformanceIndicatorRoutes.GET("/entry-types", entryTypeController.GetAllEntryTypees)
	keyPerformanceIndicatorRoutes.GET("/entry-types/:id", entryTypeController.GetEntryTypeDetails)

	// data source
	keyPerformanceIndicatorRoutes.GET("/data-sources", dataSourceController.GetAllDataSourcees)
	keyPerformanceIndicatorRoutes.GET("/data-sources/:id", dataSourceController.GetDataSourceDetails)

	// kpi register
	keyPerformanceIndicatorRoutes.GET("/kpi-register", kpiRegisterController.GetAllKPIRegister)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/:kpiRegisterId", kpiRegisterController.GetKPIRegisterDetails)
	keyPerformanceIndicatorRoutes.POST("/kpi-register", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.CreateKPIRegister)
	keyPerformanceIndicatorRoutes.PUT("/kpi-register/:kpiRegisterId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.UpdateKPIRegister)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/delete", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.DeleteKPIRegister)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/force-delete", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.ForceDeleteKPIRegister)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/:kpiRegisterId/trends", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiRegisterController.GetKPITrendAnalysis)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/:kpiRegisterId/latest-value-details", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiRegisterController.GetKPILatestValueUpdate)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/category-summary", kpiRegisterController.GetKPICategorySummary)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/status-summary", kpiRegisterController.GetKPIStatusSummary)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/overall-summary", kpiRegisterController.GetKPIOverallSummary)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/review-status-summary", kpiRegisterController.GetKPIByReviewStatus)
	keyPerformanceIndicatorRoutes.PUT("/kpi-register/activate", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.ActivateKPIRegisters)
	keyPerformanceIndicatorRoutes.PUT("/kpi-register/archive", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.ArchiveKPIRegisters)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/responsible-users", kpiRegisterController.GetAllKPIResponsibleUsers)

	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/summary", kpiDashboardController.GetKPISummary)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/summary-by-status", kpiDashboardController.GetKPIStatusSummary)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/pass-fail-summary", kpiDashboardController.GetKPIPassFailSummary)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/top-performing-kpis", kpiDashboardController.GetTopPerformingKPIs)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/bottom-performing-kpis", kpiDashboardController.GetBottomPerformingKPIs)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/kpis-by-category", kpiDashboardController.GetCategorywiseKPIStats)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/kpis-by-department", kpiDashboardController.GetDepartmentwiseKPIStats)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/overdue-for-achievement", kpiDashboardController.GetOverdueForAchievementDetails)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/overdue-for-review", kpiDashboardController.GetOverdueReviewDetails)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/status-overview", kpiDashboardController.GetKPIStatusOverview)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/achievement-summary", kpiDashboardController.GetAchievementSummary)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/kpis-by-review-status", kpiDashboardController.GetKPIByReviewStatus)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/achievement-by-domain", kpiDashboardController.GetKPIAchievementByDomain)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/kpis", kpiDashboardController.GetAllKPIs)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/top-kpis", kpiDashboardController.GetTopKPIs)
	keyPerformanceIndicatorRoutes.GET("/kpi-dashboard/bottom-kpis", kpiDashboardController.GetBottomKPIs)

	keyPerformanceIndicatorRoutes.GET("/kpi-register/template", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiRegisterController.GenerateKPIRegisterTemplate)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/template/validate", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.ValidateKPIRegisterTemplate)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/import/validate", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.ValidateImportKPIRegisters)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/import", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.ImportKPIRegisters)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/export", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.GetKPIExport)

	unauthenticatedRoutes.GET("/kpi-value-updates/:id/document-preview", kpiValueUpdateController.GetKPIValueUpdateDocumentPreview)
	unauthenticatedRoutes.GET("/kpi-value-updates/:id/document-download", kpiValueUpdateController.DownloadKPIValueUpdateDocument)
	keyPerformanceIndicatorRoutes.GET("/kpi-value-updates/export", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiValueUpdateController.GetKPIValueUpdateExport)
	keyPerformanceIndicatorRoutes.GET("/kpi-value-updates", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiValueUpdateController.GetAllKPIValueUpdates)
	keyPerformanceIndicatorRoutes.GET("/kpi-value-updates/:id", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiValueUpdateController.GetKPIValueUpdateDetails)
	keyPerformanceIndicatorRoutes.PUT("/kpi-value-updates/:id", kpiValueUpdateController.UpdateKPIValueUpdate)
	keyPerformanceIndicatorRoutes.POST("/kpi-value-updates/:id/evaluate", kpiValueUpdateController.EvaluateKPIValueUpdateInputs)
	keyPerformanceIndicatorRoutes.PUT("/kpi-value-updates/:id/submit", kpiValueUpdateController.SubmitKPIValueUpdate)
	keyPerformanceIndicatorRoutes.PUT("/kpi-value-updates/:id/approve", kpiValueUpdateController.ApproveKPIValueUpdate)
	keyPerformanceIndicatorRoutes.PUT("/kpi-value-updates/:id/revert", kpiValueUpdateController.RevertKPIValueUpdate)

	// Preserve ADF's KPI associated-item APIs.
	keyPerformanceIndicatorRoutes.GET("/kpi-register/:kpiRegisterId/controls", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIControls)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/:kpiRegisterId/controls", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIControls)
	keyPerformanceIndicatorRoutes.DELETE("/kpi-register/:kpiRegisterId/controls/:controlId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIControl)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/:kpiRegisterId/processes", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIProcesses)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/:kpiRegisterId/processes", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIProcesses)
	keyPerformanceIndicatorRoutes.DELETE("/kpi-register/:kpiRegisterId/processes/:processId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIProcess)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/:kpiRegisterId/compliances", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPICompliances)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/:kpiRegisterId/compliances", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPICompliances)
	keyPerformanceIndicatorRoutes.DELETE("/kpi-register/:kpiRegisterId/compliances/:complianceId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPICompliance)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/:kpiRegisterId/assets", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIAssets)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/:kpiRegisterId/assets", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIAssets)
	keyPerformanceIndicatorRoutes.DELETE("/kpi-register/:kpiRegisterId/assets/:assetId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIAsset)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/:kpiRegisterId/risks", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIRisks)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/:kpiRegisterId/risks", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIRisks)
	keyPerformanceIndicatorRoutes.DELETE("/kpi-register/:kpiRegisterId/risks/:riskId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIRisk)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/:kpiRegisterId/third-party-register", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIThirdParties)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/:kpiRegisterId/third-party-register", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIThirdParties)
	keyPerformanceIndicatorRoutes.DELETE("/kpi-register/:kpiRegisterId/third-party-register/:thirdPartyRegisterId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIThirdParty)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/:kpiRegisterId/documents", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIPolicies)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/:kpiRegisterId/documents", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIPolicies)
	keyPerformanceIndicatorRoutes.DELETE("/kpi-register/:kpiRegisterId/documents/:documentId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIPolicy)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/:kpiRegisterId/clauses", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllClauses)
	keyPerformanceIndicatorRoutes.PUT("/kpi-register/:kpiRegisterId/clauses", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.UpdateClauses)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/:kpiRegisterId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIStrategicObjectives)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/:kpiRegisterId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIStrategicObjective)
	keyPerformanceIndicatorRoutes.DELETE("/kpi-register/:kpiRegisterId/strategic-objectives/:strategicObjectiveId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIStrategicObjective)
	keyPerformanceIndicatorRoutes.GET("/kpi-register/:kpiRegisterId/incidents", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIIncidents)
	keyPerformanceIndicatorRoutes.POST("/kpi-register/:kpiRegisterId/incidents", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIIncident)
	keyPerformanceIndicatorRoutes.DELETE("/kpi-register/:kpiRegisterId/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIIncident)
}
