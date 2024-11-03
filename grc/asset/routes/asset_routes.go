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

}
