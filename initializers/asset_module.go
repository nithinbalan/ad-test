package initializers

import (
	"jethur/grc/asset/controllers"
	"jethur/grc/asset/repositories"
	"jethur/grc/asset/services"

	"gorm.io/gorm"
)

type AssetModule struct {
	AssetTypeController     *controllers.AssetTypeController
	AssetCategoryController *controllers.AssetCategoryController
	AssetStatusController   *controllers.AssetStatusController
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

	return &AssetModule{
		AssetTypeController:     assetTypeController,
		AssetCategoryController: assetCategoryController,
		AssetStatusController:   assetStatusController,
	}
}
