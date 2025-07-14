package initializers

import (
	"jethur/grc/project/controllers"
	"jethur/grc/project/repositories"
	"jethur/grc/project/services"

	"gorm.io/gorm"
)

type ProjectModule struct {
	ProjectTypeController       *controllers.ProjectTypeController
	ProjectStatusController     *controllers.ProjectStatusController
	ProjectController           *controllers.ProjectController
	ProjectMilestoneController  *controllers.ProjectMilestoneController
	ProjectActionPlanController *controllers.ProjectActionPlanController
}

func NewProjectModule(db *gorm.DB, sharedModule *SharedModule, organizationModule *OrganizationModule, masterModule *MasterModule) *ProjectModule {
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

	projectMilestoneRepo := repositories.NewProjectMilestoneRepository(db)
	projectMilestoneSvc := services.NewProjectMilestoneService(
		projectMilestoneRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		projectSvc,
		masterModule.ProjectMilestoneStatusSvc,
	)
	projectMilestoneController := controllers.NewProjectMilestoneController(projectMilestoneSvc)

	projectActionPlanRepo := repositories.NewProjectActionPlanRepository(db)
	projectActionPlanSvc := services.NewProjectActionPlanService(
		projectActionPlanRepo,
		sharedModule.ActivityLogService,
		masterModule.ActionPlanStatusSvc,
		organizationModule.UserSvc,
		projectSvc,
		projectMilestoneSvc,
	)
	projectActionPlanController := controllers.NewProjectActionPlanController(projectActionPlanSvc)

	return &ProjectModule{
		ProjectTypeController:       projectTypeController,
		ProjectStatusController:     projectStatusController,
		ProjectController:           projectController,
		ProjectMilestoneController:  projectMilestoneController,
		ProjectActionPlanController: projectActionPlanController,
	}
}
