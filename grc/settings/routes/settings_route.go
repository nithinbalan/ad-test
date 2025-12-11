package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	settingRoutes := router.Group("/")

	settingRoutes.Use(middlewares.AuthMiddleware())
	unauthenticatedRoutes := router.Group("/")
	db := app.DB

	organizationLevelController := app.SettingsModule.OrganizationLevelController
	systemPreferenceController := app.SettingsModule.SystemPreferenceController
	roleController := app.SettingsModule.RoleController
	organizationUnitController := app.SettingsModule.OrganizationUnitController
	logSettingsController := app.SettingsModule.LogSettingsController
	workflowLevelTypeController := app.SettingsModule.WorkflowLevelTypeController
	workflowController := app.SettingsModule.WorkflowController
	rolePermissionTypeController := app.SettingsModule.RolePermissionTypeController
	moduleController := app.SettingsModule.ModuleController
	emailTemplateSettingsController := app.SettingsModule.EmailTemplateSettingsController
	integrationController := app.SettingsModule.IntegrationController
	auditTemplateSettingController := app.SettingsModule.AuditTemplateSettingController
	riskTemplateSettingController := app.SettingsModule.RiskTemplateSettingController
	reportTemplateSettingController := app.SettingsModule.ReportTemplateSettingController
	activeDirectorySettingController := app.SettingsModule.ActiveDirectorySettingController
	siemSettingController := app.SettingsModule.SIEMSettingController
	aiSettingController := app.SettingsModule.AISettingController
	triggerConditionController := app.SettingsModule.TriggerConditionController
	reminderController := app.SettingsModule.ReminderController
	reminderRuleController := app.SettingsModule.ReminderRuleController
	escalationController := app.SettingsModule.EscalationController
	escalationLevelController := app.SettingsModule.EscalationLevelController
	reminderContextController := app.SettingsModule.ReminderContextController
	ragSettingController := app.SettingsModule.RAGSettingController
	tenableSettingsController := app.SettingsModule.TenableSettingsController
	healthCheckController := app.SettingsModule.HealthCheckController

	// Organization Level
	settingRoutes.GET("/organization-levels", organizationLevelController.GetAllOrganizationLevels)
	settingRoutes.POST("/organization-levels", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), organizationLevelController.CreateOrganizationLevels)
	settingRoutes.POST("/organization-level", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), organizationLevelController.CreateOrganizationLevel)
	settingRoutes.PUT("/organization-levels/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), organizationLevelController.UpdateOrganizationLevel)
	settingRoutes.DELETE("/organization-levels/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), organizationLevelController.DeleteOrganizationLevel)

	// System Preference
	settingRoutes.GET("/system-preferences", systemPreferenceController.GetSystemPreference)
	settingRoutes.PUT("/system-preferences/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), systemPreferenceController.UpdateSystemPreference)
	unauthenticatedRoutes.PUT("/update-license-key", systemPreferenceController.UpdateLicenseKey)
	unauthenticatedRoutes.GET("/license-key", systemPreferenceController.GetSystemPreference)

	// Email template settings
	settingRoutes.GET("/email-template-settings", emailTemplateSettingsController.GetEmailTemplateSettings)
	settingRoutes.PUT("/email-template-settings/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), emailTemplateSettingsController.UpdateEmailTemplateSettings)

	settingRoutes.GET("/role-permission-types", rolePermissionTypeController.GetAllRolePermissionTypes)
	settingRoutes.GET("/role-permission-types/:id", rolePermissionTypeController.GetRolePermissionTypeDetails)

	// Role
	settingRoutes.GET("/roles", roleController.GetAllRoles)
	settingRoutes.GET("/roles/:id", roleController.GetRoleDetails)
	settingRoutes.POST("/roles", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), roleController.CreateRole)
	settingRoutes.PUT("/roles/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), roleController.UpdateRole)
	settingRoutes.POST("/roles/delete", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), roleController.DeleteRoles)
	settingRoutes.PUT("/roles/activate", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), roleController.ActivateRoles)
	settingRoutes.PUT("/roles/deactivate", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), roleController.DeactivateRoles)

	settingRoutes.GET("/roles/:id/permissions", roleController.GetRolePermissions)
	settingRoutes.PUT("/roles/:id/permissions", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), roleController.UpdateRolePermissions)

	// Organization Unit
	settingRoutes.GET("/organization-units", organizationUnitController.GetAllOrganizationUnits)
	settingRoutes.GET("/organization-units/:id", organizationUnitController.GetOrganizationUnitDetails)
	settingRoutes.POST("/organization-units", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), organizationUnitController.CreateOrganizationUnit)
	settingRoutes.PUT("/organization-units/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), organizationUnitController.UpdateOrganizationUnit)
	settingRoutes.POST("/organization-units/delete", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), organizationUnitController.DeleteOrganizationUnits)
	settingRoutes.PUT("/organization-units/activate", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), organizationUnitController.ActivateOrganizationUnits)
	settingRoutes.PUT("/organization-units/deactivate", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), organizationUnitController.DeactivateOrganizationUnits)

	// Log Settings
	settingRoutes.GET("/log-settings", logSettingsController.GetLogSettings)
	settingRoutes.PUT("/log-settings", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), logSettingsController.UpdateLogSettings)

	settingRoutes.GET("/workflow-level-types", workflowLevelTypeController.GetAllWorkflowLevelTypes)
	settingRoutes.GET("/workflow-level-types/:id", workflowLevelTypeController.GetWorkflowLevelTypeDetails)

	settingRoutes.GET("/workflow-lists", workflowController.GetAllWorkflows)
	settingRoutes.GET("/workflows/:workflowId", workflowController.GetWorkflowDetails)
	settingRoutes.POST("/workflows", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), workflowController.CreateWorkflow)
	settingRoutes.PUT("/workflows/:workflowId", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), workflowController.UpdateWorkflow)
	settingRoutes.POST("/workflows/delete", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), workflowController.DeleteWorkflows)

	settingRoutes.GET("/workflows", workflowController.GetAllWorkflowsByType)
	settingRoutes.GET("/workflows/:workflowId/levels", workflowController.GetWorkflowlevels)
	settingRoutes.POST("/workflows/:workflowId/levels", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), workflowController.CreateWorkflowLevel)
	settingRoutes.DELETE("/workflows/:workflowId/levels/:levelId", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), workflowController.DeleteWorkflowLevel)

	settingRoutes.GET("/modules", moduleController.GetAllModules)
	settingRoutes.GET("/modules/:moduleId", moduleController.GetModuleDetail)
	settingRoutes.POST("/modules/:moduleId/enable", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), moduleController.EnableModule)
	settingRoutes.POST("/modules/:moduleId/disable", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), moduleController.DisableModule)
	settingRoutes.PUT("/modules/reorder", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), moduleController.ReorderModule)
	settingRoutes.PUT("/modules/:moduleId/reorder", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), moduleController.ReorderSubModule)
	settingRoutes.PUT("/modules/:moduleId", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), moduleController.UpdateModule)
	settingRoutes.PUT("/modules/:moduleId/submodule/:subModuleId", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), moduleController.UpdateSubModule)

	settingRoutes.GET("/integrations", integrationController.GetAllIntegrations)
	settingRoutes.GET("/integrations/:integrationId/credential", integrationController.GetIntegrationDetails)
	settingRoutes.POST("/integrations/:integrationId/credential", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), integrationController.EnableIntegration)
	settingRoutes.PUT("/integrations/:integrationId/credential", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), integrationController.UpdateCredential)
	settingRoutes.DELETE("/integrations/:integrationId/credential", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), integrationController.DisableIntegration)

	settingRoutes.GET("/organization-preferences", systemPreferenceController.GetOrganizationPreference)

	settingRoutes.GET("/audit-template-settings", auditTemplateSettingController.GetAllAuditTemplateSettings)
	settingRoutes.PUT("/audit-template-settings/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), auditTemplateSettingController.UpdateAuditTemplateSetting)
	unauthenticatedRoutes.GET("/audit-template-settings/:id/image-preview", auditTemplateSettingController.GetAuditTemplateSettingPreview)
	unauthenticatedRoutes.GET("/audit-template-settings/image-preview", auditTemplateSettingController.GetAuditTemplateSettingPreview)

	settingRoutes.GET("/risk-template-settings", riskTemplateSettingController.GetAllRiskTemplateSettings)
	settingRoutes.PUT("/risk-template-settings/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), riskTemplateSettingController.UpdateRiskTemplateSetting)
	unauthenticatedRoutes.GET("/risk-template-settings/image-preview", riskTemplateSettingController.GetRiskTemplateSettingPreview)

	settingRoutes.GET("/report-template-settings", reportTemplateSettingController.GetAllReportTemplateSettings)
	settingRoutes.PUT("/report-template-settings/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), reportTemplateSettingController.UpdateReportTemplateSetting)
	unauthenticatedRoutes.GET("/report-template-settings/image-preview", reportTemplateSettingController.GetReportTemplateSettingPreview)

	settingRoutes.GET("/active-directory-settings", activeDirectorySettingController.GetActiveDirectorySettings)
	settingRoutes.POST("/active-directory-settings", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), activeDirectorySettingController.CreateActiveDirectorySetting)
	settingRoutes.POST("/active-directory-settings/check-connection", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), activeDirectorySettingController.CheckActiveDirectoryConnection)
	settingRoutes.GET("/active-directory-settings/check-saved-connection", activeDirectorySettingController.CheckSavedADConnection)
	settingRoutes.DELETE(
		"/active-directory-settings/delete",
		middlewares.RolePermissionMiddleware(db, "general_setting", "modify"),
		activeDirectorySettingController.DeleteActiveDirectorySetting,
	)

	settingRoutes.GET("/siem-settings", siemSettingController.GetSIEMSettings)
	settingRoutes.PUT("/siem-settings", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), siemSettingController.UpdateSIEMSettings)

	settingRoutes.PUT("/ai-settings/enable", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), aiSettingController.EnableAISetting)
	settingRoutes.PUT("/ai-settings/disable", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), aiSettingController.DisableAISetting)
	settingRoutes.GET("/ai-settings/modules", aiSettingController.GetAllAIModules)
	settingRoutes.GET("/ai-settings/modules/:id", aiSettingController.GetFeaturesByModuleID)

	// Trigger Condition
	settingRoutes.GET("/trigger-conditions", triggerConditionController.GetAllTriggerConditiones)
	settingRoutes.GET("/trigger-conditions/:id", triggerConditionController.GetTriggerConditionDetails)

	// Reminder
	settingRoutes.GET("/reminders", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), reminderController.GetAllReminderes)
	settingRoutes.GET("/reminders/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), reminderController.GetReminderDetails)
	settingRoutes.GET("/reminder-modules", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), reminderController.GetAllEnabledModules)
	settingRoutes.GET("/modules/:moduleId/reminders", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), reminderController.GetReminderByModuleID)

	// Reminder Context
	settingRoutes.GET("/reminder-contexts", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), reminderContextController.GetAllReminderContexts)

	// Reminder Rule
	settingRoutes.GET("/reminder-rules", reminderRuleController.GetAllReminderRules)
	settingRoutes.GET("/reminder-rules/:reminderRuleId", reminderRuleController.GetReminderRuleDetails)
	settingRoutes.POST("/reminder-rules", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), reminderRuleController.CreateReminderRule)
	settingRoutes.PUT("/reminder-rules/:reminderRuleId", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), reminderRuleController.UpdateReminderRule)
	settingRoutes.POST("/reminder-rules/delete", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), reminderRuleController.DeleteReminderRules)

	// Escalations
	settingRoutes.GET("/escalations", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), escalationController.GetAllEscalations)
	settingRoutes.GET("/escalations/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), escalationController.GetEscalationDetails)
	settingRoutes.GET("/escalation-modules", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), escalationController.GetAllEnabledModules)
	settingRoutes.GET("/modules/:moduleId/escalations", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), escalationController.GetEscalationsByModuleID)

	// Escalation Level
	settingRoutes.GET("/escalation-levels", escalationLevelController.GetAllEscalationLevels)
	settingRoutes.GET("/escalation-levels/:escalationLevelId", escalationLevelController.GetEscalationLevelDetails)
	settingRoutes.POST("/escalation-levels", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), escalationLevelController.CreateEscalationLevel)
	settingRoutes.PUT("/escalation-levels/:escalationLevelId", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), escalationLevelController.UpdateEscalationLevel)
	settingRoutes.POST("/escalation-levels/delete", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), escalationLevelController.DeleteEscalationLevels)

	settingRoutes.GET("/rag-setting", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), ragSettingController.GetAllRAGSetting)

	// --- Microsoft Entra ID (SSO + Directory Sync) ---
	settingRoutes.POST("/active-directory-settings/test-oauth-connection",
		middlewares.RolePermissionMiddleware(db, "general_setting", "modify"),
		activeDirectorySettingController.TestOAuthConnection)

	settingRoutes.GET("/active-directory-settings/oauth-config",
		activeDirectorySettingController.GetOAuthConfiguration)

	settingRoutes.POST("/active-directory-settings/oauth-config",
		middlewares.RolePermissionMiddleware(db, "general_setting", "modify"),
		activeDirectorySettingController.SaveOAuthConfiguration)

	settingRoutes.DELETE("/active-directory-settings/oauth-config",
		middlewares.RolePermissionMiddleware(db, "general_setting", "modify"),
		activeDirectorySettingController.DeleteOAuthConfiguration)

	settingRoutes.GET("/active-directory-settings/users/:id/groups",
		middlewares.RolePermissionMiddleware(db, "general_setting", "modify"),
		activeDirectorySettingController.GetAzureUserGroups)

	//tenable settings
	settingRoutes.GET("/tenable-settings", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), tenableSettingsController.GetTenableSettings)
	settingRoutes.POST("/tenable-settings", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), tenableSettingsController.SaveTenableSettings)
	settingRoutes.DELETE("/tenable-settings", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), tenableSettingsController.DeleteTenableSettings)
	settingRoutes.POST("/tenable-settings/test-connection", middlewares.RolePermissionMiddleware(db, "general_setting", "view"), tenableSettingsController.TestConnection)

	unauthenticatedRoutes.GET("/health", healthCheckController.Health)
	unauthenticatedRoutes.GET("/health/live", healthCheckController.Liveness)
	unauthenticatedRoutes.GET("/health/ready", healthCheckController.Readiness)
}
