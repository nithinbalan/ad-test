package initializers

import (
	"jethur/grc/incident/controllers"
	"jethur/grc/incident/repositories"
	"jethur/grc/incident/services"

	"gorm.io/gorm"
)

type IncidentModule struct {
	IncidentCategoryController                *controllers.IncidentCategoryController
	IncidentSubCategoryController             *controllers.IncidentSubCategoryController
	IncidentTypeController                    *controllers.IncidentTypeController
	IncidentStatusController                  *controllers.IncidentStatusController
	IncidentRegisterController                *controllers.IncidentRegisterController
	IncidentInvestigationStatusController     *controllers.IncidentInvestigationStatusController
	IncidentInvestigationController           *controllers.IncidentInvestigationController
	IncidentActionPlanController              *controllers.IncidentActionPlanController
	IncidentActionPlanStatusUpdateController  *controllers.IncidentActionPlanStatusUpdateController
	IncidentRootCauseAnalysisController       *controllers.IncidentRootCauseAnalysisController
	IncidentCommunicationLogController        *controllers.IncidentCommunicationLogController
	IncidentClosureController                 *controllers.IncidentClosureController
	IncidentPostReviewEffectivenessController *controllers.IncidentPostReviewEffectivenessController
	IncidentPostReviewController              *controllers.IncidentPostReviewController
	IncidentAssociatedItemController          *controllers.IncidentAssociatedItemController
	IncidentTimelineController                *controllers.IncidentTimelineController

	IncidentRegisterSvc services.IncidentRegisterServiceInterface
}

func NewIncidentModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	complianceModule *ComplianceModule,
	riskModule *RiskModule) *IncidentModule {

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

	incidentTimelineRepo := repositories.NewIncidentTimelineRepository(db)
	incidentTimelineSvc := services.NewIncidentTimelineService(incidentTimelineRepo, incidentRegisterRepo)
	incidentTimelineController := controllers.NewIncidentTimelineController(incidentTimelineSvc)

	incidentRegisterSvc := services.NewIncidentRegisterService(incidentRegisterRepo,
		sharedModule.ActivityLogService,
		incidentCategorySvc,
		incidentTypeSvc,
		incidentStatusSvc,
		masterModule.SeverityLevelSvc,
		organizationModule.ControlSvc,
		complianceModule.ComplianceRequirementSvc,
		riskModule.RiskRegisterSvc,
		organizationModule.ProcessSvc,
		complianceModule.SlaRegisterSvc,
		incidentTimelineSvc,
	)
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
		incidentTimelineSvc,
	)
	incidentInvestigationController := controllers.NewIncidentInvestigationController(incidentInvestigationSvc)

	incidentActionPlanRepo := repositories.NewIncidentActionPlanRepository(db)
	incidentActionPlanSvc := services.NewIncidentActionPlanService(
		incidentActionPlanRepo,
		sharedModule.ActivityLogService,
		incidentRegisterSvc,
		masterModule.ActionPlanStatusSvc,
		masterModule.SeverityLevelSvc,
		incidentTimelineSvc,
	)
	incidentActionPlanController := controllers.NewIncidentActionPlanController(incidentActionPlanSvc)

	incidentActionPlanStatusUpdateRepo := repositories.NewIncidentActionPlanStatusUpdateRepository(db)
	incidentActionPlanStatusUpdateSvc := services.NewIncidentActionPlanStatusUpdateService(incidentActionPlanStatusUpdateRepo, sharedModule.ActivityLogService, incidentActionPlanSvc, masterModule.ActionPlanStatusSvc)
	incidentActionPlanStatusUpdateController := controllers.NewIncidentActionPlanStatusUpdateController(incidentActionPlanStatusUpdateSvc)

	incidentRootCauseAnalysisRepo := repositories.NewIncidentRootCauseAnalysisRepository(db)
	incidentRootCauseAnalysisSvc := services.NewIncidentRootCauseAnalysisService(incidentRootCauseAnalysisRepo, sharedModule.ActivityLogService, incidentRegisterSvc, incidentTimelineSvc)
	incidentRootCauseAnalysisController := controllers.NewIncidentRootCauseAnalysisController(incidentRootCauseAnalysisSvc)

	incidentCommunicationLogRepo := repositories.NewIncidentCommunicationLogRepository(db)
	incidentCommunicationLogSvc := services.NewIncidentCommunicationLogService(incidentCommunicationLogRepo, sharedModule.ActivityLogService, incidentRegisterSvc, incidentTimelineSvc)
	incidentCommunicationLogController := controllers.NewIncidentCommunicationLogController(incidentCommunicationLogSvc)

	incidentClosureRepo := repositories.NewIncidentClosureRepository(db)
	incidentClosureSvc := services.NewIncidentClosureService(incidentClosureRepo, sharedModule.ActivityLogService, incidentRegisterSvc, incidentTimelineSvc, incidentStatusSvc)
	incidentClosureController := controllers.NewIncidentClosureController(incidentClosureSvc)

	incidentPostReviewEffectivenessRepo := repositories.NewIncidentPostReviewEffectivenessRepository(db)
	incidentPostReviewEffectivenessSvc := services.NewIncidentPostReviewEffectivenessService(incidentPostReviewEffectivenessRepo)
	incidentPostReviewEffectivenessController := controllers.NewIncidentPostReviewEffectivenessController(incidentPostReviewEffectivenessSvc)

	incidentPostReviewRepo := repositories.NewIncidentPostReviewRepository(db)
	incidentPostReviewSvc := services.NewIncidentPostReviewService(
		incidentPostReviewRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		incidentPostReviewEffectivenessSvc,
		incidentRegisterSvc,
		incidentTimelineSvc,
	)
	incidentPostReviewController := controllers.NewIncidentPostReviewController(incidentPostReviewSvc)

	incidentAssociatedItemRepo := repositories.NewIncidentAssociatedItemRepository(db)
	incidentAssociatedItemSvc := services.NewIncidentAssociatedItemService(incidentAssociatedItemRepo, incidentRegisterSvc)
	incidentAssociatedItemController := controllers.NewIncidentAssociatedItemController(incidentAssociatedItemSvc)

	return &IncidentModule{
		IncidentCategoryController:                incidentCategoryController,
		IncidentSubCategoryController:             incidentSubCategoryController,
		IncidentTypeController:                    incidentTypeController,
		IncidentStatusController:                  incidentStatusController,
		IncidentRegisterController:                incidentRegisterController,
		IncidentInvestigationStatusController:     incidentInvestigationStatusController,
		IncidentInvestigationController:           incidentInvestigationController,
		IncidentActionPlanController:              incidentActionPlanController,
		IncidentActionPlanStatusUpdateController:  incidentActionPlanStatusUpdateController,
		IncidentRootCauseAnalysisController:       incidentRootCauseAnalysisController,
		IncidentCommunicationLogController:        incidentCommunicationLogController,
		IncidentClosureController:                 incidentClosureController,
		IncidentPostReviewEffectivenessController: incidentPostReviewEffectivenessController,
		IncidentPostReviewController:              incidentPostReviewController,
		IncidentAssociatedItemController:          incidentAssociatedItemController,
		IncidentTimelineController:                incidentTimelineController,

		IncidentRegisterSvc: incidentRegisterSvc,
	}
}
