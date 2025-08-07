package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	projectRoutes := router.Group("/")

	projectRoutes.Use(middlewares.AuthMiddleware())

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	projectTypeController := app.ProjectModule.ProjectTypeController
	projectStatusController := app.ProjectModule.ProjectStatusController
	projectController := app.ProjectModule.ProjectController
	projectMilestoneController := app.ProjectModule.ProjectMilestoneController
	projectActionPlanController := app.ProjectModule.ProjectActionPlanController
	projectActionPlanStatusUpdateController := app.ProjectModule.ProjectActionPlanStatusUpdateController
	projectActionPlanCommentController := app.ProjectModule.ProjectActionPlanCommentController
	projectAssociatedItemController := app.ProjectModule.ProjectAssociatedItemController

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
	projectRoutes.POST("/projects", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectController.CreateProject)
	projectRoutes.PUT("/projects/:projectId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectController.UpdateProject)
	projectRoutes.POST("/projects/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectController.DeleteProjects)

	// Project Milestones
	projectRoutes.GET("/project-milestones", projectMilestoneController.GetAllProjectMilestones)
	projectRoutes.GET("/projects/:projectId/milestones", projectMilestoneController.GetAllProjectMilestones)
	projectRoutes.GET("/projects/:projectId/milestones/:milestoneId", projectMilestoneController.GetProjectMilestoneDetails)
	projectRoutes.POST("/projects/:projectId/milestones", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectMilestoneController.CreateProjectMilestone)
	projectRoutes.PUT("/projects/:projectId/milestones/:milestoneId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectMilestoneController.UpdateProjectMilestone)
	projectRoutes.POST("/projects/:projectId/milestones/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectMilestoneController.DeleteProjectMilestones)

	// Project Action Plans
	projectRoutes.GET("/project-action-plans", projectActionPlanController.GetAllProjectActionPlans)
	projectRoutes.GET("/projects/:projectId/action-plans", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectActionPlanController.GetAllProjectActionPlans)
	projectRoutes.GET("/projects/:projectId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectActionPlanController.GetProjectActionPlanDetails)
	projectRoutes.POST("/projects/:projectId/action-plans", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectActionPlanController.CreateProjectActionPlan)
	projectRoutes.PUT("/projects/:projectId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectActionPlanController.UpdateProjectActionPlan)
	projectRoutes.POST("/projects/:projectId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectActionPlanController.DeleteProjectActionPlans)

	// Project Action Plan	 Status Update
	unauthenticatedRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/status-updates/:id/document-download", projectActionPlanStatusUpdateController.DownloadProjectActionPlanStatusUpdateDocument)
	unauthenticatedRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/status-updates/:id/document-preview", projectActionPlanStatusUpdateController.GetProjectActionPlanStatusUpdateDocumentPreview)
	projectRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectActionPlanStatusUpdateController.GetAllProjectActionPlanStatusUpdates)
	projectRoutes.POST("/projects/:projectId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectActionPlanStatusUpdateController.CreateProjectActionPlanStatusUpdate)
	projectRoutes.PUT("/projects/:projectId/action-plans/:actionPlanId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectActionPlanStatusUpdateController.UpdateProjectActionPlanStatusUpdate)
	projectRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectActionPlanStatusUpdateController.GetLatestProjectActionPlanStatusUpdate)
	projectRoutes.PUT("/projects/:projectId/action-plans/:actionPlanId/approve", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectActionPlanStatusUpdateController.ApproveProjectActionPlan)
	projectRoutes.PUT("/projects/:projectId/action-plans/:actionPlanId/revert", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectActionPlanStatusUpdateController.RevertProjectActionPlan)

	// Project Action Plan Comments
	unauthenticatedRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/comments/:commentId/document-download", projectActionPlanCommentController.DownloadProjectActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/comments/:commentId/document-preview", projectActionPlanCommentController.GetProjectActionPlanCommentDocumentPreview)
	projectRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectActionPlanCommentController.GetAllComments)
	projectRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectActionPlanCommentController.GetCommentDetails)
	projectRoutes.POST("/projects/:projectId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectActionPlanCommentController.CreateComment)
	projectRoutes.PUT("/projects/:projectId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectActionPlanCommentController.UpdateComment)
	projectRoutes.POST("/projects/:projectId/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectActionPlanCommentController.DeleteComment)

	projectRoutes.GET("/projects/:projectId/documents", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectAssociatedItemController.GetAllProjectDocuments)
	projectRoutes.POST("/projects/:projectId/documents", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.CreateProjectDocument)
	projectRoutes.DELETE("/projects/:projectId/documents/:documentId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.DeleteProjectDocument)

}
