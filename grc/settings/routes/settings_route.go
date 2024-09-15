package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	settingRoutes := router.Group("/")

	settingRoutes.Use(middlewares.AuthMiddleware())

	organizationLevelController := app.SettingsModule.OrganizationLevelController
	systemPreferenceController := app.SettingsModule.SystemPreferenceController
	roleController := app.SettingsModule.RoleController
	organizationUnitController := app.SettingsModule.OrganizationUnitController
	logSettingsController := app.SettingsModule.LogSettingsController
	workflowLevelTypeController := app.SettingsModule.WorkflowLevelTypeController
	workflowController := app.SettingsModule.WorkflowController
	rolePermissionTypeController := app.SettingsModule.RolePermissionTypeController
	moduleController := app.SettingsModule.ModuleController

	//Organization Level
	settingRoutes.GET("/organization-levels", organizationLevelController.GetAllOrganizationLevels)
	settingRoutes.POST("/organization-levels", organizationLevelController.CreateOrganizationLevel)

	//System Preference
	settingRoutes.GET("/system-preferences", systemPreferenceController.GetSystemPreference)
	settingRoutes.PUT("/system-preferences/:id", systemPreferenceController.UpdateSystemPreference)

	settingRoutes.GET("/role-permission-types", rolePermissionTypeController.GetAllRolePermissionTypes)
	settingRoutes.GET("/role-permission-types/:id", rolePermissionTypeController.GetRolePermissionTypeDetails)

	//Role
	settingRoutes.GET("/roles", roleController.GetAllRoles)
	settingRoutes.GET("/roles/:id", roleController.GetRoleDetails)
	settingRoutes.POST("/roles", roleController.CreateRole)
	settingRoutes.PUT("/roles/:id", roleController.UpdateRole)
	settingRoutes.POST("/roles/delete", roleController.DeleteRoles)
	settingRoutes.PUT("/roles/activate", roleController.ActivateRoles)
	settingRoutes.PUT("/roles/deactivate", roleController.DeactivateRoles)

	// settingRoutes.GET("/roles/:id/permissions", roleController.GetRolePermissions)
	settingRoutes.PUT("/roles/:id/permissions", roleController.UpdateRolePermissions)

	// Organization Unit
	settingRoutes.GET("/organization-units", organizationUnitController.GetAllOrganizationUnits)
	settingRoutes.GET("/organization-units/:id", organizationUnitController.GetOrganizationUnitDetails)
	settingRoutes.POST("/organization-units", organizationUnitController.CreateOrganizationUnit)
	settingRoutes.PUT("/organization-units/:id", organizationUnitController.UpdateOrganizationUnit)
	settingRoutes.POST("/organization-units/delete", organizationUnitController.DeleteOrganizationUnits)
	settingRoutes.PUT("/organization-units/activate", organizationUnitController.ActivateOrganizationUnits)
	settingRoutes.PUT("/organization-units/deactivate", organizationUnitController.DeactivateOrganizationUnits)

	// Log Settings
	settingRoutes.GET("/log-settings", logSettingsController.GetLogSettings)
	settingRoutes.PUT("/log-settings", logSettingsController.UpdateLogSettings)

	settingRoutes.GET("/workflow-level-types", workflowLevelTypeController.GetAllWorkflowLevelTypes)
	settingRoutes.GET("/workflow-level-types/:id", workflowLevelTypeController.GetWorkflowLevelTypeDetails)

	settingRoutes.GET("/workflows", workflowController.GetAllWorkflows)
	settingRoutes.GET("/workflows/:workflowId/levels", workflowController.GetWorkflowlevels)
	settingRoutes.POST("/workflows/:workflowId/levels", workflowController.CreateWorkflowLevel)
	settingRoutes.DELETE("/workflows/:workflowId/levels/:levelId", workflowController.DeleteWorkflowLevel)

	settingRoutes.GET("/modules", moduleController.GetAllModules)

}
