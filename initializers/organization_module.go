package initializers

import (
	"jethur/grc/organization/controllers"
	"jethur/grc/organization/repositories"
	"jethur/grc/organization/services"

	"gorm.io/gorm"
)

type OrganizationModule struct {
	DesignationController                     *controllers.DesignationController
	UserStatusController                      *controllers.UserStatusController
	UserController                            *controllers.UserController
	UserSettingController                     *controllers.UserSettingController
	StrategicProfileController                *controllers.StrategicProfileController
	StrategicObjectiveController              *controllers.StrategicOBjectiveController
	ControlCategoryController                 *controllers.ControlCategoryController
	ControlTypeController                     *controllers.ControlTypeController
	ControlImplementationStatusController     *controllers.ControlImplementationStatusController
	ControlEfficiencyStatusController         *controllers.ControlEfficiencyStatusController
	ControlReviewStatusController             *controllers.ControlReviewStatusController
	ControlController                         *controllers.ControlController
	ControlTestController                     *controllers.ControlTestController
	ProcessGroupController                    *controllers.ProcessGroupController
	ControlActionPlanController               *controllers.ControlActionPlanController
	ControlActionPlanCommentController        *controllers.ControlActionPlanCommentController
	ControlActionPlanStatusUpdateController   *controllers.ControlActionPlanStatusUpdateController
	ProcessController                         *controllers.ProcessController
	UserTaskController                        *controllers.UserTaskController
	StrategicProfileValueController           *controllers.StrategicProfileValueController
	ProcessAssociatedItemController           *controllers.ProcessAssociatedItemController
	ControlAssociatedItemController           *controllers.ControlAssociatedItemController
	ControlReviewController                   *controllers.ControlReviewController
	UserProfileController                     *controllers.UserProfileController
	ProcessDependencyTypeController           *controllers.ProcessDependencyTypeController
	ProcessDependencySubTypeController        *controllers.ProcessDependencySubTypeController
	GovernanceReportController                *controllers.GovernanceReportController
	ProcessDependencyController               *controllers.ProcessDependencyController
	ProcessCriticalityLevelController         *controllers.ProcessCriticalityLevelController
	ProcessActivityController                 *controllers.ProcessActivityController
	CompetencyGroupController                 *controllers.CompetencyGroupController
	CompetencyController                      *controllers.CompetencyController
	CompetencyLevelController                 *controllers.CompetencyLevelController
	JobDescriptionController                  *controllers.JobDescriptionController
	JobDescriptionCompetencyController        *controllers.JobDescriptionCompetencyController
	UserDocumentTypeController                *controllers.UserDocumentTypeController
	UserDocumentSubTypeController             *controllers.UserDocumentSubTypeController
	UserDocumentController                    *controllers.UserDocumentController
	UserJobDescriptionController              *controllers.UserJobDescriptionController
	UserDelegationHistoryController           *controllers.UserDelegationController
	OrganizationUnitsController               *controllers.OrganizationUnitsController
	OrganizationUnitBusinessContextController *controllers.OrganizationUnitBusinessContextController
	OrganizationUnitAttributeController       *controllers.OrganizationUnitAttributeController
	OrganizationCartController                *controllers.OrganizationChartController
	MyDashboardController                     *controllers.MyDashboardController
	CommitteeTermController                   *controllers.CommitteeTermController
	CommitteeTypeController                   *controllers.CommitteeTypeController
	CommitteeClassificationController         *controllers.CommitteeClassificationController
	InternalCommitteeStatusController         *controllers.InternalCommitteeStatusController
	InternalCommitteeController               *controllers.InternalCommitteeController
	InternalCommitteeWorkflowController       *controllers.InternalCommitteeWorkflowController
	InternalCommitteeWorkflowReviewController *controllers.InternalCommitteeWorkflowReviewController
	InternalCommitteeReviewCommentController  *controllers.InternalCommitteeReviewCommentController
	ControlEffectivenessRatingController      *controllers.ControlEffectivenessRatingController
	ControlStatusController                   *controllers.ControlStatusController

	UserSvc                        services.UserServiceInterface
	UserSettingSvc                 services.UserSettingServiceInterface
	ControlSvc                     services.ControlServiceInterface
	ProcessSvc                     services.ProcessServiceInterface
	StrategicObjectiveSvc          services.StrategicObjectiveServiceInterface
	ControlTestSvc                 services.ControlTestServiceInterface
	ControlImplementationStatusSvc services.ControlImplementationStatusServiceInterface
	ControlEfficiencyStatusSvc     services.ControlEfficiencyStatusServiceInterface
	ProcessAssociatedItemSvc       services.ProcessAssociatedItemServiceInterface
	ProcessDependencySvc           services.ProcessDependencyServiceInterface
	JobDescriptionSvc              services.JobDescriptionServiceInterface
	GovernanceReportSvc            services.GovernanceReportServiceInterface
	ControlStatusSvc               services.ControlStatusServiceInterface
}

func NewOrganizationModule(
	db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	taskModule *TaskModule,
	commonModule *CommonModule,
) *OrganizationModule {
	// Competency Group Setup
	competencyGroupRepo := repositories.NewCompetencyGroupRepository(db)
	competencyGroupSvc := services.NewCompetencyGroupService(competencyGroupRepo, sharedModule.ActivityLogService)
	competencyGroupController := controllers.NewCompetencyGroupController(competencyGroupSvc)

	// Competency Setup
	competencyRepo := repositories.NewCompetencyRepository(db)
	competencySvc := services.NewCompetencyService(competencyRepo, sharedModule.ActivityLogService)
	competencyController := controllers.NewCompetencyController(competencySvc)

	competencyLevelRepo := repositories.NewCompetencyLevelRepository(db)
	competencyLevelSvc := services.NewCompetencyLevelService(competencyLevelRepo)
	competencyLevelController := controllers.NewCompetencyLevelController(competencyLevelSvc)

	jobDescriptionRepo := repositories.NewJobDescriptionRepository(db)
	jobDescriptionSvc := services.NewJobDescriptionService(jobDescriptionRepo, sharedModule.ActivityLogService)
	jobDescriptionController := controllers.NewJobDescriptionController(jobDescriptionSvc)

	jobDescriptionCompetencyRepo := repositories.NewJobDescriptionCompetencyRepository(db)
	jobDescriptionCompetencySvc := services.NewJobDescriptionCompetencyService(jobDescriptionCompetencyRepo, jobDescriptionSvc, competencyLevelSvc, competencySvc, sharedModule.ActivityLogService)
	jobDescriptionCompetencyController := controllers.NewJobDescriptionCompetencyController(jobDescriptionCompetencySvc)

	designationRepo := repositories.NewDesignationRepository(db)
	designationSvc := services.NewDesignationService(designationRepo, sharedModule.ActivityLogService, jobDescriptionSvc)
	designationController := controllers.NewDesignationController(designationSvc)

	userJobDescriptionRepo := repositories.NewUserJobDescriptionRepository(db)
	userJobDescriptionSvc := services.NewUserJobDescriptionService(userJobDescriptionRepo, sharedModule.ActivityLogService, jobDescriptionSvc, designationSvc)
	userJobDescriptionController := controllers.NewUserJobDescriptionController(userJobDescriptionSvc)

	userStatusRepo := repositories.NewUserStatusRepository(db)
	userStatusSvc := services.NewUserStatusService(userStatusRepo)
	userStatusController := controllers.NewUserStatusController(userStatusSvc)

	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(
		userRepo,
		designationSvc,
		userStatusSvc,
		sharedModule.ActivityLogService,
		settingsModule.RoleSvc,
		settingsModule.SystemPreferenceSvc,
		settingsModule.OrganizationLevelSvc,
		settingsModule.ModuleSvc,
		settingsModule.EmailTemplateSettingsSvc,
		commonModule.NotificationService,
		settingsModule.ActiveDirectorySettingSvc,
		settingsModule.OrganizationUnitSvc,
		userJobDescriptionSvc,
		masterModule.LicenseTypeSvc,
	)
	userSettingRepo := repositories.NewUserSettingRepository(db)
	userSettingSvc := services.NewUserSettingService(
		userSettingRepo,
		sharedModule.ActivityLogService,
		commonModule.NotificationService,
		settingsModule.EmailTemplateSettingsSvc)
	userSvc.SetUserSettingService(userSettingSvc)
	userSettingSvc.SetUserService(userSvc)
	userController := controllers.NewUserController(userSvc)
	userSettingController := controllers.NewUserSettingController(userSettingSvc)

	strategicProfileRepo := repositories.NewStrategicProfileRepository(db)
	strategicProfileSvc := services.NewStrategicProfileService(strategicProfileRepo, commonModule.FileRepo, sharedModule.ActivityLogService)
	strategicProfileController := controllers.NewStrategicProfileController(strategicProfileSvc)

	strategicObjectiveRepo := repositories.NewStrategicObjectiveRepository(db)
	strategicObjectiveSvc := services.NewStrategicObjectiveService(strategicObjectiveRepo, sharedModule.ActivityLogService)
	strategicObjectiveController := controllers.NewStrategicObjectiveController(strategicObjectiveSvc)

	controlCategoryRepo := repositories.NewControlCategoryRepository(db)
	controlCategorySvc := services.NewControlCategoryService(controlCategoryRepo, sharedModule.ActivityLogService)
	controlCategoryController := controllers.NewControlCategoryController(controlCategorySvc)

	controlTypeRepo := repositories.NewControlTypeRepository(db)
	controlTypeSvc := services.NewControlTypeService(controlTypeRepo)
	controlTypeController := controllers.NewControlTypeController(controlTypeSvc)

	controlImplementationStatusRepo := repositories.NewControlImplementationStatusRepository(db)
	controlImplementationStatusSvc := services.NewControlImplementationStatusService(controlImplementationStatusRepo)
	controlImplementationStatusController := controllers.NewControlImplementationStatusController(controlImplementationStatusSvc)

	controlEfficiencyStatusRepo := repositories.NewControlEfficiencyStatusRepository(db)
	controlEfficiencyStatusSvc := services.NewControlEfficiencyStatusService(controlEfficiencyStatusRepo)
	controlEfficiencyStatusController := controllers.NewControlEfficiencyStatusController(controlEfficiencyStatusSvc)

	strategicProfileValueRepo := repositories.NewStrategicProfileValueRepository(db)
	strategicProfileValueSvc := services.NewStrategicProfileValueService(strategicProfileValueRepo, sharedModule.ActivityLogService, strategicProfileRepo)
	strategicProfileValueController := controllers.NewStrategicProfileValueController(strategicProfileValueSvc)

	controlEffectivenessRatingRepo := repositories.NewControlEffectivenessRatingRepository(db)
	controlEffectivenessRatingSvc := services.NewControlEffectivenessRatingService(controlEffectivenessRatingRepo)
	controlEffectivenessRatingController := controllers.NewControlEffectivenessRatingController(controlEffectivenessRatingSvc)

	controlReviewStatusRepo := repositories.NewControlReviewStatusRepository(db)
	controlReviewStatusSvc := services.NewControlReviewStatusService(controlReviewStatusRepo)
	controlReviewStatusController := controllers.NewControlReviewStatusController(controlReviewStatusSvc)

	controlRepo := repositories.NewControlRepository(db)
	controlSvc := services.NewControlService(
		controlRepo,
		sharedModule.ActivityLogService,
		controlTypeSvc,
		controlCategorySvc,
		controlImplementationStatusSvc,
		controlEfficiencyStatusSvc,
		userSvc,
		masterModule.ReviewFrequencySvc,
		masterModule.LifecycleStatusSvc,
		commonModule.NotificationService, userSettingSvc)
	controlController := controllers.NewControlController(controlSvc)

	controlTestRepo := repositories.NewControlTestRepository(db)
	controlTestSvc := services.NewControlTestService(
		controlTestRepo,
		sharedModule.ActivityLogService,
		controlSvc,
		controlImplementationStatusSvc,
		controlEfficiencyStatusSvc,
		masterModule.ReviewFrequencySvc,
		controlEffectivenessRatingSvc)
	controlTestController := controllers.NewControlTestController(controlTestSvc)

	controlAssociatedItemRepo := repositories.NewControlAssociatedItemRepository(db)
	controlAssociatedItemSvc := services.NewControlAssociatedItemService(controlAssociatedItemRepo, controlSvc, controlImplementationStatusSvc, controlEfficiencyStatusRepo, commonModule.AssociatedItemValidationService)
	controlAssociatedItemController := controllers.NewControlAssociatedItemController(controlAssociatedItemSvc)

	organizationUnitsRepo := repositories.NewOrganizationUnitsRepository(db)
	organizationUnitsSvc := services.NewOrganizationUnitsService(
		organizationUnitsRepo,
		sharedModule.ActivityLogService,
		settingsModule.OrganizationLevelSvc,
		userRepo,
		userSvc,
		settingsModule.OrganizationUnitSvc,
		settingsModule.ReportTemplateSettingSvc,
	)
	organizationUnitsController := controllers.NewOrganizationUnitsController(organizationUnitsSvc)

	// Organization Unit Business Context Setup
	organizationUnitBusinessContextRepo := repositories.NewOrganizationUnitBusinessContextRepository(db)
	organizationUnitBusinessContextSvc := services.NewOrganizationUnitBusinessContextService(organizationUnitBusinessContextRepo, sharedModule.ActivityLogService)
	organizationUnitBusinessContextController := controllers.NewOrganizationUnitBusinessContextController(organizationUnitBusinessContextSvc)

	// Organization Unit Attribute Setup
	organizationUnitAttributeRepo := repositories.NewOrganizationUnitAttributeRepository(db)
	organizationUnitAttributeSvc := services.NewOrganizationUnitAttributeService(organizationUnitAttributeRepo, sharedModule.ActivityLogService)
	organizationUnitAttributeController := controllers.NewOrganizationUnitAttributeController(organizationUnitAttributeSvc)

	controlReviewRepo := repositories.NewControlReviewRepository(db)
	controlReviewSvc := services.NewControlReviewService(controlReviewRepo, sharedModule.ActivityLogService, controlSvc, masterModule.ReviewFrequencySvc, controlReviewStatusSvc)
	controlReviewController := controllers.NewControlReviewController(controlReviewSvc)

	processGroupRepo := repositories.NewProcessGroupRepository(db)
	processGroupSvc := services.NewProcessGroupService(processGroupRepo, sharedModule.ActivityLogService)
	processGroupController := controllers.NewProcessGroupController(processGroupSvc)

	controlActionPlanRepo := repositories.NewControlActionPlanRepository(db)
	controlActionPlanSvc := services.NewControlActionPlanService(controlActionPlanRepo, sharedModule.ActivityLogService, controlSvc, masterModule.SeverityLevelSvc, masterModule.ActionPlanStatusSvc, userSvc, commonModule.NotificationService, taskModule.TaskService, commonModule.DelegateeUserService)
	controlActionPlanController := controllers.NewControlActionPlanController(controlActionPlanSvc)

	controlActionPlanStatusUpdateRepo := repositories.NewControlActionPlanStatusUpdateRepository(db)
	controlActionPlanStatusUpdateSvc := services.NewControlActionPlanStatusUpdateService(controlActionPlanStatusUpdateRepo, sharedModule.ActivityLogService, controlActionPlanSvc, masterModule.ActionPlanStatusSvc, controlSvc, userSvc, commonModule.NotificationService, taskModule.TaskService, taskModule.TaskStatusService)
	controlActionPlanStatusUpdateController := controllers.NewControlActionPlanStatusUpdateController(controlActionPlanStatusUpdateSvc)

	controlActionPlanCommentRepository := repositories.NewControlActionPlanCommentRepository(db)
	controlActionPlanCommentSvc := services.NewControlActionPlanCommentService(controlActionPlanCommentRepository, userSvc, sharedModule.ActivityLogService, controlActionPlanSvc, commonModule.NotificationService)
	controlActionPlanCommentController := controllers.NewControlActionPlanCommentController(controlActionPlanCommentSvc)

	processRepo := repositories.NewProcessRepository(db)
	processSvc := services.NewProcessService(processRepo, sharedModule.ActivityLogService, processGroupSvc, userSvc, masterModule.LifecycleStatusSvc, controlSvc, commonModule.NotificationService, settingsModule.OrganizationUnitSvc, userSettingSvc, organizationUnitsSvc)
	processController := controllers.NewProcessController(processSvc)

	userTaskRepo := repositories.NewUserTaskRepository(db)
	userTaskSvc := services.NewUserTaskService(userTaskRepo, taskModule.TaskService)
	userTaskController := controllers.NewUserTaskController(userTaskSvc)

	processAssociatedItemRepo := repositories.NewProcessAssociatedItemRepository(db)
	processAssociatedItemSvc := services.NewProcessAssociatedItemService(processAssociatedItemRepo, processSvc, controlImplementationStatusSvc, controlEfficiencyStatusRepo, commonModule.AssociatedItemValidationService)
	processAssociatedItemController := controllers.NewProcessAssociatedItemController(processAssociatedItemSvc)

	userDocumentTypeRepo := repositories.NewUserDocumentTypeRepository(db)
	userDocumentTypeSvc := services.NewUserDocumentTypeService(userDocumentTypeRepo, sharedModule.ActivityLogService)
	userDocumentTypeController := controllers.NewUserDocumentTypeController(userDocumentTypeSvc)

	userDocumentSubTypeRepo := repositories.NewUserDocumentSubTypeRepository(db)
	userDocumentSubTypeSvc := services.NewUserDocumentSubTypeService(userDocumentSubTypeRepo, userDocumentTypeSvc, sharedModule.ActivityLogService)
	userDocumentSubTypeController := controllers.NewUserDocumentSubTypeController(userDocumentSubTypeSvc)

	userDocumentRepo := repositories.NewUserDocumentRepository(db)
	userDocumentSvc := services.NewUserDocumentService(userDocumentRepo, userDocumentTypeSvc, userDocumentSubTypeSvc, sharedModule.ActivityLogService, userSvc)
	userDocumentController := controllers.NewUserDocumentController(userDocumentSvc)

	userProfileController := controllers.NewUserProfileController(userSettingSvc, userTaskSvc, userSvc, userDocumentSvc, userJobDescriptionSvc)

	processDependencyTypeRepo := repositories.NewProcessDependencyTypeRepository(db)
	processDependencyTypeSvc := services.NewProcessDependencyTypeService(processDependencyTypeRepo)
	processDependencyTypeController := controllers.NewProcessDependencyTypeController(processDependencyTypeSvc)

	processDependencySubTypeRepo := repositories.NewProcessDependencySubTypeRepository(db)
	processDependencySubTypeSvc := services.NewProcessDependencySubTypeService(processDependencySubTypeRepo, processDependencyTypeSvc, sharedModule.ActivityLogService)
	processDependencySubTypeController := controllers.NewProcessDependencySubTypeController(processDependencySubTypeSvc)

	governanceReportRepo := repositories.NewGovernanceReportRepository(db)
	governanceReportSvc := services.NewGovernanceReportService(governanceReportRepo, sharedModule.ActivityLogService, settingsModule.OrganizationLevelSvc)
	governanceReportController := controllers.NewGovernanceReportController(governanceReportSvc)

	processDependencyRepo := repositories.NewProcessDependencyRepository(db)
	processDependencySvc := services.NewProcessDependencyService(processDependencyRepo, processDependencySubTypeSvc, sharedModule.ActivityLogService, processSvc, settingsModule.OrganizationUnitSvc)
	processDependencyController := controllers.NewProcessDependencyController(processDependencySvc)
	processCriticalityLevelRepo := repositories.NewProcessCriticalityLevelRepository(db)
	processCriticalityLevelSvc := services.NewProcessCriticalityLevelService(processCriticalityLevelRepo)
	processCriticalityLevelController := controllers.NewProcessCriticalityLevelController(processCriticalityLevelSvc)

	processActivityRepo := repositories.NewProcessActivityRepository(db)
	processActivitySvc := services.NewProcessActivityService(
		processActivityRepo,
		sharedModule.ActivityLogService,
		processSvc,
		userSvc,
		commonModule.NotificationService,
	)
	processActivityController := controllers.NewProcessActivityController(processActivitySvc)

	userDelegationHistoryRepo := repositories.NewUserDelegationRepository(db)
	userDelegationHistorySvc := services.NewUserDelegationService(userDelegationHistoryRepo, sharedModule.ActivityLogService, userSvc, commonModule.NotificationService, userSettingSvc, taskModule.TaskService)
	userDelegationHistoryController := controllers.NewUserDelegationController(userDelegationHistorySvc)

	organizationChartRepo := repositories.NewOrganizationChartRepository(db)
	organizationChartSvc := services.NewOrganizationChartService(organizationChartRepo, settingsModule.OrganizationLevelSvc)
	organizationChartController := controllers.NewOrganizationChartController(organizationChartSvc)

	myDashboardRepo := repositories.NewMyDashboardRepository(db)
	myDashboardSvc := services.NewMyDashboardService(myDashboardRepo)
	myDashboardController := controllers.NewMyDashboardController(myDashboardSvc, db)

	committeeTermRepo := repositories.NewCommitteeTermRepository(db)
	committeeTermSvc := services.NewCommitteeTermService(committeeTermRepo)
	committeeTermController := controllers.NewCommitteeTermController(committeeTermSvc)

	internalCommitteeStatusRepo := repositories.NewInternalCommitteeStatusRepository(db)
	internalCommitteeStatusSvc := services.NewInternalCommitteeStatusService(internalCommitteeStatusRepo)
	internalCommitteeStatusController := controllers.NewInternalCommitteeStatusController(internalCommitteeStatusSvc)

	committeeTypeRepo := repositories.NewCommitteeTypeRepository(db)
	committeeTypeSvc := services.NewCommitteeTypeService(committeeTypeRepo)
	committeeTypeController := controllers.NewCommitteeTypeController(committeeTypeSvc)

	committeeClassificationRepo := repositories.NewCommitteeClassificationRepository(db)
	committeeClassificationSvc := services.NewCommitteeClassificationService(committeeClassificationRepo, sharedModule.ActivityLogService)
	committeeClassificationController := controllers.NewCommitteeClassificationController(committeeClassificationSvc)

	internalCommitteeRepo := repositories.NewInternalCommitteeRepository(db)
	internalCommitteeWorkflowReviewRepo := repositories.NewInternalCommitteeWorkflowReviewRepository(db)
	internalCommitteeWorkflowRepo := repositories.NewInternalCommitteeWorkflowRepository(db)

	// Create workflow service first (without committee service dependency initially)
	internalCommitteeWorkflowSvc := services.NewInternalCommitteeWorkflowService(
		internalCommitteeWorkflowRepo,
		sharedModule.ActivityLogService,
		settingsModule.WorkflowSvc,
		masterModule.WorkflowStatusSvc,
		userSvc,
		nil, // Will be set later via SetInternalCommitteeService
		masterModule.WorkflowContextSvc,
	)

	// Create workflow review service (without committee service dependency initially)
	internalCommitteeWorkflowReviewSvc := services.NewInternalCommitteeWorkflowReviewService(
		internalCommitteeWorkflowReviewRepo,
		sharedModule.ActivityLogService,
		nil, // Will be set later
		internalCommitteeWorkflowSvc,
		masterModule.WorkflowStatusSvc,
		internalCommitteeStatusSvc,
		userSvc,
		commonModule.NotificationService,
	)

	// Create committee service with all dependencies including workflow services
	internalCommitteeSvc := services.NewInternalCommitteeService(
		internalCommitteeRepo,
		sharedModule.ActivityLogService,
		userSvc,
		settingsModule.OrganizationUnitSvc,
		committeeTypeSvc,
		committeeTermSvc,
		committeeClassificationSvc,
		internalCommitteeStatusSvc,
		commonModule.NotificationService,
		organizationUnitsSvc,
		internalCommitteeWorkflowReviewSvc,
		internalCommitteeWorkflowSvc,
		masterModule.WorkflowStatusSvc,
	)
	internalCommitteeController := controllers.NewInternalCommitteeController(internalCommitteeSvc)

	// Set the circular references
	internalCommitteeWorkflowSvc.SetInternalCommitteeService(internalCommitteeSvc)
	internalCommitteeWorkflowReviewSvc.SetInternalCommitteeService(internalCommitteeSvc)

	internalCommitteeWorkflowController := controllers.NewInternalCommitteeWorkflowController(internalCommitteeWorkflowSvc)
	internalCommitteeWorkflowReviewController := controllers.NewInternalCommitteeWorkflowReviewController(internalCommitteeWorkflowReviewSvc)

	internalCommitteeReviewCommentRepo := repositories.NewInternalCommitteeReviewCommentRepository(db)
	internalCommitteeReviewCommentSvc := services.NewInternalCommitteeReviewCommentService(internalCommitteeReviewCommentRepo, userSvc, sharedModule.ActivityLogService, internalCommitteeSvc, masterModule.CommentStatusSvc)
	internalCommitteeReviewCommentController := controllers.NewInternalCommitteeReviewCommentController(internalCommitteeReviewCommentSvc)

	controlStatusRepo := repositories.NewControlStatusRepository(db)
	controlStatusSvc := services.NewControlStatusService(controlStatusRepo)
	controlStatusController := controllers.NewControlStatusController(controlStatusSvc)

	return &OrganizationModule{
		DesignationController:                   designationController,
		UserStatusController:                    userStatusController,
		UserController:                          userController,
		UserSettingController:                   userSettingController,
		StrategicProfileController:              strategicProfileController,
		StrategicObjectiveController:            strategicObjectiveController,
		ControlCategoryController:               controlCategoryController,
		ControlTypeController:                   controlTypeController,
		ControlImplementationStatusController:   controlImplementationStatusController,
		ControlEfficiencyStatusController:       controlEfficiencyStatusController,
		ControlReviewStatusController:           controlReviewStatusController,
		ControlController:                       controlController,
		ControlTestController:                   controlTestController,
		ProcessGroupController:                  processGroupController,
		ControlActionPlanController:             controlActionPlanController,
		ControlActionPlanStatusUpdateController: controlActionPlanStatusUpdateController,
		ControlActionPlanCommentController:      controlActionPlanCommentController,

		ProcessController:                         processController,
		UserTaskController:                        userTaskController,
		StrategicProfileValueController:           strategicProfileValueController,
		ProcessAssociatedItemController:           processAssociatedItemController,
		ControlAssociatedItemController:           controlAssociatedItemController,
		ControlReviewController:                   controlReviewController,
		UserProfileController:                     userProfileController,
		ProcessDependencyTypeController:           processDependencyTypeController,
		ProcessDependencySubTypeController:        processDependencySubTypeController,
		GovernanceReportController:                governanceReportController,
		ProcessDependencyController:               processDependencyController,
		ProcessCriticalityLevelController:         processCriticalityLevelController,
		ProcessActivityController:                 processActivityController,
		CompetencyGroupController:                 competencyGroupController,
		CompetencyController:                      competencyController,
		CompetencyLevelController:                 competencyLevelController,
		JobDescriptionController:                  jobDescriptionController,
		JobDescriptionCompetencyController:        jobDescriptionCompetencyController,
		UserDocumentTypeController:                userDocumentTypeController,
		UserDocumentSubTypeController:             userDocumentSubTypeController,
		UserDocumentController:                    userDocumentController,
		UserJobDescriptionController:              userJobDescriptionController,
		UserDelegationHistoryController:           userDelegationHistoryController,
		OrganizationUnitsController:               organizationUnitsController,
		OrganizationUnitBusinessContextController: organizationUnitBusinessContextController,
		OrganizationUnitAttributeController:       organizationUnitAttributeController,
		OrganizationCartController:                organizationChartController,
		MyDashboardController:                     myDashboardController,
		CommitteeTermController:                   committeeTermController,
		CommitteeTypeController:                   committeeTypeController,
		CommitteeClassificationController:         committeeClassificationController,
		InternalCommitteeStatusController:         internalCommitteeStatusController,
		InternalCommitteeController:               internalCommitteeController,
		InternalCommitteeWorkflowController:       internalCommitteeWorkflowController,
		InternalCommitteeWorkflowReviewController: internalCommitteeWorkflowReviewController,
		InternalCommitteeReviewCommentController:  internalCommitteeReviewCommentController,
		ControlEffectivenessRatingController:      controlEffectivenessRatingController,
		ControlStatusController:                   controlStatusController,

		UserSvc:                        userSvc,
		UserSettingSvc:                 userSettingSvc,
		ControlSvc:                     controlSvc,
		ProcessSvc:                     processSvc,
		StrategicObjectiveSvc:          strategicObjectiveSvc,
		ControlTestSvc:                 controlTestSvc,
		ControlImplementationStatusSvc: controlImplementationStatusSvc,
		ControlEfficiencyStatusSvc:     controlEfficiencyStatusSvc,
		ProcessAssociatedItemSvc:       processAssociatedItemSvc,
		ProcessDependencySvc:           processDependencySvc,
		JobDescriptionSvc:              jobDescriptionSvc,
		GovernanceReportSvc:            governanceReportSvc,
		ControlStatusSvc:               controlStatusSvc,
	}
}
