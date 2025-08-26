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
	ProjectAssociatedItemController         *controllers.ProjectAssociatedItemController

	projectSvc           services.ProjectServiceInterface
	projectMilestoneSvc  services.ProjectMilestoneServiceInterface
	projectActionPlanSvc services.ProjectActionPlanServiceInterface
}

func NewProjectModule(db *gorm.DB, sharedModule *SharedModule, organizationModule *OrganizationModule, masterModule *MasterModule, commonModule *CommonModule) *ProjectModule {
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
		masterModule.SeverityLevelSvc, commonModule.DelegateeUserService,
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

	projectAssociatedItemRepo := repositories.NewProjectAssociatedItemRepository(db)
	projectAssociatedItemSvc := services.NewProjectAssociatedItemService(projectAssociatedItemRepo, projectSvc, commonModule.AssociatedItemValidationService)
	projectAssociatedItemController := controllers.NewProjectAssociatedItemController(projectAssociatedItemSvc)

	return &ProjectModule{
		ProjectTypeController:                   projectTypeController,
		ProjectStatusController:                 projectStatusController,
		ProjectController:                       projectController,
		ProjectMilestoneController:              projectMilestoneController,
		ProjectActionPlanController:             projectActionPlanController,
		ProjectActionPlanStatusUpdateController: projectActionPlanStatusUpdateController,
		ProjectActionPlanCommentController:      projectActionPlanCommentController,
		ProjectAssociatedItemController:         projectAssociatedItemController,

		projectSvc:           projectSvc,
		projectMilestoneSvc:  projectMilestoneSvc,
		projectActionPlanSvc: projectActionPlanSvc,
	}
}
