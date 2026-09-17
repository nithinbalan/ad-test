package initializers

import (
	"jethur/grc/incident/controllers"
	"jethur/grc/incident/repositories"
	"jethur/grc/incident/services"

	"gorm.io/gorm"
)

type IncidentModule struct {
	IncidentCategoryController                  *controllers.IncidentCategoryController
	IncidentSubCategoryController               *controllers.IncidentSubCategoryController
	IncidentTypeController                      *controllers.IncidentTypeController
	IncidentStatusController                    *controllers.IncidentStatusController
	IncidentRegisterController                  *controllers.IncidentRegisterController
	IncidentInvestigationStatusController       *controllers.IncidentInvestigationStatusController
	IncidentInvestigationController             *controllers.IncidentInvestigationController
	IncidentActionPlanController                *controllers.IncidentActionPlanController
	IncidentActionPlanStatusUpdateController    *controllers.IncidentActionPlanStatusUpdateController
	IncidentRootCauseAnalysisController         *controllers.IncidentRootCauseAnalysisController
	IncidentCommunicationLogController          *controllers.IncidentCommunicationLogController
	IncidentClosureController                   *controllers.IncidentClosureController
	IncidentClosureWorkflowController           *controllers.IncidentClosureWorkflowController
	IncidentClosureReviewController             *controllers.IncidentClosureReviewController
	IncidentClosureReviewCommentController      *controllers.IncidentClosureReviewCommentController
	IncidentPostReviewEffectivenessController   *controllers.IncidentPostReviewEffectivenessController
	IncidentPostReviewController                *controllers.IncidentPostReviewController
	IncidentAssociatedItemController            *controllers.IncidentAssociatedItemController
	IncidentTimelineController                  *controllers.IncidentTimelineController
	IncidentActionPlanCommentController         *controllers.IncidentActionPlanCommentController
	IncidentReportController                    *controllers.IncidentReportController
	IncidentResponseTeamController              *controllers.IncidentResponseTeamController
	IncidentImpactController                    *controllers.IncidentImpactController
	IncidentScopeController                     *controllers.IncidentScopeController
	IncidentSensitivityClassificationController *controllers.IncidentSensitivityClassificationController
	IncidentEvidenceTypeController              *controllers.IncidentEvidenceTypeController
	IncidentEvidenceSourceController            *controllers.IncidentEvidenceSourceController
	IncidentEvidenceController                  *controllers.IncidentEvidenceController
	IncidentSLAEscalationMatrixController       *controllers.IncidentSLAEscalationMatrixController
	IncidentWorkflowLevelController             *controllers.IncidentWorkflowLevelController

	IncidentRegisterSvc   services.IncidentRegisterServiceInterface
	IncidentActionPlanSvc services.IncidentActionPlanServiceInterface
}

func NewIncidentModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	complianceModule *ComplianceModule,
	riskModule *RiskModule,
	taskModule *TaskModule,
	commonModule *CommonModule,
	assetModule *AssetModule,
) *IncidentModule {
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

	incidentResponseTeamRepo := repositories.NewIncidentResponseTeamRepository(db)
	incidentResponseTeamSvc := services.NewIncidentResponseTeamService(incidentResponseTeamRepo, sharedModule.ActivityLogService, organizationModule.UserSvc)

	incidentTimelineRepo := repositories.NewIncidentTimelineRepository(db)
	incidentTimelineSvc := services.NewIncidentTimelineService(incidentTimelineRepo, incidentRegisterRepo)
	incidentTimelineController := controllers.NewIncidentTimelineController(incidentTimelineSvc)

	incidentRegisterSvc := services.NewIncidentRegisterService(incidentRegisterRepo,
		sharedModule.ActivityLogService,
		incidentCategorySvc,
		incidentTypeSvc,
		incidentStatusSvc,
		incidentResponseTeamSvc,
		masterModule.SeverityLevelSvc,
		organizationModule.ControlSvc,
		complianceModule.ComplianceRequirementSvc,
		riskModule.RiskRegisterSvc,
		organizationModule.ProcessSvc,
		complianceModule.SlaRegisterSvc,
		incidentTimelineSvc,
		organizationModule.UserSvc,
		organizationModule.UserSettingSvc,
		commonModule.NotificationService,
		assetModule.AssetRegisterSvc,
		settingsModule.SystemPreferenceSvc,
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
		incidentStatusSvc,
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
		organizationModule.UserSettingSvc,
		taskModule.TaskService,
		organizationModule.UserSvc, commonModule.NotificationService, commonModule.DelegateeUserService, masterModule.actionPlanMasterSvc,
	)
	incidentActionPlanController := controllers.NewIncidentActionPlanController(incidentActionPlanSvc)

	incidentActionPlanStatusUpdateRepo := repositories.NewIncidentActionPlanStatusUpdateRepository(db)
	incidentActionPlanStatusUpdateSvc := services.NewIncidentActionPlanStatusUpdateService(incidentActionPlanStatusUpdateRepo, sharedModule.ActivityLogService, incidentActionPlanSvc, masterModule.ActionPlanStatusSvc, incidentRegisterSvc, taskModule.TaskService, taskModule.TaskStatusService, organizationModule.UserSvc, commonModule.NotificationService)
	incidentActionPlanStatusUpdateController := controllers.NewIncidentActionPlanStatusUpdateController(incidentActionPlanStatusUpdateSvc)

	incidentActionPlanCommentRepository := repositories.NewIncidentActionPlanCommentRepository(db)
	incidentActionPlanCommentSvc := services.NewIncidentActionPlanCommentService(incidentActionPlanCommentRepository, organizationModule.UserSvc, sharedModule.ActivityLogService, incidentActionPlanSvc, commonModule.NotificationService)
	incidentActionPlanCommentController := controllers.NewIncidentActionPlanCommentController(incidentActionPlanCommentSvc)

	incidentRootCauseAnalysisRepo := repositories.NewIncidentRootCauseAnalysisRepository(db)
	incidentRootCauseAnalysisSvc := services.NewIncidentRootCauseAnalysisService(incidentRootCauseAnalysisRepo, sharedModule.ActivityLogService, incidentRegisterSvc, incidentTimelineSvc)
	incidentRootCauseAnalysisController := controllers.NewIncidentRootCauseAnalysisController(incidentRootCauseAnalysisSvc)

	incidentCommunicationLogRepo := repositories.NewIncidentCommunicationLogRepository(db)
	incidentCommunicationLogSvc := services.NewIncidentCommunicationLogService(incidentCommunicationLogRepo, sharedModule.ActivityLogService, incidentRegisterSvc, incidentTimelineSvc, organizationModule.UserSvc)
	incidentCommunicationLogController := controllers.NewIncidentCommunicationLogController(incidentCommunicationLogSvc)

	incidentClosureRepo := repositories.NewIncidentClosureRepository(db)
	incidentClosureSvc := services.NewIncidentClosureService(incidentClosureRepo, sharedModule.ActivityLogService, incidentRegisterSvc, incidentTimelineSvc, incidentStatusSvc)
	incidentClosureController := controllers.NewIncidentClosureController(incidentClosureSvc)

	incidentClosureWorkflowRepo := repositories.NewIncidentClosureWorkflowRepository(db)
	incidentClosureWorkflowSvc := services.NewIncidentClosureWorkflowService(
		incidentClosureWorkflowRepo,
		sharedModule.ActivityLogService,
		settingsModule.WorkflowSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		masterModule.WorkflowContextSvc,
		commonModule.NotificationService,
		organizationModule.OrganizationUnitsSvc,
	)
	incidentClosureWorkflowSvc.SetIncidentClosureService(incidentClosureSvc)
	incidentClosureWorkflowController := controllers.NewIncidentClosureWorkflowController(incidentClosureWorkflowSvc)

	incidentClosureReviewRepo := repositories.NewIncidentClosureReviewRepository(db)
	incidentClosureReviewSvc := services.NewIncidentClosureReviewService(
		incidentClosureReviewRepo,
		sharedModule.ActivityLogService,
		incidentClosureSvc,
		incidentClosureWorkflowSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
		incidentRegisterSvc,
		incidentStatusSvc,
		incidentActionPlanSvc,
		incidentRootCauseAnalysisSvc,
	)
	incidentClosureReviewController := controllers.NewIncidentClosureReviewController(incidentClosureReviewSvc)

	incidentClosureReviewCommentRepo := repositories.NewIncidentClosureReviewCommentRepository(db)
	incidentClosureReviewCommentSvc := services.NewIncidentClosureReviewCommentService(
		incidentClosureReviewCommentRepo,
		organizationModule.UserSvc,
		sharedModule.ActivityLogService,
		incidentClosureSvc,
		masterModule.CommentStatusSvc,
	)
	incidentClosureReviewCommentController := controllers.NewIncidentClosureReviewCommentController(incidentClosureReviewCommentSvc)

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
	incidentAssociatedItemSvc := services.NewIncidentAssociatedItemService(incidentAssociatedItemRepo, incidentRegisterSvc, complianceModule.ComplianceClausevc, commonModule.AssociatedItemValidationService)
	incidentAssociatedItemController := controllers.NewIncidentAssociatedItemController(incidentAssociatedItemSvc)

	incidentReportrepo := repositories.NewIncidentReportRepository(db)
	incidentReportSvc := services.NewIncidentReportService(incidentReportrepo, settingsModule.ReportTemplateSettingSvc, incidentRootCauseAnalysisSvc, incidentActionPlanSvc, sharedModule.ActivityLogService, organizationModule.UserSettingSvc, settingsModule.SystemPreferenceSvc)
	incidentReportController := controllers.NewIncidentReportController(incidentReportSvc)

	incidentResponseTeamController := controllers.NewIncidentResponseTeamController(incidentResponseTeamSvc)

	incidentImpactRepo := repositories.NewIncidentImpactRepository(db)
	incidentImpactSvc := services.NewIncidentImpactService(incidentImpactRepo, sharedModule.ActivityLogService)
	incidentImpactController := controllers.NewIncidentImpactController(incidentImpactSvc)

	incidentScopeRepo := repositories.NewIncidentScopeRepository(db)
	incidentScopeSvc := services.NewIncidentScopeService(incidentScopeRepo, sharedModule.ActivityLogService)
	incidentScopeController := controllers.NewIncidentScopeController(incidentScopeSvc)

	incidentSensitivityClassificationRepo := repositories.NewIncidentSensitivityClassificationRepository(db)
	incidentSensitivityClassificationSvc := services.NewIncidentSensitivityClassificationService(incidentSensitivityClassificationRepo, sharedModule.ActivityLogService)
	incidentSensitivityClassificationController := controllers.NewIncidentSensitivityClassificationController(incidentSensitivityClassificationSvc)

	incidentEvidenceTypeRepo := repositories.NewIncidentEvidenceTypeRepository(db)
	incidentEvidenceTypeSvc := services.NewIncidentEvidenceTypeService(incidentEvidenceTypeRepo, sharedModule.ActivityLogService)
	incidentEvidenceTypeController := controllers.NewIncidentEvidenceTypeController(incidentEvidenceTypeSvc)

	incidentEvidenceSourceRepo := repositories.NewIncidentEvidenceSourceRepository(db)
	incidentEvidenceSourceSvc := services.NewIncidentEvidenceSourceService(incidentEvidenceSourceRepo, sharedModule.ActivityLogService)
	incidentEvidenceSourceController := controllers.NewIncidentEvidenceSourceController(incidentEvidenceSourceSvc)

	incidentEvidenceRepo := repositories.NewIncidentEvidenceRepository(db)
	incidentEvidenceSvc := services.NewIncidentEvidenceService(incidentEvidenceRepo, sharedModule.ActivityLogService, incidentRegisterSvc, incidentEvidenceSourceSvc, incidentEvidenceTypeSvc)
	incidentEvidenceController := controllers.NewIncidentEvidenceController(incidentEvidenceSvc)

	incidentSLAEscalationMatrixRepo := repositories.NewIncidentSLAEscalationMatrixRepository(db)
	incidentSLAEscalationMatrixSvc := services.NewIncidentSLAEscalationMatrixService(incidentSLAEscalationMatrixRepo, sharedModule.ActivityLogService)
	incidentSLAEscalationMatrixController := controllers.NewIncidentSLAEscalationMatrixController(incidentSLAEscalationMatrixSvc)

	incidentWorkflowLevelController := controllers.NewIncidentWorkflowLevelController(settingsModule.WorkflowSvc)

	return &IncidentModule{
		IncidentCategoryController:                  incidentCategoryController,
		IncidentSubCategoryController:               incidentSubCategoryController,
		IncidentTypeController:                      incidentTypeController,
		IncidentStatusController:                    incidentStatusController,
		IncidentRegisterController:                  incidentRegisterController,
		IncidentInvestigationStatusController:       incidentInvestigationStatusController,
		IncidentInvestigationController:             incidentInvestigationController,
		IncidentActionPlanController:                incidentActionPlanController,
		IncidentActionPlanStatusUpdateController:    incidentActionPlanStatusUpdateController,
		IncidentRootCauseAnalysisController:         incidentRootCauseAnalysisController,
		IncidentCommunicationLogController:          incidentCommunicationLogController,
		IncidentClosureController:                   incidentClosureController,
		IncidentClosureWorkflowController:           incidentClosureWorkflowController,
		IncidentClosureReviewController:             incidentClosureReviewController,
		IncidentClosureReviewCommentController:      incidentClosureReviewCommentController,
		IncidentPostReviewEffectivenessController:   incidentPostReviewEffectivenessController,
		IncidentPostReviewController:                incidentPostReviewController,
		IncidentAssociatedItemController:            incidentAssociatedItemController,
		IncidentTimelineController:                  incidentTimelineController,
		IncidentActionPlanCommentController:         incidentActionPlanCommentController,
		IncidentReportController:                    incidentReportController,
		IncidentResponseTeamController:              incidentResponseTeamController,
		IncidentImpactController:                    incidentImpactController,
		IncidentScopeController:                     incidentScopeController,
		IncidentSensitivityClassificationController: incidentSensitivityClassificationController,
		IncidentEvidenceTypeController:              incidentEvidenceTypeController,
		IncidentEvidenceSourceController:            incidentEvidenceSourceController,
		IncidentEvidenceController:                  incidentEvidenceController,
		IncidentSLAEscalationMatrixController:       incidentSLAEscalationMatrixController,
		IncidentWorkflowLevelController:             incidentWorkflowLevelController,

		IncidentRegisterSvc:   incidentRegisterSvc,
		IncidentActionPlanSvc: incidentActionPlanSvc,
	}
}
