package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	kpiRegisterRoutes := router.Group("/")

	kpiRegisterRoutes.Use(middlewares.AuthMiddleware(app.DB))

	unauthenticatedRoutes := router.Group("/")
	
	db := app.DB

	kpiRegisterController := app.KPIRegisterModule.KPIRegisterController
	kpiValueUpdateController := app.KPIRegisterModule.KPIValueUpdateController
	kpiAssociatedItemController := app.KPIRegisterModule.KPIAssociatedItemController

	kpiRegisterRoutes.GET("/kpi-register", kpiRegisterController.GetAllKPIRegister)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId", kpiRegisterController.GetKPIRegisterDetails)
	kpiRegisterRoutes.POST("/kpi-register", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.CreateKPIRegister)
	kpiRegisterRoutes.PUT("/kpi-register/:kpiRegisterId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.UpdateKPIRegister)
	kpiRegisterRoutes.POST("/kpi-register/delete", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiRegisterController.DeleteKPIRegister)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/trends", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiRegisterController.GetKPITrendAnalysis)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/latest-value-details", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiRegisterController.GetKPILatestValueUpdate)

	unauthenticatedRoutes.GET("/kpi-register/:kpiRegisterId/value-updates/:updateId/document-preview", kpiValueUpdateController.GetKPIValueUpdateDocumentPreview)
	unauthenticatedRoutes.GET("/kpi-register/:kpiRegisterId/value-updates/:updateId/document-download", kpiValueUpdateController.DownloadKPIValueUpdateDocument)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/value-updates", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiValueUpdateController.GetAllKPIValueUpdates)
	kpiRegisterRoutes.POST("/kpi-register/:kpiRegisterId/value-updates", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiValueUpdateController.CreateKPIValueUpdate)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/value-updates/:updateId", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiValueUpdateController.GetKPIValueUpdateDetails)
	kpiRegisterRoutes.PUT("/kpi-register/:kpiRegisterId/value-updates/:updateId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiValueUpdateController.UpdateKPIValueUpdate)
	kpiRegisterRoutes.POST("/kpi-register/:kpiRegisterId/value-updates/delete", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiValueUpdateController.DeleteKPIValueUpdates)

	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/controls", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIControls)
	kpiRegisterRoutes.POST("/kpi-register/:kpiRegisterId/controls", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIControls)
	kpiRegisterRoutes.DELETE("/kpi-register/:kpiRegisterId/controls/:controlId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIControl)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/processes", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIProcesses)
	kpiRegisterRoutes.POST("/kpi-register/:kpiRegisterId/processes", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIProcesses)
	kpiRegisterRoutes.DELETE("/kpi-register/:kpiRegisterId/processes/:processId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIProcess)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/compliances", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPICompliances)
	kpiRegisterRoutes.POST("/kpi-register/:kpiRegisterId/compliances", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPICompliances)
	kpiRegisterRoutes.DELETE("/kpi-register/:kpiRegisterId/compliances/:complianceId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPICompliance)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/assets", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIAssets)
	kpiRegisterRoutes.POST("/kpi-register/:kpiRegisterId/assets", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIAssets)
	kpiRegisterRoutes.DELETE("/kpi-register/:kpiRegisterId/assets/:assetId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIAsset)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/risks", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIRisks)
	kpiRegisterRoutes.POST("/kpi-register/:kpiRegisterId/risks", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIRisks)
	kpiRegisterRoutes.DELETE("/kpi-register/:kpiRegisterId/risks/:riskId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIRisk)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/third-party-register", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIThirdParties)
	kpiRegisterRoutes.POST("/kpi-register/:kpiRegisterId/third-party-register", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIThirdParties)
	kpiRegisterRoutes.DELETE("/kpi-register/:kpiRegisterId/third-party-register/:thirdPartyRegisterId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIThirdParty)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/documents", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIPolicies)
	kpiRegisterRoutes.POST("/kpi-register/:kpiRegisterId/documents", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIPolicies)
	kpiRegisterRoutes.DELETE("/kpi-register/:kpiRegisterId/documents/:documentId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIPolicy)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/clauses", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllClauses)
	kpiRegisterRoutes.PUT("/kpi-register/:kpiRegisterId/clauses", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.UpdateClauses)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIStrategicObjectives)
	kpiRegisterRoutes.POST("/kpi-register/:kpiRegisterId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIStrategicObjective)
	kpiRegisterRoutes.DELETE("/kpi-register/:kpiRegisterId/strategic-objectives/:strategicObjectiveId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIStrategicObjective)
	kpiRegisterRoutes.GET("/kpi-register/:kpiRegisterId/incidents", middlewares.RolePermissionMiddleware(db, "kpi", "view"), kpiAssociatedItemController.GetAllKPIIncidents)
	kpiRegisterRoutes.POST("/kpi-register/:kpiRegisterId/incidents", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.CreateKPIIncident)
	kpiRegisterRoutes.DELETE("/kpi-register/:kpiRegisterId/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "kpi", "modify"), kpiAssociatedItemController.DeleteKPIIncident)
}
