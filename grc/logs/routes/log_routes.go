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
	emailLogController := app.SharedModule.EmailLogController

	//Activity Log
	logRoutes.GET("/activity-logs", middlewares.RolePermissionMiddleware(db, "log", "view"), activityLogController.GetAllActivityLogs)
	logRoutes.GET("/activity-logs/:logId", middlewares.RolePermissionMiddleware(db, "log", "view"), activityLogController.GetActivityLogDetails)

	logRoutes.GET("/task-schedular-logs", middlewares.RolePermissionMiddleware(db, "log", "view"), taskSchedularLogController.GetAllTaskSchedularLogs)

	logRoutes.GET("/email-logs", middlewares.RolePermissionMiddleware(db, "log", "view"), emailLogController.GetAllEmailLogs)
	logRoutes.GET("/email-logs/summary", middlewares.RolePermissionMiddleware(db, "log", "view"), emailLogController.GetSummary)
	logRoutes.GET("/email-logs/:logId", middlewares.RolePermissionMiddleware(db, "log", "view"), emailLogController.GetEmailLogDetails)

}
