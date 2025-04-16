package initializers

import (
	"jethur/grc/strategy/controllers"
	"jethur/grc/strategy/repositories"
	"jethur/grc/strategy/services"

	"gorm.io/gorm"
)

type StrategyModule struct {
	StrategyProfileController           *controllers.StrategyProfileController
	StrategyProfileFocusAreaController  *controllers.StrategyProfileFocusAreaController
	StrategyProfileGoalController       *controllers.StrategyProfileGoalController
	StrategyProfileObjectiveController  *controllers.StrategyProfileObjectiveController
	StrategyProfileInitiativeController *controllers.StrategyProfileInitiativeController
}

func NewStrategyModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule,

) *StrategyModule {

	strategyProfileRepo := repositories.NewStrategyProfileRepository(db)
	strategyProfileSvc := services.NewStrategyProfileService(strategyProfileRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, settingsModule.OrganizationUnitSvc)
	strategyProfileController := controllers.NewStrategyProfileController(strategyProfileSvc)

	strategyProfileFocusAreaRepo := repositories.NewStrategyProfileFocusAreaRepository(db)
	strategyProfileFocusAreaSvc := services.NewStrategyProfileFocusAreaService(strategyProfileFocusAreaRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc)
	strategyProfileFocusAreaController := controllers.NewStrategyProfileFocusAreaController(strategyProfileFocusAreaSvc)

	strategyProfileGoalRepo := repositories.NewStrategyProfileGoalRepository(db)
	strategyProfileGoalSvc := services.NewStrategyProfileGoalService(strategyProfileGoalRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc)
	strategyProfileGoalController := controllers.NewStrategyProfileGoalController(strategyProfileGoalSvc)

	strategyProfileObjectiveRepo := repositories.NewStrategyProfileObjectiveRepository(db)
	strategyProfileObjectiveSvc := services.NewStrategyProfileObjectiveService(strategyProfileObjectiveRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc)
	strategyProfileObjectiveController := controllers.NewStrategyProfileObjectiveController(strategyProfileObjectiveSvc)

	strategyProfileInitiativeRepo := repositories.NewStrategyProfileInitiativeRepository(db)
	strategyProfileInitiativeSvc := services.NewStrategyProfileInitiativeService(strategyProfileInitiativeRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, strategyProfileSvc)
	strategyProfileInitiativeController := controllers.NewStrategyProfileInitiativeController(strategyProfileInitiativeSvc)

	return &StrategyModule{
		StrategyProfileController:           strategyProfileController,
		StrategyProfileFocusAreaController:  strategyProfileFocusAreaController,
		StrategyProfileGoalController:       strategyProfileGoalController,
		StrategyProfileObjectiveController:  strategyProfileObjectiveController,
		StrategyProfileInitiativeController: strategyProfileInitiativeController,
	}
}
