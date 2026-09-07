package initializers

import (
	"jethur/grc/project/controllers"
	"jethur/grc/project/repositories"
	"jethur/grc/project/services"

	"gorm.io/gorm"
)

type ProjectModule struct {
	ProjectTypeController                                *controllers.ProjectTypeController
	ProjectStatusController                              *controllers.ProjectStatusController
	ProjectController                                    *controllers.ProjectController
	ProjectMilestoneController                           *controllers.ProjectMilestoneController
	ProjectActionPlanController                          *controllers.ProjectActionPlanController
	ProjectActionPlanStatusUpdateController              *controllers.ProjectActionPlanStatusUpdateController
	ProjectActionPlanCommentController                   *controllers.ProjectActionPlanCommentController
	ProjectAssociatedItemController                      *controllers.ProjectAssociatedItemController
	ProjectClassificationController                      *controllers.ProjectClassificationController
	ProjectPhaseController                               *controllers.ProjectPhaseController
	ProjectChangeRequestController                       *controllers.ProjectChangeRequestController
	ProjectChangeRequestWorkflowController               *controllers.ProjectChangeRequestWorkflowController
	ProjectChangeRequestReviewController                 *controllers.ProjectChangeRequestReviewController
	ProjectChangeRequestReviewCommentController          *controllers.ProjectChangeRequestReviewCommentController
	ProjectImpactAssessmentTypeController                *controllers.ProjectImpactAssessmentTypeController
	ProjectSecurityGateClassificationController          *controllers.ProjectSecurityGateClassificationController
	ProjectSecurityGateStageController                   *controllers.ProjectSecurityGateStageController
	ProjectSecurityGateChecklistTemplateController       *controllers.ProjectSecurityGateChecklistTemplateController
	ProjectSecurityGateChecklistResponseStatusController *controllers.ProjectSecurityGateChecklistResponseStatusController
	ProjectWorkflowLevelController                       *controllers.ProjectWorkflowLevelController
	ProjectSecurityGateReadinessStatusController         *controllers.ProjectSecurityGateReadinessStatusController
	ProjectSecurityGateController                        *controllers.ProjectSecurityGateController
	ProjectSecurityGateChecklistController               *controllers.ProjectSecurityGateChecklistController

	projectSvc                           services.ProjectServiceInterface
	projectMilestoneSvc                  services.ProjectMilestoneServiceInterface
	projectActionPlanSvc                 services.ProjectActionPlanServiceInterface
	projectChangeRequestSvc              services.ProjectChangeRequestServiceInterface
	projectChangeRequestWorkflowSvc      services.ProjectChangeRequestWorkflowServiceInterface
	projectChangeRequestReviewSvc        services.ProjectChangeRequestReviewServiceInterface
	projectChangeRequestReviewCommentSvc services.ProjectChangeRequestReviewCommentServiceInterface
	projectImpactAssessmentTypeSvc       services.ProjectImpactAssessmentTypeServiceInterface
}

func NewProjectModule(db *gorm.DB, sharedModule *SharedModule, organizationModule *OrganizationModule, masterModule *MasterModule, commonModule *CommonModule, taskModule *TaskModule, settingsModule *SettingsModule) *ProjectModule {
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
		organizationModule.UserSettingSvc,
		commonModule.NotificationService,
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

	projectChangeRequestRepo := repositories.NewProjectChangeRequestRepository(db)
	projectChangeRequestSvc := services.NewProjectChangeRequestService(
		projectChangeRequestRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		projectSvc,
	)
	projectChangeRequestController := controllers.NewProjectChangeRequestController(projectChangeRequestSvc)

	projectChangeRequestWorkflowRepo := repositories.NewProjectChangeRequestWorkflowRepository(db)
	projectChangeRequestReviewRepo := repositories.NewProjectChangeRequestReviewRepository(db)
	projectChangeRequestWorkflowSvc := services.NewProjectChangeRequestWorkflowService(
		projectChangeRequestWorkflowRepo,
		projectChangeRequestReviewRepo,
		sharedModule.ActivityLogService,
		settingsModule.WorkflowSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		masterModule.WorkflowContextSvc,
		commonModule.NotificationService,
		organizationModule.OrganizationUnitsSvc,
	)
	projectChangeRequestWorkflowSvc.SetProjectChangeRequestService(projectChangeRequestSvc)
	projectChangeRequestSvc.SetWorkflowService(projectChangeRequestWorkflowSvc)
	projectChangeRequestWorkflowController := controllers.NewProjectChangeRequestWorkflowController(projectChangeRequestWorkflowSvc, db)

	projectChangeRequestReviewSvc := services.NewProjectChangeRequestReviewService(
		projectChangeRequestReviewRepo,
		sharedModule.ActivityLogService,
		projectChangeRequestSvc,
		projectChangeRequestWorkflowSvc,
		masterModule.WorkflowStatusSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
	)
	projectChangeRequestReviewController := controllers.NewProjectChangeRequestReviewController(projectChangeRequestReviewSvc, db)

	projectChangeRequestReviewCommentRepo := repositories.NewProjectChangeRequestReviewCommentRepository(db)
	projectChangeRequestReviewCommentSvc := services.NewProjectChangeRequestReviewCommentService(
		projectChangeRequestReviewCommentRepo,
		organizationModule.UserSvc,
		sharedModule.ActivityLogService,
		projectChangeRequestSvc,
		masterModule.CommentStatusSvc,
	)
	projectChangeRequestReviewCommentController := controllers.NewProjectChangeRequestReviewCommentController(projectChangeRequestReviewCommentSvc, db)

	projectActionPlanRepo := repositories.NewProjectActionPlanRepository(db)
	projectActionPlanSvc := services.NewProjectActionPlanService(
		projectActionPlanRepo,
		sharedModule.ActivityLogService,
		masterModule.ActionPlanStatusSvc,
		organizationModule.UserSvc,
		projectSvc,
		projectMilestoneSvc,
		masterModule.SeverityLevelSvc, commonModule.DelegateeUserService,
		taskModule.TaskService,
		commonModule.NotificationService,
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
		commonModule.NotificationService,
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

	projectClassificationRepo := repositories.NewProjectClassificationRepository(db)
	projectClassificationSvc := services.NewProjectClassificationService(projectClassificationRepo, sharedModule.ActivityLogService)
	projectClassificationController := controllers.NewProjectClassificationController(projectClassificationSvc)

	projectPhaseRepo := repositories.NewProjectPhaseRepository(db)
	projectPhaseSvc := services.NewProjectPhaseService(projectPhaseRepo)
	projectPhaseController := controllers.NewProjectPhaseController(projectPhaseSvc)

	projectImpactAssessmentTypeRepo := repositories.NewProjectImpactAssessmentTypeRepository(db)
	projectImpactAssessmentTypeSvc := services.NewProjectImpactAssessmentTypeService(projectImpactAssessmentTypeRepo)
	projectImpactAssessmentTypeController := controllers.NewProjectImpactAssessmentTypeController(projectImpactAssessmentTypeSvc)

	projectSecurityGateClassificationRepo := repositories.NewProjectSecurityGateClassificationRepository(db)
	projectSecurityGateClassificationSvc := services.NewProjectSecurityGateClassificationService(projectSecurityGateClassificationRepo)
	projectSecurityGateClassificationController := controllers.NewProjectSecurityGateClassificationController(projectSecurityGateClassificationSvc)

	projectSecurityGateStageRepo := repositories.NewProjectSecurityGateStageRepository(db)
	projectSecurityGateStageSvc := services.NewProjectSecurityGateStageService(projectSecurityGateStageRepo)
	projectSecurityGateStageController := controllers.NewProjectSecurityGateStageController(projectSecurityGateStageSvc)

	projectSecurityGateChecklistTemplateRepo := repositories.NewProjectSecurityGateChecklistTemplateRepository(db)
	projectSecurityGateChecklistTemplateSvc := services.NewProjectSecurityGateChecklistTemplateService(
		projectSecurityGateChecklistTemplateRepo,
		sharedModule.ActivityLogService,
	)
	projectSecurityGateChecklistTemplateController := controllers.NewProjectSecurityGateChecklistTemplateController(projectSecurityGateChecklistTemplateSvc)

	projectWorkflowLevelController := controllers.NewProjectWorkflowLevelController(settingsModule.WorkflowSvc, db)

	projectSecurityGateChecklistResponseStatusRepo := repositories.NewProjectSecurityGateChecklistResponseStatusRepository(db)
	projectSecurityGateChecklistResponseStatusSvc := services.NewProjectSecurityGateChecklistResponseStatusService(projectSecurityGateChecklistResponseStatusRepo)
	projectSecurityGateChecklistResponseStatusController := controllers.NewProjectSecurityGateChecklistResponseStatusController(projectSecurityGateChecklistResponseStatusSvc)

	projectSecurityGateReadinessStatusRepo := repositories.NewProjectSecurityGateReadinessStatusRepository(db)
	projectSecurityGateReadinessStatusSvc := services.NewProjectSecurityGateReadinessStatusService(projectSecurityGateReadinessStatusRepo)
	projectSecurityGateReadinessStatusController := controllers.NewProjectSecurityGateReadinessStatusController(projectSecurityGateReadinessStatusSvc)

	projectSecurityGateRepo := repositories.NewProjectSecurityGateRepository(db)
	projectSecurityGateSvc := services.NewProjectSecurityGateService(
		projectSecurityGateRepo,
		projectSecurityGateStageSvc,
		projectSecurityGateClassificationSvc,
		projectSecurityGateChecklistTemplateSvc,
		projectSecurityGateReadinessStatusSvc,
		organizationModule.UserSvc,
		sharedModule.ActivityLogService,
	)
	projectSecurityGateController := controllers.NewProjectSecurityGateController(projectSecurityGateSvc)

	projectSecurityGateChecklistRepo := repositories.NewProjectSecurityGateChecklistRepository(db)
	projectSecurityGateChecklistSvc := services.NewProjectSecurityGateChecklistService(
		projectSecurityGateChecklistRepo,
		projectSecurityGateChecklistResponseStatusSvc,
		sharedModule.ActivityLogService,
	)
	projectSecurityGateChecklistController := controllers.NewProjectSecurityGateChecklistController(projectSecurityGateChecklistSvc)

	projectSecurityGateSvc.SetChecklistService(projectSecurityGateChecklistSvc)
	projectSecurityGateChecklistSvc.SetSecurityGateService(projectSecurityGateSvc)

	return &ProjectModule{
		ProjectTypeController:                                projectTypeController,
		ProjectStatusController:                              projectStatusController,
		ProjectController:                                    projectController,
		ProjectMilestoneController:                           projectMilestoneController,
		ProjectActionPlanController:                          projectActionPlanController,
		ProjectActionPlanStatusUpdateController:              projectActionPlanStatusUpdateController,
		ProjectActionPlanCommentController:                   projectActionPlanCommentController,
		ProjectAssociatedItemController:                      projectAssociatedItemController,
		ProjectClassificationController:                      projectClassificationController,
		ProjectPhaseController:                               projectPhaseController,
		ProjectChangeRequestController:                       projectChangeRequestController,
		ProjectChangeRequestWorkflowController:               projectChangeRequestWorkflowController,
		ProjectChangeRequestReviewController:                 projectChangeRequestReviewController,
		ProjectChangeRequestReviewCommentController:          projectChangeRequestReviewCommentController,
		ProjectImpactAssessmentTypeController:                projectImpactAssessmentTypeController,
		ProjectSecurityGateClassificationController:          projectSecurityGateClassificationController,
		ProjectSecurityGateStageController:                   projectSecurityGateStageController,
		ProjectSecurityGateChecklistTemplateController:       projectSecurityGateChecklistTemplateController,
		ProjectSecurityGateChecklistResponseStatusController: projectSecurityGateChecklistResponseStatusController,
		ProjectWorkflowLevelController:                       projectWorkflowLevelController,
		ProjectSecurityGateReadinessStatusController:         projectSecurityGateReadinessStatusController,
		ProjectSecurityGateController:                        projectSecurityGateController,
		ProjectSecurityGateChecklistController:               projectSecurityGateChecklistController,

		projectSvc:                           projectSvc,
		projectMilestoneSvc:                  projectMilestoneSvc,
		projectActionPlanSvc:                 projectActionPlanSvc,
		projectChangeRequestSvc:              projectChangeRequestSvc,
		projectChangeRequestWorkflowSvc:      projectChangeRequestWorkflowSvc,
		projectChangeRequestReviewSvc:        projectChangeRequestReviewSvc,
		projectChangeRequestReviewCommentSvc: projectChangeRequestReviewCommentSvc,
		projectImpactAssessmentTypeSvc:       projectImpactAssessmentTypeSvc,
	}
}
