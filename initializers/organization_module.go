package initializers

import (
	"jethur/grc/organization/controllers"
	"jethur/grc/organization/repositories"
	"jethur/grc/organization/services"

	"gorm.io/gorm"
)

type OrganizationModule struct {
	DesignationController                 *controllers.DesignationController
	UserStatusController                  *controllers.UserStatusController
	UserController                        *controllers.UserController
	UserSettingController                 *controllers.UserSettingController
	StrategicProfileController            *controllers.StrategicProfileController
	StrategicObjectiveController          *controllers.StrategicOBjectiveController
	ControlCategoryController             *controllers.ControlCategoryController
	ControlTypeController                 *controllers.ControlTypeController
	ControlImplementationStatusController *controllers.ControlImplementationStatusController
	ControlEfficiencyStatusController     *controllers.ControlEfficiencyStatusController
	ControlController                     *controllers.ControlController
	ControlTestController                 *controllers.ControlTestController
	ProcessGroupController                *controllers.ProcessGroupController

	UserSvc        services.UserServiceInterface
	UserSettingSvc services.UserSettingServiceInterface
}

func NewOrganizationModule(
	db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule) *OrganizationModule {

	designationRepo := repositories.NewDesignationRepository(db)
	designationSvc := services.NewDesignationService(designationRepo, sharedModule.ActivityLogService)
	designationController := controllers.NewDesignationController(designationSvc)

	userStatusRepo := repositories.NewUserStatusRepository(db)
	userStatusSvc := services.NewUserStatusService(userStatusRepo)
	userStatusController := controllers.NewUserStatusController(userStatusSvc)

	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo, designationSvc, userStatusSvc, sharedModule.ActivityLogService, settingsModule.RoleSvc, settingsModule.SystemPreferenceSvc)
	userSettingRepo := repositories.NewUserSettingRepository(db)
	userSettingSvc := services.NewUserSettingService(userSettingRepo, sharedModule.ActivityLogService)
	userSvc.SetUserSettingService(userSettingSvc)
	userSettingSvc.SetUserService(userSvc)
	userController := controllers.NewUserController(userSvc)
	userSettingController := controllers.NewUserSettingController(userSettingSvc)

	strategicProfileRepo := repositories.NewStrategicProfileRepository(db)
	strategicProfileSvc := services.NewStrategicProfileService(strategicProfileRepo)
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
		masterModule.LifecycleStatusSvc)
	controlController := controllers.NewControlController(controlSvc)

	controlTestRepo := repositories.NewControlTestRepository(db)
	controlTestSvc := services.NewControlTestService(
		controlTestRepo,
		sharedModule.ActivityLogService,
		controlSvc,
		controlImplementationStatusSvc,
		controlEfficiencyStatusSvc)
	controlTestController := controllers.NewControlTestController(controlTestSvc)

	processGroupRepo := repositories.NewProcessGroupRepository(db)
	processGroupSvc := services.NewProcessGroupService(processGroupRepo, sharedModule.ActivityLogService)
	processGroupController := controllers.NewProcessGroupController(processGroupSvc)

	return &OrganizationModule{
		DesignationController:                 designationController,
		UserStatusController:                  userStatusController,
		UserController:                        userController,
		UserSettingController:                 userSettingController,
		StrategicProfileController:            strategicProfileController,
		StrategicObjectiveController:          strategicObjectiveController,
		ControlCategoryController:             controlCategoryController,
		ControlTypeController:                 controlTypeController,
		ControlImplementationStatusController: controlImplementationStatusController,
		ControlEfficiencyStatusController:     controlEfficiencyStatusController,
		ControlController:                     controlController,
		ControlTestController:                 controlTestController,
		ProcessGroupController:                processGroupController,

		UserSvc:        userSvc,
		UserSettingSvc: userSettingSvc,
	}
}
