package initializers

import (
	"jethur/grc/business-continuity/controllers"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/business-continuity/services"

	"gorm.io/gorm"
)

type BusinessContinuityModule struct {
	BusinessFunctionController *controllers.BusinessFunctionController

	BusinessFunctionSvc services.BusinessFunctionServiceInterface
}

func NewBusinessContinuityModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
) *BusinessContinuityModule {

	businessFunctionRepo := repositories.NewBusinessFunctionRepository(db)

	businessFunctionSvc := services.NewBusinessFunctionService(businessFunctionRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
	)
	businessFunctionController := controllers.NewBusinessFunctionController(businessFunctionSvc)

	return &BusinessContinuityModule{
		BusinessFunctionController: businessFunctionController,

		BusinessFunctionSvc: businessFunctionSvc,
	}
}
