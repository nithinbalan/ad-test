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
	projectStatusController := app.ProjectModule.ProjectStatusController
	projectController := app.ProjectModule.ProjectController
	projectMilestoneController := app.ProjectModule.ProjectMilestoneController
	projectActionPlanController := app.ProjectModule.ProjectActionPlanController

	// Project Type
	projectRoutes.GET("/project-types", projectTypeController.GetAllProjectTypes)
	projectRoutes.GET("/project-types/:id", projectTypeController.GetProjectTypeDetails)
	projectRoutes.POST("/project-types", projectTypeController.CreateProjectType)
	projectRoutes.PUT("/project-types/:id", projectTypeController.UpdateProjectType)
	projectRoutes.POST("/project-types/delete", projectTypeController.DeleteProjectTypes)
	projectRoutes.PUT("/project-types/activate", projectTypeController.ActivateProjectTypes)
	projectRoutes.PUT("/project-types/deactivate", projectTypeController.DeactivateProjectTypes)

	// Project Status
	projectRoutes.GET("/project-statuses", projectStatusController.GetAllProjectStatuses)
	projectRoutes.GET("/project-statuses/:id", projectStatusController.GetProjectStatusDetails)

	// Projects
	projectRoutes.GET("/projects", projectController.GetAllProjects)
	projectRoutes.GET("/projects/:id", projectController.GetProjectDetails)
	projectRoutes.POST("/projects", projectController.CreateProject)
	projectRoutes.PUT("/projects/:id", projectController.UpdateProject)
	projectRoutes.POST("/projects/delete", projectController.DeleteProjects)

	// Project Milestones
	projectRoutes.GET("/project-milestones", projectMilestoneController.GetAllProjectMilestones)
	projectRoutes.GET("/project-milestones/:id", projectMilestoneController.GetProjectMilestoneDetails)
	projectRoutes.POST("/project-milestones", projectMilestoneController.CreateProjectMilestone)
	projectRoutes.PUT("/project-milestones/:id", projectMilestoneController.UpdateProjectMilestone)
	projectRoutes.POST("/project-milestones/delete", projectMilestoneController.DeleteProjectMilestones)

	// Project Action Plans
	projectRoutes.GET("/project-action-plans", projectActionPlanController.GetAllProjectActionPlans)
	projectRoutes.GET("/project-action-plans/:id", projectActionPlanController.GetProjectActionPlanDetails)
	projectRoutes.POST("/project-action-plans", projectActionPlanController.CreateProjectActionPlan)
	projectRoutes.PUT("/project-action-plans/:id", projectActionPlanController.UpdateProjectActionPlan)
	projectRoutes.POST("/project-action-plans/delete", projectActionPlanController.DeleteProjectActionPlans)

}
