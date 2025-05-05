package initializers

import (
	"jethur/grc/reports/controllers"
	"jethur/grc/reports/repositories"
	"jethur/grc/reports/services"

	"gorm.io/gorm"
)

type ReportsModule struct {
	ComplianceReportController   *controllers.ComplianceReportController
	OrganizationReportController *controllers.OrganizationReportController
	RiskReportController         *controllers.RiskReportController
	AuditReportController        *controllers.AuditReportController
	IncidentReportController     *controllers.IncidentReportController
	AuditInsightReportController *controllers.AuditInsightReportController
	MeetingReportController      *controllers.MeetingReportController
	DocumentReportController     *controllers.DocumentReportController
	GRCReportController          *controllers.GRCReportController
}

func NewReportsModule(db *gorm.DB,
	complianceModule *ComplianceModule,
	organizationModule *OrganizationModule,
	riskModule *RiskModule,
	auditModule *AuditModule,
	auditInsightsModule *AuditInsightsModule,
	incidentModule *IncidentModule,
	meetingModule *MeetingModule,
	settingsModule *SettingsModule,
	documentModule *DocumentModule,
	sharedModule *SharedModule) *ReportsModule {

	organizationReportRepo := repositories.NewOrganizationReportRepository(db)
	organizationReportSvc := services.NewOrganizationReportService(organizationReportRepo)
	organizationReportController := controllers.NewOrganizationReportController(organizationReportSvc, organizationModule.UserSvc, organizationModule.ProcessSvc, organizationModule.ControlSvc)

	complianceReportRepo := repositories.NewComplianceReportRepository(db)
	complianceReportSvc := services.NewComplianceReportService(complianceReportRepo)
	complianceReportController := controllers.NewComplianceReportController(complianceReportSvc,
		complianceModule.ComplianceRegisterSvc,
		complianceModule.SlaRegisterSvc,
		complianceModule.ComplianceRequirementSvc,
		complianceModule.ComplianceActionPlanSvc)

	riskReportRepo := repositories.NewRiskReportRepository(db)
	riskReportSvc := services.NewRiskReportService(riskReportRepo)
	riskReportController := controllers.NewRiskReportController(riskReportSvc,
		riskModule.KriRegisterSvc,
		riskModule.RiskRegisterSvc,
		riskModule.RiskActionPlanSvc)

	auditReportRepo := repositories.NewAuditReportRepository(db)
	auditReportSvc := services.NewAuditReportService(auditReportRepo)
	auditReportController := controllers.NewAuditReportController(auditReportSvc,
		auditModule.AuditPlanSvc,
		auditModule.AuditSvc,
		auditModule.AuditFindingSvc,
		auditModule.AuditFindingActionPlanSvc)

	incidentReportRepo := repositories.NewIncidentReportRepository(db)
	incidentReportSvc := services.NewIncidentReportService(incidentReportRepo)
	incidentReportController := controllers.NewIncidentReportController(incidentReportSvc,
		incidentModule.IncidentRegisterSvc,
		incidentModule.IncidentActionPlanSvc)

	auditInsightReportRepo := repositories.NewAuditInsightReportRepository(db)
	auditInsightReportSvc := services.NewAuditInsightReportService(auditInsightReportRepo, organizationModule.UserSvc)
	auditInsightReportController := controllers.NewAuditInsightReportController(auditInsightReportSvc,
		auditInsightsModule.AuditInsightsAuditSvc,
		auditInsightsModule.AuditInsightsAuditFindingSvc,
		auditInsightsModule.AuditInsightsFindingActionPlanSvc,
		auditInsightsModule.AuditInsightsAuditPlanSvc,
	)
	meetingReportRepo := repositories.NewMeetingReportRepository(db)
	meetingReportSvc := services.NewMeetingReportService(meetingReportRepo, settingsModule.OrganizationLevelSvc)
	meetingReportController := controllers.NewMeetingReportController(meetingReportSvc,
		meetingModule.MeetingSvc,
		meetingModule.MeetingActionPlanSvc)

	documentReportRepo := repositories.NewDocumentReportRepository(db)
	documentReportSvc := services.NewDocumentReportService(documentReportRepo)
	documentReportController := controllers.NewDocumentReportController(documentReportSvc,
		documentModule.DocumentSvc, documentModule.DocumentChangeRequestSvc)

	grcReportRepo := repositories.NewGRCReportRepository(db)
	grcReportSvc := services.NewGRCReportService(grcReportRepo, sharedModule.ActivityLogService, settingsModule.OrganizationLevelSvc, riskModule.RiskResponseStrategySvc)
	grcReportController := controllers.NewGRCReportController(grcReportSvc)

	return &ReportsModule{
		ComplianceReportController:   complianceReportController,
		OrganizationReportController: organizationReportController,
		RiskReportController:         riskReportController,
		AuditReportController:        auditReportController,
		IncidentReportController:     incidentReportController,
		AuditInsightReportController: auditInsightReportController,
		MeetingReportController:      meetingReportController,
		DocumentReportController:     documentReportController,
		GRCReportController:          grcReportController,
	}
}
