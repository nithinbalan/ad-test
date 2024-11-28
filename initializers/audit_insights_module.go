package initializers

import (
	"jethur/grc/audit-insights/controllers"
	"jethur/grc/audit-insights/repositories"
	"jethur/grc/audit-insights/services"

	"gorm.io/gorm"
)

type AuditInsightsModule struct {
	AuditInsightsAuditPlanController                     *controllers.AuditInsightsAuditPlanController
	AuditInsightsAuditController                         *controllers.AuditInsightsAuditController
	AuditInsightsDocumentRequestController               *controllers.AuditInsightsDocumentRequestController
	AuditInsightsDraftReportController                   *controllers.AuditInsightsDraftReportController
	AuditInsightsDocumentRequestConversationController   *controllers.AuditInsightsDocumentRequestConversationController
	AuditInsightsDraftFindingController                  *controllers.AuditInsightsDraftFindingController
	AuditInsightsDraftFindingConversationController      *controllers.AuditInsightsDraftFindingConversationController
	AuditInsightsDraftFindingActionPlanController        *controllers.AuditInsightsDraftFindingActionPlanController
	AuditInsightsAuditFindingController                  *controllers.AuditInsightsAuditFindingController
	AuditInsightsFindingRootCauseAnalysisController      *controllers.AuditInsightsFindingRootCauseAnalysisController
	AuditInsightsFindingActionPlanController             *controllers.AuditInsightsFindingActionPlanController
	AuditInsightsFindingActionPlanStatusUpdateController *controllers.AuditInsightsFindingActionPlanStatusUpdateController
}

func NewAuditInsightsModule(db *gorm.DB, sharedModule *SharedModule, organizationModule *OrganizationModule, auditModule *AuditModule) *AuditInsightsModule {
	auditInsightsAuditPlanRepo := repositories.NewAuditInsightsAuditPlanRepository(db)
	auditInsightsAuditPlanSvc := services.NewAuditInsightsAuditPlanService(
		auditInsightsAuditPlanRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		auditModule.AuditPlanSvc,
	)
	auditInsightsAuditPlanController := controllers.NewAuditInsightsAuditPlanController(auditInsightsAuditPlanSvc)

	auditInsightsAuditRepo := repositories.NewAuditInsightsAuditRepository(db)
	auditInsightsAuditSvc := services.NewAuditInsightsAuditService(
		auditInsightsAuditRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		auditModule.AuditSvc,
	)
	auditInsightsAuditController := controllers.NewAuditInsightsAuditController(auditInsightsAuditSvc)

	auditInsightsDocumentRequestRepo := repositories.NewAuditInsightsDocumentRequestRepository(db)
	auditInsightsDocumentRequestSvc := services.NewAuditInsightsDocumentRequestService(
		auditInsightsDocumentRequestRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		auditModule.AuditDocumentRequestSvc,
	)
	auditInsightsDocumentRequestController := controllers.NewAuditInsightsDocumentRequestController(auditInsightsDocumentRequestSvc)

	//repo not using
	auditInsightsDocumentRequestConversationSvc := services.NewAuditInsightsDocumentRequestConversationService(
		sharedModule.ActivityLogService,
		auditModule.AuditDocumentRequestConversationSvc,
		auditModule.AuditDocumentRequestSvc,
	)
	auditInsightsDocumentRequestConversationController := controllers.NewAuditInsightsDocumentRequestConversationController(auditInsightsDocumentRequestConversationSvc)

	auditInsightsDraftReportRepo := repositories.NewAuditInsightsDraftReportRepository(db)
	auditInsightsDraftReportSvc := services.NewAuditInsightsDraftReportService(
		auditInsightsDraftReportRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		auditModule.AuditSvc,
		auditModule.AuditDraftReportSvc,
	)
	auditInsightsDraftReportController := controllers.NewAuditInsightsDraftReportController(auditInsightsDraftReportSvc)

	auditInsightsDraftFindingRepo := repositories.NewAuditInsightsDraftFindingRepository(db)
	auditInsightsDraftFindingSvc := services.NewAuditInsightsDraftFindingService(
		auditInsightsDraftFindingRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		auditModule.AuditDraftReportSvc,
		auditModule.AuditDraftReportStatusSvc,
		auditModule.AuditDraftFindingSvc,
		auditModule.AuditDraftFindingStatusSvc,
		auditModule.AuditDraftFindingConversationSvc,
		auditModule.AuditDraftFindingActionPlanSvc,
	)
	auditInsightsDraftFindingController := controllers.NewAuditInsightsDraftFindingController(auditInsightsDraftFindingSvc)

	auditInsightsDraftFindingConversationRepo := repositories.NewAuditInsightsDraftFindingConversationRepository(db)
	auditInsightsDraftFindingConversationSvc := services.NewAuditInsightsDraftFindingConversationService(
		auditInsightsDraftFindingConversationRepo,
		sharedModule.ActivityLogService,
		auditModule.AuditDraftFindingStatusSvc,
		auditModule.AuditDraftFindingConversationSvc,
		auditModule.AuditDraftFindingSvc,
		auditModule.AuditSvc,
		organizationModule.UserSvc,
	)
	auditInsightsDraftFindingConversationController := controllers.NewAuditInsightsDraftFindingConversationController(auditInsightsDraftFindingConversationSvc)

	auditInsightsDraftFindingActionPlanRepo := repositories.NewAuditInsightsDraftFindingActionPlanRepository(db)
	auditInsightsDraftFindingActionPlanSvc := services.NewAuditInsightsDraftFindingActionPlanService(
		auditInsightsDraftFindingActionPlanRepo,
		sharedModule.ActivityLogService,
		auditModule.AuditDraftFindingSvc,
		auditModule.AuditDraftFindingStatusSvc,
		auditModule.AuditDraftFindingActionPlanSvc,
		organizationModule.UserSvc,
		auditModule.AuditDraftReportSvc,
		auditModule.AuditDraftReportStatusSvc,
	)
	auditInsightsDraftFindingActionPlanController := controllers.NewAuditInsightsDraftFindingActionPlanController(auditInsightsDraftFindingActionPlanSvc)

	auditInsightsAuditFindingRepo := repositories.NewAuditInsightsAuditFindingRepository(db)
	auditInsightsAuditFindingSvc := services.NewAuditInsightsAuditFindingService(
		auditInsightsAuditFindingRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		auditModule.AuditFindingSvc,
	)
	auditInsightsAuditFindingController := controllers.NewAuditInsightsAuditFindingController(auditInsightsAuditFindingSvc)

	auditInsightsFindingRootCauseAnalysisSvc := services.NewAuditInsightsFindingRootCauseAnalysisService(auditModule.AuditFindingRootCauseAnalysisSvc)
	auditInsightsFindingRootCauseAnalysisController := controllers.NewAuditInsightsFindingRootCauseAnalysisController(auditInsightsFindingRootCauseAnalysisSvc)

	auditInsightsFindingActionPlanRepo := repositories.NewAuditInsightsFindingActionPlanRepository(db)
	auditInsightsFindingActionPlanSvc := services.NewAuditInsightsFindingActionPlanService(
		auditInsightsFindingActionPlanRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		auditModule.AuditFindingActionPlanSvc,
	)
	auditInsightsFindingActionPlanController := controllers.NewAuditInsightsFindingActionPlanController(auditInsightsFindingActionPlanSvc)

	auditInsightsFindingActionPlanStatusUpdateSvc := services.NewAuditInsightsFindingActionPlanStatusUpdateService(
		auditInsightsFindingActionPlanSvc,
		auditModule.AuditFindingActionPlanStatusUpdateSvc,
	)
	auditInsightsFindingActionPlanStatusUpdateController := controllers.NewAuditInsightsFindingActionPlanStatusUpdateController(auditInsightsFindingActionPlanStatusUpdateSvc)

	return &AuditInsightsModule{
		AuditInsightsAuditPlanController:                     auditInsightsAuditPlanController,
		AuditInsightsAuditController:                         auditInsightsAuditController,
		AuditInsightsDocumentRequestController:               auditInsightsDocumentRequestController,
		AuditInsightsDraftReportController:                   auditInsightsDraftReportController,
		AuditInsightsDocumentRequestConversationController:   auditInsightsDocumentRequestConversationController,
		AuditInsightsDraftFindingController:                  auditInsightsDraftFindingController,
		AuditInsightsDraftFindingConversationController:      auditInsightsDraftFindingConversationController,
		AuditInsightsDraftFindingActionPlanController:        auditInsightsDraftFindingActionPlanController,
		AuditInsightsAuditFindingController:                  auditInsightsAuditFindingController,
		AuditInsightsFindingRootCauseAnalysisController:      auditInsightsFindingRootCauseAnalysisController,
		AuditInsightsFindingActionPlanController:             auditInsightsFindingActionPlanController,
		AuditInsightsFindingActionPlanStatusUpdateController: auditInsightsFindingActionPlanStatusUpdateController,
	}
}
