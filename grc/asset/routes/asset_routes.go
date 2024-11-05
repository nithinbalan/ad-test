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

}
