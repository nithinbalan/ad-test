package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	organizationRoutes := router.Group("/")

	organizationRoutes.Use(middlewares.AuthMiddleware())

	designationController := app.OrganizationModule.DesignationController
	userStatusController := app.OrganizationModule.UserStatusController
	userController := app.OrganizationModule.UserController
	userSettingController := app.OrganizationModule.UserSettingController
	strategicProfileController := app.OrganizationModule.StrategicProfileController
	strategicProfileValueController := app.OrganizationModule.StrategicProfileValueController
	strategicObjectiveController := app.OrganizationModule.StrategicObjectiveController
	controlCategoryController := app.OrganizationModule.ControlCategoryController
	controlTypeController := app.OrganizationModule.ControlTypeController
	controlImplementationStatusController := app.OrganizationModule.ControlImplementationStatusController
	controlEfficiencyStatusController := app.OrganizationModule.ControlEfficiencyStatusController
	controlController := app.OrganizationModule.ControlController
	controlTestController := app.OrganizationModule.ControlTestController
	processGroupController := app.OrganizationModule.ProcessGroupController
	controlActionPlanController := app.OrganizationModule.ControlActionPlanController
	controlActionPlanStatusUpdateController := app.OrganizationModule.ControlActionPlanStatusUpdateController
	processController := app.OrganizationModule.ProcessController
	userTaskController := app.OrganizationModule.UserTaskController
	processAssociatedItemController := app.OrganizationModule.ProcessAssociatedItemController
	controlAssociatedItemController := app.OrganizationModule.ControlAssociatedItemController
	controlReviewStatusController := app.OrganizationModule.ControlReviewStatusController
	controlReviewController := app.OrganizationModule.ControlReviewController

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
	organizationRoutes.GET("/users/:userId", userController.GetUserDetails)
	organizationRoutes.POST("/users", userController.CreateUser)
	organizationRoutes.PUT("/users/:userId", userController.UpdateUser)
	organizationRoutes.POST("/users/delete", userController.DeleteUsers)
	organizationRoutes.GET("/users/:userId/login-history", userController.GetUserLoginHistory)

	organizationRoutes.GET("/users/me", userController.GetCurrentUserDetails)

	//Strategic Profile
	organizationRoutes.GET("/strategic-profiles", strategicProfileController.GetStrategicProfile)
	organizationRoutes.PUT("/strategic-profiles", strategicProfileController.UpdateStrategicProfile)
	organizationRoutes.PUT("/strategic-profiles/vision", strategicProfileController.UpdateVision)
	organizationRoutes.PUT("/strategic-profiles/mission", strategicProfileController.UpdateMission)
	organizationRoutes.PUT("/strategic-profiles/values", strategicProfileController.UpdateValues)

	organizationRoutes.POST("/strategic-profiles/upload-image", strategicProfileController.UploadProfileImage)
	organizationRoutes.GET("/strategic-profiles/preview-image", strategicProfileController.PreviewProfileImage)

	//Strategic Profile Values
	organizationRoutes.GET("/strategic-profiles/values", strategicProfileValueController.GetAllStrategicProfileValues)
	organizationRoutes.GET("/strategic-profiles/values/:id", strategicProfileValueController.GetStrategicProfileValueDetails)
	organizationRoutes.POST("/strategic-profiles/values", strategicProfileValueController.CreateStrategicProfileValue)
	organizationRoutes.PUT("/strategic-profiles/values/:id", strategicProfileValueController.UpdateStrategicProfileValue)
	organizationRoutes.POST("/strategic-profiles/values/delete", strategicProfileValueController.DeleteStrategicProfileValues)

	//Strategic Objective
	organizationRoutes.GET("/strategic-objectives", strategicObjectiveController.GetAllStrategicObjectives)
	organizationRoutes.GET("/strategic-objectives/:id", strategicObjectiveController.GetStrategicObjectiveDetails)
	organizationRoutes.POST("/strategic-objectives", strategicObjectiveController.CreateStrategicObjective)
	organizationRoutes.PUT("/strategic-objectives/:id", strategicObjectiveController.UpdateStrategicObjective)
	organizationRoutes.POST("/strategic-objectives/delete", strategicObjectiveController.DeleteStrategicObjectives)
	organizationRoutes.PUT("/strategic-objectives/activate", strategicObjectiveController.ActivateStrategicObjectives)
	organizationRoutes.PUT("/strategic-objectives/deactivate", strategicObjectiveController.DeactivateStrategicObjectives)

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

	//User Tasks
	organizationRoutes.GET("/users/:userId/tasks", userTaskController.GetAllUserTasks)
	organizationRoutes.GET("/users/:userId/task-summary", userTaskController.GetUserTaskSummary)

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

	//Control Review Statuses
	organizationRoutes.GET("/control-review-statuses", controlReviewStatusController.GetAllControlReviewStatuses)
	organizationRoutes.GET("/control-review-statuses/:id", controlReviewStatusController.GetControlReviewStatusDetails)

	//Control Review
	organizationRoutes.GET("/controls/:controlId/reviews", controlReviewController.GetAllControlReviews)
	organizationRoutes.POST("/controls/:controlId/reviews", controlReviewController.CreateControlReview)
	organizationRoutes.PUT("/controls/:controlId/reviews/:id", controlReviewController.UpdateControlReview)

	//Control
	organizationRoutes.GET("/controls", controlController.GetAllControls)
	organizationRoutes.GET("/controls/:controlId", controlController.GetControlDetails)
	organizationRoutes.POST("/controls", controlController.CreateControl)
	organizationRoutes.PUT("/controls/:controlId", controlController.UpdateControl)
	organizationRoutes.POST("/controls/delete", controlController.DeleteControls)
	organizationRoutes.PUT("/controls/archive", controlController.ArchiveControls)
	organizationRoutes.PUT("/controls/restore", controlController.RestoreControls)

	//control Test
	organizationRoutes.GET("/controls/:controlId/tests", controlTestController.GetAllControlTests)
	organizationRoutes.GET("/controls/:controlId/tests/:controlTestId", controlTestController.GetControlTestDetails)
	organizationRoutes.POST("/controls/:controlId/tests", controlTestController.CreateControlTest)
	organizationRoutes.PUT("/controls/:controlId/tests/:controlTestId", controlTestController.UpdateControlTest)
	organizationRoutes.POST("/controls/:controlId/tests/delete", controlTestController.DeleteControlTest)
	organizationRoutes.GET("/controls/:controlId/trends", controlTestController.GetControlTestTrendAnalysis)

	//Control Associated Item
	organizationRoutes.GET("/controls/:controlId/processes", controlAssociatedItemController.GetAllControlProcesses)
	organizationRoutes.GET("/controls/:controlId/compliance-requirements", controlAssociatedItemController.GetAllControlComplianceRequirements)
	organizationRoutes.GET("/controls/:controlId/risks", controlAssociatedItemController.GetAllControlRisks)
	organizationRoutes.GET("/controls/:controlId/incidents", controlAssociatedItemController.GetAllControlIncidents)
	organizationRoutes.GET("/controls/:controlId/associated-item-summary", controlAssociatedItemController.GetAllAssociatedCounts)

	//Process Group
	organizationRoutes.GET("/process-groups", processGroupController.GetAllProcessGroups)
	organizationRoutes.GET("/process-groups/:id", processGroupController.GetProcessGroupDetails)
	organizationRoutes.POST("/process-groups", processGroupController.CreateProcessGroup)
	organizationRoutes.PUT("/process-groups/:id", processGroupController.UpdateProcessGroup)
	organizationRoutes.POST("/process-groups/delete", processGroupController.DeleteProcessGroups)
	organizationRoutes.PUT("/process-groups/activate", processGroupController.ActivateProcessGroups)
	organizationRoutes.PUT("/process-groups/deactivate", processGroupController.DeactivateProcessGroups)

	//control action plan
	organizationRoutes.GET("/controls/:controlId/action-plans", controlActionPlanController.GetAllControlActionPlans)
	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId", controlActionPlanController.GetControlActionPlanDetails)
	organizationRoutes.POST("/controls/:controlId/action-plans", controlActionPlanController.CreateControlActionPlan)
	organizationRoutes.PUT("/controls/:controlId/action-plans/:controlActionPlanId", controlActionPlanController.UpdateControlActionPlan)
	organizationRoutes.POST("/controls/:controlId/action-plans/delete", controlActionPlanController.DeleteControlActionPlans)

	//control action plan status-updates
	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId/status-updates", controlActionPlanStatusUpdateController.GetAllControlActionPlanStatusUpdates)
	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId/status-updates/:controlActionPlanUpdateId", controlActionPlanStatusUpdateController.GetControlActionPlanStatusUpdateDetails)
	organizationRoutes.POST("/controls/:controlId/action-plans/:controlActionPlanId/status-updates", controlActionPlanStatusUpdateController.CreateControlActionPlanStatusUpdate)
	organizationRoutes.PUT("/controls/:controlId/action-plans/:controlActionPlanId/status-updates/:controlActionPlanUpdateId", controlActionPlanStatusUpdateController.UpdateControlActionPlanStatusUpdate)

	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId/status-updates/latest", controlActionPlanStatusUpdateController.GetLatestControlActionPlanStatusUpdate)

	//process
	organizationRoutes.GET("/processes", processController.GetAllProcesses)
	organizationRoutes.GET("/processes/:processId", processController.GetProcessDetails)
	organizationRoutes.POST("/processes", processController.CreateProcess)
	organizationRoutes.PUT("/processes/:processId", processController.UpdateProcess)
	organizationRoutes.POST("/processes/delete", processController.DeleteProcesses)
	organizationRoutes.PUT("/processes/archive", processController.ArchiveProcesses)
	organizationRoutes.PUT("/processes/restore", processController.RestoreProcesses)
	organizationRoutes.GET("/processes/:processId/control-summary", processController.GetControlSummary)

	organizationRoutes.GET("/processes/:processId/controls", processAssociatedItemController.GetAllProcessControls)
	organizationRoutes.GET("/processes/:processId/risks", processAssociatedItemController.GetAllProcessRisks)
	organizationRoutes.GET("/processes/:processId/incidents", processAssociatedItemController.GetAllProcessIncidents)
	organizationRoutes.GET("/processes/:processId/documents", processAssociatedItemController.GetAllProcessDocuments)
	organizationRoutes.GET("/processes/:processId/associated-item-summary", processAssociatedItemController.GetAllAssociatedCounts)

}
