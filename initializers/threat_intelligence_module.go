package initializers

import (
	"jethur/grc/threat-intelligence/controllers"
	"jethur/grc/threat-intelligence/repositories"
	"jethur/grc/threat-intelligence/services"

	"gorm.io/gorm"
)

type ThreatIntelligenceModule struct {
	PlannedAdvisoryStatusController                 *controllers.PlannedAdvisoryStatusController
	PlannedAdvisorySourceController                 *controllers.PlannedAdvisorySourceController
	PlannedAdvisoryClassificationController         *controllers.PlannedAdvisoryClassificationController
	PlannedAdvisoryTypeController                   *controllers.PlannedAdvisoryTypeController
	PlannedAdvisoryImportanceController             *controllers.PlannedAdvisoryImportanceController
	PlannedAdvisorySecurityClassificationController *controllers.PlannedAdvisorySecurityClassificationController
	PlannedAdvisoryInternalStatusController         *controllers.PlannedAdvisoryInternalStatusController
}

func NewThreatIntelligenceModule(db *gorm.DB, sharedModule *SharedModule) *ThreatIntelligenceModule {
	plannedAdvisoryStatusRepo := repositories.NewPlannedAdvisoryStatusRepository(db)
	plannedAdvisoryStatusSvc := services.NewPlannedAdvisoryStatusService(plannedAdvisoryStatusRepo)
	plannedAdvisoryStatusController := controllers.NewPlannedAdvisoryStatusController(plannedAdvisoryStatusSvc)

	plannedAdvisorySourceRepo := repositories.NewPlannedAdvisorySourceRepository(db)
	plannedAdvisorySourceSvc := services.NewPlannedAdvisorySourceService(plannedAdvisorySourceRepo, sharedModule.ActivityLogService)
	plannedAdvisorySourceController := controllers.NewPlannedAdvisorySourceController(plannedAdvisorySourceSvc)

	plannedAdvisoryClassificationRepo := repositories.NewPlannedAdvisoryClassificationRepository(db)
	plannedAdvisoryClassificationSvc := services.NewPlannedAdvisoryClassificationService(plannedAdvisoryClassificationRepo, sharedModule.ActivityLogService)
	plannedAdvisoryClassificationController := controllers.NewPlannedAdvisoryClassificationController(plannedAdvisoryClassificationSvc)

	plannedAdvisoryTypeRepo := repositories.NewPlannedAdvisoryTypeRepository(db)
	plannedAdvisoryTypeSvc := services.NewPlannedAdvisoryTypeService(plannedAdvisoryTypeRepo, sharedModule.ActivityLogService)
	plannedAdvisoryTypeController := controllers.NewPlannedAdvisoryTypeController(plannedAdvisoryTypeSvc)

	plannedAdvisoryImportanceRepo := repositories.NewPlannedAdvisoryImportanceRepository(db)
	plannedAdvisoryImportanceSvc := services.NewPlannedAdvisoryImportanceService(plannedAdvisoryImportanceRepo, sharedModule.ActivityLogService)
	plannedAdvisoryImportanceController := controllers.NewPlannedAdvisoryImportanceController(plannedAdvisoryImportanceSvc)

	plannedAdvisorySecurityClassificationRepo := repositories.NewPlannedAdvisorySecurityClassificationRepository(db)
	plannedAdvisorySecurityClassificationSvc := services.NewPlannedAdvisorySecurityClassificationService(plannedAdvisorySecurityClassificationRepo, sharedModule.ActivityLogService)
	plannedAdvisorySecurityClassificationController := controllers.NewPlannedAdvisorySecurityClassificationController(plannedAdvisorySecurityClassificationSvc)

	plannedAdvisoryInternalStatusRepo := repositories.NewPlannedAdvisoryInternalStatusRepository(db)
	plannedAdvisoryInternalStatusSvc := services.NewPlannedAdvisoryInternalStatusService(plannedAdvisoryInternalStatusRepo, sharedModule.ActivityLogService)
	plannedAdvisoryInternalStatusController := controllers.NewPlannedAdvisoryInternalStatusController(plannedAdvisoryInternalStatusSvc)

	return &ThreatIntelligenceModule{
		PlannedAdvisoryStatusController:                 plannedAdvisoryStatusController,
		PlannedAdvisorySourceController:                 plannedAdvisorySourceController,
		PlannedAdvisoryClassificationController:         plannedAdvisoryClassificationController,
		PlannedAdvisoryTypeController:                   plannedAdvisoryTypeController,
		PlannedAdvisoryImportanceController:             plannedAdvisoryImportanceController,
		PlannedAdvisorySecurityClassificationController: plannedAdvisorySecurityClassificationController,
		PlannedAdvisoryInternalStatusController:         plannedAdvisoryInternalStatusController,
	}
}
