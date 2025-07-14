package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	projectRoutes := router.Group("/")

	projectRoutes.Use(middlewares.AuthMiddleware())

	db := app.DB

	projectTypeController := app.ProjectModule.ProjectTypeController
	projectStatusController := app.ProjectModule.ProjectStatusController
	projectController := app.ProjectModule.ProjectController
	projectMilestoneController := app.ProjectModule.ProjectMilestoneController
	projectActionPlanController := app.ProjectModule.ProjectActionPlanController

	// Project Type
	projectRoutes.GET("/project-types", projectTypeController.GetAllProjectTypes)
	projectRoutes.GET("/project-types/:id", projectTypeController.GetProjectTypeDetails)
	projectRoutes.POST("/project-types", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectTypeController.CreateProjectType)
	projectRoutes.PUT("/project-types/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectTypeController.UpdateProjectType)
	projectRoutes.POST("/project-types/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectTypeController.DeleteProjectTypes)
	projectRoutes.PUT("/project-types/activate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectTypeController.ActivateProjectTypes)
	projectRoutes.PUT("/project-types/deactivate", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectTypeController.DeactivateProjectTypes)

	// Project Status
	projectRoutes.GET("/project-statuses", projectStatusController.GetAllProjectStatuses)
	projectRoutes.GET("/project-statuses/:id", projectStatusController.GetProjectStatusDetails)

	// Projects
	projectRoutes.GET("/projects", projectController.GetAllProjects)
	projectRoutes.GET("/projects/:projectId", projectController.GetProjectDetails)
	projectRoutes.POST("/projects", projectController.CreateProject)
	projectRoutes.PUT("/projects/:projectId", projectController.UpdateProject)
	projectRoutes.POST("/projects/delete", projectController.DeleteProjects)

	// Project Milestones
	projectRoutes.GET("/project-milestones", projectMilestoneController.GetAllProjectMilestones)
	projectRoutes.GET("/project-milestones/:id", projectMilestoneController.GetProjectMilestoneDetails)
	projectRoutes.POST("/project-milestones", projectMilestoneController.CreateProjectMilestone)
	projectRoutes.PUT("/project-milestones/:id", projectMilestoneController.UpdateProjectMilestone)
	projectRoutes.POST("/project-milestones/delete", projectMilestoneController.DeleteProjectMilestones)

	// Project Action Plans
	projectRoutes.GET("/projects/:projectId/action-plans", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), projectActionPlanController.GetAllProjectActionPlans)
	projectRoutes.GET("/projects/:projectId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), projectActionPlanController.GetProjectActionPlanDetails)
	projectRoutes.POST("/projects/:projectId/action-plans", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), projectActionPlanController.CreateProjectActionPlan)
	projectRoutes.PUT("/projects/:projectId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), projectActionPlanController.UpdateProjectActionPlan)
	projectRoutes.POST("/projects/:projectId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), projectActionPlanController.DeleteProjectActionPlans)

}
