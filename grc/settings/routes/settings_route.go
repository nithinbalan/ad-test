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
	currencyController := app.SettingsModule.CurrencyController
	systemPreferenceController := app.SettingsModule.SystemPreferenceController
	timezoneController := app.SettingsModule.TimezoneController

	//Organization Level
	settingRoutes.GET("/organization-levels", organizationLevelController.GetAllOrganizationLevels)
	settingRoutes.POST("/organization-levels", organizationLevelController.CreateOrganizationLevel)

	//Currency
	settingRoutes.GET("/currencies", currencyController.GetAllCurrencies)
	settingRoutes.GET("/currencies/:id", currencyController.GetCurrencyDetails)

	// Timezones
	settingRoutes.GET("/timezones", timezoneController.GetAllTimezones)
	settingRoutes.GET("/timezones/:id", timezoneController.GetTimezoneDetails)

	//System Preference
	settingRoutes.GET("/system-preferences", systemPreferenceController.GetSystemPreference)
	settingRoutes.PUT("/system-preferences/:id", systemPreferenceController.UpdateSystemPreference)
}
