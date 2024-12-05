package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	logRoutes := router.Group("/")

	logRoutes.Use(middlewares.AuthMiddleware())

	activityLogController := app.SharedModule.ActivityLogController

	//Activity Log
	logRoutes.GET("/activity-logs", activityLogController.GetAllActivityLogs)
	logRoutes.GET("/activity-logs/:logId", activityLogController.GetActivityLogDetails)

}
