package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	taskRoutes := router.Group("/")

	taskRoutes.Use(middlewares.AuthMiddleware())

	taskStatusController := app.TaskModule.TaskStatusController

	//Task Status
	taskRoutes.GET("/task-statuses", taskStatusController.GetAllTaskStatuses)
	taskRoutes.GET("/task-statuses/:id", taskStatusController.GetTaskStatusDetails)

}
