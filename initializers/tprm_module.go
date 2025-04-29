package initializers

import (
	"jethur/grc/tprm/controllers"
	"jethur/grc/tprm/repositories"
	"jethur/grc/tprm/services"

	"gorm.io/gorm"
)

type TprmModule struct {
	ThirdPartyTypeController             *controllers.ThirdPartyTypeController
	ThirdPartyRegisterController         *controllers.ThirdPartyRegisterController
	ThirdPartyQuestionCategoryController *controllers.ThirdPartyQuestionCategoryController
	ThirdPartyQuestionTypeController     *controllers.ThirdPartyQuestionTypeController

	ThirdPartyTypeSvc             services.ThirdPartyTypeServiceInterface
	ThirdPartyRegisterSvc         services.ThirdPartyRegisterServiceInterface
	ThirdPartyQuestionCategorySvc services.ThirdPartyQuestionCategoryServiceInterface
	ThirdPartyQuestionTypeSvc     services.ThirdPartyQuestionTypeServiceInterface
}

func NewTprmModule(
	db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule,
) *TprmModule {

	// === Third Party Type ===
	thirdPartyTypeRepo := repositories.NewThirdPartyTypeRepository(db)
	thirdPartyTypeSvc := services.NewThirdPartyTypeService(thirdPartyTypeRepo, sharedModule.ActivityLogService)
	thirdPartyTypeController := controllers.NewThirdPartyTypeController(thirdPartyTypeSvc)

	// === Third Party Register ===
	thirdPartyRegisterRepo := repositories.NewThirdPartyRegisterRepository(db)
	thirdPartyRegisterSvc := services.NewThirdPartyRegisterService(
		thirdPartyRegisterRepo,
		sharedModule.ActivityLogService,
		masterModule.SeverityLevelSvc,
		organizationModule.UserSvc,
		settingsModule.OrganizationUnitSvc,
		thirdPartyTypeSvc,
		masterModule.CountrySvc,
		organizationModule.UserSettingSvc,
	)
	thirdPartyRegisterController := controllers.NewThirdPartyRegisterController(thirdPartyRegisterSvc)

	thirdPartyQuestionCategoryRepo := repositories.NewThirdPartyQuestionCategoryRepository(db)
	thirdPartyQuestionCategorySvc := services.NewThirdPartyQuestionCategoryService(thirdPartyQuestionCategoryRepo, sharedModule.ActivityLogService)
	thirdPartyQuestionCategoryController := controllers.NewThirdPartyQuestionCategoryController(thirdPartyQuestionCategorySvc)

	thirdPartyQuestionTypeRepo := repositories.NewThirdPartyQuestionTypeRepository(db)
	thirdPartyQuestionTypeSvc := services.NewThirdPartyQuestionTypeService(thirdPartyQuestionTypeRepo)
	thirdPartyQuestionTypeController := controllers.NewThirdPartyQuestionTypeController(thirdPartyQuestionTypeSvc)

	return &TprmModule{
		ThirdPartyTypeController:             thirdPartyTypeController,
		ThirdPartyRegisterController:         thirdPartyRegisterController,
		ThirdPartyQuestionCategoryController: thirdPartyQuestionCategoryController,
		ThirdPartyQuestionTypeController:     thirdPartyQuestionTypeController,

		ThirdPartyTypeSvc:             thirdPartyTypeSvc,
		ThirdPartyRegisterSvc:         thirdPartyRegisterSvc,
		ThirdPartyQuestionCategorySvc: thirdPartyQuestionCategorySvc,
		ThirdPartyQuestionTypeSvc:     thirdPartyQuestionTypeSvc,
	}
}
