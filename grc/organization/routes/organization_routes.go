package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	organizationRoutes := router.Group("/")

	organizationRoutes.Use(middlewares.AuthMiddleware(app.DB))

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	designationController := app.OrganizationModule.DesignationController
	userStatusController := app.OrganizationModule.UserStatusController
	userController := app.OrganizationModule.UserController
	userSettingController := app.OrganizationModule.UserSettingController
	strategicProfileController := app.OrganizationModule.StrategicProfileController
	strategicProfileValueController := app.OrganizationModule.StrategicProfileValueController
	// strategicObjectiveController := app.OrganizationModule.StrategicObjectiveController
	controlCategoryController := app.OrganizationModule.ControlCategoryController
	controlTypeController := app.OrganizationModule.ControlTypeController
	controlImplementationStatusController := app.OrganizationModule.ControlImplementationStatusController
	controlEfficiencyStatusController := app.OrganizationModule.ControlEfficiencyStatusController
	controlController := app.OrganizationModule.ControlController
	controlTestController := app.OrganizationModule.ControlTestController
	processGroupController := app.OrganizationModule.ProcessGroupController
	controlActionPlanController := app.OrganizationModule.ControlActionPlanController
	controlActionPlanStatusUpdateController := app.OrganizationModule.ControlActionPlanStatusUpdateController
	controlActionPlanCommentController := app.OrganizationModule.ControlActionPlanCommentController
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
	processActivityController := app.OrganizationModule.ProcessActivityController
	competencyGroupController := app.OrganizationModule.CompetencyGroupController
	competencyController := app.OrganizationModule.CompetencyController
	jobDescriptionController := app.OrganizationModule.JobDescriptionController
	jobDescriptionCompetencyController := app.OrganizationModule.JobDescriptionCompetencyController
	competencyLevelController := app.OrganizationModule.CompetencyLevelController
	userDocumentTypeController := app.OrganizationModule.UserDocumentTypeController
	userDocumentSubTypeController := app.OrganizationModule.UserDocumentSubTypeController
	userDocumentController := app.OrganizationModule.UserDocumentController
	userJobDescriptionController := app.OrganizationModule.UserJobDescriptionController
	userDelegationHistoryController := app.OrganizationModule.UserDelegationHistoryController
	organizationUnitsController := app.OrganizationModule.OrganizationUnitsController
	organizationUnitBusinessContextController := app.OrganizationModule.OrganizationUnitBusinessContextController
	organizationUnitAttributeController := app.OrganizationModule.OrganizationUnitAttributeController
	organizationChartUsers := app.OrganizationModule.OrganizationCartController
	myDashboardController := app.OrganizationModule.MyDashboardController
	committeeTermController := app.OrganizationModule.CommitteeTermController
	committeeMemberTypeController := app.OrganizationModule.CommitteeMemberTypeController
	committeeTypeController := app.OrganizationModule.CommitteeTypeController
	committeeClassificationController := app.OrganizationModule.CommitteeClassificationController
	internalCommitteeStatusController := app.OrganizationModule.InternalCommitteeStatusController
	internalCommitteeController := app.OrganizationModule.InternalCommitteeController
	internalCommitteeWorkflowReviewController := app.OrganizationModule.InternalCommitteeWorkflowReviewController
	internalCommitteeWorkflowController := app.OrganizationModule.InternalCommitteeWorkflowController
	internalCommitteeReviewCommentController := app.OrganizationModule.InternalCommitteeReviewCommentController
	controlEffectivenessRatingController := app.OrganizationModule.ControlEffectivenessRatingController
	controlStatusController := app.OrganizationModule.ControlStatusController
	governanceSettingsController := app.OrganizationModule.GovernanceSettingsController
	processDependencyTierController := app.OrganizationModule.ProcessDependencyTierController
	processExternalDependencyEntityTypeController := app.OrganizationModule.ProcessExternalDependencyEntityTypeController
	processInternalDependencyController := app.OrganizationModule.ProcessInternalDependencyController
	processExternalDependencyController := app.OrganizationModule.ProcessExternalDependencyController
	processEmployeeAnalysisController := app.OrganizationModule.ProcessEmployeeAnalysisController
	processEquipmentToolController := app.OrganizationModule.ProcessEquipmentToolController
	// Designation
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

	// User Status
	organizationRoutes.GET("/user-statuses", userStatusController.GetAllUserStatuses)
	organizationRoutes.GET("/user-statuses/:id", userStatusController.GetUserStatusDetails)

	// User
	organizationRoutes.GET("/users/export", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userController.GetUsersExport)
	organizationRoutes.GET("/users/template", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userController.GenerateUserTemplate)
	organizationRoutes.POST("/users/template/validate", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userController.ValidateUserTemplate)
	organizationRoutes.POST("/users/import/validate", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userController.ValidateImportUsers)
	organizationRoutes.POST("/users/import", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userController.ImportUsers)
	organizationRoutes.GET("/users", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userController.GetAllUsers)
	organizationRoutes.GET("/users/:userId", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userController.GetUserDetails)
	organizationRoutes.POST("/users", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userController.CreateUser)
	organizationRoutes.PUT("/users/:userId", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userController.UpdateUser)
	organizationRoutes.POST("/users/delete", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userController.DeleteUsers)
	organizationRoutes.GET("/users/:userId/login-history", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userController.GetUserLoginHistory)
	organizationRoutes.PUT("/users/:userId/change-unit", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userController.UpdateUserOrganizationUnit)
	organizationRoutes.GET("/users/lookup", userController.GetAllUsersWithPrimaryUnit)

	organizationRoutes.GET("/users/me", userController.GetCurrentUserDetails)
	organizationRoutes.GET("/users/me/permissions", userController.GetUserPermissions)
	organizationRoutes.GET("/users/me/notifications", userController.GetUserNotifications)
	organizationRoutes.GET("/users/me/notifications/new-count", userController.GetUserNewNotificationCount)
	organizationRoutes.PUT("/users/me/notifications/mark-as-read", userController.MarkwNotificationAsRead)
	organizationRoutes.PUT("/users/me/notifications/mark-as-unread", userController.MarkwNotificationAsUnread)

	// User Profile Settings
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

	organizationRoutes.GET("/users/me/documents", userProfileController.GetAllUserDocuments)
	organizationRoutes.GET("/users/me/documents/:id", userProfileController.GetUserDocumentDetails)
	organizationRoutes.POST("/users/me/documents", userProfileController.CreateUserDocument)
	organizationRoutes.PUT("/users/me/documents/:id", userProfileController.UpdateUserDocument)
	organizationRoutes.POST("/users/me/documents/delete", userProfileController.DeleteUserDocuments)
	organizationRoutes.GET("/users/me/job-descriptions", userProfileController.GetAllUserJobDescriptions)

	// Strategic Profile
	organizationRoutes.GET("/strategic-profiles", middlewares.RolePermissionMiddleware(db, "strategic_profile", "view"), strategicProfileController.GetStrategicProfile)
	organizationRoutes.PUT("/strategic-profiles", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileController.UpdateStrategicProfile)
	organizationRoutes.PUT("/strategic-profiles/vision", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileController.UpdateVision)
	organizationRoutes.PUT("/strategic-profiles/mission", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileController.UpdateMission)
	organizationRoutes.PUT("/strategic-profiles/values", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileController.UpdateValues)

	organizationRoutes.POST("/strategic-profiles/upload-image", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileController.UploadProfileImage)
	organizationRoutes.GET("/strategic-profiles/preview-image", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileController.PreviewProfileImage)

	// Strategic Profile Values
	organizationRoutes.GET("/strategic-profiles/values", middlewares.RolePermissionMiddleware(db, "strategic_profile", "view"), strategicProfileValueController.GetAllStrategicProfileValues)
	organizationRoutes.GET("/strategic-profiles/values/:id", middlewares.RolePermissionMiddleware(db, "strategic_profile", "view"), strategicProfileValueController.GetStrategicProfileValueDetails)
	organizationRoutes.POST("/strategic-profiles/values", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileValueController.CreateStrategicProfileValue)
	organizationRoutes.PUT("/strategic-profiles/values/:id", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileValueController.UpdateStrategicProfileValue)
	organizationRoutes.POST("/strategic-profiles/values/delete", middlewares.RolePermissionMiddleware(db, "strategic_profile", "modify"), strategicProfileValueController.DeleteStrategicProfileValues)

	// Strategic Objective
	// organizationRoutes.GET("/strategic-objectives", strategicObjectiveController.GetAllStrategicObjectives)
	// organizationRoutes.GET("/strategic-objectives/:id", strategicObjectiveController.GetStrategicObjectiveDetails)
	// organizationRoutes.POST("/strategic-objectives", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), strategicObjectiveController.CreateStrategicObjective)
	// organizationRoutes.PUT("/strategic-objectives/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), strategicObjectiveController.UpdateStrategicObjective)
	// organizationRoutes.POST("/strategic-objectives/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), strategicObjectiveController.DeleteStrategicObjectives)
	// organizationRoutes.PUT("/strategic-objectives/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), strategicObjectiveController.ActivateStrategicObjectives)
	// organizationRoutes.PUT("/strategic-objectives/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), strategicObjectiveController.DeactivateStrategicObjectives)

	// User Settings
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

	// User Tasks
	organizationRoutes.GET("/users/:userId/tasks", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userTaskController.GetAllUserTasks)
	organizationRoutes.GET("/users/:userId/task-summary", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userTaskController.GetUserTaskSummary)

	// Control Category
	organizationRoutes.GET("/control-categories", controlCategoryController.GetAllControlCategories)
	organizationRoutes.GET("/control-categories/:id", controlCategoryController.GetControlCategoryDetails)
	organizationRoutes.POST("/control-categories", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), controlCategoryController.CreateControlCategory)
	organizationRoutes.PUT("/control-categories/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), controlCategoryController.UpdateControlCategory)
	organizationRoutes.POST("/control-categories/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), controlCategoryController.DeleteControlCategory)
	organizationRoutes.PUT("/control-categories/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), controlCategoryController.ActivateControlCategories)
	organizationRoutes.PUT("/control-categories/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), controlCategoryController.DeactivateControlCategories)

	// Control Type
	organizationRoutes.GET("/control-types", controlTypeController.GetAllControlTypes)

	// Control Implementation Status
	organizationRoutes.GET("/control-implementation-statuses", controlImplementationStatusController.GetAllControlImplementationStatuses)
	organizationRoutes.GET("/control-implementation-statuses/:id", controlImplementationStatusController.GetControlImplementationStatusDetails)

	// Control Efficiency Status
	organizationRoutes.GET("/control-efficiency-statuses", controlEfficiencyStatusController.GetAllControlEfficiencyStatuses)
	organizationRoutes.GET("/control-efficiency-statuses/:id", controlEfficiencyStatusController.GetControlEfficiencyStatusDetails)

	// Control Review Statuses
	organizationRoutes.GET("/control-review-statuses", controlReviewStatusController.GetAllControlReviewStatuses)
	organizationRoutes.GET("/control-review-statuses/:id", controlReviewStatusController.GetControlReviewStatusDetails)

	// Control Review
	organizationRoutes.GET("/controls/:controlId/reviews", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlReviewController.GetAllControlReviews)
	organizationRoutes.POST("/controls/:controlId/reviews", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlReviewController.CreateControlReview)
	organizationRoutes.PUT("/controls/:controlId/reviews/:id", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlReviewController.UpdateControlReview)

	// Control
	organizationRoutes.POST("/controls/direct-import", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.DirectImportControls)
	organizationRoutes.GET("/controls/template", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.GenerateControlTemplate)
	organizationRoutes.POST("/controls/template/validate", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.ValidateControlTemplate)
	organizationRoutes.POST("/controls/import/validate", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.ValidateImportControls)
	organizationRoutes.POST("/controls/import", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.ImportControls)

	organizationRoutes.GET("/controls/export", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlController.GetControlsExport)
	organizationRoutes.GET("/controls", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlController.GetAllControls)
	organizationRoutes.GET("/controls/:controlId", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlController.GetControlDetails)
	organizationRoutes.POST("/controls", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.CreateControl)
	organizationRoutes.PUT("/controls/:controlId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.UpdateControl)
	organizationRoutes.POST("/controls/delete", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.DeleteControls)
	organizationRoutes.POST("/controls/force-delete", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.DeleteForceControls)
	organizationRoutes.PUT("/controls/archive", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.ArchiveControls)
	organizationRoutes.PUT("/controls/restore", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.RestoreControls)

	organizationRoutes.PUT("/controls/:controlId/mark-as-applicable", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.MarkControlAsApplicable)
	organizationRoutes.PUT("/controls/:controlId/mark-as-not-applicable", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlController.MarkControlAsNotApplicable)

	// control Test
	unauthenticatedRoutes.GET("/controls/:controlId/tests/:controlTestId/documents/:id/document-download", controlTestController.DownloadControlTestDocument)
	unauthenticatedRoutes.GET("/controls/:controlId/tests/:controlTestId/documents/:id/document-preview", controlTestController.GetControlTestDocumentPreview)
	organizationRoutes.GET("/controls/:controlId/tests", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlTestController.GetAllControlTests)
	organizationRoutes.GET("/controls/:controlId/tests/:controlTestId", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlTestController.GetControlTestDetails)
	organizationRoutes.POST("/controls/:controlId/tests", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlTestController.CreateControlTest)
	organizationRoutes.PUT("/controls/:controlId/tests/:controlTestId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlTestController.UpdateControlTest)
	organizationRoutes.POST("/controls/:controlId/tests/delete", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlTestController.DeleteControlTest)
	organizationRoutes.GET("/controls/:controlId/trends", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlTestController.GetControlTestTrendAnalysis)

	// Control Associated Item
	organizationRoutes.GET("/controls/:controlId/processes", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllControlProcesses)
	organizationRoutes.GET("/controls/:controlId/compliance-requirements", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllControlComplianceRequirements)
	organizationRoutes.GET("/controls/:controlId/risks", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllControlRisks)
	organizationRoutes.GET("/controls/:controlId/incidents", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllControlIncidents)
	organizationRoutes.GET("/controls/:controlId/associated-item-summary", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllAssociatedCounts)
	organizationRoutes.GET("/controls/:controlId/clauses", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllClauses)
	organizationRoutes.PUT("/controls/:controlId/clauses", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.UpdateClauses)

	organizationRoutes.POST("/controls/:controlId/risks", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.CreateControlRisk)
	organizationRoutes.DELETE("/controls/:controlId/risks/:riskId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.DeleteControlRisk)

	organizationRoutes.GET("/controls/:controlId/documents", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllControlDocument)
	organizationRoutes.POST("/controls/:controlId/documents", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.CreateControlDocument)
	organizationRoutes.DELETE("/controls/:controlId/documents/:documentId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.DeleteControlDocument)

	organizationRoutes.POST("/controls/:controlId/processes", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.CreateControlPocess)
	organizationRoutes.DELETE("/controls/:controlId/processes/:processId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.DeleteControlProcess)

	organizationRoutes.GET("/controls/:controlId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllControlStrategicObjectives)
	organizationRoutes.POST("/controls/:controlId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.CreateControlStrategicObjectives)
	organizationRoutes.DELETE("/controls/:controlId/strategic-objectives/:strategicObjectiveId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.DeleteControlStrategicObjectives)

	organizationRoutes.GET("/controls/:controlId/kpis", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllControlKPI)
	organizationRoutes.POST("/controls/:controlId/kpis", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.CreateControlKPI)
	organizationRoutes.DELETE("/controls/:controlId/kpis/:kpiRegisterId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.DeleteControlKPI)

	organizationRoutes.GET("/controls/:controlId/assets", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlAssociatedItemController.GetAllControlAssets)
	organizationRoutes.POST("/controls/:controlId/assets", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.CreateControlAsset)
	organizationRoutes.DELETE("/controls/:controlId/assets/:assetId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.DeleteControlAsset)

	organizationRoutes.POST("/controls/:controlId/incidents", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.CreateControlIncident)
	organizationRoutes.DELETE("/controls/:controlId/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlAssociatedItemController.DeleteControlIncident)

	//Process Group
	// Process Group
	// Process Group
	organizationRoutes.GET("/process-groups", processGroupController.GetAllProcessGroups)
	organizationRoutes.GET("/process-groups/:id", processGroupController.GetProcessGroupDetails)
	organizationRoutes.POST("/process-groups", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processGroupController.CreateProcessGroup)
	organizationRoutes.PUT("/process-groups/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processGroupController.UpdateProcessGroup)
	organizationRoutes.POST("/process-groups/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processGroupController.DeleteProcessGroups)
	organizationRoutes.PUT("/process-groups/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processGroupController.ActivateProcessGroups)
	organizationRoutes.PUT("/process-groups/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processGroupController.DeactivateProcessGroups)

	// control action plan
	organizationRoutes.GET("/controls/:controlId/action-plans", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlActionPlanController.GetAllControlActionPlans)
	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlActionPlanController.GetControlActionPlanDetails)
	organizationRoutes.POST("/controls/:controlId/action-plans", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlActionPlanController.CreateControlActionPlan)
	organizationRoutes.PUT("/controls/:controlId/action-plans/:controlActionPlanId", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlActionPlanController.UpdateControlActionPlan)
	organizationRoutes.POST("/controls/:controlId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "controls", "modify"), controlActionPlanController.DeleteControlActionPlans)

	// control action plan status-updates
	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlActionPlanStatusUpdateController.GetAllControlActionPlanStatusUpdates)
	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId/status-updates/:controlActionPlanUpdateId", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlActionPlanStatusUpdateController.GetControlActionPlanStatusUpdateDetails)
	organizationRoutes.POST("/controls/:controlId/action-plans/:controlActionPlanId/status-updates", controlActionPlanStatusUpdateController.CreateControlActionPlanStatusUpdate)
	organizationRoutes.PUT("/controls/:controlId/action-plans/:controlActionPlanId/status-updates/:controlActionPlanUpdateId", controlActionPlanStatusUpdateController.UpdateControlActionPlanStatusUpdate)

	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "controls", "view"), controlActionPlanStatusUpdateController.GetLatestControlActionPlanStatusUpdate)
	organizationRoutes.PUT("/controls/:controlId/action-plans/:controlActionPlanId/approve", controlActionPlanStatusUpdateController.ApproveControlActionPlan)
	organizationRoutes.PUT("/controls/:controlId/action-plans/:controlActionPlanId/revert", controlActionPlanStatusUpdateController.RevertControlActionPlan)

	// control Action Plan Comments
	unauthenticatedRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId/comments/:commentId/document-download", controlActionPlanCommentController.DownloadControlActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId/comments/:commentId/document-preview", controlActionPlanCommentController.GetControlActionPlanCommentDocumentPreview)
	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId/comments", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), controlActionPlanCommentController.GetAllComments)
	organizationRoutes.GET("/controls/:controlId/action-plans/:controlActionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), controlActionPlanCommentController.GetCommentDetails)
	organizationRoutes.POST("/controls/:controlId/action-plans/:controlActionPlanId/comments", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), controlActionPlanCommentController.CreateComment)
	organizationRoutes.PUT("/controls/:controlId/action-plans/:controlActionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), controlActionPlanCommentController.UpdateComment)
	organizationRoutes.POST("/controls/:controlId/action-plans/:controlActionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), controlActionPlanCommentController.UpdateComment)
	organizationRoutes.POST("/controls/:controlId/action-plans/:controlActionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), controlActionPlanCommentController.DeleteComment)

	// process
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
	organizationRoutes.POST("/processes/:processId/incidents", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.CreateProcessIncident)
	organizationRoutes.DELETE("/processes/:processId/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.DeleteProcessIncident)
	organizationRoutes.GET("/processes/:processId/assets", middlewares.RolePermissionMiddleware(db, "process", "view"), processAssociatedItemController.GetAllProcessAssets)
	organizationRoutes.POST("/processes/:processId/assets", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.CreateProcessAsset)
	organizationRoutes.DELETE("/processes/:processId/assets/:assetId", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.DeleteProcessAsset)
	organizationRoutes.GET("/processes/:processId/documents", middlewares.RolePermissionMiddleware(db, "process", "view"), processAssociatedItemController.GetAllProcessDocuments)
	organizationRoutes.GET("/processes/:processId/associated-item-summary", middlewares.RolePermissionMiddleware(db, "process", "view"), processAssociatedItemController.GetAllAssociatedCounts)
	organizationRoutes.GET("/processes/:processId/clauses", middlewares.RolePermissionMiddleware(db, "process", "view"), processAssociatedItemController.GetAllClauses)
	organizationRoutes.PUT("/processes/:processId/clauses", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.UpdateClauses)

	organizationRoutes.POST("/processes/:processId/risks", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.CreateProcessRisk)
	organizationRoutes.DELETE("/processes/:processId/risks/:riskId", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.DeleteProcessRisk)

	organizationRoutes.POST("/processes/:processId/documents", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.CreateProcessDocument)
	organizationRoutes.DELETE("/processes/:processId/documents/:documentId", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.DeleteProcessDocument)

	organizationRoutes.POST("/processes/:processId/controls", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.CreateProcessControl)
	organizationRoutes.DELETE("/processes/:processId/controls/:controlId", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.DeleteProcessControl)

	organizationRoutes.GET("/processes/:processId/kpis", middlewares.RolePermissionMiddleware(db, "process", "view"), processAssociatedItemController.GetAllProcessKPI)
	organizationRoutes.POST("/processes/:processId/kpis", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.CreateProcessKPI)
	organizationRoutes.DELETE("/processes/:processId/kpis/:kpiRegisterId", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.DeleteProcessKPI)

	organizationRoutes.GET("/processes/:processId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "process", "view"), processAssociatedItemController.GetAllProcessStrategicObjectives)
	organizationRoutes.POST("/processes/:processId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.CreateProcessStrategicObjective)
	organizationRoutes.DELETE("/processes/:processId/strategic-objectives/:strategicObjectiveId", middlewares.RolePermissionMiddleware(db, "process", "modify"), processAssociatedItemController.DeleteProcessStrategicObjective)

	organizationRoutes.GET("/processes/:processId/associated-controls", middlewares.RolePermissionMiddleware(db, "process", "view"), processAssociatedItemController.GetAllProcessAssociatedControls)

	// process dependency type
	organizationRoutes.GET("/process-dependency-types", processDependencyTypeController.GetAllProcessDependencyTypes)
	organizationRoutes.GET("/process-dependency-types/:id", processDependencyTypeController.GetProcessDependencyTypeDetails)

	organizationRoutes.GET("/process-dependency-sub-types", processDependencySubTypeController.GetAllProcessDependencySubTypes)
	organizationRoutes.GET("/process-dependency-sub-types/:id", processDependencySubTypeController.GetProcessDependencySubTypeDetails)
	organizationRoutes.POST("/process-dependency-sub-types", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processDependencySubTypeController.CreateProcessDependencySubType)
	organizationRoutes.PUT("/process-dependency-sub-types/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processDependencySubTypeController.UpdateProcessDependencySubType)
	organizationRoutes.POST("/process-dependency-sub-types/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processDependencySubTypeController.DeleteProcessDependencySubTypes)
	organizationRoutes.PUT("/process-dependency-sub-types/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processDependencySubTypeController.ActivateProcessDependencySubTypes)
	organizationRoutes.PUT("/process-dependency-sub-types/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processDependencySubTypeController.DeactivateProcessDependencySubTypes)

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

	organizationRoutes.GET("/active-directory/ous", userController.ListOUs)
	organizationRoutes.GET("/active-directory/ou-groups", userController.ListGroupsInOU)
	organizationRoutes.POST("/active-directory/ou-group-users", userController.ListUsersInGroups)
	organizationRoutes.GET("/active-directory-users", userController.GetADUsers)

	organizationRoutes.GET("/active-directory-settings/sync-users", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userController.SyncAzureUsers)
	organizationRoutes.GET("/active-directory-settings/sync-groups", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userController.SyncAzureGroups)

	organizationRoutes.GET("/processes/:processId/activities", processActivityController.GetAllProcessActivities)
	organizationRoutes.GET("/processes/:processId/activities/:activityId", processActivityController.GetProcessActivityDetails)
	organizationRoutes.POST("/processes/:processId/activities", middlewares.RolePermissionMiddleware(db, "process", "modify"), processActivityController.CreateProcessActivity)
	organizationRoutes.PUT("/processes/:processId/activities/:activityId", middlewares.RolePermissionMiddleware(db, "process", "modify"), processActivityController.UpdateProcessActivity)
	organizationRoutes.POST("/processes/:processId/activities/delete", middlewares.RolePermissionMiddleware(db, "process", "modify"), processActivityController.DeleteProcessActivities)
	organizationRoutes.PUT("/processes/:processId/activities/:activityId/raci", middlewares.RolePermissionMiddleware(db, "process", "modify"), processActivityController.UpdateProcessActivityUsers)
	organizationRoutes.DELETE("/processes/:processId/activities/:activityId/raci/:activityUserId", middlewares.RolePermissionMiddleware(db, "process", "modify"), processActivityController.DeleteActivityUser)

	organizationRoutes.GET("/competency-groups", competencyGroupController.GetAllCompetencyGroups)
	organizationRoutes.GET("/competency-groups/:competencyGroupId", competencyGroupController.GetCompetencyGroupDetails)
	organizationRoutes.POST("/competency-groups", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), competencyGroupController.CreateCompetencyGroup)
	organizationRoutes.PUT("/competency-groups/:competencyGroupId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), competencyGroupController.UpdateCompetencyGroup)
	organizationRoutes.POST("/competency-groups/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), competencyGroupController.DeleteCompetencyGroups)
	organizationRoutes.PUT("/competency-groups/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), competencyGroupController.ActivateCompetencyGroups)
	organizationRoutes.PUT("/competency-groups/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), competencyGroupController.DeactivateCompetencyGroups)

	organizationRoutes.GET("/competencies", competencyController.GetAllCompetencies)
	organizationRoutes.GET("/competencies/:id", competencyController.GetCompetencyDetails)
	organizationRoutes.POST("/competencies", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), competencyController.CreateCompetency)
	organizationRoutes.PUT("/competencies/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), competencyController.UpdateCompetency)
	organizationRoutes.POST("/competencies/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), competencyController.DeleteCompetencies)
	organizationRoutes.PUT("/competencies/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), competencyController.ActivateCompetencies)
	organizationRoutes.PUT("/competencies/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), competencyController.DeactivateCompetencies)

	organizationRoutes.GET("/competency-levels", competencyLevelController.GetAllCompetencyLevels)
	organizationRoutes.GET("/competency-levels/:id", competencyLevelController.GetCompetencyLevelDetails)

	organizationRoutes.GET("/job-descriptions", jobDescriptionController.GetAllJobDescriptions)
	organizationRoutes.GET("/job-descriptions/:id", jobDescriptionController.GetJobDescriptionDetails)
	organizationRoutes.POST("/job-descriptions", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), jobDescriptionController.CreateJobDescription)
	organizationRoutes.PUT("/job-descriptions/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), jobDescriptionController.UpdateJobDescription)
	organizationRoutes.POST("/job-descriptions/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), jobDescriptionController.DeleteJobDescription)
	organizationRoutes.PUT("/job-descriptions/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), jobDescriptionController.ActivateJobDescriptions)
	organizationRoutes.PUT("/job-descriptions/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), jobDescriptionController.DeactivateJobDescriptions)

	organizationRoutes.POST("/job-descriptions/:id/competencies", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), jobDescriptionCompetencyController.MapCompetencies)
	organizationRoutes.PUT("/job-descriptions/:id/competencies/:cId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), jobDescriptionCompetencyController.UpdateCompetencyDetails)
	organizationRoutes.DELETE("/job-descriptions/:id/competencies/:cId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), jobDescriptionCompetencyController.DeleteCompetency)

	organizationRoutes.GET("/users/:userId/job-descriptions", userJobDescriptionController.GetAllUserJobDescriptions)
	organizationRoutes.GET("/users/:userId/job-descriptions/:id", userJobDescriptionController.GetUserJobDescriptionDetailsByID)
	organizationRoutes.POST("/users/:userId/job-descriptions", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userJobDescriptionController.AssignJobDescriptionsToUser)
	organizationRoutes.PUT("/users/:userId/job-descriptions/:id/competencies/:cId", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userJobDescriptionController.EvaluateCompetency)
	organizationRoutes.DELETE("/users/:userId/job-descriptions/:id", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userJobDescriptionController.DeleteJobDescriptionUser)
	unauthenticatedRoutes.GET("/users/:userId/job-descriptions/:id/competencies/:cId/document/:docId/document-download", userJobDescriptionController.DownloadUserJobDescriptionDocument)
	unauthenticatedRoutes.GET("/users/:userId/job-descriptions/:id/competencies/:cId/document/:docId/document-preview", userJobDescriptionController.GetUserJobDescriptionDocumentPreview)

	// user document type
	organizationRoutes.GET("/user-document-types", userDocumentTypeController.GetAllUserDocumentTypes)
	organizationRoutes.GET("/user-document-types/:id", userDocumentTypeController.GetUserDocumentTypeDetails)
	organizationRoutes.POST("/user-document-types", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), userDocumentTypeController.CreateUserDocumentType)
	organizationRoutes.PUT("/user-document-types/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), userDocumentTypeController.UpdateUserDocumentType)
	organizationRoutes.POST("/user-document-types/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), userDocumentTypeController.DeleteUserDocumentTypes)
	organizationRoutes.PUT("/user-document-types/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), userDocumentTypeController.ActivateUserDocumentTypes)
	organizationRoutes.PUT("/user-document-types/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), userDocumentTypeController.DeactivateUserDocumentTypes)

	// user document sub type
	organizationRoutes.GET("/user-document-sub-types", userDocumentSubTypeController.GetAllUserDocumentSubTypes)
	organizationRoutes.GET("/user-document-sub-types/:id", userDocumentSubTypeController.GetUserDocumentSubTypeDetails)
	organizationRoutes.POST("/user-document-sub-types", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), userDocumentSubTypeController.CreateUserDocumentSubType)
	organizationRoutes.PUT("/user-document-sub-types/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), userDocumentSubTypeController.UpdateUserDocumentSubType)
	organizationRoutes.POST("/user-document-sub-types/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), userDocumentSubTypeController.DeleteUserDocumentSubTypes)
	organizationRoutes.PUT("/user-document-sub-types/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), userDocumentSubTypeController.ActivateUserDocumentSubTypes)
	organizationRoutes.PUT("/user-document-sub-types/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), userDocumentSubTypeController.DeactivateUserDocumentSubTypes)

	// user document
	unauthenticatedRoutes.GET("/users/:userId/documents/:id/document-download", userDocumentController.DownloadUserDocument)
	unauthenticatedRoutes.GET("/users/:userId/documents/:id/document-preview", userDocumentController.GetUserDocumentPreview)
	organizationRoutes.GET("/users/:userId/documents", userDocumentController.GetAllUserDocuments)
	organizationRoutes.GET("/users/:userId/documents/:id", userDocumentController.GetUserDocumentDetails)
	organizationRoutes.POST("/users/:userId/documents", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userDocumentController.CreateUserDocument)
	organizationRoutes.PUT("/users/:userId/documents/:id", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userDocumentController.UpdateUserDocument)
	organizationRoutes.POST("/users/:userId/documents/delete", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userDocumentController.DeleteUserDocuments)

	// user delegation
	organizationRoutes.GET("/users/:userId/modules", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userController.GetUserModulesWithPermissions)
	organizationRoutes.GET("/users/:userId/task-modules", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userController.GetUserModulesWithOpenTasks)
	organizationRoutes.POST("/users/:userId/user-delegate", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userDelegationHistoryController.CreateUserDelegation)
	organizationRoutes.GET("/users/:userId/user-delegate-history", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userDelegationHistoryController.GetUserDelegationHistory)
	organizationRoutes.GET("/users/:userId/user-delegate-history/latest", middlewares.RolePermissionMiddleware(db, "user_management", "view"), userDelegationHistoryController.GetLatestUserDelegation)
	organizationRoutes.PUT("/users/:userId/end-delegation", middlewares.RolePermissionMiddleware(db, "user_management", "modify"), userDelegationHistoryController.EndUserDelegation)

	// organization structure
	organizationRoutes.GET("/organization-units/hierarchies", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationUnitsController.GetOrganizationUnitsHierarchies)
	organizationRoutes.GET("/organization-units/:id/users", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationUnitsController.GetUsersByOrganizationUnitID)
	organizationRoutes.PUT("/organization-units/:id/users/:userId", middlewares.RolePermissionMiddleware(db, "organization_unit", "modify"), organizationUnitsController.UpdatePrimaryOrganizationUnitHead)
	organizationRoutes.GET("/organization-units/:id/summary", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationUnitsController.GetOrganizationUnitSummary)
	organizationRoutes.GET("/organization-units/:id/policies", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationUnitsController.GetAllPoliciesByOrganizationUnit)
	organizationRoutes.GET("/organization-units/:id/processes", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationUnitsController.GetAllProcessesByOrganizationUnit)
	organizationRoutes.GET("/organization-units/:id/risks", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationUnitsController.GetAllRiskByOrganizationUnit)
	organizationRoutes.GET("/organization-units/:id/details", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationUnitsController.GetParentOrganizationUnitAndCurrentUnitHead)

	// Organization unit reports
	organizationRoutes.GET("/organization-units/reports", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationUnitsController.GetAllOrganizationUnitReports)
	organizationRoutes.GET("/organization-units/reports/:id", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationUnitsController.GetOrganizationUnitReportDetails)
	organizationRoutes.POST("/organization-units/reports", middlewares.RolePermissionMiddleware(db, "organization_unit", "modify"), organizationUnitsController.CreateOrganizationUnitReport)
	organizationRoutes.DELETE("/organization-units/reports/:id", middlewares.RolePermissionMiddleware(db, "organization_unit", "modify"), organizationUnitsController.DeleteOrganizationUnitReport)
	organizationRoutes.GET("/organization-units/reports/:id/export", organizationUnitsController.ExportOrganizationUnitReportToPDF)

	// Organization Unit Business Context
	organizationRoutes.GET("/organization-units/:id/business-context", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationUnitBusinessContextController.GetAllBusinessContexts)
	organizationRoutes.GET("/organization-units/:id/business-context/details", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationUnitBusinessContextController.GetBusinessContextDetails)
	organizationRoutes.POST("/organization-units/:id/business-context", middlewares.RolePermissionMiddleware(db, "organization_unit", "modify"), organizationUnitBusinessContextController.CreateBusinessContext)
	organizationRoutes.PUT("/organization-units/:id/business-context", middlewares.RolePermissionMiddleware(db, "organization_unit", "modify"), organizationUnitBusinessContextController.UpdateBusinessContext)
	organizationRoutes.DELETE("/organization-units/:id/business-context/delete", middlewares.RolePermissionMiddleware(db, "organization_unit", "modify"), organizationUnitBusinessContextController.DeleteBusinessContext)

	// Organization Unit Attributes
	organizationRoutes.GET("/organization-units/:id/attributes", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationUnitAttributeController.GetAllAttributes)
	organizationRoutes.GET("/organization-units/:id/attributes/:attributeId", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationUnitAttributeController.GetAttributeDetails)
	organizationRoutes.POST("/organization-units/:id/attributes", middlewares.RolePermissionMiddleware(db, "organization_unit", "modify"), organizationUnitAttributeController.CreateAttribute)
	organizationRoutes.PUT("/organization-units/:id/attributes/:attributeId", middlewares.RolePermissionMiddleware(db, "organization_unit", "modify"), organizationUnitAttributeController.UpdateAttribute)
	organizationRoutes.POST("/organization-units/:id/attributes/delete", middlewares.RolePermissionMiddleware(db, "organization_unit", "modify"), organizationUnitAttributeController.DeleteAttributes)

	organizationRoutes.GET("/organization-chart/users", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationChartUsers.GetUsers)
	organizationRoutes.GET("/organization-chart/users/:id", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationChartUsers.GetUserFullDetails)
	organizationRoutes.GET("/organization-charts/users", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationChartUsers.GetUsers)
	organizationRoutes.GET("/organization-charts/users/:id", middlewares.RolePermissionMiddleware(db, "organization_unit", "view"), organizationChartUsers.GetUserFullDetails)
	organizationRoutes.PUT("/organization-charts/users/:id/change-reporting-user", middlewares.RolePermissionMiddleware(db, "organization_unit", "modify"), organizationChartUsers.ChangeReportingUser)

	// My Dashboard
	organizationRoutes.GET("/my-dashboard/summary", myDashboardController.GetDashboardSummary)
	organizationRoutes.GET("/my-dashboard/kpis", myDashboardController.GetKPIs)
	organizationRoutes.GET("/my-dashboard/tasks", myDashboardController.GetUserTasks)
	organizationRoutes.GET("/my-dashboard/overdue-tasks", myDashboardController.GetOverdueTasks)
	organizationRoutes.GET("/my-dashboard/assigned-risks", myDashboardController.GetAssignedRisks)
	organizationRoutes.GET("/my-dashboard/recent-activity", myDashboardController.GetRecentActivity)

	// Committee Term
	organizationRoutes.GET("/committee-terms", committeeTermController.GetAllCommitteeTerms)
	organizationRoutes.GET("/committee-terms/:id", committeeTermController.GetCommitteeTermDetails)

	// Committee Member Type
	organizationRoutes.GET("/committee-member-types", committeeMemberTypeController.GetAllCommitteeMemberTypes)
	organizationRoutes.GET("/committee-member-types/:id", committeeMemberTypeController.GetCommitteeMemberTypeDetails)

	// Committee Term Status
	organizationRoutes.GET("/internal-committee-statuses", internalCommitteeStatusController.GetAllInternalCommitteeStatuses)
	organizationRoutes.GET("/internal-committee-statuses/:id", internalCommitteeStatusController.GetInternalCommitteeStatusDetails)

	// Committee Type
	organizationRoutes.GET("/committee-types", committeeTypeController.GetAllCommitteeTypes)
	organizationRoutes.GET("/committee-types/:id", committeeTypeController.GetCommitteeTypeDetails)
	organizationRoutes.POST("/committee-types", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), committeeTypeController.CreateCommitteeType)
	organizationRoutes.PUT("/committee-types/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), committeeTypeController.UpdateCommitteeType)
	organizationRoutes.POST("/committee-types/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), committeeTypeController.DeleteCommitteeTypes)

	// Committee Classification
	organizationRoutes.GET("/committee-classifications", committeeClassificationController.GetAllCommitteeClassifications)
	organizationRoutes.GET("/committee-classifications/:id", committeeClassificationController.GetCommitteeClassificationDetails)
	organizationRoutes.POST("/committee-classifications", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), committeeClassificationController.CreateCommitteeClassification)
	organizationRoutes.PUT("/committee-classifications/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), committeeClassificationController.UpdateCommitteeClassification)
	organizationRoutes.POST("/committee-classifications/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), committeeClassificationController.DeleteCommitteeClassifications)
	organizationRoutes.PUT("/committee-classifications/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), committeeClassificationController.ActivateCommitteeClassifications)
	organizationRoutes.PUT("/committee-classifications/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), committeeClassificationController.DeactivateCommitteeClassifications)

	//Internal Committee
	organizationRoutes.GET("/internal-committees", internalCommitteeController.GetAllInternalCommittees)
	organizationRoutes.GET("/internal-committees/export", internalCommitteeController.GetInternalCommitteesExport)
	organizationRoutes.GET("/internal-committees/:committeeId", internalCommitteeController.GetInternalCommitteeDetails)
	organizationRoutes.POST("/internal-committees", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeController.CreateInternalCommittee)
	organizationRoutes.PUT("/internal-committees/:committeeId", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeController.UpdateInternalCommittee)
	organizationRoutes.POST("/internal-committees/reform", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeController.ReformInternalCommittee)
	organizationRoutes.POST("/internal-committees/delete", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeController.DeleteInternalCommittee)
	organizationRoutes.POST("/internal-committees/:committeeId/disband", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeController.DisbandInternalCommittee)
	organizationRoutes.PUT("/internal-committees/:committeeId/details", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeController.UpdateCommitteeDetails)
	organizationRoutes.PUT("/internal-committees/:committeeId/cancel", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeController.CancelInternalCommittee)
	organizationRoutes.PUT("/internal-committees/:committeeId/complete", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeController.CompleteInternalCommittee)

	organizationRoutes.GET("/internal-committees/:committeeId/committee-formation/export-pdf", internalCommitteeController.ExportCommitteeFormationToPDF)

	// Internal Committee Participants
	organizationRoutes.POST("/internal-committee-nominee-departments/:nomineeDepartmentId/participants", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeController.CreateInternalCommitteeParticipants)
	organizationRoutes.DELETE("/internal-committee-nominee-departments/:nomineeDepartmentId/participants/:participantId", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeController.DeleteInternalCommitteeParticipant)
	organizationRoutes.POST("/internal-committee-nominee-departments/:nomineeDepartmentId/external-participants", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeController.CreateExternalCommitteeParticipant)

	//Internal Committee Workflow Review
	organizationRoutes.POST("/internal-committees/:committeeId/submit", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeWorkflowReviewController.SubmitInternalCommittee)
	organizationRoutes.POST("/internal-committees/:committeeId/approve", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeWorkflowReviewController.ApproveInternalCommittee)
	organizationRoutes.POST("/internal-committees/:committeeId/revert", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeWorkflowReviewController.RevertInternalCommittee)
	organizationRoutes.GET("/internal-committees/:committeeId/workflow-list", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeWorkflowController.GetWorkflowList)
	organizationRoutes.GET("/internal-committees/:committeeId/previous-version-history", middlewares.RolePermissionMiddleware(db, "committee_management", "view"), internalCommitteeWorkflowController.GetPreviousVersionWorkflowHistory)

	// Internal Committee Review Comments
	organizationRoutes.GET("/internal-committees/:committeeId/review-comments", middlewares.RolePermissionMiddleware(db, "committee_management", "view"), internalCommitteeReviewCommentController.GetAllComments)
	organizationRoutes.GET("/internal-committees/:committeeId/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "committee_management", "view"), internalCommitteeReviewCommentController.GetCommentDetails)
	organizationRoutes.POST("/internal-committees/:committeeId/review-comments", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeReviewCommentController.CreateComment)
	organizationRoutes.PUT("/internal-committees/:committeeId/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeReviewCommentController.UpdateComment)
	organizationRoutes.POST("/internal-committees/:committeeId/review-comments/delete", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeReviewCommentController.DeleteComment)
	organizationRoutes.PUT("/internal-committees/:committeeId/review-comments/:commentId/resolve", middlewares.RolePermissionMiddleware(db, "committee_management", "modify"), internalCommitteeReviewCommentController.ResolveComment)

	unauthenticatedRoutes.GET("/internal-committees/:committeeId/review-comments/:commentId/document-preview", internalCommitteeReviewCommentController.GetDocumentPreview)
	unauthenticatedRoutes.GET("/internal-committees/:committeeId/review-comments/:commentId/document-download", internalCommitteeReviewCommentController.DownloadDocument)

	// Control Effectiveness Rating Master
	organizationRoutes.GET("/control-effectiveness-ratings", controlEffectivenessRatingController.GetAllControlEffectivenessRatings)
	organizationRoutes.GET("/control-effectiveness-ratings/:id", controlEffectivenessRatingController.GetControlEffectivenessRatingDetails)

	organizationRoutes.GET("/governance-settings", governanceSettingsController.GetAllGovernanceSettings)
	organizationRoutes.PUT("/governance-settings", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), governanceSettingsController.UpdateGovernanceSettings)
	// Control Status Master
	organizationRoutes.GET("/control-statuses", controlStatusController.GetAllControlStatuses)
	organizationRoutes.GET("/control-statuses/:id", controlStatusController.GetControlStatusDetails)

	// Process Dependency Tiers
	organizationRoutes.GET("/process-dependency-tiers", processDependencyTierController.GetAllProcessDependencyTiers)
	organizationRoutes.GET("/process-dependency-tiers/:id", processDependencyTierController.GetProcessDependencyTierDetails)
	organizationRoutes.POST("/process-dependency-tiers", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processDependencyTierController.CreateProcessDependencyTier)
	organizationRoutes.PUT("/process-dependency-tiers/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processDependencyTierController.UpdateProcessDependencyTier)
	organizationRoutes.POST("/process-dependency-tiers/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processDependencyTierController.DeleteProcessDependencyTiers)

	// Process Internal Dependencies
	organizationRoutes.GET("/processes/:processId/internal-dependencies", middlewares.RolePermissionMiddleware(db, "process", "view"), processInternalDependencyController.GetAllProcessInternalDependencies)
	organizationRoutes.GET("/processes/:processId/internal-dependencies/:internalDependencyId", middlewares.RolePermissionMiddleware(db, "process", "view"), processInternalDependencyController.GetProcessInternalDependencyDetails)
	organizationRoutes.POST("/processes/:processId/internal-dependencies", middlewares.RolePermissionMiddleware(db, "process", "modify"), processInternalDependencyController.CreateProcessInternalDependency)
	organizationRoutes.PUT("/processes/:processId/internal-dependencies/:internalDependencyId", middlewares.RolePermissionMiddleware(db, "process", "modify"), processInternalDependencyController.UpdateProcessInternalDependency)
	organizationRoutes.POST("/processes/:processId/internal-dependencies/delete", middlewares.RolePermissionMiddleware(db, "process", "modify"), processInternalDependencyController.DeleteProcessInternalDependencies)

	// Process External Dependencies
	organizationRoutes.GET("/processes/:processId/external-dependencies", middlewares.RolePermissionMiddleware(db, "process", "view"), processExternalDependencyController.GetAllProcessExternalDependencies)
	organizationRoutes.GET("/processes/:processId/external-dependencies/:externalDependencyId", middlewares.RolePermissionMiddleware(db, "process", "view"), processExternalDependencyController.GetProcessExternalDependencyDetails)
	organizationRoutes.POST("/processes/:processId/external-dependencies", middlewares.RolePermissionMiddleware(db, "process", "modify"), processExternalDependencyController.CreateProcessExternalDependency)
	organizationRoutes.PUT("/processes/:processId/external-dependencies/:externalDependencyId", middlewares.RolePermissionMiddleware(db, "process", "modify"), processExternalDependencyController.UpdateProcessExternalDependency)
	organizationRoutes.POST("/processes/:processId/external-dependencies/delete", middlewares.RolePermissionMiddleware(db, "process", "modify"), processExternalDependencyController.DeleteProcessExternalDependencies)

	// Process External Dependency Entity Types
	organizationRoutes.GET("/process-external-dependency-entity-types", processExternalDependencyEntityTypeController.GetAllProcessExternalDependencyEntityTypes)
	organizationRoutes.GET("/process-external-dependency-entity-types/:id", processExternalDependencyEntityTypeController.GetProcessExternalDependencyEntityTypeDetails)
	organizationRoutes.POST("/process-external-dependency-entity-types", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processExternalDependencyEntityTypeController.CreateProcessExternalDependencyEntityType)
	organizationRoutes.PUT("/process-external-dependency-entity-types/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processExternalDependencyEntityTypeController.UpdateProcessExternalDependencyEntityType)
	organizationRoutes.POST("/process-external-dependency-entity-types/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processExternalDependencyEntityTypeController.DeleteProcessExternalDependencyEntityTypes)

	// Process Employee Analysis
	organizationRoutes.GET("/processes/:processId/employee-analyses", middlewares.RolePermissionMiddleware(db, "process", "view"), processEmployeeAnalysisController.GetAllProcessEmployeeAnalyses)
	organizationRoutes.GET("/processes/:processId/employee-analyses/:id", middlewares.RolePermissionMiddleware(db, "process", "view"), processEmployeeAnalysisController.GetProcessEmployeeAnalysisDetails)
	organizationRoutes.POST("/processes/:processId/employee-analyses", middlewares.RolePermissionMiddleware(db, "process", "modify"), processEmployeeAnalysisController.CreateProcessEmployeeAnalysis)
	organizationRoutes.PUT("/processes/:processId/employee-analyses/:id", middlewares.RolePermissionMiddleware(db, "process", "modify"), processEmployeeAnalysisController.UpdateProcessEmployeeAnalysis)
	organizationRoutes.POST("/processes/:processId/employee-analyses/delete", middlewares.RolePermissionMiddleware(db, "process", "modify"), processEmployeeAnalysisController.DeleteProcessEmployeeAnalyses)

	// Process Equipment Tools
	organizationRoutes.GET("/process-equipment-tools", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), processEquipmentToolController.GetAllProcessEquipmentTools)
	organizationRoutes.GET("/process-equipment-tools/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), processEquipmentToolController.GetProcessEquipmentToolDetails)
	organizationRoutes.POST("/process-equipment-tools", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processEquipmentToolController.CreateProcessEquipmentTool)
	organizationRoutes.PUT("/process-equipment-tools/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processEquipmentToolController.UpdateProcessEquipmentTool)
	organizationRoutes.POST("/process-equipment-tools/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), processEquipmentToolController.DeleteProcessEquipmentTools)
}
