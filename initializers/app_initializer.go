package initializers

import (
	"gorm.io/gorm"
)

type App struct {
	SettingsModule      *SettingsModule
	CommonModule        *CommonModule
	MasterModule        *MasterModule
	MeetingModule       *MeetingModule
	LibraryModule       *LibraryModule
	OrganizationModule  *OrganizationModule
	AuthModule          *AuthModule
	ComplianceModule    *ComplianceModule
	RiskModule          *RiskModule
	IncidentModule      *IncidentModule
	TaskModule          *TaskModule
	DocumentModule      *DocumentModule
	AuditModule         *AuditModule
	KPIModule           *KPIModule
	AuditInsightsModule *AuditInsightsModule
}

func NewApp(db *gorm.DB) *App {
	sharedModule := NewSharedModule(db)
	masterModule := NewMasterModule(db, sharedModule)
	settingsModule := NewSettingsModule(db, sharedModule, masterModule)
	taskModule := NewTaskModule(db, settingsModule, masterModule)
	commonModule := NewCommonModule(db, settingsModule)
	libraryModule := NewLibraryModule(db, sharedModule, masterModule)
	organizationModule := NewOrganizationModule(db, sharedModule, settingsModule, masterModule, taskModule, commonModule)
	authModule := NewAuthModule(db, organizationModule, settingsModule, commonModule)
	complianceModule := NewComplianceModule(db, sharedModule, masterModule, organizationModule, taskModule, settingsModule)
	riskModule := NewRiskModule(db, sharedModule, settingsModule, masterModule, libraryModule, organizationModule, complianceModule, commonModule, taskModule)
	incidentModule := NewIncidentModule(db, sharedModule, settingsModule, masterModule, organizationModule, complianceModule, riskModule, taskModule)
	documentModule := NewDocumentModule(db, sharedModule, settingsModule, masterModule, organizationModule)
	meetingModule := NewMeetingModule(db, sharedModule, organizationModule, masterModule, commonModule, complianceModule, riskModule, incidentModule)
	auditModule := NewAuditModule(db, sharedModule, organizationModule, complianceModule, riskModule, incidentModule, documentModule, settingsModule, masterModule, taskModule)
	kpiModule := NewKPIModule(db)
	auditInsightsModule := NewAuditInsightsModule(db, sharedModule, organizationModule, auditModule)

	return &App{
		SettingsModule:      settingsModule,
		CommonModule:        commonModule,
		MasterModule:        masterModule,
		LibraryModule:       libraryModule,
		OrganizationModule:  organizationModule,
		AuthModule:          authModule,
		TaskModule:          taskModule,
		ComplianceModule:    complianceModule,
		RiskModule:          riskModule,
		IncidentModule:      incidentModule,
		DocumentModule:      documentModule,
		MeetingModule:       meetingModule,
		AuditModule:         auditModule,
		KPIModule:           kpiModule,
		AuditInsightsModule: auditInsightsModule,
	}
}
