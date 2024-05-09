package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	organizationRoutes := router.Group("/")

	organizationRoutes.Use(middlewares.AuthMiddleware())

	// activityLogRepo := logRepositories.NewActivityLogRepository(db)
	// activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	// currencyRepo := settingRepositories.NewCurrencyRepository(db)
	// currencySvc := settingServices.NewCurrencyService(currencyRepo)

	// systemPreferenceRepo := settingRepositories.NewSystemPreferenceRepository(db)
	// systemPreferenceSvc := settingServices.NewSystemPreferenceService(systemPreferenceRepo, activityLogSvc, currencySvc)

	// reviewFrequiencyRepo := masterRepositories.NewReviewFrequencyRepository(db)
	// reviewFrequiencySvc := masterServices.NewReviewFrequencyService(reviewFrequiencyRepo)

	// lifecycleStatusRepo := masterRepositories.NewLifecycleStatusRepository(db)
	// lifecycleStatusSvc := masterServices.NewLifecycleStatusService(lifecycleStatusRepo)

	// roleRepo := settingRepositories.NewRoleRepository(db)
	// roleSvc := settingServices.NewRoleService(roleRepo, activityLogSvc)

	// designationRepo := organizationRepositories.NewDesignationRepository(db)
	// designationSvc := organizationServices.NewDesignationService(designationRepo, activityLogSvc)
	// designationController := controllers.NewDesignationController(designationSvc)

	// userStatusRepo := organizationRepositories.NewUserStatusRepository(db)
	// userStatusSvc := organizationServices.NewUserStatusService(userStatusRepo)
	// userStatusController := controllers.NewUserStatusController(userStatusSvc)

	// userRepo := organizationRepositories.NewUserRepository(db)
	// userSvc := organizationServices.NewUserService(userRepo, designationSvc, userStatusSvc, activityLogSvc, roleSvc, systemPreferenceSvc)
	// userSettingRepo := organizationRepositories.NewUserSettingRepository(db)
	// userSettingSvc := organizationServices.NewUserSettingService(userSettingRepo, activityLogSvc)
	// userSvc.SetUserSettingService(userSettingSvc)
	// userSettingSvc.SetUserService(userSvc)
	// userController := controllers.NewUserController(userSvc)
	// userSettingController := controllers.NewUserSettingController(userSettingSvc)

	// strategicProfileRepo := organizationRepositories.NewStrategicProfileRepository(db)
	// strategicProfileSvc := organizationServices.NewStrategicProfileService(strategicProfileRepo)
	// strategicProfileController := controllers.NewStrategicProfileController(strategicProfileSvc)

	// strategicObjectiveRepo := organizationRepositories.NewStrategicObjectiveRepository(db)
	// strategicObjectiveSvc := organizationServices.NewStrategicObjectiveService(strategicObjectiveRepo, activityLogSvc)
	// strategicObjectiveController := controllers.NewStrategicObjectiveController(strategicObjectiveSvc)

	// controlCategoryRepo := organizationRepositories.NewControlCategoryRepository(db)
	// controlCategorySvc := organizationServices.NewControlCategoryService(controlCategoryRepo, activityLogSvc)
	// controlCategoryController := controllers.NewControlCategoryController(controlCategorySvc)

	// controlTypeRepo := organizationRepositories.NewControlTypeRepository(db)
	// controlTypeSvc := organizationServices.NewControlTypeService(controlTypeRepo)
	// controlTypeController := controllers.NewControlTypeController(controlTypeSvc)

	// controlImplementationStatusRepo := organizationRepositories.NewControlImplementationStatusRepository(db)
	// controlImplementationStatusSvc := organizationServices.NewControlImplementationStatusService(controlImplementationStatusRepo)
	// controlImplementationStatusController := controllers.NewControlImplementationStatusController(controlImplementationStatusSvc)

	// controlEfficiencyStatusRepo := organizationRepositories.NewControlEfficiencyStatusRepository(db)
	// controlEfficiencyStatusSvc := organizationServices.NewControlEfficiencyStatusService(controlEfficiencyStatusRepo)
	// controlEfficiencyStatusController := controllers.NewControlEfficiencyStatusController(controlEfficiencyStatusSvc)

	// controlRepo := organizationRepositories.NewControlRepository(db)
	// controlSvc := organizationServices.NewControlService(
	// 	controlRepo,
	// 	activityLogSvc,
	// 	controlTypeSvc,
	// 	controlCategorySvc,
	// 	controlImplementationStatusSvc,
	// 	controlEfficiencyStatusSvc,
	// 	userSvc,
	// 	reviewFrequiencySvc,
	// 	lifecycleStatusSvc)
	// controlController := controllers.NewControlController(controlSvc)

	// controlTestRepo := organizationRepositories.NewControlTestRepository(db)
	// controlTestSvc := organizationServices.NewControlTestService(
	// 	controlTestRepo,
	// 	activityLogSvc,
	// 	controlSvc,
	// 	controlImplementationStatusSvc,
	// 	controlEfficiencyStatusSvc)
	// controlTestController := controllers.NewControlTestController(controlTestSvc)

	// processGroupRepo := organizationRepositories.NewProcessGroupRepository(db)
	// processGroupSvc := organizationServices.NewProcessGroupService(processGroupRepo, activityLogSvc)
	// processGroupController := controllers.NewProcessGroupController(processGroupSvc)
	userController := app.UserController
	userSettingController := app.UserSettingController
	//Designation
	// organizationRoutes.GET("/designations", designationController.GetAllDesignations)
	// organizationRoutes.GET("/designations/:id", designationController.GetDesignationDetails)
	// organizationRoutes.POST("/designations", designationController.CreateDesignation)
	// organizationRoutes.PUT("/designations/:id", designationController.UpdateDesignation)
	// organizationRoutes.POST("/designations/delete", designationController.DeleteDesignations)
	// organizationRoutes.PUT("/designations/activate", designationController.ActivateDesignations)
	// organizationRoutes.PUT("/designations/deactivate", designationController.DeactivateDesignations)

	//User Status
	// organizationRoutes.GET("/user-statuses", userStatusController.GetAllUserStatuses)
	// organizationRoutes.GET("/user-statuses/:id", userStatusController.GetUserStatusDetails)

	//User
	organizationRoutes.GET("/users", userController.GetAllUsers)
	organizationRoutes.GET("/users/:userId", userController.GetUserDetails)
	organizationRoutes.POST("/users", userController.CreateUser)
	organizationRoutes.PUT("/users/:userId", userController.UpdateUser)
	organizationRoutes.POST("/users/delete", userController.DeleteUsers)

	//Strategic Profile
	// organizationRoutes.GET("/strategic-profiles", strategicProfileController.GetStrategicProfile)
	// organizationRoutes.PUT("/strategic-profiles", strategicProfileController.UpdateStrategicProfile)
	// organizationRoutes.PUT("/strategic-profiles/vision", strategicProfileController.UpdateVision)
	// organizationRoutes.PUT("/strategic-profiles/mission", strategicProfileController.UpdateMission)
	// organizationRoutes.PUT("/strategic-profiles/values", strategicProfileController.UpdateValues)

	//Strategic Objective
	// organizationRoutes.GET("/strategic-objectives", strategicObjectiveController.GetAllStrategicObjectives)
	// organizationRoutes.GET("/strategic-objectives/:id", strategicObjectiveController.GetStrategicObjectiveDetails)
	// organizationRoutes.POST("/strategic-objectives", strategicObjectiveController.CreateStrategicObjective)
	// organizationRoutes.PUT("/strategic-objectives/:id", strategicObjectiveController.UpdateStrategicObjective)
	// organizationRoutes.POST("/strategic-objectives/delete", strategicObjectiveController.DeleteStrategicObjectives)
	// organizationRoutes.PUT("/strategic-objectives/activate", strategicObjectiveController.ActivateStrategicObjectives)
	// organizationRoutes.PUT("/strategic-objectives/deactivate", strategicObjectiveController.DeactivateStrategicObjectives)

	//User Settings
	organizationRoutes.GET("/users/:userId/settings/preferences", userSettingController.GetSetting)
	organizationRoutes.PUT("/users/:userId/settings/preferences", userSettingController.UpdateSetting)

	organizationRoutes.PUT("/users/:userId/settings/security/enable-mfa", userSettingController.EnableMFA)
	organizationRoutes.PUT("/users/:userId/settings/security/disable-mfa", userSettingController.DisableMFA)
	organizationRoutes.POST("/users/:userId/settings/security/reset-password", userSettingController.NewPassword)

	organizationRoutes.PUT("/users/:userId/settings/enable-notification", userSettingController.EnableMainNotification)
	organizationRoutes.PUT("/users/:userId/settings/disable-notification", userSettingController.DisableMainNotification)

	organizationRoutes.GET("/users/:userId/settings/notifications", userSettingController.GetNotificationSettings)
	organizationRoutes.POST("/users/:userId/settings/notifications/:notificationMasterId/enable", userSettingController.EnableNotification)
	organizationRoutes.POST("/users/:userId/settings/notifications/:notificationMasterId/disable", userSettingController.DisableNotification)

	//Control Category
	// organizationRoutes.GET("/control-categories", controlCategoryController.GetAllControlCategories)
	// organizationRoutes.GET("/control-categories/:id", controlCategoryController.GetControlCategoryDetails)
	// organizationRoutes.POST("/control-categories", controlCategoryController.CreateControlCategory)
	// organizationRoutes.PUT("/control-categories/:id", controlCategoryController.UpdateControlCategory)
	// organizationRoutes.POST("/control-categories/delete", controlCategoryController.DeleteControlCategory)
	// organizationRoutes.PUT("/control-categories/activate", controlCategoryController.ActivateControlCategories)
	// organizationRoutes.PUT("/control-categories/deactivate", controlCategoryController.DeactivateControlCategories)

	//Control Type
	//organizationRoutes.GET("/control-types", controlTypeController.GetAllControlTypes)

	//Control Implementation Status
	// organizationRoutes.GET("/control-implementation-statuses", controlImplementationStatusController.GetAllControlImplementationStatuses)
	// organizationRoutes.GET("/control-implementation-statuses/:id", controlImplementationStatusController.GetControlImplementationStatusDetails)

	//Control Efficiency Status
	// organizationRoutes.GET("/control-efficiency-statuses", controlEfficiencyStatusController.GetAllControlEfficiencyStatuses)
	// organizationRoutes.GET("/control-efficiency-statuses/:id", controlEfficiencyStatusController.GetControlEfficiencyStatusDetails)

	//Control
	// organizationRoutes.GET("/controls", controlController.GetAllControls)
	// organizationRoutes.GET("/controls/:controlId", controlController.GetControlDetails)
	// organizationRoutes.POST("/controls", controlController.CreateControl)
	// organizationRoutes.PUT("/controls/:controlId", controlController.UpdateControl)
	// organizationRoutes.POST("/controls/delete", controlController.DeleteControls)
	// organizationRoutes.PUT("/controls/archive", controlController.ArchiveControls)
	// organizationRoutes.PUT("/controls/restore", controlController.RestoreControls)

	//control Test
	// organizationRoutes.GET("/controls/:controlId/test", controlTestController.GetAllControlTests)
	// organizationRoutes.GET("/controls/:controlId/test/:controlTestId", controlTestController.GetControlTestDetails)
	// organizationRoutes.POST("/controls/:controlId/test", controlTestController.CreateControlTest)
	// organizationRoutes.PUT("/controls/:controlId/test/:controlTestId", controlTestController.UpdateControlTest)
	// organizationRoutes.POST("/controls/:controlId/test/delete", controlTestController.DeleteControlTest)

	//Process Group
	// organizationRoutes.GET("/process-groups", processGroupController.GetAllProcessGroups)
	// organizationRoutes.GET("/process-groups/:id", processGroupController.GetProcessGroupDetails)
	// organizationRoutes.POST("/process-groups", processGroupController.CreateProcessGroup)
	// organizationRoutes.PUT("/process-groups/:id", processGroupController.UpdateProcessGroup)
	// organizationRoutes.POST("/process-groups/delete", processGroupController.DeleteProcessGroups)
	// organizationRoutes.PUT("/process-groups/activate", processGroupController.ActivateProcessGroups)
	// organizationRoutes.PUT("/process-groups/deactivate", processGroupController.DeactivateProcessGroups)

}
