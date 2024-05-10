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
}

func NewApp(db *gorm.DB) *App {
	sharedModule := NewSharedModule(db)
	settingsModule := NewSettingsModule(db, sharedModule)
	masterModule := NewMasterModule(db, sharedModule)
	libraryModule := NewLibraryModule(db, sharedModule, MasterModule)
	organizationModule := NewOrganizationModule(db, sharedModule, settingsModule, masterModule)
	authModule := NewAuthModule(db, organizationModule)
	complianceModule := NewComplianceModule(db, sharedModule, masterModule, organizationModule)

	return &App{
		SettingsModule:     settingsModule,
		MasterModule:       masterModule,
		LibraryModule:      libraryModule,
		OrganizationModule: organizationModule,
		AuthModule:         authModule,
		ComplianceModule:   complianceModule,
	}
}
