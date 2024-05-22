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

	//Organization Level
	settingRoutes.GET("/organization-levels", organizationLevelController.GetAllOrganizationLevels)
	settingRoutes.POST("/organization-levels", organizationLevelController.CreateOrganizationLevel)

	//System Preference
	settingRoutes.GET("/system-preferences", systemPreferenceController.GetSystemPreference)
	settingRoutes.PUT("/system-preferences/:id", systemPreferenceController.UpdateSystemPreference)

	//Role
	settingRoutes.GET("/roles", roleController.GetAllRoles)
	settingRoutes.GET("/roles/:id", roleController.GetRoleDetails)
	settingRoutes.POST("/roles", roleController.CreateRole)
	settingRoutes.PUT("/roles/:id", roleController.UpdateRole)
	settingRoutes.POST("/roles/delete", roleController.DeleteRoles)
	settingRoutes.PUT("/roles/activate", roleController.ActivateRoles)
	settingRoutes.PUT("/roles/deactivate", roleController.DeactivateRoles)

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
	settingRoutes.PUT("/log-settings/:id", logSettingsController.UpdateLogSettings)

}
