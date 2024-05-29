package initializers

import (
	"gorm.io/gorm"
)

type App struct {
	SettingsModule     *SettingsModule
	MasterModule       *MasterModule
	MeetingModule      *MeetingModule
	LibraryModule      *LibraryModule
	OrganizationModule *OrganizationModule
	AuthModule         *AuthModule
	ComplianceModule   *ComplianceModule
	RiskModule         *RiskModule
	IncidentModule     *IncidentModule
	TaskModule         *TaskModule
}

func NewApp(db *gorm.DB) *App {
	sharedModule := NewSharedModule(db)
	masterModule := NewMasterModule(db, sharedModule)
	settingsModule := NewSettingsModule(db, sharedModule, masterModule)
	libraryModule := NewLibraryModule(db, sharedModule, masterModule)
	organizationModule := NewOrganizationModule(db, sharedModule, settingsModule, masterModule)
	meetingModule := NewMeetingModule(db, sharedModule, organizationModule, masterModule)
	authModule := NewAuthModule(db, organizationModule)
	taskModule := NewTaskModule(db)
	complianceModule := NewComplianceModule(db, sharedModule, masterModule, organizationModule,taskModule)
	riskModule := NewRiskModule(db, sharedModule, settingsModule, masterModule, libraryModule, organizationModule)
	incidentModule := NewIncidentModule(db, sharedModule, settingsModule, masterModule, organizationModule, complianceModule, riskModule)

	return &App{
		SettingsModule:     settingsModule,
		MasterModule:       masterModule,
		MeetingModule:      meetingModule,
		LibraryModule:      libraryModule,
		OrganizationModule: organizationModule,
		AuthModule:         authModule,
		TaskModule:         taskModule,
		ComplianceModule:   complianceModule,
		RiskModule:         riskModule,
		IncidentModule:     incidentModule,
	}
}
