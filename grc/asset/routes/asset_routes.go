package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	assetRoutes := router.Group("/")

	assetRoutes.Use(middlewares.AuthMiddleware(app.DB))

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
	vulnerabilityAssociatedItemController := app.AssetModule.VulnerabilityAssociatedItemController
	threatTypeController := app.AssetModule.ThreatTypeController
	threatStatusController := app.AssetModule.ThreatStatusController
	threatController := app.AssetModule.ThreatController
	vulnerabilityActionPlanController := app.AssetModule.VulnerabilityActionPlanController
	vulnerabilityActionPlanStatusUpdateController := app.AssetModule.VulnerabilityActionPlanStatusUpdateController
	vulnerabilityActionPlanCommentController := app.AssetModule.VulnerabilityActionPlanCommentController
	threatActionPlanController := app.AssetModule.ThreatActionPlanController
	threatActionPlanCommentController := app.AssetModule.ThreatActionPlanCommentController

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
	tenableController := app.AssetModule.TenableController
	qualysKBController := app.AssetModule.QualysController
	systemRegisterDisruptionLevelController := app.AssetModule.SystemRegisterDisruptionLevelController
	systemRegisterServiceImportanceController := app.AssetModule.SystemRegisterServiceImportanceController
	systemRegisterController := app.AssetModule.SystemRegisterController
	assetDeviceCategoryController := app.AssetModule.AssetDeviceCategoryController
	confidentialityClassificationController := app.AssetModule.ConfidentialityClassificationController
	assetEnvironmentController := app.AssetModule.AssetEnvironmentController
	applicationController := app.AssetModule.ApplicationController
	slaEscalationMatrixController := app.AssetModule.SLAEscalationMatrixController
	vulnerabilitySettingsController := app.AssetModule.VulnerabilitySettingsController
	vulnerabilityReportController := app.AssetModule.VulnerabilityReportController
	qualysVulnerabilityHistoryController := app.AssetModule.QualysVulnerabilityHistoryController
	assetRiskRegisterController := app.AssetModule.AssetRiskRegisterController
	assetRiskWorkflowController := app.AssetModule.AssetRiskWorkflowController
	assetRiskReviewController := app.AssetModule.AssetRiskReviewController
	assetRiskReviewCommentController := app.AssetModule.AssetRiskReviewCommentController
	vulnerabilityWorkflowController := app.AssetModule.VulnerabilityClosedWorkflowController
	vulnerabilityClosedReviewController := app.AssetModule.VulnerabilityClosedReviewController
	vulnerabilityClosedReviewCommentController := app.AssetModule.VulnerabilityClosedReviewCommentController
	assetRiskRootCauseAnalysisController := app.AssetModule.AssetRiskRootCauseAnalysisController
	assetRiskAssociatedItemController := app.AssetModule.AssetRiskAssociatedItemController
	assetRiskMatrixController := app.AssetModule.AssetRiskMatrixController
	assetRiskInherentEvaluationController := app.AssetModule.AssetRiskInherentEvaluationController
	assetRiskInherentControlAnalysisController := app.AssetModule.AssetRiskInherentControlAnalysisController
	assetRiskResidualEvaluationController := app.AssetModule.AssetRiskResidualEvaluationController
	assetRiskKriController := app.AssetModule.AssetRiskKriController
	assetRiskClosureController := app.AssetModule.AssetRiskClosureController
	assetRiskResponseController := app.AssetModule.AssetRiskResponseController
	assetRiskActionPlanController := app.AssetModule.AssetRiskActionPlanController
	assetRiskActionPlanStatusUpdateController := app.AssetModule.AssetRiskActionPlanStatusUpdateController
	assetRiskActionPlanCommentController := app.AssetModule.AssetRiskActionPlanCommentController

	// Asset Type
	assetRoutes.GET("/asset-types", assetTypeController.GetAllAssetTypes)
	assetRoutes.GET("/asset-types/:id", assetTypeController.GetAssetTypeDetails)
	assetRoutes.POST("/asset-types", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetTypeController.CreateAssetType)
	assetRoutes.PUT("/asset-types/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetTypeController.UpdateAssetType)
	assetRoutes.POST("/asset-types/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetTypeController.DeleteAssetTypes)
	assetRoutes.PUT("/asset-types/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetTypeController.ActivateAssetTypes)
	assetRoutes.PUT("/asset-types/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetTypeController.DeactivateAssetTypes)

	// Asset category
	assetRoutes.GET("/asset-categories", assetCategoryController.GetAllAssetCategories)
	assetRoutes.GET("/asset-categories/:id", assetCategoryController.GetAssetCategoryDetails)
	assetRoutes.POST("/asset-categories", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetCategoryController.CreateAssetCategory)
	assetRoutes.PUT("/asset-categories/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetCategoryController.UpdateAssetCategory)
	assetRoutes.POST("/asset-categories/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetCategoryController.DeleteAssetCategories)
	assetRoutes.PUT("/asset-categories/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetCategoryController.ActivateAssetCategories)
	assetRoutes.PUT("/asset-categories/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetCategoryController.DeactivateAssetCategories)

	// Asset Sub category
	assetRoutes.GET("/asset-sub-categories", assetSubCategoryController.GetAllAssetSubCategories)
	assetRoutes.GET("/asset-sub-categories/:id", assetSubCategoryController.GetAssetSubCategoryDetails)
	assetRoutes.POST("/asset-sub-categories", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetSubCategoryController.CreateAssetSubCategory)
	assetRoutes.PUT("/asset-sub-categories/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetSubCategoryController.UpdateAssetSubCategory)
	assetRoutes.POST("/asset-sub-categories/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetSubCategoryController.DeleteAssetSubCategories)
	assetRoutes.PUT("/asset-sub-categories/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetSubCategoryController.ActivateAssetSubCategories)
	assetRoutes.PUT("/asset-sub-categories/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetSubCategoryController.DeactivateAssetSubCategories)

	// asset status
	assetRoutes.GET("/asset-statuses", assetStatusController.GetAllAssetStatuses)
	assetRoutes.GET("/asset-statuses/:id", assetStatusController.GetAssetStatusDetails)

	// asset register
	assetRoutes.GET("/asset-register", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetRegisterController.GetAllAssets)
	assetRoutes.GET("/asset-register/:id", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetRegisterController.GetAssetDetails)
	assetRoutes.POST("/asset-register", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetRegisterController.CreateAsset)
	assetRoutes.PUT("/asset-register/:id", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetRegisterController.UpdateAsset)
	assetRoutes.POST("/asset-register/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetRegisterController.DeleteAssets)
	assetRoutes.POST("/asset-register/force-delete", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetRegisterController.DeleteForceAssets)
	assetRoutes.POST("/asset-register/archive", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetRegisterController.ArchiveAssets)
	assetRoutes.POST("/asset-register/unarchive", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetRegisterController.UnarchiveAssets)
	assetRoutes.GET("/asset-register/:id/control-summary", middlewares.RolePermissionMiddleware(db, "process", "view"), assetRegisterController.GetControlSummary)
	unauthenticatedRoutes.GET("/asset-register/:id/document-download", assetRegisterController.DownloadAssetDocument)
	assetRoutes.GET("/asset-register/export", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetRegisterController.GetAssetsExport)

	// Asset Matrix
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

	// vulnerability types
	assetRoutes.GET("/vulnerability-types", vulnerabilityTypeController.GetAllVulnerabilityTypes)
	assetRoutes.GET("/vulnerability-types/:id", vulnerabilityTypeController.GetVulnerabilityTypeDetails)
	assetRoutes.POST("/vulnerability-types", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), vulnerabilityTypeController.CreateVulnerabilityType)
	assetRoutes.PUT("/vulnerability-types/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), vulnerabilityTypeController.UpdateVulnerabilityType)
	assetRoutes.POST("/vulnerability-types/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), vulnerabilityTypeController.DeleteVulnerabilityTypes)
	assetRoutes.PUT("/vulnerability-types/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), vulnerabilityTypeController.ActivateVulnerabilityTypes)
	assetRoutes.PUT("/vulnerability-types/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), vulnerabilityTypeController.DeactivateVulnerabilityTypes)

	// vulnerability status
	assetRoutes.GET("/vulnerability-statuses", vulnerabilityStatusController.GetAllVulnerabilityStatuses)
	assetRoutes.GET("/vulnerability-statuses/:id", vulnerabilityStatusController.GetVulnerabilityStatusDetails)

	assetRoutes.GET("/vulnerabilities", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityController.GetAllVulnerabilities)
	assetRoutes.GET("/vulnerabilities/:id", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityController.GetVulnerabilityDetails)
	assetRoutes.POST("/vulnerabilities/:id/qualys-sync", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityController.SyncVulnerabilityFromQualys)
	assetRoutes.POST("/vulnerabilities", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityController.CreateVulnerability)
	assetRoutes.PUT("/vulnerabilities/:id", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityController.UpdateVulnerability)
	assetRoutes.POST("/vulnerabilities/delete", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityController.DeleteVulnerabilities)
	assetRoutes.POST("/vulnerabilities/force-delete", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityController.DeleteForceVulnerabilities)
	assetRoutes.GET("/assets/:id/vulnerability-summary", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityController.GetVulnerabilitySummary)
	assetRoutes.GET("/vulnerabilities/export", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityController.GetVulnerabilitiesExport)
	assetRoutes.PUT("/vulnerabilities/:id/reject", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityController.RejectVulnerability)
	assetRoutes.PUT("/vulnerabilities/:id/approve", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityController.ApproveVulnerability)
	assetRoutes.PUT("/vulnerabilities/:id/deferred", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityController.DeferVulnerability)
	assetRoutes.POST("/vulnerability-remediation/:id", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityController.CreateVulnerabilityRemediation)
	assetRoutes.GET("/vulnerabilities/:id/incidents", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityAssociatedItemController.GetAllVulnerabilityIncidents)
	assetRoutes.POST("/vulnerabilities/:id/incidents", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityAssociatedItemController.CreateVulnerabilityIncident)
	assetRoutes.DELETE("/vulnerabilities/:id/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityAssociatedItemController.DeleteVulnerabilityIncident)

	// Threat types
	assetRoutes.GET("/threat-types", threatTypeController.GetAllThreatTypes)
	assetRoutes.GET("/threat-types/:id", threatTypeController.GetThreatTypeDetails)
	assetRoutes.POST("/threat-types", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), threatTypeController.CreateThreatType)
	assetRoutes.PUT("/threat-types/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), threatTypeController.UpdateThreatType)
	assetRoutes.POST("/threat-types/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), threatTypeController.DeleteThreatTypes)
	assetRoutes.PUT("/threat-types/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), threatTypeController.ActivateThreatTypes)
	assetRoutes.PUT("/threat-types/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), threatTypeController.DeactivateThreatTypes)

	// threat status
	assetRoutes.GET("/threat-statuses", threatStatusController.GetAllThreatStatuses)
	assetRoutes.GET("/threat-statuses/:id", threatStatusController.GetThreatStatusDetails)

	// Threat
	assetRoutes.GET("/threats", threatController.GetAllThreats)
	assetRoutes.GET("/threats/:id", threatController.GetThreatDetails)
	assetRoutes.POST("/threats", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatController.CreateThreat)
	assetRoutes.PUT("/threats/:id", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatController.UpdateThreat)
	assetRoutes.POST("/threats/delete", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatController.DeleteThreats)
	assetRoutes.POST("/threats/force-delete", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatController.DeleteForceThreats)
	assetRoutes.GET("/assets/:id/threat-summary", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatController.GetThreatSummary)
	assetRoutes.GET("/threats/export", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatController.GetThreatsExport)

	// vulnerability action plan
	assetRoutes.GET("/vulnerability-action-plans/export", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanController.GetVulnerabilityActionPlanExport)
	assetRoutes.GET("/vulnerability-action-plans", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanController.GetAllVulnerabilityActionPlans)
	assetRoutes.GET("/vulnerabilities/:id/action-plans", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanController.GetAllVulnerabilityActionPlans)
	assetRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanController.GetVulnerabilityActionPlanDetails)
	assetRoutes.POST("/vulnerabilities/:id/action-plans", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityActionPlanController.CreateVulnerabilityActionPlan)
	assetRoutes.PUT("/vulnerabilities/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityActionPlanController.UpdateVulnerabilityActionPlan)
	assetRoutes.POST("/vulnerabilities/:id/action-plans/delete", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityActionPlanController.DeleteVulnerabilityActionPlans)

	assetRoutes.GET("/vulnerabilities/:id/action-plan-summary", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanController.GetVulnerabilityActionPlanSummary)

	// vulnerability action plan comment
	unauthenticatedRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/comments/:commentId/document-download", vulnerabilityActionPlanCommentController.DownloadVulnerabilityActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/comments/:commentId/document-preview", vulnerabilityActionPlanCommentController.GetVulnerabilityActionPlanCommentDocumentPreview)
	assetRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "vulnerabilities", "view"), vulnerabilityActionPlanCommentController.GetAllComments)
	assetRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "vulnerabilities", "view"), vulnerabilityActionPlanCommentController.GetCommentDetails)
	assetRoutes.POST("/vulnerabilities/:id/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "vulnerabilitie", "modify"), vulnerabilityActionPlanCommentController.CreateComment)
	assetRoutes.PUT("/vulnerabilities/:id/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "threat", "modify"), vulnerabilityActionPlanCommentController.UpdateComment)
	assetRoutes.POST("/vulnerabilities/:id/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "threat", "modify"), vulnerabilityActionPlanCommentController.DeleteComment)

	// Document Preview and Download for Vulnerability Action Plan Status Update
	unauthenticatedRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-preview", vulnerabilityActionPlanStatusUpdateController.GetDocumentPreview)
	unauthenticatedRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-download", vulnerabilityActionPlanStatusUpdateController.DownloadVulnerabilityActionPlanDocument)

	assetRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanStatusUpdateController.GetAllVulnerabilityActionPlanStatusUpdates)
	assetRoutes.POST("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityActionPlanStatusUpdateController.CreateVulnerabilityActionPlanStatusUpdate)
	assetRoutes.PUT("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates/:actionPlanStatusUpdateId", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityActionPlanStatusUpdateController.UpdateVulnerabilityActionPlanStatusUpdate)
	assetRoutes.GET("/vulnerabilities/:id/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityActionPlanStatusUpdateController.GetLatestVulnerabilityActionPlanStatusUpdate)
	assetRoutes.PUT("/vulnerabilities/:id/action-plans/:actionPlanId/approve", vulnerabilityActionPlanStatusUpdateController.ApproveVulnerabilityActionPlan)
	assetRoutes.PUT("/vulnerabilities/:id/action-plans/:actionPlanId/revert", vulnerabilityActionPlanStatusUpdateController.RevertVulnerabilityActionPlan)

	// Threat action plan
	assetRoutes.GET("/threat-action-plans/export", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanController.GetThreatActionPlanExport)
	assetRoutes.GET("/threat-action-plans", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanController.GetAllThreatActionPlans)
	assetRoutes.GET("/threats/:id/action-plans", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanController.GetAllThreatActionPlans)
	assetRoutes.GET("/threats/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanController.GetThreatActionPlanDetails)
	assetRoutes.POST("/threats/:id/action-plans", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatActionPlanController.CreateThreatActionPlan)
	assetRoutes.PUT("/threats/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatActionPlanController.UpdateThreatActionPlan)
	assetRoutes.POST("/threats/:id/action-plans/delete", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatActionPlanController.DeleteThreatActionPlans)

	assetRoutes.GET("/threats/:id/action-plan-summary", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanController.GetThreatActionPlanSummary)

	// Threat action plan comment
	unauthenticatedRoutes.GET("/threats/:id/action-plans/:actionPlanId/comments/:commentId/document-download", threatActionPlanCommentController.DownloadThreatActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/threats/:id/action-plans/:actionPlanId/comments/:commentId/document-preview", threatActionPlanCommentController.GetThreatActionPlanCommentDocumentPreview)
	assetRoutes.GET("/threats/:id/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanCommentController.GetAllComments)
	assetRoutes.GET("/threats/:id/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanCommentController.GetCommentDetails)
	assetRoutes.POST("/threats/:id/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatActionPlanCommentController.CreateComment)
	assetRoutes.PUT("/threats/:id/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatActionPlanCommentController.UpdateComment)
	assetRoutes.POST("/threats/:id/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatActionPlanCommentController.DeleteComment)

	// Document Preview and Download for Threat Action Plan Status Update
	unauthenticatedRoutes.GET("/threats/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-preview", threatActionPlanStatusUpdateController.GetDocumentPreview)
	unauthenticatedRoutes.GET("/threats/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-download", threatActionPlanStatusUpdateController.DownloadThreatActionPlanDocument)

	assetRoutes.GET("/threats/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanStatusUpdateController.GetAllThreatActionPlanStatusUpdates)
	assetRoutes.POST("/threats/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatActionPlanStatusUpdateController.CreateThreatActionPlanStatusUpdate)
	assetRoutes.PUT("/threats/:id/action-plans/:actionPlanId/status-updates/:actionPlanStatusUpdateId", middlewares.RolePermissionMiddleware(db, "threat", "modify"), threatActionPlanStatusUpdateController.UpdateThreatActionPlanStatusUpdate)
	assetRoutes.GET("/threats/:id/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "threat", "view"), threatActionPlanStatusUpdateController.GetLatestThreatActionPlanStatusUpdate)

	assetRoutes.PUT("/threats/:id/action-plans/:actionPlanId/approve", threatActionPlanStatusUpdateController.ApproveThreatActionPlan)
	assetRoutes.PUT("/threats/:id/action-plans/:actionPlanId/revert", threatActionPlanStatusUpdateController.RevertThreatActionPlan)

	assetRoutes.GET("/assets/:id/maintenance", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetMaintenanceController.GetAllMaintenances)
	assetRoutes.GET("/assets/:id/maintenance/:maintenanceId", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetMaintenanceController.GetMaintenanceDetails)
	assetRoutes.POST("/assets/:id/maintenance", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetMaintenanceController.CreateMaintenance)
	assetRoutes.PUT("/assets/:id/maintenance/:maintenanceId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetMaintenanceController.UpdateMaintenance)
	assetRoutes.POST("/assets/:id/maintenance/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetMaintenanceController.DeleteMaintenances)

	// asset associated item list
	assetRoutes.GET("assets/:id/controls", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetControls)
	assetRoutes.GET("assets/:id/incidents", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetIncidents)
	assetRoutes.GET("assets/:id/risks", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetRisks)
	assetRoutes.GET("assets/:id/vulnerabilities", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetVulnerabilities)
	assetRoutes.GET("assets/:id/clauses", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetClauses)
	assetRoutes.PUT("assets/:id/clauses", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.UpdateClauses)

	assetRoutes.POST("assets/:id/vulnerabilities", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.CreateAssetVulnerabilities)
	assetRoutes.POST("assets/:id/vulnerabilities/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.DeleteAssetVulnerabilities)

	assetRoutes.POST("assets/:id/compliance-clauses", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.CreateAssetComplianceClauses)
	assetRoutes.POST("assets/:id/compliance-clauses/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.DeleteAssetComplianceClauses)

	assetRoutes.POST("vulnerabilities/:id/threats", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.CreateAssetThreats)
	assetRoutes.POST("vulnerabilities/:id/threats/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.DeleteAssetThreats)

	assetRoutes.POST("threats/:id/risks", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.CreateAssetThreatRisks)
	assetRoutes.POST("threats/:id/risks/delete", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.DeleteAssetThreatRisks)

	assetRoutes.GET("assets/:id/risk-analysis", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetRiskAnalysis)

	assetRoutes.POST("assets/:id/risks", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.CreateAssetRisk)
	assetRoutes.DELETE("assets/:id/risks/:riskId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.DeleteAssetRisk)

	assetRoutes.GET("assets/:id/documents", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetDocuments)
	assetRoutes.POST("assets/:id/documents", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.CreateAssetDocument)
	assetRoutes.DELETE("assets/:id/documents/:documentId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.DeleteAssetDocument)

	assetRoutes.POST("assets/:id/controls", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.CreateAssetControl)
	assetRoutes.DELETE("assets/:id/controls/:controlId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.DeleteAssetControl)

	assetRoutes.GET("assets/:id/processes", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetProcess)
	assetRoutes.POST("assets/:id/processes", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.CreateAssetProcess)
	assetRoutes.DELETE("assets/:id/processes/:processId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.DeleteAssetProcess)

	assetRoutes.GET("assets/:id/strategic-objectives", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetStrategicObjectives)
	assetRoutes.POST("assets/:id/strategic-objectives", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.CreateAssetStrategicObjective)
	assetRoutes.DELETE("assets/:id/strategic-objectives/:strategicObjectiveId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.DeleteAssetStrategicObjective)

	assetRoutes.GET("assets/:id/kpis", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), assetAssociatedItemController.GetAllAssetKPI)
	assetRoutes.POST("assets/:id/kpis", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.CreateAssetKPI)
	assetRoutes.DELETE("assets/:id/kpis/:kpiRegisterId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.DeleteAssetKPI)

	assetRoutes.POST("assets/:id/incidents", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.CreateAssetIncident)
	assetRoutes.DELETE("assets/:id/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), assetAssociatedItemController.DeleteAssetIncident)

	//asset criticality
	// asset criticality
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

	// Asset Document category
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

	//tenable
	assetRoutes.GET("/tenable/assets", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), tenableController.GetAllAssets)
	assetRoutes.POST("/tenable/assets/validate", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), tenableController.ValidateTenableAssets)
	assetRoutes.POST("/tenable/assets/import", middlewares.RolePermissionMiddleware(db, "asset_register", "modify"), tenableController.ImportAssets)

	//qualys knowledge base
	assetRoutes.GET("/qualys/:qid", middlewares.RolePermissionMiddleware(db, "asset_register", "view"), qualysKBController.GetKBDetailsByQID)

	// System Register Disruption Levels
	assetRoutes.GET("/system-register-disruption-levels", systemRegisterDisruptionLevelController.GetAllSystemRegisterDisruptionLevels)
	assetRoutes.GET("/system-register-disruption-levels/:id", systemRegisterDisruptionLevelController.GetSystemRegisterDisruptionLevelDetails)
	assetRoutes.POST("/system-register-disruption-levels", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), systemRegisterDisruptionLevelController.CreateSystemRegisterDisruptionLevel)
	assetRoutes.PUT("/system-register-disruption-levels/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), systemRegisterDisruptionLevelController.UpdateSystemRegisterDisruptionLevel)
	assetRoutes.POST("/system-register-disruption-levels/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), systemRegisterDisruptionLevelController.DeleteSystemRegisterDisruptionLevels)
	assetRoutes.PUT("/system-register-disruption-levels/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), systemRegisterDisruptionLevelController.ActivateSystemRegisterDisruptionLevels)
	assetRoutes.PUT("/system-register-disruption-levels/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), systemRegisterDisruptionLevelController.DeactivateSystemRegisterDisruptionLevels)

	// System Register Service Importances
	assetRoutes.GET("/system-register-service-importances", systemRegisterServiceImportanceController.GetAllSystemRegisterServiceImportances)
	assetRoutes.GET("/system-register-service-importances/:id", systemRegisterServiceImportanceController.GetSystemRegisterServiceImportanceDetails)
	assetRoutes.POST("/system-register-service-importances", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), systemRegisterServiceImportanceController.CreateSystemRegisterServiceImportance)
	assetRoutes.PUT("/system-register-service-importances/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), systemRegisterServiceImportanceController.UpdateSystemRegisterServiceImportance)
	assetRoutes.POST("/system-register-service-importances/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), systemRegisterServiceImportanceController.DeleteSystemRegisterServiceImportances)
	assetRoutes.PUT("/system-register-service-importances/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), systemRegisterServiceImportanceController.ActivateSystemRegisterServiceImportances)
	assetRoutes.PUT("/system-register-service-importances/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), systemRegisterServiceImportanceController.DeactivateSystemRegisterServiceImportances)

	// System Register
	assetRoutes.GET("/system-registers", middlewares.RolePermissionMiddleware(db, "system_register", "view"), systemRegisterController.GetAllSystemRegisters)
	assetRoutes.GET("/system-registers/:id", middlewares.RolePermissionMiddleware(db, "system_register", "view"), systemRegisterController.GetSystemRegisterDetails)
	assetRoutes.POST("/system-registers", middlewares.RolePermissionMiddleware(db, "system_register", "modify"), systemRegisterController.CreateSystemRegister)
	assetRoutes.PUT("/system-registers/:id", middlewares.RolePermissionMiddleware(db, "system_register", "modify"), systemRegisterController.UpdateSystemRegister)
	assetRoutes.POST("/system-registers/delete", middlewares.RolePermissionMiddleware(db, "system_register", "modify"), systemRegisterController.DeleteSystemRegisters)
	assetRoutes.GET("/system-registers/export", middlewares.RolePermissionMiddleware(db, "system_register", "view"), systemRegisterController.GetSystemRegistersExport)

	// Asset Device Category
	assetRoutes.GET("/asset-device-categories", assetDeviceCategoryController.GetAllAssetDeviceCategories)
	assetRoutes.GET("/asset-device-categories/:id", assetDeviceCategoryController.GetAssetDeviceCategoryDetails)
	assetRoutes.POST("/asset-device-categories", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetDeviceCategoryController.CreateAssetDeviceCategory)
	assetRoutes.PUT("/asset-device-categories/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetDeviceCategoryController.UpdateAssetDeviceCategory)
	assetRoutes.POST("/asset-device-categories/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetDeviceCategoryController.DeleteAssetDeviceCategories)
	assetRoutes.PUT("/asset-device-categories/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetDeviceCategoryController.ActivateAssetDeviceCategories)
	assetRoutes.PUT("/asset-device-categories/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), assetDeviceCategoryController.DeactivateAssetDeviceCategories)

	// Confidentiality Classification
	assetRoutes.GET("/confidentiality-classifications", confidentialityClassificationController.GetAllConfidentialityClassifications)
	assetRoutes.GET("/confidentiality-classifications/:id", confidentialityClassificationController.GetConfidentialityClassificationDetails)

	// Asset Environments
	assetRoutes.GET("/asset-environments", assetEnvironmentController.GetAllAssetEnvironments)
	assetRoutes.GET("/asset-environments/:id", assetEnvironmentController.GetAssetEnvironmentDetails)

	// Applications
	assetRoutes.GET("/applications", applicationController.GetAllApplications)
	assetRoutes.GET("/applications/:id", applicationController.GetApplicationDetails)
	assetRoutes.POST("/applications", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), applicationController.CreateApplication)
	assetRoutes.PUT("/applications/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), applicationController.UpdateApplication)
	assetRoutes.POST("/applications/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), applicationController.DeleteApplications)
	assetRoutes.PUT("/applications/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), applicationController.ActivateApplications)
	assetRoutes.PUT("/applications/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), applicationController.DeactivateApplications)

	// SLA Escalation Matrix
	assetRoutes.GET("/sla-escalation-matrix", slaEscalationMatrixController.GetAllMatrices)
	assetRoutes.GET("/sla-escalation-matrix/:id", slaEscalationMatrixController.GetMatrixDetails)
	assetRoutes.POST("/sla-escalation-matrix", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), slaEscalationMatrixController.CreateMatrix)
	assetRoutes.PUT("/sla-escalation-matrix/:id", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), slaEscalationMatrixController.UpdateMatrix)
	assetRoutes.POST("/sla-escalation-matrix/delete", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), slaEscalationMatrixController.DeleteMatrices)
	assetRoutes.PUT("/sla-escalation-matrix/activate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), slaEscalationMatrixController.ActivateMatrices)
	assetRoutes.PUT("/sla-escalation-matrix/deactivate", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), slaEscalationMatrixController.DeactivateMatrices)

	// Vulnerability Settings
	assetRoutes.GET("/vulnerability-settings", vulnerabilitySettingsController.GetAllVulnerabilitySettings)
	assetRoutes.PUT("/vulnerability-settings", middlewares.RolePermissionMiddleware(db, "asset_setting", "modify"), vulnerabilitySettingsController.UpdateVulnerabilitySettings)

	// Vulnerability Report
	assetRoutes.GET("/vulnerability-reports", middlewares.RolePermissionMiddleware(db, "asset_report", "view"), vulnerabilityReportController.GetAllVulnerabilityReports)
	assetRoutes.GET("/vulnerability-reports/:vulnerabilityReportId", middlewares.RolePermissionMiddleware(db, "asset_report", "view"), vulnerabilityReportController.GetVulnerabilityReportDetails)
	assetRoutes.POST("/vulnerability-reports", middlewares.RolePermissionMiddleware(db, "asset_report", "view"), vulnerabilityReportController.CreateVulnerabilityReport)
	assetRoutes.DELETE("/vulnerability-reports/:vulnerabilityReportId", middlewares.RolePermissionMiddleware(db, "asset_report", "view"), vulnerabilityReportController.DeleteVulnerabilityReport)
	assetRoutes.GET("/vulnerability-reports/:vulnerabilityReportId/export-pdf", middlewares.RolePermissionMiddleware(db, "asset_report", "view"), vulnerabilityReportController.GetVulnerabilityReportExportZip)

	//Qulays Vulnerability History
	assetRoutes.GET("/vulnerabilities/:id/qualys-history", qualysVulnerabilityHistoryController.GetAllQualysVulnerabilityHistories)
	assetRoutes.POST("/vulnerabilities/:id/qualys-history/:historyId/documents", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), qualysVulnerabilityHistoryController.CreateQualysVulnerabilityHistoryDocuments)
	unauthenticatedRoutes.GET("/vulnerabilities/:id/qualys-history/:historyId/document-preview", qualysVulnerabilityHistoryController.GetQualysVulnerabilityHistoryDocumentPreview)

	// Asset Risk Register
	assetRoutes.GET("/asset-risk-register", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskRegisterController.GetAllAssetRiskRegisters)
	assetRoutes.GET("/asset-risk-register/export", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskRegisterController.GetAssetRiskRegisterExport)
	assetRoutes.GET("/asset-risk-register/:id", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskRegisterController.GetAssetRiskRegisterDetails)
	assetRoutes.POST("/asset-risk-register", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), assetRiskRegisterController.CreateAssetRiskRegister)
	assetRoutes.PUT("/asset-risk-register/:id", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), assetRiskRegisterController.UpdateAssetRiskRegister)
	assetRoutes.POST("/asset-risk-register/delete", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskRegisterController.DeleteAssetRiskRegisters)

	// Risk Acceptance Workflow
	assetRoutes.POST("/asset-risk-register/:id/submit", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), assetRiskReviewController.SubmitAssetRiskRegister)
	assetRoutes.POST("/asset-risk-register/:id/approve", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), assetRiskReviewController.ApproveAssetRiskRegister)
	assetRoutes.POST("/asset-risk-register/:id/reject", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), assetRiskReviewController.RejectAssetRiskRegister)
	assetRoutes.POST("/asset-risk-register/:id/resubmit", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), assetRiskReviewController.ResubmitAssetRiskRegister)
	assetRoutes.POST("/asset-risk-register/:id/revert", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), assetRiskReviewController.RevertAssetRiskRegister)
	assetRoutes.GET("/asset-risk-register/:id/workflow-list", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), assetRiskWorkflowController.GetWorkflowList)

	// Risk Acceptance Workflow Review Comment
	assetRoutes.GET("/asset-risk-register/:id/review-comments", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), assetRiskReviewCommentController.GetAllComments)
	assetRoutes.GET("/asset-risk-register/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), assetRiskReviewCommentController.GetCommentDetails)
	assetRoutes.POST("/asset-risk-register/:id/review-comments", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), assetRiskReviewCommentController.CreateComment)
	assetRoutes.PUT("/asset-risk-register/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), assetRiskReviewCommentController.UpdateComment)
	assetRoutes.POST("/asset-risk-register/:id/review-comments/delete", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), assetRiskReviewCommentController.DeleteComment)
	assetRoutes.PUT("/asset-risk-register/:id/review-comments/:commentId/resolve", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), assetRiskReviewCommentController.ResolveComment)

	// Vulnerability Closure Workflow
	assetRoutes.POST("/vulnerabilities/:id/submit", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityClosedReviewController.SubmitVulnerability)
	assetRoutes.POST("/vulnerabilities/:id/approve", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityClosedReviewController.ApproveVulnerability)
	assetRoutes.POST("/vulnerabilities/:id/reject", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityClosedReviewController.RejectVulnerability)
	assetRoutes.POST("/vulnerabilities/:id/resubmit", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityClosedReviewController.ResubmitVulnerability)
	assetRoutes.POST("/vulnerabilities/:id/revert", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityClosedReviewController.RevertVulnerability)
	assetRoutes.GET("/vulnerabilities/:id/workflow-list", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityWorkflowController.GetWorkflowList)

	assetRoutes.GET("/vulnerabilities/:id/review-comments", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityClosedReviewCommentController.GetAllComments)
	assetRoutes.GET("/vulnerabilities/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "vulnerability", "view"), vulnerabilityClosedReviewCommentController.GetCommentDetails)
	assetRoutes.POST("/vulnerabilities/:id/review-comments", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityClosedReviewCommentController.CreateComment)
	assetRoutes.PUT("/vulnerabilities/:id/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityClosedReviewCommentController.UpdateComment)
	assetRoutes.POST("/vulnerabilities/:id/review-comments/delete", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityClosedReviewCommentController.DeleteComments)
	assetRoutes.PUT("/vulnerabilities/:id/review-comments/:commentId/resolve", middlewares.RolePermissionMiddleware(db, "vulnerability", "modify"), vulnerabilityClosedReviewCommentController.ResolveComment)

	// Risk Root Cause Analysis
	assetRoutes.GET("/asset-risk-register/:id/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskRootCauseAnalysisController.GetAllAssetRiskRootCauseAnalyses)
	assetRoutes.POST("/asset-risk-register/:id/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskRootCauseAnalysisController.CreateAssetRiskRootCauseAnalysis)
	assetRoutes.POST("/asset-risk-register/:id/root-cause-analyses", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskRootCauseAnalysisController.CreateAssetRiskRootCauseAnalyses)
	assetRoutes.PUT("/asset-risk-register/:id/root-cause-analysis/:rcaId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskRootCauseAnalysisController.UpdateAssetRiskRootCauseAnalysis)
	assetRoutes.DELETE("/asset-risk-register/:id/root-cause-analysis/:rcaId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskRootCauseAnalysisController.DeleteAssetRiskRootCauseAnalysis)

	// Risk Associated Items
	assetRoutes.GET("/asset-risk-register/:id/controls", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskAssociatedItemController.GetAllAssetRiskControls)
	assetRoutes.POST("/asset-risk-register/:id/controls", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.CreateAssetRiskControl)
	assetRoutes.DELETE("/asset-risk-register/:id/controls/:controlId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.DeleteAssetRiskControl)

	assetRoutes.GET("/asset-risk-register/:id/processes", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskAssociatedItemController.GetAllAssetRiskProcesses)
	assetRoutes.POST("/asset-risk-register/:id/processes", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.CreateAssetRiskProcess)
	assetRoutes.DELETE("/asset-risk-register/:id/processes/:processId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.DeleteAssetRiskProcess)

	assetRoutes.GET("/asset-risk-register/:id/strategic-objectives", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskAssociatedItemController.GetAllAssetRiskStrategicObjectives)
	assetRoutes.POST("/asset-risk-register/:id/strategic-objectives", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.CreateAssetRiskStrategicObjective)
	assetRoutes.DELETE("/asset-risk-register/:id/strategic-objectives/:strategicObjectiveId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.DeleteAssetRiskStrategicObjective)

	assetRoutes.GET("/asset-risk-register/:id/incidents", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskAssociatedItemController.GetAllAssetRiskIncidents)
	assetRoutes.POST("/asset-risk-register/:id/incidents", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.CreateAssetRiskIncident)
	assetRoutes.DELETE("/asset-risk-register/:id/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.DeleteAssetRiskIncident)

	assetRoutes.GET("/asset-risk-register/:id/documents", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskAssociatedItemController.GetAllAssetRiskDocuments)
	assetRoutes.POST("/asset-risk-register/:id/documents", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.CreateAssetRiskDocument)
	assetRoutes.DELETE("/asset-risk-register/:id/documents/:documentId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.DeleteAssetRiskDocument)

	assetRoutes.GET("/asset-risk-register/:id/assets", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskAssociatedItemController.GetAllAssetRiskAssets)
	assetRoutes.POST("/asset-risk-register/:id/assets", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.CreateAssetRiskAsset)
	assetRoutes.DELETE("/asset-risk-register/:id/assets/:assetId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.DeleteAssetRiskAsset)

	assetRoutes.GET("/asset-risk-register/:id/clauses", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskAssociatedItemController.GetAllAssetRiskClauses)
	assetRoutes.PUT("/asset-risk-register/:id/clauses", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.UpdateAssetRiskClauses)

	assetRoutes.GET("/asset-risk-register/:id/kpis", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskAssociatedItemController.GetAllAssetRiskKPI)
	assetRoutes.POST("/asset-risk-register/:id/kpis", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.CreateAssetRiskKPI)
	assetRoutes.DELETE("/asset-risk-register/:id/kpis/:kpiId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskAssociatedItemController.DeleteAssetRiskKPI)

	// Risk Matrix (ISMS)
	assetRoutes.GET("/asset-risk-matrix/impacts", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskMatrixController.GetAllAssetRiskMatrixImpacts)
	assetRoutes.GET("/asset-risk-matrix/likelihoods", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskMatrixController.GetAllAssetRiskMatrixLikelihoods)
	assetRoutes.GET("/asset-risk-matrix/risk-ratings", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskMatrixController.GetAllAssetRiskMatrixRiskRatings)
	assetRoutes.GET("/asset-risk-matrix/risk-rating-values", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskMatrixController.GetAllAssetRiskMatrixRiskRatingValues)
	assetRoutes.GET("/asset-risk-matrix-configuration", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskMatrixController.GetAssetRiskMatrixConfiguration)

	// Risk Inherent Evaluation
	assetRoutes.GET("/asset-risk-register/:id/inherent-evaluation", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskInherentEvaluationController.GetAssetRiskInherentEvaluation)
	assetRoutes.POST("/asset-risk-register/:id/inherent-evaluation", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskInherentEvaluationController.CreateAssetRiskInherentEvaluation)
	assetRoutes.PUT("/asset-risk-register/:id/inherent-evaluation/:evaluationId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskInherentEvaluationController.UpdateAssetRiskInherentEvaluation)
	assetRoutes.DELETE("/asset-risk-register/:id/inherent-evaluation/:evaluationId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskInherentEvaluationController.DeleteAssetRiskInherentEvaluation)

	// Risk Inherent Control Analysis
	assetRoutes.GET("/asset-risk-register/:id/inherent-control-analysis", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskInherentControlAnalysisController.GetAllAssetRiskControlAnalyses)
	assetRoutes.GET("/asset-risk-register/:id/inherent-control-analysis/control-summary", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskInherentControlAnalysisController.GetAssetRiskControlSummary)
	assetRoutes.GET("/asset-risk-register/:id/inherent-control-analysis/:analysisId", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskInherentControlAnalysisController.GetAssetRiskControlAnalysisDetails)
	assetRoutes.POST("/asset-risk-register/:id/inherent-control-analysis", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskInherentControlAnalysisController.CreateAssetRiskControlAnalysis)
	assetRoutes.PUT("/asset-risk-register/:id/inherent-control-analysis/:analysisId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskInherentControlAnalysisController.UpdateAssetRiskControlAnalysis)
	assetRoutes.POST("/asset-risk-register/:id/inherent-control-analysis/delete", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskInherentControlAnalysisController.DeleteAssetRiskControlAnalyses)

	// Risk Residual Evaluation
	assetRoutes.GET("/asset-risk-register/:id/residual-evaluation", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskResidualEvaluationController.GetAssetRiskResidualEvaluation)
	assetRoutes.POST("/asset-risk-register/:id/residual-evaluation", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskResidualEvaluationController.CreateAssetRiskResidualEvaluation)
	assetRoutes.PUT("/asset-risk-register/:id/residual-evaluation/:evaluationId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskResidualEvaluationController.UpdateAssetRiskResidualEvaluation)
	assetRoutes.DELETE("/asset-risk-register/:id/residual-evaluation/:evaluationId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskResidualEvaluationController.DeleteAssetRiskResidualEvaluation)

	// Risk Response
	assetRoutes.GET("/asset-risk-register/:id/response", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskResponseController.GetAssetRiskResponse)
	assetRoutes.POST("/asset-risk-register/:id/response", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskResponseController.CreateAssetRiskResponse)
	assetRoutes.PUT("/asset-risk-register/:id/response/:responseId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskResponseController.UpdateAssetRiskResponse)
	assetRoutes.DELETE("/asset-risk-register/:id/response/:responseId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskResponseController.DeleteAssetRiskResponse)

	// Risk Action Plan
	assetRoutes.GET("/asset-risk-register/:id/action-plans", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskActionPlanController.GetAllAssetRiskActionPlans)
	assetRoutes.GET("/asset-risk-register/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskActionPlanController.GetAssetRiskActionPlanDetails)
	assetRoutes.POST("/asset-risk-register/:id/action-plans", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskActionPlanController.CreateAssetRiskActionPlan)
	assetRoutes.PUT("/asset-risk-register/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskActionPlanController.UpdateAssetRiskActionPlan)
	assetRoutes.POST("/asset-risk-register/:id/action-plans/delete", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskActionPlanController.DeleteAssetRiskActionPlans)

	// Passive Routes
	assetRoutes.GET("/asset-risk-action-plans/export", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "view"), assetRiskActionPlanController.GetAssetRiskActionPlanExport)
	assetRoutes.GET("/asset-risk-action-plans", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "view"), assetRiskActionPlanController.GetAllAssetRiskActionPlans)
	assetRoutes.GET("/asset-risk-action-plans/:actionPlanId", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "view"), assetRiskActionPlanController.GetAssetRiskActionPlanDetails)

	// Risk Action Plan Status Updates
	unauthenticatedRoutes.GET("/asset-risk-register/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-preview", assetRiskActionPlanStatusUpdateController.GetAssetRiskActionPlanStatusUpdateDocumentPreview)
	unauthenticatedRoutes.GET("/asset-risk-register/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-download", assetRiskActionPlanStatusUpdateController.DownloadAssetRiskActionPlanDocument)

	assetRoutes.GET("/asset-risk-register/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskActionPlanStatusUpdateController.GetAllAssetRiskActionPlanStatusUpdates)
	assetRoutes.POST("/asset-risk-register/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskActionPlanStatusUpdateController.CreateAssetRiskActionPlanStatusUpdate)
	assetRoutes.PUT("/asset-risk-register/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskActionPlanStatusUpdateController.UpdateAssetRiskActionPlanStatusUpdate)
	assetRoutes.GET("/asset-risk-register/:id/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskActionPlanStatusUpdateController.GetLatestAssetRiskActionPlanStatusUpdate)
	assetRoutes.PUT("/asset-risk-register/:id/action-plans/:actionPlanId/approve", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskActionPlanStatusUpdateController.ApproveAssetRiskActionPlan)
	assetRoutes.PUT("/asset-risk-register/:id/action-plans/:actionPlanId/revert", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskActionPlanStatusUpdateController.RevertAssetRiskActionPlan)

	// Passive Routes
	assetRoutes.GET("/asset-risk-action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "view"), assetRiskActionPlanStatusUpdateController.GetAllAssetRiskActionPlanStatusUpdates)
	assetRoutes.POST("/asset-risk-action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "modify"), assetRiskActionPlanStatusUpdateController.CreateAssetRiskActionPlanStatusUpdate)
	assetRoutes.PUT("/asset-risk-action-plans/:actionPlanId/status-updates/:statusUpdateId", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "modify"), assetRiskActionPlanStatusUpdateController.UpdateAssetRiskActionPlanStatusUpdate)
	assetRoutes.GET("/asset-risk-action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "view"), assetRiskActionPlanStatusUpdateController.GetLatestAssetRiskActionPlanStatusUpdate)
	assetRoutes.PUT("/asset-risk-action-plans/:actionPlanId/approve", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "modify"), assetRiskActionPlanStatusUpdateController.ApproveAssetRiskActionPlan)
	assetRoutes.PUT("/asset-risk-action-plans/:actionPlanId/revert", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "modify"), assetRiskActionPlanStatusUpdateController.RevertAssetRiskActionPlan)

	// Risk Action Plan Comments
	unauthenticatedRoutes.GET("/asset-risk-register/:id/action-plans/:actionPlanId/comments/:commentId/document-download", assetRiskActionPlanCommentController.DownloadAssetRiskActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/asset-risk-register/:id/action-plans/:actionPlanId/comments/:commentId/document-preview", assetRiskActionPlanCommentController.GetAssetRiskActionPlanCommentDocumentPreview)

	assetRoutes.GET("/asset-risk-register/:id/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskActionPlanCommentController.GetAllAssetRiskActionPlanComments)
	assetRoutes.GET("/asset-risk-register/:id/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskActionPlanCommentController.GetAssetRiskActionPlanCommentDetails)
	assetRoutes.POST("/asset-risk-register/:id/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskActionPlanCommentController.CreateAssetRiskActionPlanComment)
	assetRoutes.PUT("/asset-risk-register/:id/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskActionPlanCommentController.UpdateAssetRiskActionPlanComment)
	assetRoutes.POST("/asset-risk-register/:id/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskActionPlanCommentController.DeleteAssetRiskActionPlanComments)

	// Passive Routes
	assetRoutes.GET("/asset-risk-action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "view"), assetRiskActionPlanCommentController.GetAllAssetRiskActionPlanComments)
	assetRoutes.GET("/asset-risk-action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "view"), assetRiskActionPlanCommentController.GetAssetRiskActionPlanCommentDetails)
	assetRoutes.POST("/asset-risk-action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "modify"), assetRiskActionPlanCommentController.CreateAssetRiskActionPlanComment)
	assetRoutes.PUT("/asset-risk-action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "modify"), assetRiskActionPlanCommentController.UpdateAssetRiskActionPlanComment)
	assetRoutes.POST("/asset-risk-action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddlewareWithPassive(db, "asset_risk_mitigation", "modify"), assetRiskActionPlanCommentController.DeleteAssetRiskActionPlanComments)

	// Risk KRI
	assetRoutes.GET("/asset-risk-register/:id/kris", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskKriController.GetAllAssetRiskKris)
	assetRoutes.GET("/asset-risk-register/:id/kris/:kriId", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskKriController.GetAssetRiskKriDetails)
	assetRoutes.POST("/asset-risk-register/:id/kris", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskKriController.CreateAssetRiskKri)
	assetRoutes.POST("/asset-risk-register/:id/kris/delete", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskKriController.DeleteAssetRiskKris)

	// Risk Closure
	unauthenticatedRoutes.GET("/asset-risk-register/:id/closures/:closureId/document-preview", assetRiskClosureController.GetAssetRiskClosureDocumentPreview)
	unauthenticatedRoutes.GET("/asset-risk-register/:id/closures/:closureId/document-download", assetRiskClosureController.DownloadAssetRiskClosureDocument)
	assetRoutes.GET("/asset-risk-register/:id/closures", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskClosureController.GetAllAssetRiskClosures)
	assetRoutes.GET("/asset-risk-register/:id/closures/:closureId", middlewares.RolePermissionMiddleware(db, "asset_risk", "view"), assetRiskClosureController.GetAssetRiskClosureDetails)
	assetRoutes.POST("/asset-risk-register/:id/closures", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskClosureController.CreateAssetRiskClosure)
	assetRoutes.PUT("/asset-risk-register/:id/closures/:closureId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskClosureController.UpdateAssetRiskClosure)
	assetRoutes.DELETE("/asset-risk-register/:id/closures/:closureId", middlewares.RolePermissionMiddleware(db, "asset_risk", "modify"), assetRiskClosureController.DeleteAssetRiskClosure)

}
