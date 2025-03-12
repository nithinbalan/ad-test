package initializers

import (
	"jethur/grc/strategy/controllers"
	"jethur/grc/strategy/repositories"
	"jethur/grc/strategy/services"

	"gorm.io/gorm"
)

type StrategyModule struct {
	SwotPestelQuestionGroupController *controllers.SwotPestelQuestionGroupController
	SwotPestelQuestionController      *controllers.SwotPestelQuestionController
}

func NewStrategyModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,

) *StrategyModule {
	swotQuestionGroupRepo := repositories.NewSwotPestelQuestionGroupRepository(db)
	swotQuestionGroupSvc := services.NewSwotPestelQuestionGroupService(swotQuestionGroupRepo, sharedModule.ActivityLogService)
	swotQuestionGroupController := controllers.NewSwotPestelQuestionGroupController(swotQuestionGroupSvc)

	swotQuestionRepo := repositories.NewSwotPestelQuestionRepository(db)
	swotQuestionSvc := services.NewSwotPestelQuestionService(swotQuestionRepo, sharedModule.ActivityLogService)
	swotQuestionController := controllers.NewSwotPestelQuestionController(swotQuestionSvc)

	return &StrategyModule{
		SwotPestelQuestionGroupController: swotQuestionGroupController,
		SwotPestelQuestionController:      swotQuestionController,
	}
}
