package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	logRoutes := router.Group("/")

	logRoutes.Use(middlewares.AuthMiddleware(app.DB))

	db := app.DB

	activityLogController := app.SharedModule.ActivityLogController
	taskSchedularLogController := app.SharedModule.TaskSchedularLogController

	//Activity Log
	logRoutes.GET("/activity-logs", middlewares.RolePermissionMiddleware(db, "log", "view"), activityLogController.GetAllActivityLogs)
	logRoutes.GET("/activity-logs/:logId", middlewares.RolePermissionMiddleware(db, "log", "view"), activityLogController.GetActivityLogDetails)

	logRoutes.GET("/task-schedular-logs", middlewares.RolePermissionMiddleware(db, "log", "view"), taskSchedularLogController.GetAllTaskSchedularLogs)

}
