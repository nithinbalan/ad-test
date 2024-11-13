package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	assetRoutes := router.Group("/")

	assetRoutes.Use(middlewares.AuthMiddleware())

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	assetTypeController := app.AssetModule.AssetTypeController
	assetCategoryController := app.AssetModule.AssetCategoryController
	assetStatusController := app.AssetModule.AssetStatusController
	assetRegisterController := app.AssetModule.AssetRegisterController
	assetCriticalityMatrixController := app.AssetModule.AssetCriticalityMatrixController
	assetCriticalityMatrixFactorController := app.AssetModule.AssetCriticalityMatrixFactorController
	assetCriticalityMatrixFactorRatingController := app.AssetModule.AssetCriticalityMatrixFactorRatingController
	assetCriticalityMatrixRatingScaleController := app.AssetModule.AssetCriticalityMatrixRatingScaleController
	vulnerabilityTypeController := app.AssetModule.VulnerabilityTypeController
	vulnerabilityStatusController := app.AssetModule.VulnerabilityStatusController
	vulnerabilityController := app.AssetModule.VulnerabilityController
	threatTypeController := app.AssetModule.ThreatTypeController
	threatStatusController := app.AssetModule.ThreatStatusController
	threatController := app.AssetModule.ThreatController
	vulnerabilityActionPlanController := app.AssetModule.VulnerabilityActionPlanController
	vulnerabilityActionPlanStatusUpdateController := app.AssetModule.VulnerabilityActionPlanStatusUpdateController
	threatActionPlanController := app.AssetModule.ThreatActionPlanController
	threatActionPlanStatusUpdateController := app.AssetModule.ThreatActionPlanStatusUpdateController
	assetMaintenanceController := app.AssetModule.AssetMaintenanceController
	assetAssociatedItemController := app.AssetModule.AssetAssociatedItemController
	assetCriticalityController := app.AssetModule.AssetCriticalityController

	// Asset Type
	assetRoutes.GET("/asset-types", assetTypeController.GetAllAssetTypes)
	assetRoutes.GET("/asset-types/:id", assetTypeController.GetAssetTypeDetails)
	assetRoutes.POST("/asset-types", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetTypeController.CreateAssetType)
	assetRoutes.PUT("/asset-types/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetTypeController.UpdateAssetType)
	assetRoutes.POST("/asset-types/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetTypeController.DeleteAssetTypes)
	assetRoutes.PUT("/asset-types/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetTypeController.ActivateAssetTypes)
	assetRoutes.PUT("/asset-types/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetTypeController.DeactivateAssetTypes)

	//Asset category
	assetRoutes.GET("/asset-categories", assetCategoryController.GetAllAssetCategories)
	assetRoutes.GET("/asset-categories/:id", assetCategoryController.GetAssetCategoryDetails)
	assetRoutes.POST("/asset-categories", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetCategoryController.CreateAssetCategory)
	assetRoutes.PUT("/asset-categories/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetCategoryController.UpdateAssetCategory)
	assetRoutes.POST("/asset-categories/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetCategoryController.DeleteAssetCategories)
	assetRoutes.PUT("/asset-categories/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetCategoryController.ActivateAssetCategories)
	assetRoutes.PUT("/asset-categories/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetCategoryController.DeactivateAssetCategories)

	//asset status
	assetRoutes.GET("/asset-statuses", assetStatusController.GetAllAssetStatuses)
	assetRoutes.GET("/asset-statuses/:id", assetStatusController.GetAssetStatusDetails)

	//asset register
	assetRoutes.GET("/asset-register", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetRegisterController.GetAllAssets)
	assetRoutes.GET("/asset-register/:id", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetRegisterController.GetAssetDetails)
	assetRoutes.POST("/asset-register", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetRegisterController.CreateAsset)
	assetRoutes.PUT("/asset-register/:id", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetRegisterController.UpdateAsset)
	assetRoutes.POST("/asset-register/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetRegisterController.DeleteAssets)

	//Asset Matrix
	assetRoutes.GET("/asset-criticality-matrices", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixController.GetAllMatrices)
	assetRoutes.GET("/asset-criticality-matrices/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixController.GetMatrixDetails)
	assetRoutes.POST("/asset-criticality-matrices", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixController.CreateMatrix)
	assetRoutes.PUT("/asset-criticality-matrices/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixController.UpdateMatrix)
	assetRoutes.POST("/asset-criticality-matrices/delete", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixController.DeleteMatrices)

	// Asset Criticality Matrix Factors
	assetRoutes.GET("/asset-criticality-matrix/:matrixId/factors", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixFactorController.GetAllFactors)
	assetRoutes.GET("/asset-criticality-matrix/:matrixId/factors/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixFactorController.GetFactorDetails)
	assetRoutes.POST("/asset-criticality-matrix/:matrixId/factors", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixFactorController.CreateFactor)
	assetRoutes.PUT("/asset-criticality-matrix/:matrixId/factors/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixFactorController.UpdateFactor)
	assetRoutes.POST("/asset-criticality-matrix/:matrixId/factors/delete", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixFactorController.DeleteFactors)

	// Asset Criticality Matrix Factor Ratings
	assetRoutes.GET("/asset-criticality-matrix/:matrixId/factor-ratings", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixFactorRatingController.GetAllRatings)
	assetRoutes.GET("/asset-criticality-matrix/:matrixId/factor-ratings/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixFactorRatingController.GetRatingDetails)
	assetRoutes.POST("/asset-criticality-matrix/:matrixId/factor-ratings", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixFactorRatingController.CreateRating)
	assetRoutes.PUT("/asset-criticality-matrix/:matrixId/factor-ratings/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixFactorRatingController.UpdateRating)
	assetRoutes.POST("/asset-criticality-matrix/:matrixId/factor-ratings/delete", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixFactorRatingController.DeleteRatings)

	// Asset Criticality Matrix Rating Scales
	assetRoutes.GET("/asset-criticality-matrix/:matrixId/rating-scales", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixRatingScaleController.GetAllScales)
	assetRoutes.GET("/asset-criticality-matrix/:matrixId/rating-scales/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixRatingScaleController.GetScaleDetails)
	assetRoutes.POST("/asset-criticality-matrix/:matrixId/rating-scales", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixRatingScaleController.CreateScale)
	assetRoutes.PUT("/asset-criticality-matrix/:matrixId/rating-scales/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixRatingScaleController.UpdateScale)
	assetRoutes.POST("/asset-criticality-matrix/:matrixId/rating-scales/delete", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixRatingScaleController.DeleteScales)

	//vulnerability types
	assetRoutes.GET("/vulnerability-types", vulnerabilityTypeController.GetAllVulnerabilityTypes)
	assetRoutes.GET("/vulnerability-types/:id", vulnerabilityTypeController.GetVulnerabilityTypeDetails)
	assetRoutes.POST("/vulnerability-types", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), vulnerabilityTypeController.CreateVulnerabilityType)
	assetRoutes.PUT("/vulnerability-types/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), vulnerabilityTypeController.UpdateVulnerabilityType)
	assetRoutes.POST("/vulnerability-types/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), vulnerabilityTypeController.DeleteVulnerabilityTypes)
	assetRoutes.PUT("/vulnerability-types/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), vulnerabilityTypeController.ActivateVulnerabilityTypes)
	assetRoutes.PUT("/vulnerability-types/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), vulnerabilityTypeController.DeactivateVulnerabilityTypes)

	//vulnerability status
	assetRoutes.GET("/vulnerability-statuses", vulnerabilityStatusController.GetAllVulnerabilityStatuses)
	assetRoutes.GET("/vulnerability-statuses/:id", vulnerabilityStatusController.GetVulnerabilityStatusDetails)

	assetRoutes.GET("/vulnerabilities", middlewares.RolePermissionMiddleware(db, "assets", "view"), vulnerabilityController.GetAllVulnerabilities)
	assetRoutes.GET("/vulnerabilities/:id", middlewares.RolePermissionMiddleware(db, "assets", "view"), vulnerabilityController.GetVulnerabilityDetails)
	assetRoutes.POST("/vulnerabilities", middlewares.RolePermissionMiddleware(db, "assets", "modify"), vulnerabilityController.CreateVulnerability)
	assetRoutes.PUT("/vulnerabilities/:id", middlewares.RolePermissionMiddleware(db, "assets", "modify"), vulnerabilityController.UpdateVulnerability)
	assetRoutes.POST("/vulnerabilities/delete", middlewares.RolePermissionMiddleware(db, "assets", "modify"), vulnerabilityController.DeleteVulnerabilities)
	assetRoutes.GET("assets/:id/vulnerability-summary", vulnerabilityController.GetVulnerabilitySummary)
	//Threat types
	assetRoutes.GET("/threat-types", threatTypeController.GetAllThreatTypes)
	assetRoutes.GET("/threat-types/:id", threatTypeController.GetThreatTypeDetails)
	assetRoutes.POST("/threat-types", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), threatTypeController.CreateThreatType)
	assetRoutes.PUT("/threat-types/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), threatTypeController.UpdateThreatType)
	assetRoutes.POST("/threat-types/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), threatTypeController.DeleteThreatTypes)
	assetRoutes.PUT("/threat-types/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), threatTypeController.ActivateThreatTypes)
	assetRoutes.PUT("/threat-types/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), threatTypeController.DeactivateThreatTypes)

	//threat status
	assetRoutes.GET("/threat-statuses", threatStatusController.GetAllThreatStatuses)
	assetRoutes.GET("/threat-statuses/:id", threatStatusController.GetThreatStatusDetails)

	//Threat
	assetRoutes.GET("/threats", threatController.GetAllThreats)
	assetRoutes.GET("/threats/:id", threatController.GetThreatDetails)
	assetRoutes.POST("/threats", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), threatController.CreateThreat)
	assetRoutes.PUT("/threats/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), threatController.UpdateThreat)
	assetRoutes.POST("/threats/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), threatController.DeleteThreats)
	assetRoutes.GET("assets/:id/threat-summary", threatController.GetThreatSummary)

	//vulnerability action plan
	assetRoutes.GET("/vulnerability-action-plans/export", middlewares.RolePermissionMiddleware(db, "vulnerability_mitigation", "view"), vulnerabilityActionPlanController.GetVulnerabilityActionPlanExport)
	assetRoutes.GET("/vulnerability-action-plans", middlewares.RolePermissionMiddleware(db, "vulnerability_mitigation", "view"), vulnerabilityActionPlanController.GetAllVulnerabilityActionPlans)
	assetRoutes.GET("/vulnerabilities/:id/action-plans", middlewares.RolePermissionMiddleware(db, "vulnerability_mitigation", "view"), vulnerabilityActionPlanController.GetAllVulnerabilityActionPlans)
	assetRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "vulnerability_mitigation", "view"), vulnerabilityActionPlanController.GetVulnerabilityActionPlanDetails)
	assetRoutes.POST("/vulnerabilities/:id/action-plans", middlewares.RolePermissionMiddleware(db, "vulnerability_mitigation", "modify"), vulnerabilityActionPlanController.CreateVulnerabilityActionPlan)
	assetRoutes.PUT("/vulnerabilities/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "vulnerability_mitigation", "modify"), vulnerabilityActionPlanController.UpdateVulnerabilityActionPlan)
	assetRoutes.POST("/vulnerabilities/:id/action-plans/delete", middlewares.RolePermissionMiddleware(db, "vulnerability_mitigation", "modify"), vulnerabilityActionPlanController.DeleteVulnerabilityActionPlans)

	assetRoutes.GET("/vulnerabilities/:id/action-plan-summary", middlewares.RolePermissionMiddleware(db, "vulnerability_mitigation", "view"), vulnerabilityActionPlanController.GetVulnerabilityActionPlanSummary)

	// Document Preview and Download for Vulnerability Action Plan Status Update
	unauthenticatedRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-preview", vulnerabilityActionPlanStatusUpdateController.GetDocumentPreview)
	unauthenticatedRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-download", vulnerabilityActionPlanStatusUpdateController.DownloadVulnerabilityActionPlanDocument)

	assetRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "vulnerability_management", "view"), vulnerabilityActionPlanStatusUpdateController.GetAllVulnerabilityActionPlanStatusUpdates)
	assetRoutes.POST("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "vulnerability_management", "modify"), vulnerabilityActionPlanStatusUpdateController.CreateVulnerabilityActionPlanStatusUpdate)
	assetRoutes.PUT("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates/:actionPlanStatusUpdateId", middlewares.RolePermissionMiddleware(db, "vulnerability_management", "modify"), vulnerabilityActionPlanStatusUpdateController.UpdateVulnerabilityActionPlanStatusUpdate)
	assetRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "vulnerability_management", "view"), vulnerabilityActionPlanStatusUpdateController.GetLatestVulnerabilityActionPlanStatusUpdate)

	// Threat action plan
	assetRoutes.GET("/threat-action-plans/export", middlewares.RolePermissionMiddleware(db, "threat_mitigation", "view"), threatActionPlanController.GetThreatActionPlanExport)
	assetRoutes.GET("/threat-action-plans", middlewares.RolePermissionMiddleware(db, "threat_mitigation", "view"), threatActionPlanController.GetAllThreatActionPlans)
	assetRoutes.GET("/threats/:id/action-plans", middlewares.RolePermissionMiddleware(db, "threat_mitigation", "view"), threatActionPlanController.GetAllThreatActionPlans)
	assetRoutes.GET("/threats/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "threat_mitigation", "view"), threatActionPlanController.GetThreatActionPlanDetails)
	assetRoutes.POST("/threats/:id/action-plans", middlewares.RolePermissionMiddleware(db, "threat_mitigation", "modify"), threatActionPlanController.CreateThreatActionPlan)
	assetRoutes.PUT("/threats/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "threat_mitigation", "modify"), threatActionPlanController.UpdateThreatActionPlan)
	assetRoutes.POST("/threats/:id/action-plans/delete", middlewares.RolePermissionMiddleware(db, "threat_mitigation", "modify"), threatActionPlanController.DeleteThreatActionPlans)

	assetRoutes.GET("/threats/:id/action-plan-summary", middlewares.RolePermissionMiddleware(db, "threat_mitigation", "view"), threatActionPlanController.GetThreatActionPlanSummary)

	// Document Preview and Download for Threat Action Plan Status Update
	unauthenticatedRoutes.GET("/threats/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-preview", threatActionPlanStatusUpdateController.GetDocumentPreview)
	unauthenticatedRoutes.GET("/threats/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-download", threatActionPlanStatusUpdateController.DownloadThreatActionPlanDocument)

	assetRoutes.GET("/threats/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "threat_management", "view"), threatActionPlanStatusUpdateController.GetAllThreatActionPlanStatusUpdates)
	assetRoutes.POST("/threats/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "threat_management", "modify"), threatActionPlanStatusUpdateController.CreateThreatActionPlanStatusUpdate)
	assetRoutes.PUT("/threats/:id/action-plans/:actionPlanId/status-updates/:actionPlanStatusUpdateId", middlewares.RolePermissionMiddleware(db, "threat_management", "modify"), threatActionPlanStatusUpdateController.UpdateThreatActionPlanStatusUpdate)
	assetRoutes.GET("/threats/:id/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "threat_management", "view"), threatActionPlanStatusUpdateController.GetLatestThreatActionPlanStatusUpdate)

	assetRoutes.GET("/assets/:id/maintenance", middlewares.RolePermissionMiddleware(db, "asset_maintenance", "view"), assetMaintenanceController.GetAllMaintenances)
	assetRoutes.GET("/assets/:id/maintenance/:maintenanceId", middlewares.RolePermissionMiddleware(db, "asset_maintenance", "view"), assetMaintenanceController.GetMaintenanceDetails)
	assetRoutes.POST("/assets/:id/maintenance", middlewares.RolePermissionMiddleware(db, "asset_maintenance", "modify"), assetMaintenanceController.CreateMaintenance)
	assetRoutes.PUT("/assets/:id/maintenance/:maintenanceId", middlewares.RolePermissionMiddleware(db, "asset_maintenance", "modify"), assetMaintenanceController.UpdateMaintenance)
	assetRoutes.POST("/assets/:id/maintenance/delete", middlewares.RolePermissionMiddleware(db, "asset_maintenance", "modify"), assetMaintenanceController.DeleteMaintenances)

	//asset associated item list
	assetRoutes.GET("assets/:id/controls", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetControls)
	assetRoutes.GET("assets/:id/incidents", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetIncidents)
	assetRoutes.GET("assets/:id/risks", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetRisks)

	//asset criticality
	assetRoutes.GET("assets/:id/criticality", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetCriticalityController.GetAssetCriticality)
	assetRoutes.POST("assets/:id/criticality", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetCriticalityController.CreateAssetCriticality)
	assetRoutes.PUT("assets/:id/criticality/:criticalityId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetCriticalityController.UpdateAssetCriticality)
	assetRoutes.DELETE("assets/:id/criticality", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetCriticalityController.DeleteAssetCriticality)
}
