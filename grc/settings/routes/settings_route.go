package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	settingRoutes := router.Group("/")

	settingRoutes.Use(middlewares.AuthMiddleware())

	// activityLogRepo := logRepositories.NewActivityLogRepository(db)
	// activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	// organizationLevelRepo := settingRepositories.NewOrganizationLevelRepository(db)
	// organizationLevelSvc := settingServices.NewOrganizationLevelService(organizationLevelRepo, activityLogSvc)
	// OrganizationLevelController := controllers.NewOrganizationLevelController(organizationLevelSvc)

	// currencyRepo := settingRepositories.NewCurrencyRepository(db)
	// currencySvc := settingServices.NewCurrencyService(currencyRepo)
	// currencyController := controllers.NewCurrencyController(currencySvc)

	// systemPreferenceRepo := settingRepositories.NewSystemPreferenceRepository(db)
	// systemPreferenceSvc := settingServices.NewSystemPreferenceService(systemPreferenceRepo, activityLogSvc, currencySvc)
	// systemPreferenceController := controllers.NewSystemPreferenceController(systemPreferenceSvc)

	organizationLevelController := app.SettingsModule.OrganizationLevelController
	currencyController := app.SettingsModule.CurrencyController
	systemPreferenceController := app.SettingsModule.SystemPreferenceController

	//Organization Level
	settingRoutes.GET("/organization-levels", organizationLevelController.GetAllOrganizationLevels)
	settingRoutes.POST("/organization-levels", organizationLevelController.CreateOrganizationLevel)

	//Currency
	settingRoutes.GET("/currencies", currencyController.GetAllCurrencies)
	settingRoutes.GET("/currencies/:id", currencyController.GetCurrencyDetails)

	//System Preference
	settingRoutes.GET("/system-preferences", systemPreferenceController.GetSystemPreference)
	settingRoutes.PUT("/system-preferences/:id", systemPreferenceController.UpdateSystemPreference)
}
