package initializers

import (
	"jethur/grc/incident/controllers"
	"jethur/grc/incident/repositories"
	"jethur/grc/incident/services"

	"gorm.io/gorm"
)

type IncidentModule struct {
	IncidentCategoryController               *controllers.IncidentCategoryController
	IncidentSubCategoryController            *controllers.IncidentSubCategoryController
	IncidentTypeController                   *controllers.IncidentTypeController
	IncidentStatusController                 *controllers.IncidentStatusController
	IncidentRegisterController               *controllers.IncidentRegisterController
	IncidentInvestigationStatusController    *controllers.IncidentInvestigationStatusController
	IncidentInvestigationController          *controllers.IncidentInvestigationController
	IncidentActionPlanController             *controllers.IncidentActionPlanController
	IncidentActionPlanStatusUpdateController *controllers.IncidentActionPlanStatusUpdateController
}

func NewIncidentModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	complianceModule *ComplianceModule) *IncidentModule {

	incidentCategoryRepo := repositories.NewIncidentCategoryRepository(db)
	incidentCategorySvc := services.NewIncidentCategoryService(incidentCategoryRepo, sharedModule.ActivityLogService)
	incidentCategoryController := controllers.NewIncidentCategoryController(incidentCategorySvc)

	incidentSubCategoryRepo := repositories.NewIncidentSubCategoryRepository(db)
	incidentSubCategorySvc := services.NewIncidentSubCategoryService(incidentSubCategoryRepo, sharedModule.ActivityLogService, incidentCategorySvc)
	incidentSubCategoryController := controllers.NewIncidentSubCategoryController(incidentSubCategorySvc)

	incidentTypeRepo := repositories.NewIncidentTypeRepository(db)
	incidentTypeSvc := services.NewIncidentTypeService(incidentTypeRepo, sharedModule.ActivityLogService)
	incidentTypeController := controllers.NewIncidentTypeController(incidentTypeSvc)

	incidentStatusRepo := repositories.NewIncidentStatusRepository(db)
	incidentStatusSvc := services.NewIncidentStatusService(incidentStatusRepo)
	incidentStatusController := controllers.NewIncidentStatusController(incidentStatusSvc)

	incidentRegisterRepo := repositories.NewIncidentRegisterRepository(db)
	incidentRegisterSvc := services.NewIncidentRegisterService(incidentRegisterRepo,
		sharedModule.ActivityLogService,
		incidentCategorySvc,
		incidentSubCategorySvc,
		incidentTypeSvc,
		incidentStatusSvc,
		masterModule.SeverityLevelSvc,
		organizationModule.ControlSvc,
		complianceModule.ComplianceRequirementSvc)
	incidentRegisterController := controllers.NewIncidentRegisterController(incidentRegisterSvc)

	incidentInvestigationStatusRepo := repositories.NewIncidentInvestigationStatusRepository(db)
	incidentInvestigationStatusSvc := services.NewIncidentInvestigationStatusService(incidentInvestigationStatusRepo)
	incidentInvestigationStatusController := controllers.NewIncidentInvestigationStatusController(incidentInvestigationStatusSvc)

	incidentInvestigationRepo := repositories.NewIncidentInvestigationRepository(db)
	incidentInvestigationSvc := services.NewIncidentInvestigationService(
		incidentInvestigationRepo,
		sharedModule.ActivityLogService,
		incidentRegisterSvc,
		organizationModule.UserSvc,
		incidentInvestigationStatusSvc,
	)
	incidentInvestigationController := controllers.NewIncidentInvestigationController(incidentInvestigationSvc)

	incidentActionPlanRepo := repositories.NewIncidentActionPlanRepository(db)
	incidentActionPlanSvc := services.NewIncidentActionPlanService(
		incidentActionPlanRepo,
		sharedModule.ActivityLogService,
		incidentRegisterSvc,
		masterModule.ActionPlanStatusSvc,
		masterModule.SeverityLevelSvc,
	)
	incidentActionPlanController := controllers.NewIncidentActionPlanController(incidentActionPlanSvc)

	incidentActionPlanStatusUpdateRepo := repositories.NewIncidentActionPlanStatusUpdateRepository(db)
	incidentActionPlanStatusUpdateSvc := services.NewIncidentActionPlanStatusUpdateService(incidentActionPlanStatusUpdateRepo, sharedModule.ActivityLogService, incidentActionPlanSvc, masterModule.ActionPlanStatusSvc)
	incidentActionPlanStatusUpdateController := controllers.NewIncidentActionPlanStatusUpdateController(incidentActionPlanStatusUpdateSvc)

	return &IncidentModule{
		IncidentCategoryController:               incidentCategoryController,
		IncidentSubCategoryController:            incidentSubCategoryController,
		IncidentTypeController:                   incidentTypeController,
		IncidentStatusController:                 incidentStatusController,
		IncidentRegisterController:               incidentRegisterController,
		IncidentInvestigationStatusController:    incidentInvestigationStatusController,
		IncidentInvestigationController:          incidentInvestigationController,
		IncidentActionPlanController:             incidentActionPlanController,
		IncidentActionPlanStatusUpdateController: incidentActionPlanStatusUpdateController,
	}
}
