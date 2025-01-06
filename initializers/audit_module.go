package initializers

import (
	"jethur/grc/audit/controllers"
	"jethur/grc/audit/repositories"
	"jethur/grc/audit/services"

	"gorm.io/gorm"
)

type AuditModule struct {
	AuditStatusController                        *controllers.AuditStatusController
	AuditorRiskRatingController                  *controllers.AuditorRiskRatingController
	AuditUniverseController                      *controllers.AuditUniverseController
	AuditableItemController                      *controllers.AuditableItemController
	AuditTeamController                          *controllers.AuditTeamController
	AuditCategoryController                      *controllers.AuditCategoryController
	AuditPlanController                          *controllers.AuditPlanController
	AuditPlanChecklistController                 *controllers.AuditPlanChecklistController
	AuditPlanToolController                      *controllers.AuditPlanToolController
	AuditPlanAuditableItemController             *controllers.AuditPlanAuditableItemController
	AuditScheduleStatusController                *controllers.AuditScheduleStatusController
	AuditActivityTypeController                  *controllers.AuditActivityTypeController
	AuditPlanBudgetController                    *controllers.AuditPlanBudgetController
	AuditPlanScheduleController                  *controllers.AuditPlanScheduleController
	AuditEngagementTypeController                *controllers.AuditEngagementTypeController
	AuditPlanStakeholderEngagementController     *controllers.AuditPlanStakeholderEngagementController
	AuditController                              *controllers.AuditController
	AuditWorkpaperController                     *controllers.AuditWorkpaperController
	AuditMeetingController                       *controllers.AuditMeetingController
	AuditTimelineController                      *controllers.AuditTimelineController
	AuditFindingStatusController                 *controllers.AuditFindingStatusController
	AuditDraftFindingStatusController            *controllers.AuditDraftFindingStatusController
	AuditDraftReportStatusController             *controllers.AuditDraftReportStatusController
	AuditFindingCategoryController               *controllers.AuditFindingCategoryController
	AuditFindingController                       *controllers.AuditFindingController
	AuditDraftFindingController                  *controllers.AuditDraftFindingController
	AuditDocumentRequestStatusController         *controllers.AuditDocumentRequestStatusController
	AuditDocumentRequestController               *controllers.AuditDocumentRequestController
	AuditFindingRootCauseAnalysisController      *controllers.AuditFindingRootCauseAnalysisController
	AuditFindingActionPlanController             *controllers.AuditFindingActionPlanController
	AuditFindingActionPlanStatusUpdateController *controllers.AuditFindingActionPlanStatusUpdateController
	AuditDocumentRequestConversationController   *controllers.AuditDocumentRequestConversationController
	AuditChecklistController                     *controllers.AuditChecklistController
	AuditScheduleController                      *controllers.AuditScheduleController
	AuditDraftReportController                   *controllers.AuditDraftReportController
	RCMController                                *controllers.RCMController
	AuditDraftFindingConversationController      *controllers.AuditDraftFindingConversationController
	AuditDraftFindingActionPlanController        *controllers.AuditDraftFindingActionPlanController
	AuditScheduleStatusUpdateController          *controllers.AuditScheduleStatusUpdateController
	AuditDraftReportReviewController             *controllers.AuditDraftReportReviewController
	AuditDraftReportCommentController            *controllers.AuditDraftReportCommentController
	AuditPlanCommencementLetterStatusController  *controllers.AuditPlanCommencementLetterStatusController
	AuditPlanCommencementLetterController        *controllers.AuditPlanCommencementLetterController
	AuditFinalReportStatusController             *controllers.AuditFinalReportStatusController
	AuditFinalReportController                   *controllers.AuditFinalReportController
	AuditFinalReportSectionController            *controllers.AuditFinalReportSectionController
	AuditReportWorkflowController                *controllers.AuditReportWorkflowController
	AuditReportReviewController                  *controllers.AuditReportReviewController

	AuditPlanSvc                          services.AuditPlanServiceInterface
	AuditSvc                              services.AuditServiceInterface
	AuditFindingSvc                       services.AuditFindingServiceInterface
	AuditDocumentRequestSvc               services.AuditDocumentRequestServiceInterface
	AuditDocumentRequestConversationSvc   services.AuditDocumentRequestConversationServiceInterface
	AuditDraftReportStatusSvc             services.AuditDraftReportStatusServiceInterface
	AuditDraftReportSvc                   services.AuditDraftReportServiceInterface
	AuditDraftFindingSvc                  services.AuditDraftFindingServiceInterface
	AuditDraftFindingStatusSvc            services.AuditDraftFindingStatusServiceInterface
	AuditDraftFindingConversationSvc      services.AuditDraftFindingConversationServiceInterface
	AuditDraftFindingActionPlanSvc        services.AuditDraftFindingActionPlanServiceInterface
	AuditFindingRootCauseAnalysisSvc      services.AuditFindingRootCauseAnalysisServiceInterface
	AuditFindingActionPlanSvc             services.AuditFindingActionPlanServiceInterface
	AuditFindingActionPlanStatusUpdateSvc services.AuditFindingActionPlanStatusUpdateServiceInterface
}

func NewAuditModule(db *gorm.DB, sharedModule *SharedModule, organizationModule *OrganizationModule, complianceModule *ComplianceModule, riskModule *RiskModule, incidentModule *IncidentModule, documentModule *DocumentModule, settingModule *SettingsModule, masterModule *MasterModule, taskModule *TaskModule, commonModule *CommonModule) *AuditModule {
	auditStatusRepo := repositories.NewAuditStatusRepository(db)
	auditStatusSvc := services.NewAuditStatusService(auditStatusRepo)
	auditStatusController := controllers.NewAuditStatusController(auditStatusSvc)

	auditorRiskRatingRepo := repositories.NewAuditorRiskRatingRepository(db)
	auditorRiskRatingSvc := services.NewAuditorRiskRatingService(auditorRiskRatingRepo)
	auditorRiskRatingController := controllers.NewAuditorRiskRatingController(auditorRiskRatingSvc)

	auditUniverseRepo := repositories.NewAuditUniverseRepository(db)
	auditUniverseSvc := services.NewAuditUniverseService(
		auditUniverseRepo,
		sharedModule.ActivityLogService,
		auditStatusSvc, auditorRiskRatingSvc,
		organizationModule.ProcessSvc,
		organizationModule.ControlSvc,
		organizationModule.UserSvc,
		complianceModule.ComplianceRequirementSvc,
		riskModule.RiskRegisterSvc,
		incidentModule.IncidentRegisterSvc,
		documentModule.DocumentSvc,
	)
	auditUniverseController := controllers.NewAuditUniverseController(auditUniverseSvc)

	auditableItemRepo := repositories.NewAuditableItemRepository(db)
	auditableItemSvc := services.NewAuditableItemService(
		auditableItemRepo,
		sharedModule.ActivityLogService,
		auditStatusSvc, auditorRiskRatingSvc,
		settingModule.OrganizationUnitSvc,
		organizationModule.UserSvc,
	)
	auditableItemController := controllers.NewAuditableItemController(auditableItemSvc)

	auditTeamRepo := repositories.NewAuditTeamRepository(db)
	auditTeamSvc := services.NewAuditTeamService(
		auditTeamRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
	)
	auditTeamController := controllers.NewAuditTeamController(auditTeamSvc)

	auditCategoryRepo := repositories.NewAuditCategoryRepository(db)
	auditCategorySvc := services.NewAuditCategoryService(auditCategoryRepo, sharedModule.ActivityLogService)
	auditCategoryController := controllers.NewAuditCategoryController(auditCategorySvc)

	auditPlanRepo := repositories.NewAuditPlanRepository(db)
	auditPlanSvc := services.NewAuditPlanService(
		auditPlanRepo,
		sharedModule.ActivityLogService,
		auditCategorySvc,
		auditStatusSvc,
		auditTeamSvc,
		organizationModule.UserSvc,
		auditableItemSvc,
		commonModule.NotificationService,
	)
	auditPlanController := controllers.NewAuditPlanController(auditPlanSvc)

	auditPlanChecklistRepo := repositories.NewAuditPlanChecklistRepository(db)
	auditPlanChecklistSvc := services.NewAuditPlanChecklistService(
		auditPlanChecklistRepo,
		sharedModule.ActivityLogService,
		auditPlanSvc,
		auditableItemSvc,
	)
	auditPlanChecklistController := controllers.NewAuditPlanChecklistController(auditPlanChecklistSvc)

	auditPlanToolRepo := repositories.NewAuditPlanToolRepository(db)
	auditPlanToolSvc := services.NewAuditPlanToolService(
		auditPlanToolRepo,
		sharedModule.ActivityLogService,
		auditPlanSvc,
		masterModule.ToolTypeSvc,
	)
	auditPlanToolController := controllers.NewAuditPlanToolController(auditPlanToolSvc)

	auditPlanAuditableItemRepo := repositories.NewAuditPlanAuditableItemRepository(db)
	auditPlanAuditableItemSvc := services.NewAuditPlanAuditableItemService(auditPlanAuditableItemRepo, sharedModule.ActivityLogService, auditPlanSvc, auditableItemSvc)
	auditPlanAuditableItemController := controllers.NewAuditPlanAuditableItemController(auditPlanAuditableItemSvc)

	auditPlanBudgetRepo := repositories.NewAuditPlanBudgetRepository(db)
	auditPlanBudgetSvc := services.NewAuditPlanBudgetService(
		auditPlanBudgetRepo,
		sharedModule.ActivityLogService,
		auditPlanSvc,
		masterModule.ExpenseCategorySvc,
	)
	auditPlanBudgetController := controllers.NewAuditPlanBudgetController(auditPlanBudgetSvc)

	auditActivityTypeRepo := repositories.NewAuditActivityTypeRepository(db)
	auditActivityTypeSvc := services.NewAuditActivityTypeService(auditActivityTypeRepo, sharedModule.ActivityLogService)
	auditActivityTypeController := controllers.NewAuditActivityTypeController(auditActivityTypeSvc)

	auditScheduleStatusRepo := repositories.NewAuditScheduleStatusRepository(db)
	auditScheduleStatusSvc := services.NewAuditScheduleStatusService(auditScheduleStatusRepo)
	auditScheduleStatusController := controllers.NewAuditScheduleStatusController(auditScheduleStatusSvc)

	auditPlanScheduleRepo := repositories.NewAuditPlanScheduleRepository(db)
	auditPlanScheduleSvc := services.NewAuditPlanScheduleService(
		auditPlanScheduleRepo,
		sharedModule.ActivityLogService,
		auditPlanSvc,
		auditActivityTypeSvc,
		organizationModule.UserSvc,
		auditScheduleStatusSvc,
		auditableItemSvc,
	)
	auditPlanScheduleController := controllers.NewAuditPlanScheduleController(auditPlanScheduleSvc)

	auditEngagementTypeRepo := repositories.NewAuditEngagementTypeRepository(db)
	auditEngagementTypeSvc := services.NewAuditEngagementTypeService(auditEngagementTypeRepo, sharedModule.ActivityLogService)
	auditEngagementTypeController := controllers.NewAuditEngagementTypeController(auditEngagementTypeSvc)

	auditPlanStakeholderEngagementRepo := repositories.NewAuditPlanStakeholderEngagementRepository(db)
	auditPlanStakeholderEngagementSvc := services.NewAuditPlanStakeholderEngagementService(
		auditPlanStakeholderEngagementRepo,
		sharedModule.ActivityLogService,
		auditPlanSvc,
		auditEngagementTypeSvc,
		organizationModule.UserSvc,
		auditScheduleStatusSvc,
	)
	auditPlanStakeholderEngagementController := controllers.NewAuditPlanStakeholderEngagementController(auditPlanStakeholderEngagementSvc)

	auditRepo := repositories.NewAuditRepository(db)

	auditTimelineRepo := repositories.NewAuditTimelineRepository(db)
	auditTimelineSvc := services.NewAuditTimelineService(auditTimelineRepo, auditRepo)
	auditTimelineController := controllers.NewAuditTimelineController(auditTimelineSvc)

	auditSvc := services.NewAuditService(
		auditRepo,
		sharedModule.ActivityLogService,
		auditPlanSvc,
		auditTimelineSvc,
		auditPlanScheduleSvc,
		masterModule.StatusSvc,
		auditStatusSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
	)
	auditController := controllers.NewAuditController(auditSvc)

	auditWorkpaperRepo := repositories.NewAuditWorkpaperRepository(db)
	auditWorkpaperSvc := services.NewAuditWorkpaperService(
		auditWorkpaperRepo,
		sharedModule.ActivityLogService,
		auditSvc,
		masterModule.SourceSvc,
	)
	auditWorkpaperController := controllers.NewAuditWorkpaperController(auditWorkpaperSvc)

	auditMeetingRepo := repositories.NewAuditMeetingRepository(db)
	auditMeetingSvc := services.NewAuditMeetingService(
		auditMeetingRepo,
		sharedModule.ActivityLogService,
		auditSvc,
	)
	auditMeetingController := controllers.NewAuditMeetingController(auditMeetingSvc)

	auditFindingStatusRepo := repositories.NewAuditFindingStatusRepository(db)
	auditFindingStatusSvc := services.NewAuditFindingStatusService(auditFindingStatusRepo)
	auditFindingStatusController := controllers.NewAuditFindingStatusController(auditFindingStatusSvc)

	auditDraftFindingStatusRepo := repositories.NewAuditDraftFindingStatusRepository(db)
	auditDraftFindingStatusSvc := services.NewAuditDraftFindingStatusService(auditDraftFindingStatusRepo)
	auditDraftFindingStatusController := controllers.NewAuditDraftFindingStatusController(auditDraftFindingStatusSvc)

	auditDraftReportStatusRepo := repositories.NewAuditDraftReportStatusRepository(db)
	auditDraftReportStatusSvc := services.NewAuditDraftReportStatusService(auditDraftReportStatusRepo)
	auditDraftReportStatusController := controllers.NewAuditDraftReportStatusController(auditDraftReportStatusSvc)

	auditFindingCategoryRepo := repositories.NewAuditFindingCategoryRepository(db)
	auditFindingCategorySvc := services.NewAuditFindingCategoryService(auditFindingCategoryRepo, sharedModule.ActivityLogService)
	auditFindingCategoryController := controllers.NewAuditFindingCategoryController(auditFindingCategorySvc)

	auditChecklistRepo := repositories.NewAuditChecklistRepository(db)
	auditChecklistSvc := services.NewAuditChecklistService(auditChecklistRepo, sharedModule.ActivityLogService, auditSvc, auditableItemSvc, masterModule.ConformanceSvc, masterModule.StatusSvc, auditPlanAuditableItemSvc)
	auditChecklistController := controllers.NewAuditChecklistController(auditChecklistSvc)

	auditFindingRepo := repositories.NewAuditFindingRepository(db)

	auditFindingSvc := services.NewAuditFindingService(
		auditFindingRepo,
		sharedModule.ActivityLogService,
		auditSvc,
		auditFindingStatusSvc,
		auditTimelineSvc,
		auditableItemSvc,
		auditFindingCategorySvc,
		masterModule.SeverityLevelSvc,
	)
	auditFindingController := controllers.NewAuditFindingController(auditFindingSvc)

	auditDraftFindingRepo := repositories.NewAuditDraftFindingRepository(db)
	auditDraftFindingSvc := services.NewAuditDraftFindingService(
		auditDraftFindingRepo,
		sharedModule.ActivityLogService,
		auditSvc,
		auditDraftFindingStatusSvc,
		auditTimelineSvc,
		auditableItemSvc,
		auditFindingCategorySvc,
		masterModule.SeverityLevelSvc,
		auditChecklistSvc,
	)
	auditDraftFindingController := controllers.NewAuditDraftFindingController(auditDraftFindingSvc)

	auditDraftFindingActionPlanRepo := repositories.NewAuditDraftFindingActionPlanRepository(db)
	auditDraftFindingActionPlanSvc := services.NewAuditDraftFindingActionPlanService(
		auditDraftFindingActionPlanRepo,
		sharedModule.ActivityLogService,
		auditDraftFindingSvc,
		auditDraftFindingStatusSvc,
	)
	auditDraftFindingActionPlanController := controllers.NewAuditDraftFindingActionPlanController(auditDraftFindingActionPlanSvc)

	auditDocumentRequestStatusRepo := repositories.NewAuditDocumentRequestStatusRepository(db)
	auditDocumentRequestStatusSvc := services.NewAuditDocumentRequestStatusService(auditDocumentRequestStatusRepo)
	auditDocumentRequestStatusController := controllers.NewAuditDocumentRequestStatusController(auditDocumentRequestStatusSvc)

	auditDocumentRequestRepo := repositories.NewAuditDocumentRequestRepository(db)
	auditDocumentRequestSvc := services.NewAuditDocumentRequestService(
		auditDocumentRequestRepo,
		sharedModule.ActivityLogService,
		auditSvc,
		auditDocumentRequestStatusSvc,
		auditableItemSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService, taskModule.TaskService,
	)
	auditDocumentRequestController := controllers.NewAuditDocumentRequestController(auditDocumentRequestSvc)

	auditDocumentRequestConversationRepo := repositories.NewAuditDocumentRequestConversationRepository(db)
	auditDocumentRequestConversationSvc := services.NewAuditDocumentRequestConversationService(
		auditDocumentRequestConversationRepo,
		sharedModule.ActivityLogService,
		auditSvc,
		auditDocumentRequestSvc,
		auditDocumentRequestStatusSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService, taskModule.TaskService,
	)
	auditDocumentRequestConversationController := controllers.NewAuditDocumentRequestConversationController(auditDocumentRequestConversationSvc)

	auditFindingRootCauseAnalysisRepo := repositories.NewAuditFindingRootCauseAnalysisRepository(db)
	auditFindingRootCauseAnalysisSvc := services.NewAuditFindingRootCauseAnalysisService(auditFindingRootCauseAnalysisRepo, sharedModule.ActivityLogService, auditFindingSvc, auditTimelineSvc)
	auditFindingRootCauseAnalysisController := controllers.NewAuditFindingRootCauseAnalysisController(auditFindingRootCauseAnalysisSvc)

	auditFindingActionPlanRepo := repositories.NewAuditFindingActionPlanRepository(db)
	auditFindingActionPlanSvc := services.NewAuditFindingActionPlanService(
		auditFindingActionPlanRepo,
		sharedModule.ActivityLogService,
		auditFindingSvc,
		masterModule.ActionPlanStatusSvc,
		masterModule.SeverityLevelSvc,
		auditTimelineSvc,
		taskModule.TaskService,
		organizationModule.UserSvc,
		commonModule.NotificationService,
	)
	auditFindingActionPlanController := controllers.NewAuditFindingActionPlanController(auditFindingActionPlanSvc)

	auditFindingActionPlanStatusUpdateRepo := repositories.NewAuditFindingActionPlanStatusUpdateRepository(db)
	auditFindingActionPlanStatusUpdateSvc := services.NewAuditFindingActionPlanStatusUpdateService(auditFindingActionPlanStatusUpdateRepo, sharedModule.ActivityLogService, auditFindingActionPlanSvc, masterModule.ActionPlanStatusSvc)
	auditFindingActionPlanStatusUpdateController := controllers.NewAuditFindingActionPlanStatusUpdateController(auditFindingActionPlanStatusUpdateSvc)

	auditScheduleRepo := repositories.NewAuditScheduleRepository(db)
	auditScheduleSvc := services.NewAuditScheduleService(
		auditScheduleRepo,
		sharedModule.ActivityLogService,
		auditSvc,
		auditActivityTypeSvc,
		organizationModule.UserSvc,
		auditScheduleStatusSvc,
		auditableItemSvc,
	)
	auditScheduleController := controllers.NewAuditScheduleController(auditScheduleSvc)

	auditDraftReportRepo := repositories.NewAuditDraftReportRepository(db)
	auditDraftReportSvc := services.NewAuditDraftReportService(
		auditDraftReportRepo,
		sharedModule.ActivityLogService,
		auditSvc,
		auditDraftReportStatusSvc,
		auditDraftFindingSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
	)
	auditDraftReportController := controllers.NewAuditDraftReportController(auditDraftReportSvc)

	auditDraftFindingConversationRepo := repositories.NewAuditDraftFindingConversationRepository(db)
	auditDraftFindingConversationSvc := services.NewAuditDraftFindingConversationService(
		auditDraftFindingConversationRepo,
		sharedModule.ActivityLogService,
		auditDraftFindingSvc,
		auditDraftFindingStatusSvc,
		auditSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
		auditDraftReportSvc,
	)
	auditDraftFindingConversationController := controllers.NewAuditDraftFindingConversationController(auditDraftFindingConversationSvc)


	rcmRepo := repositories.NewRCMRepository(db)
	rcmSvc := services.NewRCMService(
		rcmRepo,
		sharedModule.ActivityLogService,
		auditSvc,
	)
	rcmController := controllers.NewRCMController(rcmSvc)

	auditScheduleStatusUpdateRepo := repositories.NewAuditScheduleStatusUpdateRepository(db)
	auditScheduleStatusUpdateSvc := services.NewAuditScheduleStatusUpdateService(auditScheduleStatusUpdateRepo, sharedModule.ActivityLogService, auditScheduleSvc, auditScheduleStatusSvc)
	auditScheduleStatusUpdateController := controllers.NewAuditScheduleStatusUpdateController(auditScheduleStatusUpdateSvc)

	auditDraftReportReviewRepo := repositories.NewAuditDraftReportReviewRepository(db)
	auditDraftReportReviewSvc := services.NewAuditDraftReportReviewService(
		auditDraftReportReviewRepo,
		sharedModule.ActivityLogService,
		auditSvc,
		auditDraftReportSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		auditDraftReportStatusSvc,
		commonModule.NotificationService,
		taskModule.TaskService,
	)
	auditDraftReportReviewController := controllers.NewAuditDraftReportReviewController(auditDraftReportReviewSvc)

	auditDraftReportCommentRepo := repositories.NewAuditDraftReportCommentRepository(db)
	auditDraftReportCommentSvc := services.NewAuditDraftReportCommentService(
		auditDraftReportCommentRepo,
		sharedModule.ActivityLogService,
		auditSvc,
		auditDraftReportSvc,
	)
	auditDraftReportCommentController := controllers.NewAuditDraftReportCommentController(auditDraftReportCommentSvc)

	auditPlanCommencementLetterStatusRepo := repositories.NewAuditPlanCommencementLetterStatusRepository(db)
	auditPlanCommencementLetterStatusSvc := services.NewAuditPlanCommencementLetterStatusService(auditPlanCommencementLetterStatusRepo)
	auditPlanCommencementLetterStatusController := controllers.NewAuditPlanCommencementLetterStatusController(auditPlanCommencementLetterStatusSvc)

	auditPlanCommencementLetterRepo := repositories.NewAuditPlanCommencementLetterRepository(db)
	auditPlanCommencementLetterSvc := services.NewAuditPlanCommencementLetterService(
		auditPlanCommencementLetterRepo,
		sharedModule.ActivityLogService,
		auditPlanSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
		auditPlanCommencementLetterStatusSvc,
	)
	auditPlanCommencementLetterController := controllers.NewAuditPlanCommencementLetterController(auditPlanCommencementLetterSvc)

	auditFinalReportStatusRepo := repositories.NewAuditFinalReportStatusRepository(db)
	auditFinalReportStatusSvc := services.NewAuditFinalReportStatusService(auditFinalReportStatusRepo)
	auditFinalReportStatusController := controllers.NewAuditFinalReportStatusController(auditFinalReportStatusSvc)

	auditReportWorkflowRepo := repositories.NewAuditReportWorkflowRepository(db)
	auditReportWorkflowSvc := services.NewAuditReportWorkflowService(
		auditReportWorkflowRepo,
		sharedModule.ActivityLogService,
		settingModule.WorkflowSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		auditSvc,
	)

	auditFinalReportRepo := repositories.NewAuditFinalReportRepository(db)
	auditFinalReportSvc := services.NewAuditFinalReportService(
		auditFinalReportRepo,
		sharedModule.ActivityLogService,
		auditSvc,
		auditDraftReportStatusSvc,
		auditDraftReportSvc,
		auditDraftFindingSvc,
		organizationModule.UserSvc,
		auditFinalReportStatusSvc,
		auditFindingSvc,
		commonModule.NotificationService,
	)
	auditFinalReportController := controllers.NewAuditFinalReportController(auditFinalReportSvc)

	auditFinalReportSectionRepo := repositories.NewAuditFinalReportSectionRepository(db)
	auditFinalReportSectionSvc := services.NewAuditFinalReportSectionService(
		auditFinalReportSectionRepo,
		sharedModule.ActivityLogService,
		auditSvc,
		auditFindingSvc,
	)
	auditFinalReportSectionController := controllers.NewAuditFinalReportSectionController(auditFinalReportSectionSvc)

	auditFinalReportSvc.SetAuditReportWorkflowService(auditReportWorkflowSvc)
	auditReportWorkflowSvc.SetAuditFinalReportService(auditFinalReportSvc)
	auditReportWorkflowController := controllers.NewAuditReportWorkflowController(auditReportWorkflowSvc)

	auditReportReviewRepo := repositories.NewAuditReportReviewRepository(db)
	auditReportReviewSvc := services.NewAuditReportReviewService(
		auditReportReviewRepo,
		sharedModule.ActivityLogService,
		auditFinalReportSvc,
		auditReportWorkflowSvc,
		masterModule.WorkflowStatusSvc,
		auditFinalReportStatusSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
		taskModule.TaskService)
	auditReportReviewController := controllers.NewAuditReportReviewController(auditReportReviewSvc)

	return &AuditModule{
		AuditorRiskRatingController:                  auditorRiskRatingController,
		AuditUniverseController:                      auditUniverseController,
		AuditableItemController:                      auditableItemController,
		AuditTeamController:                          auditTeamController,
		AuditCategoryController:                      auditCategoryController,
		AuditStatusController:                        auditStatusController,
		AuditPlanController:                          auditPlanController,
		AuditPlanChecklistController:                 auditPlanChecklistController,
		AuditPlanToolController:                      auditPlanToolController,
		AuditPlanAuditableItemController:             auditPlanAuditableItemController,
		AuditActivityTypeController:                  auditActivityTypeController,
		AuditScheduleStatusController:                auditScheduleStatusController,
		AuditPlanBudgetController:                    auditPlanBudgetController,
		AuditPlanScheduleController:                  auditPlanScheduleController,
		AuditEngagementTypeController:                auditEngagementTypeController,
		AuditPlanStakeholderEngagementController:     auditPlanStakeholderEngagementController,
		AuditController:                              auditController,
		AuditWorkpaperController:                     auditWorkpaperController,
		AuditMeetingController:                       auditMeetingController,
		AuditTimelineController:                      auditTimelineController,
		AuditFindingStatusController:                 auditFindingStatusController,
		AuditDraftFindingStatusController:            auditDraftFindingStatusController,
		AuditDraftReportStatusController:             auditDraftReportStatusController,
		AuditFindingCategoryController:               auditFindingCategoryController,
		AuditFindingController:                       auditFindingController,
		AuditDraftFindingController:                  auditDraftFindingController,
		AuditDocumentRequestStatusController:         auditDocumentRequestStatusController,
		AuditDocumentRequestController:               auditDocumentRequestController,
		AuditFindingRootCauseAnalysisController:      auditFindingRootCauseAnalysisController,
		AuditFindingActionPlanController:             auditFindingActionPlanController,
		AuditFindingActionPlanStatusUpdateController: auditFindingActionPlanStatusUpdateController,
		AuditDocumentRequestConversationController:   auditDocumentRequestConversationController,
		AuditChecklistController:                     auditChecklistController,
		AuditScheduleController:                      auditScheduleController,
		AuditDraftReportController:                   auditDraftReportController,
		RCMController:                                rcmController,
		AuditDraftFindingConversationController:      auditDraftFindingConversationController,
		AuditDraftFindingActionPlanController:        auditDraftFindingActionPlanController,
		AuditScheduleStatusUpdateController:          auditScheduleStatusUpdateController,
		AuditDraftReportReviewController:             auditDraftReportReviewController,
		AuditDraftReportCommentController:            auditDraftReportCommentController,
		AuditPlanCommencementLetterStatusController:  auditPlanCommencementLetterStatusController,
		AuditPlanCommencementLetterController:        auditPlanCommencementLetterController,
		AuditFinalReportStatusController:             auditFinalReportStatusController,
		AuditFinalReportController:                   auditFinalReportController,
		AuditFinalReportSectionController:            auditFinalReportSectionController,
		AuditReportWorkflowController:                auditReportWorkflowController,
		AuditReportReviewController:                  auditReportReviewController,

		AuditPlanSvc:                          auditPlanSvc,
		AuditSvc:                              auditSvc,
		AuditFindingSvc:                       auditFindingSvc,
		AuditDocumentRequestSvc:               auditDocumentRequestSvc,
		AuditDocumentRequestConversationSvc:   auditDocumentRequestConversationSvc,
		AuditDraftReportSvc:                   auditDraftReportSvc,
		AuditDraftReportStatusSvc:             auditDraftReportStatusSvc,
		AuditDraftFindingSvc:                  auditDraftFindingSvc,
		AuditDraftFindingStatusSvc:            auditDraftFindingStatusSvc,
		AuditDraftFindingConversationSvc:      auditDraftFindingConversationSvc,
		AuditDraftFindingActionPlanSvc:        auditDraftFindingActionPlanSvc,
		AuditFindingRootCauseAnalysisSvc:      auditFindingRootCauseAnalysisSvc,
		AuditFindingActionPlanSvc:             auditFindingActionPlanSvc,
		AuditFindingActionPlanStatusUpdateSvc: auditFindingActionPlanStatusUpdateSvc,
	}
}
