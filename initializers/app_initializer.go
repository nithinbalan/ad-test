package initializers

import (
	"gorm.io/gorm"
)

type App struct {
	DB                          *gorm.DB
	SettingsModule              *SettingsModule
	CommonModule                *CommonModule
	MasterModule                *MasterModule
	MeetingModule               *MeetingModule
	LibraryModule               *LibraryModule
	OrganizationModule          *OrganizationModule
	AuthModule                  *AuthModule
	AwarenessReportingModule    *AwarenessReportingModule
	ComplianceModule            *ComplianceModule
	RiskModule                  *RiskModule
	IncidentModule              *IncidentModule
	TaskModule                  *TaskModule
	DocumentModule              *DocumentModule
	AuditModule                 *AuditModule
	KPIModule                   *KPIModule
	AuditInsightsModule         *AuditInsightsModule
	ReportsModule               *ReportsModule
	AssetModule                 *AssetModule
	SharedModule                *SharedModule
	BusinessContinuityModule    *BusinessContinuityModule
	OrganizationalContextModule *OrganizationalContextModule
	StrategyModule              *StrategyModule
	KPIRegisterModule           *KPIRegisterModule
	TprmModule                  *TprmModule
	ProjectModule               *ProjectModule
	JethurAIModule              *JethurAIModule
	PolicyModule                *PolicyModule
	ConsultationRequestModule   *ConsultationRequestModule
	ViolationReportModule       *ViolationReportModule
	IdentityAccessReviewModule  *IdentityAccessReviewModule
	SecurityOperationsModule    *SecurityOperationsModule
	PenetrationTestingModule    *PenetrationTestingModule
	ThreatIntelligenceModule    *ThreatIntelligenceModule
}

func NewApp(db *gorm.DB) *App {
	sharedModule := NewSharedModule(db)
	masterModule := NewMasterModule(db, sharedModule)
	settingsModule := NewSettingsModule(db, sharedModule, masterModule)
	taskModule := NewTaskModule(db, settingsModule, masterModule)
	commonModule := NewCommonModule(db, settingsModule, masterModule)
	organizationModule := NewOrganizationModule(db, sharedModule, settingsModule, masterModule, taskModule, commonModule)
	libraryModule := NewLibraryModule(db, sharedModule, masterModule, organizationModule)
	authModule := NewAuthModule(db, organizationModule, settingsModule, commonModule)
	awarenessReportingModule := NewAwarenessReportingModule(db, sharedModule, settingsModule, organizationModule, commonModule, masterModule)
	documentModule := NewDocumentModule(db, sharedModule, settingsModule, masterModule, organizationModule, taskModule, commonModule)
	complianceModule := NewComplianceModule(db, sharedModule, masterModule, organizationModule, taskModule, settingsModule, commonModule, documentModule)
	assetModule := NewAssetModule(db, sharedModule, settingsModule, masterModule, organizationModule, taskModule, commonModule, documentModule, complianceModule)
	riskModule := NewRiskModule(db, sharedModule, settingsModule, masterModule, libraryModule, organizationModule, complianceModule, commonModule, taskModule, assetModule)
	assetModule.SetRiskControllers(riskModule)
	incidentModule := NewIncidentModule(db, sharedModule, settingsModule, masterModule, organizationModule, complianceModule, riskModule, taskModule, commonModule, assetModule)
	meetingModule := NewMeetingModule(db, sharedModule, organizationModule, masterModule, commonModule, complianceModule, riskModule, incidentModule, taskModule, settingsModule)
	auditModule := NewAuditModule(db, sharedModule, organizationModule, complianceModule, riskModule, incidentModule, documentModule, settingsModule, masterModule, taskModule, commonModule)
	kpiModule := NewKPIModule(db, riskModule, organizationModule, documentModule)
	auditInsightsModule := NewAuditInsightsModule(db, sharedModule, organizationModule, auditModule, commonModule)
	tprmModule := NewTprmModule(db, sharedModule, settingsModule, masterModule, organizationModule, taskModule, commonModule, riskModule)
	projectModule := NewProjectModule(db, sharedModule, organizationModule, masterModule, commonModule, taskModule, settingsModule, incidentModule)
	reportsModule := NewReportsModule(db, complianceModule, organizationModule, riskModule, auditModule, auditInsightsModule, incidentModule, meetingModule, settingsModule, documentModule, sharedModule, kpiModule, commonModule, tprmModule, projectModule)
	businessFunctionModule := NewBusinessContinuityModule(db, sharedModule, settingsModule, masterModule, organizationModule, assetModule, documentModule, riskModule, commonModule, taskModule)
	organizationalContextModule := NewOrganizationalContextModule(db, sharedModule, settingsModule, masterModule, organizationModule, taskModule, commonModule)
	strategyModule := NewStrategyModule(db, sharedModule, settingsModule, masterModule, organizationModule, taskModule, commonModule, riskModule, assetModule, tprmModule, complianceModule, documentModule, projectModule)
	kpiRegisterModule := NewKPIRegisterModule(db, sharedModule, settingsModule, masterModule, organizationModule, commonModule, riskModule, assetModule, tprmModule, complianceModule, documentModule, strategyModule)
	jethurAIModule := NewJethurAIModule(db, documentModule, organizationModule, settingsModule, complianceModule)
	policyModule := NewPolicyModule(db, sharedModule, settingsModule, masterModule, organizationModule, taskModule, commonModule, documentModule)
	consultationRequestModule := NewConsultationRequestModule(db, sharedModule, masterModule, organizationModule, taskModule, settingsModule, commonModule, documentModule)
	violationReportModule := NewViolationReportModule(db, sharedModule, settingsModule, masterModule, organizationModule, complianceModule, riskModule, taskModule, commonModule, assetModule)
	identityAccessReviewModule := NewIdentityAccessReviewModule(db, sharedModule, organizationModule, settingsModule, masterModule, taskModule, commonModule)
	securityOperationsModule := NewSecurityOperationsModule(db, sharedModule, settingsModule, organizationModule, commonModule)
	penetrationTestingModule := NewPenetrationTestingModule(db, sharedModule, masterModule, organizationModule)
	threatIntelligenceModule := NewThreatIntelligenceModule(db, sharedModule)

	return &App{
		DB:                          db,
		SettingsModule:              settingsModule,
		CommonModule:                commonModule,
		MasterModule:                masterModule,
		LibraryModule:               libraryModule,
		OrganizationModule:          organizationModule,
		AuthModule:                  authModule,
		AwarenessReportingModule:    awarenessReportingModule,
		TaskModule:                  taskModule,
		ComplianceModule:            complianceModule,
		RiskModule:                  riskModule,
		IncidentModule:              incidentModule,
		DocumentModule:              documentModule,
		MeetingModule:               meetingModule,
		AuditModule:                 auditModule,
		KPIModule:                   kpiModule,
		AuditInsightsModule:         auditInsightsModule,
		ReportsModule:               reportsModule,
		AssetModule:                 assetModule,
		SharedModule:                sharedModule,
		BusinessContinuityModule:    businessFunctionModule,
		OrganizationalContextModule: organizationalContextModule,
		StrategyModule:              strategyModule,
		KPIRegisterModule:           kpiRegisterModule,
		TprmModule:                  tprmModule,
		ProjectModule:               projectModule,
		JethurAIModule:              jethurAIModule,
		PolicyModule:                policyModule,
		ConsultationRequestModule:   consultationRequestModule,
		ViolationReportModule:       violationReportModule,
		IdentityAccessReviewModule:  identityAccessReviewModule,
		SecurityOperationsModule:    securityOperationsModule,
		PenetrationTestingModule:    penetrationTestingModule,
		ThreatIntelligenceModule:    threatIntelligenceModule,
	}
}
