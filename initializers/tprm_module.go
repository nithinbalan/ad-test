package initializers

import (
	"jethur/grc/tprm/controllers"
	"jethur/grc/tprm/repositories"
	"jethur/grc/tprm/services"

	"gorm.io/gorm"
)

type TprmModule struct {
	ThirdPartyTypeController     *controllers.ThirdPartyTypeController
	ThirdPartyRegisterController *controllers.ThirdPartyRegisterController

	ThirdPartyTypeSvc     services.ThirdPartyTypeServiceInterface
	ThirdPartyRegisterSvc services.ThirdPartyRegisterServiceInterface
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
	)
	thirdPartyRegisterController := controllers.NewThirdPartyRegisterController(thirdPartyRegisterSvc)

	return &TprmModule{
		ThirdPartyTypeController:     thirdPartyTypeController,
		ThirdPartyRegisterController: thirdPartyRegisterController,

		ThirdPartyTypeSvc:     thirdPartyTypeSvc,
		ThirdPartyRegisterSvc: thirdPartyRegisterSvc,
	}
}
