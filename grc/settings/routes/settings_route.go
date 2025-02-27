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

	//Organization Level
	settingRoutes.GET("/organization-levels", organizationLevelController.GetAllOrganizationLevels)
	settingRoutes.POST("/organization-levels", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), organizationLevelController.CreateOrganizationLevels)
	settingRoutes.POST("/organization-level", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), organizationLevelController.CreateOrganizationLevel)
	settingRoutes.PUT("/organization-levels/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), organizationLevelController.UpdateOrganizationLevel)
	settingRoutes.DELETE("/organization-levels/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), organizationLevelController.DeleteOrganizationLevel)

	//System Preference
	settingRoutes.GET("/system-preferences", systemPreferenceController.GetSystemPreference)
	settingRoutes.PUT("/system-preferences/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), systemPreferenceController.UpdateSystemPreference)
	unauthenticatedRoutes.PUT("/update-license-key", systemPreferenceController.UpdateLicenseKey)
	unauthenticatedRoutes.GET("/license-key", systemPreferenceController.GetSystemPreference)

	//Email template settings
	settingRoutes.GET("/email-template-settings", emailTemplateSettingsController.GetEmailTemplateSettings)
	settingRoutes.PUT("/email-template-settings/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), emailTemplateSettingsController.UpdateEmailTemplateSettings)

	settingRoutes.GET("/role-permission-types", rolePermissionTypeController.GetAllRolePermissionTypes)
	settingRoutes.GET("/role-permission-types/:id", rolePermissionTypeController.GetRolePermissionTypeDetails)

	//Role
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

	settingRoutes.GET("/workflows", workflowController.GetAllWorkflows)
	settingRoutes.GET("/workflows/:workflowId/levels", workflowController.GetWorkflowlevels)
	settingRoutes.POST("/workflows/:workflowId/levels", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), workflowController.CreateWorkflowLevel)
	settingRoutes.DELETE("/workflows/:workflowId/levels/:levelId", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), workflowController.DeleteWorkflowLevel)

	settingRoutes.GET("/modules", moduleController.GetAllModules)
	settingRoutes.GET("/modules/:moduleId", moduleController.GetModuleDetail)
	settingRoutes.POST("/modules/:moduleId/enable", moduleController.EnableModule)
	settingRoutes.POST("/modules/:moduleId/disable", moduleController.DisableModule)
	settingRoutes.PUT("/modules/reorder", moduleController.ReorderModule)
	settingRoutes.PUT("/modules/:moduleId/reorder", moduleController.ReorderSubModule)
	settingRoutes.PUT("/modules/:moduleId", moduleController.UpdateModule)
	settingRoutes.PUT("/modules/:moduleId/submodule/:subModuleId", moduleController.UpdateSubModule)

	settingRoutes.GET("/integrations", integrationController.GetAllIntegrations)
	settingRoutes.GET("/integrations/:integrationId/credential", integrationController.GetIntegrationDetails)
	settingRoutes.POST("/integrations/:integrationId/credential", integrationController.EnableIntegration)
	settingRoutes.PUT("/integrations/:integrationId/credential", integrationController.UpdateCredential)
	settingRoutes.DELETE("/integrations/:integrationId/credential", integrationController.DisableIntegration)

	settingRoutes.GET("/organization-preferences", systemPreferenceController.GetOrganizationPreference)

	settingRoutes.GET("/audit-template-settings", auditTemplateSettingController.GetAllAuditTemplateSettings)
	settingRoutes.PUT("/audit-template-settings/:id", middlewares.RolePermissionMiddleware(db, "general_setting", "modify"), auditTemplateSettingController.UpdateAuditTemplateSetting)
	unauthenticatedRoutes.GET("/audit-template-settings/:id/image-preview", auditTemplateSettingController.GetAuditTemplateSettingPreview)

}
