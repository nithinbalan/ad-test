package initializers

import (
	"jethur/grc/organizational-context/controllers"
	"jethur/grc/organizational-context/repositories"
	"jethur/grc/organizational-context/services"

	"gorm.io/gorm"
)

type OrganizationalContextModule struct {
	SwotPestelQuestionGroupController  *controllers.SwotPestelQuestionGroupController
	SwotPestelQuestionController       *controllers.SwotPestelQuestionController
	SwotPestelCategoryController       *controllers.SwotPestelCategoryController
	SwotPestelQuestionFactorController *controllers.SwotPestelQuestionFactorController
}

func NewOrganizationalContextModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,

) *OrganizationalContextModule {
	swotQuestionGroupRepo := repositories.NewSwotPestelQuestionGroupRepository(db)
	swotQuestionGroupSvc := services.NewSwotPestelQuestionGroupService(swotQuestionGroupRepo, sharedModule.ActivityLogService)
	swotQuestionGroupController := controllers.NewSwotPestelQuestionGroupController(swotQuestionGroupSvc)

	swotQuestionRepo := repositories.NewSwotPestelQuestionRepository(db)
	swotQuestionSvc := services.NewSwotPestelQuestionService(swotQuestionRepo, sharedModule.ActivityLogService)
	swotQuestionController := controllers.NewSwotPestelQuestionController(swotQuestionSvc)

	swotPestelCategoryRepo := repositories.NewSwotPestelCategoryRepository(db)
	swotPestelCategorySvc := services.NewSwotPestelCategoryService(swotPestelCategoryRepo)
	swotPestelCategoryController := controllers.NewSwotPestelCategoryController(swotPestelCategorySvc)

	swotQuestionFactorRepo := repositories.NewSwotPestelQuestionFactorRepository(db)
	swotQuestionFactorSvc := services.NewSwotPestelQuestionFactorService(swotQuestionFactorRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc)
	swotQuestionFactorController := controllers.NewSwotPestelQuestionFactorController(swotQuestionFactorSvc)

	return &OrganizationalContextModule{
		SwotPestelQuestionGroupController:  swotQuestionGroupController,
		SwotPestelQuestionController:       swotQuestionController,
		SwotPestelCategoryController:       swotPestelCategoryController,
		SwotPestelQuestionFactorController: swotQuestionFactorController,
	}
}
