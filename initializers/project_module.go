package initializers

import (
	"jethur/grc/project/controllers"
	"jethur/grc/project/repositories"
	"jethur/grc/project/services"

	"gorm.io/gorm"
)

type ProjectModule struct {
	ProjectTypeController   *controllers.ProjectTypeController
	ProjectStatusController *controllers.ProjectStatusController
	ProjectController       *controllers.ProjectController
}

func NewProjectModule(db *gorm.DB, sharedModule *SharedModule, organizationModule *OrganizationModule) *ProjectModule {
	projectTypeRepo := repositories.NewProjectTypeRepository(db)
	projectTypeSvc := services.NewProjectTypeService(projectTypeRepo, sharedModule.ActivityLogService)
	projectTypeController := controllers.NewProjectTypeController(projectTypeSvc)

	projectStatusRepo := repositories.NewProjectStatusRepository(db)
	projectStatusSvc := services.NewProjectStatusService(projectStatusRepo)
	projectStatusController := controllers.NewProjectStatusController(projectStatusSvc)

	projectRepo := repositories.NewProjectRepository(db)
	projectSvc := services.NewProjectService(
		projectRepo,
		sharedModule.ActivityLogService,
		projectTypeSvc,
		projectStatusSvc,
		organizationModule.UserSvc,
	)
	projectController := controllers.NewProjectController(projectSvc)

	return &ProjectModule{
		ProjectTypeController:   projectTypeController,
		ProjectStatusController: projectStatusController,
		ProjectController:       projectController,
	}
}
