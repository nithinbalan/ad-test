package initializers

import (
	"jethur/grc/strategy/controllers"
	"jethur/grc/strategy/repositories"
	"jethur/grc/strategy/services"

	"gorm.io/gorm"
)

type StrategyModule struct {
	SwotQuestionGroupController *controllers.SwotQuestionGroupController
	SwotQuestionController      *controllers.SwotQuestionController
}

func NewStrategyModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,

) *StrategyModule {
	swotQuestionGroupRepo := repositories.NewSwotQuestionGroupRepository(db)
	swotQuestionGroupSvc := services.NewSwotQuestionGroupService(swotQuestionGroupRepo, sharedModule.ActivityLogService)
	swotQuestionGroupController := controllers.NewSwotQuestionGroupController(swotQuestionGroupSvc)

	swotQuestionRepo := repositories.NewSwotQuestionRepository(db)
	swotQuestionSvc := services.NewSwotQuestionService(swotQuestionRepo, sharedModule.ActivityLogService)
	swotQuestionController := controllers.NewSwotQuestionController(swotQuestionSvc)

	return &StrategyModule{
		SwotQuestionGroupController: swotQuestionGroupController,
		SwotQuestionController:      swotQuestionController,
	}
}
