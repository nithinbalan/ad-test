package initializers

import (
	"jethur/grc/tprm/controllers"
	"jethur/grc/tprm/repositories"
	"jethur/grc/tprm/services"

	"gorm.io/gorm"
)

type TprmModule struct {
	ThirdPartyTypeController *controllers.ThirdPartyTypeController

	ThirdPartyTypeSvc services.ThirdPartyTypeServiceInterface
}

func NewTprmModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule,
) *TprmModule {
	thirdPartyTypeRepo := repositories.NewThirdPartyTypeRepository(db)
	thirdPartyTypeSvc := services.NewThirdPartyTypeService(thirdPartyTypeRepo, sharedModule.ActivityLogService)
	thirdPartyTypeController := controllers.NewThirdPartyTypeController(thirdPartyTypeSvc)

	return &TprmModule{
		ThirdPartyTypeController: thirdPartyTypeController,

		ThirdPartyTypeSvc: thirdPartyTypeSvc,
	}
}
