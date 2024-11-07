package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	assetRoutes := router.Group("/")

	assetRoutes.Use(middlewares.AuthMiddleware())

	// unauthenticatedRoutes := router.Group("/")

	db := app.DB

	assetTypeController := app.AssetModule.AssetTypeController
	assetCategoryController := app.AssetModule.AssetCategoryController
	assetStatusController := app.AssetModule.AssetStatusController
	assetRegisterController := app.AssetModule.AssetRegisterController
	assetCriticalityMatrixFactorController := app.AssetModule.AssetCriticalityMatrixFactorController
	assetCriticalityMatrixFactorRatingController := app.AssetModule.AssetCriticalityMatrixFactorRatingController
	assetCriticalityMatrixRatingScaleController := app.AssetModule.AssetCriticalityMatrixRatingScaleController
	vulnerabilityTypeController := app.AssetModule.VulnerabilityTypeController
	vulnerabilityStatusController := app.AssetModule.VulnerabilityStatusController
	vulnerabilityController := app.AssetModule.VulnerabilityController
	threatTypeController := app.AssetModule.ThreatTypeController
	threatStatusController := app.AssetModule.ThreatStatusController
	threatController := app.AssetModule.ThreatController

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
	assetRoutes.GET("/assets", middlewares.RolePermissionMiddleware(db, "assets", "view"), assetRegisterController.GetAllAssets)
	assetRoutes.GET("/assets/:id", middlewares.RolePermissionMiddleware(db, "assets", "view"), assetRegisterController.GetAssetDetails)
	assetRoutes.POST("/assets", middlewares.RolePermissionMiddleware(db, "assets", "modify"), assetRegisterController.CreateAsset)
	assetRoutes.PUT("/assets/:id", middlewares.RolePermissionMiddleware(db, "assets", "modify"), assetRegisterController.UpdateAsset)
	assetRoutes.POST("/assets/delete", middlewares.RolePermissionMiddleware(db, "assets", "modify"), assetRegisterController.DeleteAssets)

	// Asset Criticality Matrix Factors
	assetRoutes.GET("/asset-criticality-matrix-factors", assetCriticalityMatrixFactorController.GetAllFactors)
	assetRoutes.GET("/asset-criticality-matrix-factors/:id", assetCriticalityMatrixFactorController.GetFactorDetails)
	assetRoutes.POST("/asset-criticality-matrix-factors", middlewares.RolePermissionMiddleware(db, "asset_criticality", "modify"), assetCriticalityMatrixFactorController.CreateFactor)
	assetRoutes.PUT("/asset-criticality-matrix-factors/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality", "modify"), assetCriticalityMatrixFactorController.UpdateFactor)
	assetRoutes.POST("/asset-criticality-matrix-factors/delete", middlewares.RolePermissionMiddleware(db, "asset_criticality", "modify"), assetCriticalityMatrixFactorController.DeleteFactors)

	// Asset Criticality Matrix Factor Ratings
	assetRoutes.GET("/asset-criticality-matrix-factor-ratings", assetCriticalityMatrixFactorRatingController.GetAllRatings)
	assetRoutes.GET("/asset-criticality-matrix-factor-ratings/:id", assetCriticalityMatrixFactorRatingController.GetRatingDetails)
	assetRoutes.POST("/asset-criticality-matrix-factor-ratings", middlewares.RolePermissionMiddleware(db, "asset_criticality", "modify"), assetCriticalityMatrixFactorRatingController.CreateRating)
	assetRoutes.PUT("/asset-criticality-matrix-factor-ratings/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality", "modify"), assetCriticalityMatrixFactorRatingController.UpdateRating)
	assetRoutes.POST("/asset-criticality-matrix-factor-ratings/delete", middlewares.RolePermissionMiddleware(db, "asset_criticality", "modify"), assetCriticalityMatrixFactorRatingController.DeleteRatings)

	// Asset Criticality Matrix Rating Scales
	assetRoutes.GET("/asset-criticality-matrix-rating-scales", assetCriticalityMatrixRatingScaleController.GetAllScales)
	assetRoutes.GET("/asset-criticality-matrix-rating-scales/:id", assetCriticalityMatrixRatingScaleController.GetScaleDetails)
	assetRoutes.POST("/asset-criticality-matrix-rating-scales", middlewares.RolePermissionMiddleware(db, "asset_criticality", "modify"), assetCriticalityMatrixRatingScaleController.CreateScale)
	assetRoutes.PUT("/asset-criticality-matrix-rating-scales/:id", middlewares.RolePermissionMiddleware(db, "asset_criticality", "modify"), assetCriticalityMatrixRatingScaleController.UpdateScale)
	assetRoutes.POST("/asset-criticality-matrix-rating-scales/delete", middlewares.RolePermissionMiddleware(db, "asset_criticality", "modify"), assetCriticalityMatrixRatingScaleController.DeleteScales)

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
	

}
