package initializers

import (
	"jethur/grc/audit-insights/controllers"
	"jethur/grc/audit-insights/repositories"
	"jethur/grc/audit-insights/services"

	"gorm.io/gorm"
)

type AuditInsightsModule struct {
	AuditInsightsAuditPlanController       *controllers.AuditInsightsAuditPlanController
	AuditInsightsAuditController           *controllers.AuditInsightsAuditController
	AuditInsightsDocumentRequestController *controllers.AuditInsightsDocumentRequestController
	AuditInsightsDraftReportController     *controllers.AuditInsightsDraftReportController


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

	auditInsightsDraftReportRepo := repositories.NewAuditInsightsDraftReportRepository(db)
	auditInsightsDraftReportSvc := services.NewAuditInsightsDraftReportService(
		auditInsightsDraftReportRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		auditModule.AuditSvc,
		auditModule.AuditDraftReportSvc,
	)
	auditInsightsDraftReportController := controllers.NewAuditInsightsDraftReportController(auditInsightsDraftReportSvc)

	return &AuditInsightsModule{
		AuditInsightsAuditPlanController:       auditInsightsAuditPlanController,
		AuditInsightsAuditController:           auditInsightsAuditController,
		AuditInsightsDocumentRequestController: auditInsightsDocumentRequestController,
		AuditInsightsDraftReportController:     auditInsightsDraftReportController,
	}
}
