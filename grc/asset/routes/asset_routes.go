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
	assetSubCategoryController := app.AssetModule.AssetSubCategoryController
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
	assetPolicyController := app.AssetModule.AssetPolicyController
	assetPolicyChecklistController := app.AssetModule.AssetPolicyChecklistController
	assetPolicyCommentController := app.AssetModule.AssetPolicyCommentController
	assetDocumentCategoryController := app.AssetModule.AssetDocumentCategoryController
	assetOtherDocumentController := app.AssetModule.AssetOtherDocumentController
	assetControlAnalysisController := app.AssetModule.AssetControlAnalysisController
	assetKriController := app.AssetModule.AssetKriController

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

	//Asset Sub category
	assetRoutes.GET("/asset-sub-categories", assetSubCategoryController.GetAllAssetSubCategories)
	assetRoutes.GET("/asset-sub-categories/:id", assetSubCategoryController.GetAssetSubCategoryDetails)
	assetRoutes.POST("/asset-sub-categories", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetSubCategoryController.CreateAssetSubCategory)
	assetRoutes.PUT("/asset-sub-categories/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetSubCategoryController.UpdateAssetSubCategory)
	assetRoutes.POST("/asset-sub-categories/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetSubCategoryController.DeleteAssetSubCategories)
	assetRoutes.PUT("/asset-sub-categories/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetSubCategoryController.ActivateAssetSubCategories)
	assetRoutes.PUT("/asset-sub-categories/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetSubCategoryController.DeactivateAssetSubCategories)

	//asset status
	assetRoutes.GET("/asset-statuses", assetStatusController.GetAllAssetStatuses)
	assetRoutes.GET("/asset-statuses/:id", assetStatusController.GetAssetStatusDetails)

	//asset register
	assetRoutes.GET("/asset-register", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetRegisterController.GetAllAssets)
	assetRoutes.GET("/asset-register/:id", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetRegisterController.GetAssetDetails)
	assetRoutes.POST("/asset-register", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetRegisterController.CreateAsset)
	assetRoutes.PUT("/asset-register/:id", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetRegisterController.UpdateAsset)
	assetRoutes.POST("/asset-register/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetRegisterController.DeleteAssets)
	assetRoutes.GET("/asset-register/:id/control-summary", middlewares.RolePermissionMiddleware(db, "process", "view"), assetRegisterController.GetControlSummary)
	unauthenticatedRoutes.GET("/asset-register/:id/document-download", assetRegisterController.DownloadAssetDocument)
	assetRoutes.GET("/asset-register/export", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetRegisterController.GetAssetsExport)

	//Asset Matrix
	assetRoutes.GET("/asset-criticality-matrices", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixController.GetAllMatrices)
	assetRoutes.GET("/asset-criticality-matrices/:matrixId", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixController.GetMatrixDetails)
	assetRoutes.POST("/asset-criticality-matrices", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixController.CreateMatrix)
	assetRoutes.PUT("/asset-criticality-matrices/:matrixId", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixController.UpdateMatrix)
	assetRoutes.POST("/asset-criticality-matrices/delete", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixController.DeleteMatrices)

	// Asset Criticality Matrix Factors
	assetRoutes.GET("/asset-criticality-matrices/:matrixId/factors", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixFactorController.GetAllFactors)
	assetRoutes.GET("/asset-criticality-matrices/:matrixId/factors/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixFactorController.GetFactorDetails)
	assetRoutes.POST("/asset-criticality-matrices/:matrixId/factors", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixFactorController.CreateFactor)
	assetRoutes.PUT("/asset-criticality-matrices/:matrixId/factors/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixFactorController.UpdateFactor)
	assetRoutes.POST("/asset-criticality-matrices/:matrixId/factors/delete", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixFactorController.DeleteFactors)

	// Asset Criticality Matrix Factor Ratings
	assetRoutes.GET("/asset-criticality-matrices/:matrixId/factor-ratings", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixFactorRatingController.GetAllRatings)
	assetRoutes.GET("/asset-criticality-matrices/:matrixId/factor-ratings/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixFactorRatingController.GetRatingDetails)
	assetRoutes.POST("/asset-criticality-matrices/:matrixId/factor-ratings", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixFactorRatingController.CreateRating)
	assetRoutes.PUT("/asset-criticality-matrices/:matrixId/factor-ratings/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixFactorRatingController.UpdateRating)
	assetRoutes.POST("/asset-criticality-matrices/:matrixId/factor-ratings/delete", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixFactorRatingController.DeleteRatings)

	// Asset Criticality Matrix Rating Scales
	assetRoutes.GET("/asset-criticality-matrices/:matrixId/rating-scales", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixRatingScaleController.GetAllScales)
	assetRoutes.GET("/asset-criticality-matrices/:matrixId/rating-scales/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "view"), assetCriticalityMatrixRatingScaleController.GetScaleDetails)
	assetRoutes.POST("/asset-criticality-matrices/:matrixId/rating-scales", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixRatingScaleController.CreateScale)
	assetRoutes.PUT("/asset-criticality-matrices/:matrixId/rating-scales/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixRatingScaleController.UpdateScale)
	assetRoutes.POST("/asset-criticality-matrices/:matrixId/rating-scales/delete", middlewares.RolePermissionMiddleware(db, "asset_criticality_matrix", "modify"), assetCriticalityMatrixRatingScaleController.DeleteScales)

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

	assetRoutes.GET("/vulnerabilities", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityController.GetAllVulnerabilities)
	assetRoutes.GET("/vulnerabilities/:id", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityController.GetVulnerabilityDetails)
	assetRoutes.POST("/vulnerabilities", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityController.CreateVulnerability)
	assetRoutes.PUT("/vulnerabilities/:id", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityController.UpdateVulnerability)
	assetRoutes.POST("/vulnerabilities/delete", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityController.DeleteVulnerabilities)
	assetRoutes.GET("/assets/:id/vulnerability-summary", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityController.GetVulnerabilitySummary)
	assetRoutes.GET("/vulnerabilities/export", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityController.GetVulnerabilitiesExport)
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
	assetRoutes.POST("/threats", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatController.CreateThreat)
	assetRoutes.PUT("/threats/:id", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatController.UpdateThreat)
	assetRoutes.POST("/threats/delete", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatController.DeleteThreats)
	assetRoutes.GET("/assets/:id/threat-summary", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatController.GetThreatSummary)
	assetRoutes.GET("/threats/export", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatController.GetThreatsExport)

	//vulnerability action plan
	assetRoutes.GET("/vulnerability-action-plans/export", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanController.GetVulnerabilityActionPlanExport)
	assetRoutes.GET("/vulnerability-action-plans", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanController.GetAllVulnerabilityActionPlans)
	assetRoutes.GET("/vulnerabilities/:id/action-plans", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanController.GetAllVulnerabilityActionPlans)
	assetRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanController.GetVulnerabilityActionPlanDetails)
	assetRoutes.POST("/vulnerabilities/:id/action-plans", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityActionPlanController.CreateVulnerabilityActionPlan)
	assetRoutes.PUT("/vulnerabilities/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityActionPlanController.UpdateVulnerabilityActionPlan)
	assetRoutes.POST("/vulnerabilities/:id/action-plans/delete", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityActionPlanController.DeleteVulnerabilityActionPlans)

	assetRoutes.GET("/vulnerabilities/:id/action-plan-summary", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanController.GetVulnerabilityActionPlanSummary)

	// Document Preview and Download for Vulnerability Action Plan Status Update
	unauthenticatedRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-preview", vulnerabilityActionPlanStatusUpdateController.GetDocumentPreview)
	unauthenticatedRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-download", vulnerabilityActionPlanStatusUpdateController.DownloadVulnerabilityActionPlanDocument)

	assetRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanStatusUpdateController.GetAllVulnerabilityActionPlanStatusUpdates)
	assetRoutes.POST("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityActionPlanStatusUpdateController.CreateVulnerabilityActionPlanStatusUpdate)
	assetRoutes.PUT("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates/:actionPlanStatusUpdateId", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityActionPlanStatusUpdateController.UpdateVulnerabilityActionPlanStatusUpdate)
	assetRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanStatusUpdateController.GetLatestVulnerabilityActionPlanStatusUpdate)

	// Threat action plan
	assetRoutes.GET("/threat-action-plans/export", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanController.GetThreatActionPlanExport)
	assetRoutes.GET("/threat-action-plans", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanController.GetAllThreatActionPlans)
	assetRoutes.GET("/threats/:id/action-plans", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanController.GetAllThreatActionPlans)
	assetRoutes.GET("/threats/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanController.GetThreatActionPlanDetails)
	assetRoutes.POST("/threats/:id/action-plans", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatActionPlanController.CreateThreatActionPlan)
	assetRoutes.PUT("/threats/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatActionPlanController.UpdateThreatActionPlan)
	assetRoutes.POST("/threats/:id/action-plans/delete", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatActionPlanController.DeleteThreatActionPlans)

	assetRoutes.GET("/threats/:id/action-plan-summary", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanController.GetThreatActionPlanSummary)

	// Document Preview and Download for Threat Action Plan Status Update
	unauthenticatedRoutes.GET("/threats/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-preview", threatActionPlanStatusUpdateController.GetDocumentPreview)
	unauthenticatedRoutes.GET("/threats/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-download", threatActionPlanStatusUpdateController.DownloadThreatActionPlanDocument)

	assetRoutes.GET("/threats/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanStatusUpdateController.GetAllThreatActionPlanStatusUpdates)
	assetRoutes.POST("/threats/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatActionPlanStatusUpdateController.CreateThreatActionPlanStatusUpdate)
	assetRoutes.PUT("/threats/:id/action-plans/:actionPlanId/status-updates/:actionPlanStatusUpdateId", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatActionPlanStatusUpdateController.UpdateThreatActionPlanStatusUpdate)
	assetRoutes.GET("/threats/:id/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanStatusUpdateController.GetLatestThreatActionPlanStatusUpdate)

	assetRoutes.GET("/assets/:id/maintenance", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetMaintenanceController.GetAllMaintenances)
	assetRoutes.GET("/assets/:id/maintenance/:maintenanceId", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetMaintenanceController.GetMaintenanceDetails)
	assetRoutes.POST("/assets/:id/maintenance", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetMaintenanceController.CreateMaintenance)
	assetRoutes.PUT("/assets/:id/maintenance/:maintenanceId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetMaintenanceController.UpdateMaintenance)
	assetRoutes.POST("/assets/:id/maintenance/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetMaintenanceController.DeleteMaintenances)

	//asset associated item list
	assetRoutes.GET("assets/:id/controls", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetControls)
	assetRoutes.GET("assets/:id/incidents", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetIncidents)
	assetRoutes.GET("assets/:id/risks", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetRisks)
	assetRoutes.GET("assets/:id/vulnerabilities", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetVulnerabilities)
	assetRoutes.GET("assets/:id/clauses", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetClauses)

	assetRoutes.POST("assets/:id/vulnerabilities", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.CreateAssetVulnerabilities)
	assetRoutes.POST("assets/:id/vulnerabilities/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.DeleteAssetVulnerabilities)

	assetRoutes.POST("assets/:id/compliance-clauses", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.CreateAssetComplianceClauses)
	assetRoutes.POST("assets/:id/compliance-clauses/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.DeleteAssetComplianceClauses)

	assetRoutes.POST("vulnerabilities/:id/threats", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.CreateAssetThreats)
	assetRoutes.POST("vulnerabilities/:id/threats/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.DeleteAssetThreats)

	assetRoutes.POST("threats/:id/risks", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.CreateAssetThreatRisks)
	assetRoutes.POST("threats/:id/risks/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.DeleteAssetThreatRisks)

	assetRoutes.GET("assets/:id/risk-analysis", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetRiskAnalysis)

	//asset criticality
	assetRoutes.POST("/asset/criticality-recommendation/ai", assetCriticalityController.GetAIRecommendation)
	assetRoutes.GET("assets/:id/criticality", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetCriticalityController.GetAssetCriticality)
	assetRoutes.POST("assets/:id/criticality", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetCriticalityController.CreateAssetCriticality)
	assetRoutes.PUT("assets/:id/criticality/:criticalityId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetCriticalityController.UpdateAssetCriticality)
	assetRoutes.DELETE("assets/:id/criticality", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetCriticalityController.DeleteAssetCriticality)

	assetRoutes.GET("/assets/:id/policies", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetPolicyController.GetAllPolicies)
	assetRoutes.GET("/assets/:id/policies/:policyId", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetPolicyController.GetPolicyDetails)
	assetRoutes.POST("/assets/:id/policies", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetPolicyController.CreatePolicy)
	assetRoutes.POST("/assets/:id/policies/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetPolicyController.DeletePolicies)
	assetRoutes.GET("/assets/:id/policies-summary", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetPolicyController.GetAssetPolicyComplianceSummary)

	assetRoutes.POST("/assets/:id/policies/:policyId/checklists", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetPolicyChecklistController.CreateAssetPolicyChecklist)
	assetRoutes.PUT("/assets/:id/policies/:policyId/checklists/:checklistId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetPolicyChecklistController.UpdateAssetPolicyChecklist)
	assetRoutes.GET("/assets/:id/policies/:policyId/checklists/:checklistId", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetPolicyChecklistController.GetAssetPolicyChecklistDetails)
	assetRoutes.POST("/assets/:id/policies/:policyId/checklists/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetPolicyChecklistController.DeleteAssetPolicyChecklists)
	assetRoutes.PUT("/assets/:id/policies/:policyId/checklists/:checklistId/answer", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetPolicyChecklistController.AnswerAssetPolicyChecklist)
	assetRoutes.POST("/assets/:id/policies/:policyId/checklists/delete-answer", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetPolicyChecklistController.DeleteAnswerAssetPolicyChecklist)
	assetRoutes.GET("/assets/:id/policies/:policyId/checklists-summary", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetPolicyChecklistController.GetAssetPolicyChecklistSummary)

	assetRoutes.GET("/assets/:id/policies/:policyId/comments", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetPolicyCommentController.GetAllComments)
	assetRoutes.GET("/assets/:id/policies/:policyId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetPolicyCommentController.GetCommentDetails)
	assetRoutes.POST("/assets/:id/policies/:policyId/comments", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetPolicyCommentController.CreateComment)
	assetRoutes.PUT("/assets/:id/policies/:policyId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetPolicyCommentController.UpdateComment)
	assetRoutes.POST("/assets/:id/policies/:policyId/comments/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetPolicyCommentController.DeleteComment)

	//Asset Document category
	assetRoutes.GET("/asset-document-categories", assetDocumentCategoryController.GetAllAssetDocumentCategories)
	assetRoutes.GET("/asset-document-categories/:id", assetDocumentCategoryController.GetAssetDocumentCategoryDetails)
	assetRoutes.POST("/asset-document-categories", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetDocumentCategoryController.CreateAssetDocumentCategory)
	assetRoutes.PUT("/asset-document-categories/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetDocumentCategoryController.UpdateAssetDocumentCategory)
	assetRoutes.POST("/asset-document-categories/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetDocumentCategoryController.DeleteAssetDocumentCategories)
	assetRoutes.PUT("/asset-document-categories/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetDocumentCategoryController.ActivateAssetDocumentCategories)
	assetRoutes.PUT("/asset-document-categories/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetDocumentCategoryController.DeactivateAssetDocumentCategories)

	assetRoutes.GET("/assets/:id/other-documents", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetOtherDocumentController.GetAllDocuments)
	assetRoutes.GET("/assets/:id/other-documents/:documentId", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetOtherDocumentController.GetDocumentDetails)
	assetRoutes.POST("/assets/:id/other-documents", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetOtherDocumentController.CreateDocument)
	assetRoutes.PUT("/assets/:id/other-documents/:documentId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetOtherDocumentController.UpdateDocument)
	assetRoutes.POST("/assets/:id/other-documents/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetOtherDocumentController.DeleteDocuments)
	unauthenticatedRoutes.GET("/assets/:id/other-documents/:documentId/document-download", assetOtherDocumentController.DownloadAssetOtherDocument)

	assetRoutes.GET("assets/:id/control-analysis", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetControlAnalysisController.GetAllAssetControlAnalyses)
	assetRoutes.GET("assets/:id/control-analysis/:analysisId", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetControlAnalysisController.GetAssetControlAnalysisDetails)
	assetRoutes.POST("assets/:id/control-analysis", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetControlAnalysisController.CreateAssetControlAnalysis)
	assetRoutes.PUT("assets/:id/control-analysis/:analysisId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetControlAnalysisController.UpdateAssetControlAnalysis)
	assetRoutes.POST("assets/:id/control-analysis/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetControlAnalysisController.DeleteAssetControlAnalyses)

	assetRoutes.GET("assets/:id/control-analysis/control-summary", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetControlAnalysisController.GetAssetControlSummary)

	assetRoutes.GET("/assets/:id/kris", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetKriController.GetAllAssetKris)
	assetRoutes.GET("/assets/:id/kris/:kriId", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetKriController.GetAssetKriDetails)
	assetRoutes.POST("/assets/:id/kris", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetKriController.CreateAssetKri)
	assetRoutes.PUT("/assets/:id/kris", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetKriController.UpdateAssetKri)
	assetRoutes.POST("/assets/:id/kris/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetKriController.DeleteAssetKris)

}
