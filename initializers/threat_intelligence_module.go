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
	ThreatIntelligenceLikelihoodController          *controllers.ThreatIntelligenceLikelihoodController
	ThreatIntelligenceSeverityController            *controllers.ThreatIntelligenceSeverityController
	ThreatIntelligenceRiskLevelController           *controllers.ThreatIntelligenceRiskLevelController
	ThreatIntelligenceRiskLevelValueController      *controllers.ThreatIntelligenceRiskLevelValueController
	ThreatIntelligencePlannedAdvisoryController     *controllers.ThreatIntelligencePlannedAdvisoryController
}

func NewThreatIntelligenceModule(db *gorm.DB, sharedModule *SharedModule, organizationModule *OrganizationModule, settingsModule *SettingsModule, commonModule *CommonModule) *ThreatIntelligenceModule {
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

	threatIntelligenceLikelihoodRepo := repositories.NewThreatIntelligenceLikelihoodRepository(db)
	threatIntelligenceLikelihoodSvc := services.NewThreatIntelligenceLikelihoodService(threatIntelligenceLikelihoodRepo, sharedModule.ActivityLogService)
	threatIntelligenceLikelihoodController := controllers.NewThreatIntelligenceLikelihoodController(threatIntelligenceLikelihoodSvc)

	threatIntelligenceSeverityRepo := repositories.NewThreatIntelligenceSeverityRepository(db)
	threatIntelligenceSeveritySvc := services.NewThreatIntelligenceSeverityService(threatIntelligenceSeverityRepo, sharedModule.ActivityLogService)
	threatIntelligenceSeverityController := controllers.NewThreatIntelligenceSeverityController(threatIntelligenceSeveritySvc)

	threatIntelligenceRiskLevelRepo := repositories.NewThreatIntelligenceRiskLevelRepository(db)
	threatIntelligenceRiskLevelSvc := services.NewThreatIntelligenceRiskLevelService(threatIntelligenceRiskLevelRepo, sharedModule.ActivityLogService)
	threatIntelligenceRiskLevelController := controllers.NewThreatIntelligenceRiskLevelController(threatIntelligenceRiskLevelSvc)

	threatIntelligenceRiskLevelValueRepo := repositories.NewThreatIntelligenceRiskLevelValueRepository(db)
	threatIntelligenceRiskLevelValueSvc := services.NewThreatIntelligenceRiskLevelValueService(
		threatIntelligenceRiskLevelValueRepo,
		threatIntelligenceLikelihoodSvc,
		threatIntelligenceSeveritySvc,
		threatIntelligenceRiskLevelSvc,
	)
	threatIntelligenceRiskLevelValueController := controllers.NewThreatIntelligenceRiskLevelValueController(threatIntelligenceRiskLevelValueSvc)

	threatIntelligencePlannedAdvisoryRepo := repositories.NewThreatIntelligencePlannedAdvisoryRepository(db)
	threatIntelligencePlannedAdvisorySvc := services.NewThreatIntelligencePlannedAdvisoryService(
		threatIntelligencePlannedAdvisoryRepo,
		sharedModule.ActivityLogService,
		commonModule.NotificationService,
		plannedAdvisorySourceSvc,
		plannedAdvisoryClassificationSvc,
		plannedAdvisoryTypeSvc,
		plannedAdvisoryImportanceSvc,
		plannedAdvisorySecurityClassificationSvc,
		plannedAdvisoryStatusSvc,
		plannedAdvisoryInternalStatusSvc,
		threatIntelligenceLikelihoodSvc,
		threatIntelligenceSeveritySvc,
		threatIntelligenceRiskLevelSvc,
		threatIntelligenceRiskLevelValueSvc,
		organizationModule.UserSvc,
		settingsModule.OrganizationUnitSvc,
		organizationModule.UserSettingSvc,
	)
	threatIntelligencePlannedAdvisoryController := controllers.NewThreatIntelligencePlannedAdvisoryController(threatIntelligencePlannedAdvisorySvc)

	return &ThreatIntelligenceModule{
		PlannedAdvisoryStatusController:                 plannedAdvisoryStatusController,
		PlannedAdvisorySourceController:                 plannedAdvisorySourceController,
		PlannedAdvisoryClassificationController:         plannedAdvisoryClassificationController,
		PlannedAdvisoryTypeController:                   plannedAdvisoryTypeController,
		PlannedAdvisoryImportanceController:             plannedAdvisoryImportanceController,
		PlannedAdvisorySecurityClassificationController: plannedAdvisorySecurityClassificationController,
		PlannedAdvisoryInternalStatusController:         plannedAdvisoryInternalStatusController,
		ThreatIntelligenceLikelihoodController:          threatIntelligenceLikelihoodController,
		ThreatIntelligenceSeverityController:            threatIntelligenceSeverityController,
		ThreatIntelligenceRiskLevelController:           threatIntelligenceRiskLevelController,
		ThreatIntelligenceRiskLevelValueController:      threatIntelligenceRiskLevelValueController,
		ThreatIntelligencePlannedAdvisoryController:     threatIntelligencePlannedAdvisoryController,
	}
}
