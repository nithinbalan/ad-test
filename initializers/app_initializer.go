package initializers

import (
	"gorm.io/gorm"
)

type App struct {
	SettingsModule     *SettingsModule
	CommonModule       *CommonModule
	MasterModule       *MasterModule
	MeetingModule      *MeetingModule
	LibraryModule      *LibraryModule
	OrganizationModule *OrganizationModule
	AuthModule         *AuthModule
	ComplianceModule   *ComplianceModule
	RiskModule         *RiskModule
	IncidentModule     *IncidentModule
	TaskModule         *TaskModule
	DocumentModule     *DocumentModule
	AuditModule        *AuditModule
}

func NewApp(db *gorm.DB) *App {
	sharedModule := NewSharedModule(db)
	masterModule := NewMasterModule(db, sharedModule)
	settingsModule := NewSettingsModule(db, sharedModule, masterModule)
	taskModule := NewTaskModule(db)
	commonModule := NewCommonModule(db, settingsModule)
	libraryModule := NewLibraryModule(db, sharedModule, masterModule)
	organizationModule := NewOrganizationModule(db, sharedModule, settingsModule, masterModule, taskModule)
	meetingModule := NewMeetingModule(db, sharedModule, organizationModule, masterModule, commonModule)
	authModule := NewAuthModule(db, organizationModule)
	complianceModule := NewComplianceModule(db, sharedModule, masterModule, organizationModule, taskModule, settingsModule)
	riskModule := NewRiskModule(db, sharedModule, settingsModule, masterModule, libraryModule, organizationModule, complianceModule)
	incidentModule := NewIncidentModule(db, sharedModule, settingsModule, masterModule, organizationModule, complianceModule, riskModule)
	documentModule := NewDocumentModule(db, sharedModule, settingsModule, masterModule, organizationModule)
	auditModule := NewAuditModule(db,sharedModule,organizationModule,complianceModule,riskModule,incidentModule,documentModule)

	return &App{
		SettingsModule:     settingsModule,
		CommonModule:       commonModule,
		MasterModule:       masterModule,
		MeetingModule:      meetingModule,
		LibraryModule:      libraryModule,
		OrganizationModule: organizationModule,
		AuthModule:         authModule,
		TaskModule:         taskModule,
		ComplianceModule:   complianceModule,
		RiskModule:         riskModule,
		IncidentModule:     incidentModule,
		DocumentModule:     documentModule,
		AuditModule:        auditModule,
	}
}
