package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	organizationRoutes := router.Group("/")

	organizationRoutes.Use(middlewares.AuthMiddleware())

	db := app.DB

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
	userProfileController := app.OrganizationModule.UserProfileController
	processDependencyTypeController := app.OrganizationModule.ProcessDependencyTypeController
	processDependencySubTypeController := app.OrganizationModule.ProcessDependencySubTypeController
	governanceReportController := app.OrganizationModule.GovernanceReportController
	processDependencyController := app.OrganizationModule.ProcessDependencyController
	processCriticalityLevelController := app.OrganizationModule.ProcessCriticalityLevelController

	//Designation
	organizationRoutes.GET("/designations/template", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), designationController.GetDesignationsTemplate)
	organizationRoutes.POST("/designations/template/validate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), designationController.ValidateDesignationTemplate)
	organizationRoutes.POST("/designations/import/validate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), designationController.ValidateImportDesignation)
	organizationRoutes.POST("/designations/import", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), designationController.ImportDesignations)
	organizationRoutes.GET("/designations", designationController.GetAllDesignations)
	organizationRoutes.GET("/designations/:id", designationController.GetDesignationDetails)
	organizationRoutes.POST("/designations", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), designationController.CreateDesignation)
	organizationRoutes.PUT("/designations/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), designationController.UpdateDesignation)
	organizationRoutes.POST("/designations/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), designationController.DeleteDesignations)
	organizationRoutes.PUT("/designations/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), designationController.ActivateDesignations)
	organizationRoutes.PUT("/designations/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), designationController.DeactivateDesignations)

	//User Status
	organizationRoutes.GET("/user-statuses", userStatusController.GetAllUserStatuses)
	organizationRoutes.GET("/user-statuses/:id", userStatusController.GetUserStatusDetails)

	//User
	organizationRoutes.GET("/users/export", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userController.GetUsersExport)
	organizationRoutes.GET("/users", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userController.GetAllUsers)
	organizationRoutes.GET("/users/:userId", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userController.GetUserDetails)
	organizationRoutes.POST("/users", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userController.CreateUser)
	organizationRoutes.PUT("/users/:userId", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userController.UpdateUser)
	organizationRoutes.POST("/users/delete", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userController.DeleteUsers)
	organizationRoutes.GET("/users/:userId/login-history", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userController.GetUserLoginHistory)

	organizationRoutes.GET("/users/me", userController.GetCurrentUserDetails)
	organizationRoutes.GET("/users/me/permissions", userController.GetUserPermissions)
	organizationRoutes.GET("/users/me/notifications", userController.GetUserNotifications)
	organizationRoutes.GET("/users/me/notifications/new-count", userController.GetUserNewNotificationCount)
	organizationRoutes.PUT("/users/me/notifications/mark-as-read", userController.MarkwNotificationAsRead)
	organizationRoutes.PUT("/users/me/notifications/mark-as-unread", userController.MarkwNotificationAsUnread)

	//User Profile Settings
	organizationRoutes.GET("/users/me/settings/preferences", userProfileController.GetProfileSetting)
	organizationRoutes.PUT("/users/me/settings/preferences", userProfileController.UpdateProfileSetting)

	organizationRoutes.PUT("/users/me/settings/security/enable-mfa", userProfileController.EnableMFA)
	organizationRoutes.PUT("/users/me/settings/security/disable-mfa", userProfileController.DisableMFA)
	organizationRoutes.POST("/users/me/settings/security/reset-password", userProfileController.NewPassword)

	organizationRoutes.PUT("/users/me/settings/enable-notification", userProfileController.EnableMainNotification)
	organizationRoutes.PUT("/users/me/settings/disable-notification", userProfileController.DisableMainNotification)

	organizationRoutes.GET("/users/me/settings/notifications", userProfileController.GetNotificationSettings)
	organizationRoutes.POST("/users/me/settings/notifications/:notificationMasterId/enable", userProfileController.EnableNotification)
	organizationRoutes.POST("/users/me/settings/notifications/:notificationMasterId/disable", userProfileController.DisableNotification)

	organizationRoutes.GET("/users/me/tasks", userProfileController.GetAllUserTasks)
	organizationRoutes.GET("/users/me/task-summary", userProfileController.GetUserTaskSummary)

	organizationRoutes.GET("/users/me/login-history", userProfileController.GetUserLoginHistory)

	organizationRoutes.POST("/users/me/change-password", userProfileController.ChangePassword)

	//Strategic Profile
	organizationRoutes.GET("/strategic-profiles", middlewares.RolePermissionMiddleware(db, "strategic_profile", "view"), strategicProfileController.GetStrategicProfile)
	organizationRoutes.PUT("/strategic-profiles", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileController.UpdateStrategicProfile)
	organizationRoutes.PUT("/strategic-profiles/vision", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileController.UpdateVision)
	organizationRoutes.PUT("/strategic-profiles/mission", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileController.UpdateMission)
	organizationRoutes.PUT("/strategic-profiles/values", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileController.UpdateValues)

	organizationRoutes.POST("/strategic-profiles/upload-image", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileController.UploadProfileImage)
	organizationRoutes.GET("/strategic-profiles/preview-image", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileController.PreviewProfileImage)

	//Strategic Profile Values
	organizationRoutes.GET("/strategic-profiles/values", middlewares.RolePermissionMiddleware(db, "strategic_profile", "view"), strategicProfileValueController.GetAllStrategicProfileValues)
	organizationRoutes.GET("/strategic-profiles/values/:id", middlewares.RolePermissionMiddleware(db, "strategic_profile", "view"), strategicProfileValueController.GetStrategicProfileValueDetails)
	organizationRoutes.POST("/strategic-profiles/values", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileValueController.CreateStrategicProfileValue)
	organizationRoutes.PUT("/strategic-profiles/values/:id", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileValueController.UpdateStrategicProfileValue)
	organizationRoutes.POST("/strategic-profiles/values/delete", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileValueController.DeleteStrategicProfileValues)

	//Strategic Objective
	organizationRoutes.GET("/strategic-objectives", strategicObjectiveController.GetAllStrategicObjectives)
	organizationRoutes.GET("/strategic-objectives/:id", strategicObjectiveController.GetStrategicObjectiveDetails)
	organizationRoutes.POST("/strategic-objectives", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), strategicObjectiveController.CreateStrategicObjective)
	organizationRoutes.PUT("/strategic-objectives/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), strategicObjectiveController.UpdateStrategicObjective)
	organizationRoutes.POST("/strategic-objectives/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), strategicObjectiveController.DeleteStrategicObjectives)
	organizationRoutes.PUT("/strategic-objectives/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), strategicObjectiveController.ActivateStrategicObjectives)
	organizationRoutes.PUT("/strategic-objectives/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), strategicObjectiveController.DeactivateStrategicObjectives)

	//User Settings
	organizationRoutes.GET("/users/:userId/settings/preferences", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userSettingController.GetSetting)
	organizationRoutes.PUT("/users/:userId/settings/preferences", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userSettingController.UpdateSetting)

	organizationRoutes.PUT("/users/:userId/settings/security/enable-mfa", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userSettingController.EnableMFA)
	organizationRoutes.PUT("/users/:userId/settings/security/disable-mfa", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userSettingController.DisableMFA)
	organizationRoutes.POST("/users/:userId/settings/security/reset-password", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userSettingController.NewPassword)

	organizationRoutes.PUT("/users/:userId/settings/enable-notification", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userSettingController.EnableMainNotification)
	organizationRoutes.PUT("/users/:userId/settings/disable-notification", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userSettingController.DisableMainNotification)

	organizationRoutes.GET("/users/:userId/settings/notifications", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userSettingController.GetNotificationSettings)
	organizationRoutes.POST("/users/:userId/settings/notifications/:notificationMasterId/enable", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userSettingController.EnableNotification)
	organizationRoutes.POST("/users/:userId/settings/notifications/:notificationMasterId/disable", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userSettingController.DisableNotification)

	//User Tasks
	organizationRoutes.GET("/users/:userId/tasks", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userTaskController.GetAllUserTasks)
	organizationRoutes.GET("/users/:userId/task-summary", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userTaskController.GetUserTaskSummary)

	//Control Category
	organizationRoutes.GET("/control-categories", controlCategoryController.GetAllControlCategories)
	organizationRoutes.GET("/control-categories/:id", controlCategoryController.GetControlCategoryDetails)
	organizationRoutes.POST("/control-categories", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), controlCategoryController.CreateControlCategory)
	organizationRoutes.PUT("/control-categories/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), controlCategoryController.UpdateControlCategory)
	organizationRoutes.POST("/control-categories/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), controlCategoryController.DeleteControlCategory)
	organizationRoutes.PUT("/control-categories/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), controlCategoryController.ActivateControlCategories)
	organizationRoutes.PUT("/control-categories/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), controlCategoryController.DeactivateControlCategories)

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
	organizationRoutes.GET("/controls/:controlId/reviews", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlReviewController.GetAllControlReviews)
	organizationRoutes.POST("/controls/:controlId/reviews", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlReviewController.CreateControlReview)
	organizationRoutes.PUT("/controls/:controlId/reviews/:id", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlReviewController.UpdateControlReview)

	//Control
	organizationRoutes.POST("/controls/direct-import", middlewares.RolePermissionMiddleware(db, "control", "modify"), controlController.DirectImportControls)
	organizationRoutes.GET("/controls/template", middlewares.RolePermissionMiddleware(db, "control", "modify"), controlController.GenerateControlTemplate)
	organizationRoutes.POST("/controls/template/validate", middlewares.RolePermissionMiddleware(db, "control", "modify"), controlController.ValidateControlTemplate)
	organizationRoutes.POST("/controls/import/validate", middlewares.RolePermissionMiddleware(db, "control", "modify"), controlController.ValidateImportControls)
	organizationRoutes.POST("/controls/import", middlewares.RolePermissionMiddleware(db, "control", "modify"), controlController.ImportControls)

	organizationRoutes.GET("/controls/export", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlController.GetControlsExport)
	organizationRoutes.GET("/controls", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlController.GetAllControls)
	organizationRoutes.GET("/controls/:controlId", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlController.GetControlDetails)
	organizationRoutes.POST("/controls", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.CreateControl)
	organizationRoutes.PUT("/controls/:controlId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.UpdateControl)
	organizationRoutes.POST("/controls/delete", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.DeleteControls)
	organizationRoutes.PUT("/controls/archive", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.ArchiveControls)
	organizationRoutes.PUT("/controls/restore", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.RestoreControls)

	organizationRoutes.PUT("/controls/:controlId/mark-as-applicable", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.MarkControlAsApplicable)
	organizationRoutes.PUT("/controls/:controlId/mark-as-not-applicable", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.MarkControlAsNotApplicable)

	//control Test
	organizationRoutes.GET("/controls/:controlId/tests", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlTestController.GetAllControlTests)
	organizationRoutes.GET("/controls/:controlId/tests/:controlTestId", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlTestController.GetControlTestDetails)
	organizationRoutes.POST("/controls/:controlId/tests", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlTestController.CreateControlTest)
	organizationRoutes.PUT("/controls/:controlId/tests/:controlTestId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlTestController.UpdateControlTest)
	organizationRoutes.POST("/controls/:controlId/tests/delete", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlTestController.DeleteControlTest)
	organizationRoutes.GET("/controls/:controlId/trends", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlTestController.GetControlTestTrendAnalysis)

	//Control Associated Item
	organizationRoutes.GET("/controls/:controlId/processes", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllControlProcesses)
	organizationRoutes.GET("/controls/:controlId/compliance-requirements", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllControlComplianceRequirements)
	organizationRoutes.GET("/controls/:controlId/risks", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllControlRisks)
	organizationRoutes.GET("/controls/:controlId/incidents", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllControlIncidents)
	organizationRoutes.GET("/controls/:controlId/associated-item-summary", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllAssociatedCounts)

	//Process Group
	organizationRoutes.GET("/process-groups", processGroupController.GetAllProcessGroups)
	organizationRoutes.GET("/process-groups/:id", processGroupController.GetProcessGroupDetails)
	organizationRoutes.POST("/process-groups", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processGroupController.CreateProcessGroup)
	organizationRoutes.PUT("/process-groups/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processGroupController.UpdateProcessGroup)
	organizationRoutes.POST("/process-groups/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processGroupController.DeleteProcessGroups)
	organizationRoutes.PUT("/process-groups/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processGroupController.ActivateProcessGroups)
	organizationRoutes.PUT("/process-groups/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processGroupController.DeactivateProcessGroups)

	//control action plan
	organizationRoutes.GET("/controls/:controlId/action-plans", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlActionPlanController.GetAllControlActionPlans)
	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlActionPlanController.GetControlActionPlanDetails)
	organizationRoutes.POST("/controls/:controlId/action-plans", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlActionPlanController.CreateControlActionPlan)
	organizationRoutes.PUT("/controls/:controlId/action-plans/:controlActionPlanId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlActionPlanController.UpdateControlActionPlan)
	organizationRoutes.POST("/controls/:controlId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlActionPlanController.DeleteControlActionPlans)

	//control action plan status-updates
	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlActionPlanStatusUpdateController.GetAllControlActionPlanStatusUpdates)
	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId/status-updates/:controlActionPlanUpdateId", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlActionPlanStatusUpdateController.GetControlActionPlanStatusUpdateDetails)
	organizationRoutes.POST("/controls/:controlId/action-plans/:controlActionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlActionPlanStatusUpdateController.CreateControlActionPlanStatusUpdate)
	organizationRoutes.PUT("/controls/:controlId/action-plans/:controlActionPlanId/status-updates/:controlActionPlanUpdateId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlActionPlanStatusUpdateController.UpdateControlActionPlanStatusUpdate)

	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlActionPlanStatusUpdateController.GetLatestControlActionPlanStatusUpdate)

	//process
	organizationRoutes.GET("/processes/template", middlewares.RolePermissionMiddleware(db, "process", "modify"), processController.GetProcessesTemplate)
	organizationRoutes.POST("/processes/template/validate", middlewares.RolePermissionMiddleware(db, "process", "modify"), processController.ValidateProcessTemplate)
	organizationRoutes.POST("/processes/import/validate", middlewares.RolePermissionMiddleware(db, "process", "modify"), processController.ValidateImportProcess)
	organizationRoutes.POST("/processes/import", middlewares.RolePermissionMiddleware(db, "process", "modify"), processController.ImportProcesses)
	organizationRoutes.GET("/processes/export", middlewares.RolePermissionMiddleware(db, "process", "view"), processController.GetProcessesExport)
	organizationRoutes.GET("/processes", middlewares.RolePermissionMiddleware(db, "process", "view"), processController.GetAllProcesses)
	organizationRoutes.GET("/processes/:processId", middlewares.RolePermissionMiddleware(db, "process", "view"), processController.GetProcessDetails)
	organizationRoutes.POST("/processes", middlewares.RolePermissionMiddleware(db, "process", "modify"), processController.CreateProcess)
	organizationRoutes.PUT("/processes/:processId", middlewares.RolePermissionMiddleware(db, "process", "modify"), processController.UpdateProcess)
	organizationRoutes.POST("/processes/delete", middlewares.RolePermissionMiddleware(db, "process", "modify"), processController.DeleteProcesses)
	organizationRoutes.PUT("/processes/archive", middlewares.RolePermissionMiddleware(db, "process", "modify"), processController.ArchiveProcesses)
	organizationRoutes.PUT("/processes/restore", middlewares.RolePermissionMiddleware(db, "process", "modify"), processController.RestoreProcesses)
	organizationRoutes.GET("/processes/:processId/control-summary", middlewares.RolePermissionMiddleware(db, "process", "view"), processController.GetControlSummary)

	organizationRoutes.GET("/processes/:processId/controls", middlewares.RolePermissionMiddleware(db, "process", "view"), processAssociatedItemController.GetAllProcessControls)
	organizationRoutes.GET("/processes/:processId/risks", middlewares.RolePermissionMiddleware(db, "process", "view"), processAssociatedItemController.GetAllProcessRisks)
	organizationRoutes.GET("/processes/:processId/incidents", middlewares.RolePermissionMiddleware(db, "process", "view"), processAssociatedItemController.GetAllProcessIncidents)
	organizationRoutes.GET("/processes/:processId/assets", middlewares.RolePermissionMiddleware(db, "process", "view"), processAssociatedItemController.GetAllProcessAssets)
	organizationRoutes.GET("/processes/:processId/documents", middlewares.RolePermissionMiddleware(db, "process", "view"), processAssociatedItemController.GetAllProcessDocuments)
	organizationRoutes.GET("/processes/:processId/associated-item-summary", middlewares.RolePermissionMiddleware(db, "process", "view"), processAssociatedItemController.GetAllAssociatedCounts)

	// process dependency type
	organizationRoutes.GET("/process-dependency-types", processDependencyTypeController.GetAllProcessDependencyTypes)
	organizationRoutes.GET("/process-dependency-types/:id", processDependencyTypeController.GetProcessDependencyTypeDetails)

	organizationRoutes.GET("/process-dependency-sub-types", processDependencySubTypeController.GetAllProcessDependencySubTypes)
	organizationRoutes.GET("/process-dependency-sub-types/:id", processDependencySubTypeController.GetProcessDependencySubTypeDetails)
	organizationRoutes.POST("/process-dependency-sub-types", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processDependencySubTypeController.CreateProcessDependencySubType)
	organizationRoutes.PUT("/process-dependency-sub-types/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processDependencySubTypeController.UpdateProcessDependencySubType)
	organizationRoutes.POST("/process-dependency-sub-types/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processDependencySubTypeController.DeleteProcessDependencySubTypes)
	organizationRoutes.PUT("/process-dependency-sub-types/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processDependencySubTypeController.ActivateProcessDependencySubTypes)
	organizationRoutes.PUT("/process-dependency-sub-types/deactivate", middlewares.RolePermissionMiddleware(db, "processorganization_setting_dependency_setting", "modify"), processDependencySubTypeController.DeactivateProcessDependencySubTypes)

	organizationRoutes.GET("/governance-reports", governanceReportController.GetAllGovernanceReports)
	organizationRoutes.GET("/governance-reports/:governanceReportId", governanceReportController.GetGovernanceReportDetails)
	organizationRoutes.GET("/governance-reports/:governanceReportId/export-pdf", governanceReportController.ExportGovernanceReportToPDF)
	organizationRoutes.POST("/governance-reports", governanceReportController.CreateGovernanceReport)
	organizationRoutes.GET("/governance-reports/:governanceReportId/sections/:sectionId", governanceReportController.GetGovernanceReportSection)
	organizationRoutes.PUT("/governance-reports/:governanceReportId/sections/:sectionId", governanceReportController.UpdateGovernanceReportSection)
	organizationRoutes.DELETE("/governance-reports/:governanceReportId/sections/:sectionId", governanceReportController.DeleteGovernanceReportSection)
	organizationRoutes.DELETE("/governance-reports/:governanceReportId", governanceReportController.DeleteGovernanceReport)

	organizationRoutes.GET("/processes/:processId/dependencies", processDependencyController.GetAllProcessDependencies)
	organizationRoutes.GET("/processes/:processId/dependencies/:dependencyId", processDependencyController.GetProcessDependencyDetails)
	organizationRoutes.POST("/processes/:processId/dependencies", middlewares.RolePermissionMiddleware(db, "process", "modify"), processDependencyController.CreateProcessDependency)
	organizationRoutes.PUT("/processes/:processId/dependencies/:dependencyId", middlewares.RolePermissionMiddleware(db, "process", "modify"), processDependencyController.UpdateProcessDependency)
	organizationRoutes.POST("/processes/:processId/dependencies/delete", middlewares.RolePermissionMiddleware(db, "process", "modify"), processDependencyController.DeleteProcessDependencies)
	// process criticality level
	organizationRoutes.GET("/process-criticality-levels", processCriticalityLevelController.GetAllProcessCriticalityLevels)
	organizationRoutes.GET("/process-criticality-levels/:id", processCriticalityLevelController.GetProcessCriticalityLevelDetails)
	organizationRoutes.POST("/ad-users-sync", userController.SyncADUsers)
	organizationRoutes.GET("/active-directory-users", userController.GetADUsers)

}
