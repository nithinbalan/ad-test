package initializers

import (
	"gorm.io/gorm"
)

type App struct {
	DB                       *gorm.DB
	SettingsModule           *SettingsModule
	CommonModule             *CommonModule
	MasterModule             *MasterModule
	MeetingModule            *MeetingModule
	LibraryModule            *LibraryModule
	OrganizationModule       *OrganizationModule
	AuthModule               *AuthModule
	ComplianceModule         *ComplianceModule
	RiskModule               *RiskModule
	IncidentModule           *IncidentModule
	TaskModule               *TaskModule
	DocumentModule           *DocumentModule
	AuditModule              *AuditModule
	KPIModule                *KPIModule
	AuditInsightsModule      *AuditInsightsModule
	ReportsModule            *ReportsModule
	AssetModule              *AssetModule
	SharedModule             *SharedModule
	BusinessContinuityModule *BusinessContinuityModule
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
	complianceModule := NewComplianceModule(db, sharedModule, masterModule, organizationModule, taskModule, settingsModule, commonModule)
	assetModule := NewAssetModule(db, sharedModule, settingsModule, masterModule, organizationModule, taskModule, commonModule)
	riskModule := NewRiskModule(db, sharedModule, settingsModule, masterModule, libraryModule, organizationModule, complianceModule, commonModule, taskModule, assetModule)
	incidentModule := NewIncidentModule(db, sharedModule, settingsModule, masterModule, organizationModule, complianceModule, riskModule, taskModule, commonModule, assetModule)
	documentModule := NewDocumentModule(db, sharedModule, settingsModule, masterModule, organizationModule, taskModule, commonModule)
	meetingModule := NewMeetingModule(db, sharedModule, organizationModule, masterModule, commonModule, complianceModule, riskModule, incidentModule, taskModule)
	auditModule := NewAuditModule(db, sharedModule, organizationModule, complianceModule, riskModule, incidentModule, documentModule, settingsModule, masterModule, taskModule, commonModule)
	kpiModule := NewKPIModule(db, riskModule)
	auditInsightsModule := NewAuditInsightsModule(db, sharedModule, organizationModule, auditModule, commonModule)
	reportsModule := NewReportsModule(db, complianceModule, organizationModule, riskModule)
	businessFunctionModule := NewBusinessContinuityModule(db, sharedModule, settingsModule, masterModule, organizationModule, assetModule, documentModule)

	return &App{
		DB:                       db,
		SettingsModule:           settingsModule,
		CommonModule:             commonModule,
		MasterModule:             masterModule,
		LibraryModule:            libraryModule,
		OrganizationModule:       organizationModule,
		AuthModule:               authModule,
		TaskModule:               taskModule,
		ComplianceModule:         complianceModule,
		RiskModule:               riskModule,
		IncidentModule:           incidentModule,
		DocumentModule:           documentModule,
		MeetingModule:            meetingModule,
		AuditModule:              auditModule,
		KPIModule:                kpiModule,
		AuditInsightsModule:      auditInsightsModule,
		ReportsModule:            reportsModule,
		AssetModule:              assetModule,
		SharedModule:             sharedModule,
		BusinessContinuityModule: businessFunctionModule,
	}
}
