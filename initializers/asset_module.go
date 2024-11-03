package initializers

import (
	"jethur/grc/asset/controllers"
	"jethur/grc/asset/repositories"
	"jethur/grc/asset/services"

	"gorm.io/gorm"
)

type AssetModule struct {
	AssetTypeController                          *controllers.AssetTypeController
	AssetCategoryController                      *controllers.AssetCategoryController
	AssetStatusController                        *controllers.AssetStatusController
	AssetCriticalityMatrixFactorController       *controllers.AssetCriticalityMatrixFactorController
	AssetCriticalityMatrixFactorRatingController *controllers.AssetCriticalityMatrixFactorRatingController
	AssetCriticalityRatingScaleController        *controllers.AssetCriticalityRatingScaleController
}

func NewAssetModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule) *AssetModule {

	assetTypeRepo := repositories.NewAssetTypeRepository(db)
	assetTypeSvc := services.NewAssetTypeService(assetTypeRepo, sharedModule.ActivityLogService)
	assetTypeController := controllers.NewAssetTypeController(assetTypeSvc)

	assetCategoryRepo := repositories.NewAssetCategoryRepository(db)
	assetCategorySvc := services.NewAssetCategoryService(assetCategoryRepo, assetTypeSvc, sharedModule.ActivityLogService)
	assetCategoryController := controllers.NewAssetCategoryController(assetCategorySvc)

	assetStatusRepo := repositories.NewAssetStatusRepository(db)
	assetStatusSvc := services.NewAssetStatusService(assetStatusRepo)
	assetStatusController := controllers.NewAssetStatusController(assetStatusSvc)

	assetCriticalityMatrixFactorRepo := repositories.NewAssetCriticalityMatrixFactorRepository(db)
	assetCriticalityMatrixFactorSvc := services.NewAssetCriticalityMatrixFactorService(assetCriticalityMatrixFactorRepo, sharedModule.ActivityLogService)
	assetCriticalityMatrixFactorController := controllers.NewAssetCriticalityMatrixFactorController(assetCriticalityMatrixFactorSvc)

	assetCriticalityMatrixFactorRatingRepo := repositories.NewAssetCriticalityMatrixFactorRatingRepository(db)
	assetCriticalityMatrixFactorRatingSvc := services.NewAssetCriticalityMatrixFactorRatingService(assetCriticalityMatrixFactorRatingRepo, sharedModule.ActivityLogService)
	assetCriticalityMatrixFactorRatingController := controllers.NewAssetCriticalityMatrixFactorRatingController(assetCriticalityMatrixFactorRatingSvc)

	assetCriticalityRatingScaleRepo := repositories.NewAssetCriticalityRatingScaleRepository(db)
	assetCriticalityRatingScaleSvc := services.NewAssetCriticalityRatingScaleService(assetCriticalityRatingScaleRepo, sharedModule.ActivityLogService)
	assetCriticalityRatingScaleController := controllers.NewAssetCriticalityRatingScaleController(assetCriticalityRatingScaleSvc)

	return &AssetModule{
		AssetTypeController:                          assetTypeController,
		AssetCategoryController:                      assetCategoryController,
		AssetStatusController:                        assetStatusController,
		AssetCriticalityMatrixFactorController:       assetCriticalityMatrixFactorController,
		AssetCriticalityMatrixFactorRatingController: assetCriticalityMatrixFactorRatingController,
		AssetCriticalityRatingScaleController:        assetCriticalityRatingScaleController,
	}
}
