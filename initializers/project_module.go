package initializers

import (
	"jethur/grc/project/controllers"
	"jethur/grc/project/repositories"
	"jethur/grc/project/services"

	"gorm.io/gorm"
)

type ProjectModule struct {
	ProjectTypeController *controllers.ProjectTypeController
}

func NewProjectModule(db *gorm.DB, sharedModule *SharedModule, masterModule *MasterModule, organizationModule *OrganizationModule, settingsModule *SettingsModule, commonModule *CommonModule) *ProjectModule {
	projectTypeRepo := repositories.NewProjectTypeRepository(db)
	projectTypeSvc := services.NewProjectTypeService(projectTypeRepo, sharedModule.ActivityLogService)
	projectTypeController := controllers.NewProjectTypeController(projectTypeSvc)

	return &ProjectModule{
		ProjectTypeController: projectTypeController,
	}
}
