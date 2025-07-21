package initializers

import (
	"jethur/grc/project/controllers"
	"jethur/grc/project/repositories"
	"jethur/grc/project/services"

	"gorm.io/gorm"
)

type ProjectModule struct {
	ProjectTypeController                   *controllers.ProjectTypeController
	ProjectStatusController                 *controllers.ProjectStatusController
	ProjectController                       *controllers.ProjectController
	ProjectMilestoneController              *controllers.ProjectMilestoneController
	ProjectActionPlanController             *controllers.ProjectActionPlanController
	ProjectActionPlanStatusUpdateController *controllers.ProjectActionPlanStatusUpdateController
	ProjectActionPlanCommentController      *controllers.ProjectActionPlanCommentController
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
		masterModule.SeverityLevelSvc,
	)
	projectActionPlanController := controllers.NewProjectActionPlanController(projectActionPlanSvc)

	projectActionPlanStatusUpdateRepo := repositories.NewProjectActionPlanStatusUpdateRepository(db)
	projectActionPlanStatusUpdateSvc := services.NewProjectActionPlanStatusUpdateService(
		projectActionPlanStatusUpdateRepo,
		sharedModule.ActivityLogService,
		projectSvc,
		projectActionPlanSvc,
		masterModule.ActionPlanStatusSvc,
		organizationModule.UserSvc,
	)
	projectActionPlanStatusUpdateController := controllers.NewProjectActionPlanStatusUpdateController(projectActionPlanStatusUpdateSvc)

	projectActionPlanCommentRepo := repositories.NewProjectActionPlanCommentRepository(db)
	projectActionPlanCommentSvc := services.NewProjectActionPlanCommentService(
		projectActionPlanCommentRepo,
		organizationModule.UserSvc,
		sharedModule.ActivityLogService,
		projectActionPlanSvc,
	)
	projectActionPlanCommentController := controllers.NewProjectActionPlanCommentController(projectActionPlanCommentSvc)

	return &ProjectModule{
		ProjectTypeController:                   projectTypeController,
		ProjectStatusController:                 projectStatusController,
		ProjectController:                       projectController,
		ProjectMilestoneController:              projectMilestoneController,
		ProjectActionPlanController:             projectActionPlanController,
		ProjectActionPlanStatusUpdateController: projectActionPlanStatusUpdateController,
		ProjectActionPlanCommentController:      projectActionPlanCommentController,
	}
}
