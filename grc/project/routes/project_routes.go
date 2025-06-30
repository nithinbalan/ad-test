package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	projectRoutes := router.Group("/")

	projectRoutes.Use(middlewares.AuthMiddleware())

	projectTypeController := app.ProjectModule.ProjectTypeController

	// Project Type
	projectRoutes.GET("/project-types", projectTypeController.GetAllProjectTypes)
	projectRoutes.GET("/project-types/:id", projectTypeController.GetProjectTypeDetails)
	projectRoutes.POST("/project-types", projectTypeController.CreateProjectType)
	projectRoutes.PUT("/project-types/:id", projectTypeController.UpdateProjectType)
	projectRoutes.POST("/project-types/delete", projectTypeController.DeleteProjectTypes)
	projectRoutes.PUT("/project-types/activate", projectTypeController.ActivateProjectTypes)
	projectRoutes.PUT("/project-types/deactivate", projectTypeController.DeactivateProjectTypes)

}
