package routes

import (
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/organization/controllers"
	organizationRepositories "jethur/grc/organization/repositories"
	organizationServices "jethur/grc/organization/services"
	settingRepositories "jethur/grc/settings/repositories"
	settingServices "jethur/grc/settings/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	organizationRoutes := router.Group("/")

	organizationRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	reviewFrequiencyRepo := masterRepositories.NewReviewFrequencyRepository(db)
	reviewFrequiencySvc := masterServices.NewReviewFrequencyService(reviewFrequiencyRepo)

	lifecycleStatusRepo := masterRepositories.NewLifecycleStatusRepository(db)
	lifecycleStatusSvc := masterServices.NewLifecycleStatusService(lifecycleStatusRepo)

	roleRepo := settingRepositories.NewRoleRepository(db)
	roleSvc := settingServices.NewRoleService(roleRepo, activityLogSvc)

	designationRepo := organizationRepositories.NewDesignationRepository(db)
	designationSvc := organizationServices.NewDesignationService(designationRepo, activityLogSvc)
	designationController := controllers.NewDesignationController(designationSvc)

	userStatusRepo := organizationRepositories.NewUserStatusRepository(db)
	userStatusSvc := organizationServices.NewUserStatusService(userStatusRepo)
	userStatusController := controllers.NewUserStatusController(userStatusSvc)

	userRepo := organizationRepositories.NewUserRepository(db)
	userSvc := organizationServices.NewUserService(userRepo, designationSvc, userStatusSvc, activityLogSvc, roleSvc)
	userController := controllers.NewUserController(userSvc)

	strategicProfileRepo := organizationRepositories.NewStrategicProfileRepository(db)
	strategicProfileSvc := organizationServices.NewStrategicProfileService(strategicProfileRepo)
	strategicProfileController := controllers.NewStrategicProfileController(strategicProfileSvc)

	strategicObjectiveRepo := organizationRepositories.NewStrategicObjectiveRepository(db)
	strategicObjectiveSvc := organizationServices.NewStrategicObjectiveService(strategicObjectiveRepo, activityLogSvc)
	strategicObjectiveController := controllers.NewStrategicObjectiveController(strategicObjectiveSvc)

	userSettingRepo := organizationRepositories.NewUserSettingRepository(db)
	userSettingSvc := organizationServices.NewUserSettingService(userSettingRepo, activityLogSvc)
	userSettingController := controllers.NewUserSettingController(userSettingSvc)

	controlCategoryRepo := organizationRepositories.NewControlCategoryRepository(db)
	controlCategorySvc := organizationServices.NewControlCategoryService(controlCategoryRepo, activityLogSvc)
	controlCategoryController := controllers.NewControlCategoryController(controlCategorySvc)

	controlTypeRepo := organizationRepositories.NewControlTypeRepository(db)
	controlTypeSvc := organizationServices.NewControlTypeService(controlTypeRepo)
	controlTypeController := controllers.NewControlTypeController(controlTypeSvc)

	controlImplementationStatusRepo := organizationRepositories.NewControlImplementationStatusRepository(db)
	controlImplementationStatusSvc := organizationServices.NewControlImplementationStatusService(controlImplementationStatusRepo)
	controlImplementationStatusController := controllers.NewControlImplementationStatusController(controlImplementationStatusSvc)

	controlEfficiencyStatusRepo := organizationRepositories.NewControlEfficiencyStatusRepository(db)
	controlEfficiencyStatusSvc := organizationServices.NewControlEfficiencyStatusService(controlEfficiencyStatusRepo)
	controlEfficiencyStatusController := controllers.NewControlEfficiencyStatusController(controlEfficiencyStatusSvc)

	controlRepo := organizationRepositories.NewControlRepository(db)
	controlSvc := organizationServices.NewControlService(
		controlRepo,
		activityLogSvc,
		controlTypeSvc,
		controlCategorySvc,
		controlImplementationStatusSvc,
		controlEfficiencyStatusSvc,
		userSvc,
		reviewFrequiencySvc,
		lifecycleStatusSvc)
	controlController := controllers.NewControlController(controlSvc)

	controlTestRepo := organizationRepositories.NewControlTestRepository(db)
	controlTestSvc := organizationServices.NewControlTestService(
		controlTestRepo,
		activityLogSvc,
		controlSvc,
		controlImplementationStatusSvc,
		controlEfficiencyStatusSvc)
	controlTestController := controllers.NewControlTestController(controlTestSvc)

	//Designation
	organizationRoutes.GET("/designations", designationController.GetAllDesignations)
	organizationRoutes.GET("/designations/:id", designationController.GetDesignationDetails)
	organizationRoutes.POST("/designations", designationController.CreateDesignation)
	organizationRoutes.PUT("/designations/:id", designationController.UpdateDesignation)
	organizationRoutes.POST("/designations/delete", designationController.DeleteDesignations)
	organizationRoutes.PUT("/designations/activate", designationController.ActivateDesignations)
	organizationRoutes.PUT("/designations/deactivate", designationController.DeactivateDesignations)

	//User Status
	organizationRoutes.GET("/user-statuses", userStatusController.GetAllUserStatuses)
	organizationRoutes.GET("/user-statuses/:id", userStatusController.GetUserStatusDetails)

	//User
	organizationRoutes.GET("/users", userController.GetAllUsers)
	organizationRoutes.POST("/users", userController.CreateUser)
	organizationRoutes.PUT("/users/:userId", userController.UpdateUser)
	organizationRoutes.POST("/users/delete", userController.DeleteUsers)
	organizationRoutes.POST("/users/:userId/reset-password", userController.NewPassword)

	//Strategic Profile
	organizationRoutes.GET("/strategic-profiles", strategicProfileController.GetStrategicProfile)
	organizationRoutes.PUT("/strategic-profiles", strategicProfileController.UpdateStrategicProfile)
	organizationRoutes.PUT("/strategic-profiles/vision", strategicProfileController.UpdateVision)
	organizationRoutes.PUT("/strategic-profiles/mission", strategicProfileController.UpdateMission)
	organizationRoutes.PUT("/strategic-profiles/values", strategicProfileController.UpdateValues)

	//Strategic Objective
	organizationRoutes.GET("/strategic-objectives", strategicObjectiveController.GetAllStrategicObjectives)
	organizationRoutes.GET("/strategic-objectives/:id", strategicObjectiveController.GetStrategicObjectiveDetails)
	organizationRoutes.POST("/strategic-objectives", strategicObjectiveController.CreateStrategicObjective)
	organizationRoutes.PUT("/strategic-objectives/:id", strategicObjectiveController.UpdateStrategicObjective)
	organizationRoutes.POST("/strategic-objectives/delete", strategicObjectiveController.DeleteStrategicObjectives)
	organizationRoutes.PUT("/strategic-objectives/activate", strategicObjectiveController.ActivateStrategicObjectives)
	organizationRoutes.PUT("/strategic-objectives/deactivate", strategicObjectiveController.DeactivateStrategicObjectives)

	//User Settings
	organizationRoutes.GET("/users/:userId/settings/preferences", userController.GetSetting)
	organizationRoutes.PUT("/users/:userId/settings/preferences", userController.UpdateSetting)

	organizationRoutes.PUT("/users/:userId/settings/security/enable-mfa", userController.EnableMFA)
	organizationRoutes.PUT("/users/:userId/settings/security/disable-mfa", userController.DisableMFA)
	//@review
	//you can change this notification url once individual notification enable/disable url ready if required.
	organizationRoutes.PUT("/users/:userId/settings/enable-notification", userController.EnableNotification)
	organizationRoutes.PUT("/users/:userId/settings/disable-notification", userController.DisableNotification)

	organizationRoutes.GET("/users/:userId/settings/notifications", userSettingController.GetNotificationSettings)
	organizationRoutes.POST("/users/:userId/settings/notifications/:notificationMasterId/enable", userSettingController.EnableNotification)
	organizationRoutes.POST("/users/:userId/settings/notifications/:notificationMasterId/disable", userSettingController.DisableNotification)

	//Control Category
	organizationRoutes.GET("/control-categories", controlCategoryController.GetAllControlCategories)
	organizationRoutes.GET("/control-categories/:id", controlCategoryController.GetControlCategoryDetails)
	organizationRoutes.POST("/control-categories", controlCategoryController.CreateControlCategory)
	organizationRoutes.PUT("/control-categories/:id", controlCategoryController.UpdateControlCategory)
	organizationRoutes.POST("/control-categories/delete", controlCategoryController.DeleteControlCategory)
	organizationRoutes.PUT("/control-categories/activate", controlCategoryController.ActivateControlCategories)
	organizationRoutes.PUT("/control-categories/deactivate", controlCategoryController.DeactivateControlCategories)

	//Control Type
	organizationRoutes.GET("/control-types", controlTypeController.GetAllControlTypes)

	//Control Implementation Status
	organizationRoutes.GET("/control-implementation-statuses", controlImplementationStatusController.GetAllControlImplementationStatuses)
	organizationRoutes.GET("/control-implementation-statuses/:id", controlImplementationStatusController.GetControlImplementationStatusDetails)

	//Control Efficiency Status
	organizationRoutes.GET("/control-efficiency-statuses", controlEfficiencyStatusController.GetAllControlEfficiencyStatuses)
	organizationRoutes.GET("/control-efficiency-statuses/:id", controlEfficiencyStatusController.GetControlEfficiencyStatusDetails)

	//Control
	organizationRoutes.GET("/controls", controlController.GetAllControls)
	organizationRoutes.GET("/controls/:controlId", controlController.GetControlDetails)
	organizationRoutes.POST("/controls", controlController.CreateControl)
	organizationRoutes.PUT("/controls/:controlId", controlController.UpdateControl)
	organizationRoutes.POST("/controls/delete", controlController.DeleteControls)
	organizationRoutes.PUT("/controls/archive", controlController.ArchiveControls)
	organizationRoutes.PUT("/controls/restore", controlController.RestoreControls)

	//control Test
	organizationRoutes.GET("/controls/:controlId/test", controlTestController.GetAllControlTests)
	organizationRoutes.GET("/controls/:controlId/test/:controlTestId", controlTestController.GetControlTestDetails)
	organizationRoutes.POST("/controls/:controlId/test", controlTestController.CreateControlTest)
	organizationRoutes.PUT("/controls/:controlId/test/:controlTestId", controlTestController.UpdateControlTest)
	organizationRoutes.POST("/controls/:controlId/test/delete", controlTestController.DeleteControlTest)

}
