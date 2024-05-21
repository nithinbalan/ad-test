package initializers

import (
	"gorm.io/gorm"
)

type App struct {
	SettingsModule     *SettingsModule
	MasterModule       *MasterModule
	LibraryModule      *LibraryModule
	OrganizationModule *OrganizationModule
	AuthModule         *AuthModule
	ComplianceModule   *ComplianceModule
	RiskModule         *RiskModule
	IncidentModule     *IncidentModule
}

func NewApp(db *gorm.DB) *App {
	sharedModule := NewSharedModule(db)
	masterModule := NewMasterModule(db, sharedModule)
	settingsModule := NewSettingsModule(db, sharedModule, masterModule)
	libraryModule := NewLibraryModule(db, sharedModule, masterModule)
	organizationModule := NewOrganizationModule(db, sharedModule, settingsModule, masterModule)
	authModule := NewAuthModule(db, organizationModule)
	complianceModule := NewComplianceModule(db, sharedModule, masterModule, organizationModule)
	riskModule := NewRiskModule(db, sharedModule, settingsModule, masterModule, libraryModule, organizationModule)
	incidentModule := NewIncidentModule(db, sharedModule, settingsModule, masterModule, organizationModule, complianceModule, riskModule)

	return &App{
		SettingsModule:     settingsModule,
		MasterModule:       masterModule,
		LibraryModule:      libraryModule,
		OrganizationModule: organizationModule,
		AuthModule:         authModule,
		ComplianceModule:   complianceModule,
		RiskModule:         riskModule,
		IncidentModule:     incidentModule,
	}
}
