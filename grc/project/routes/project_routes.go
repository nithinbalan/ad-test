package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	projectRoutes := router.Group("/")

	projectRoutes.Use(middlewares.AuthMiddleware(app.DB))

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
	projectClassificationController := app.ProjectModule.ProjectClassificationController
	projectPhaseController := app.ProjectModule.ProjectPhaseController
	projectChangeRequestController := app.ProjectModule.ProjectChangeRequestController
	projectChangeRequestWorkflowController := app.ProjectModule.ProjectChangeRequestWorkflowController
	projectChangeRequestReviewController := app.ProjectModule.ProjectChangeRequestReviewController
	projectChangeRequestReviewCommentController := app.ProjectModule.ProjectChangeRequestReviewCommentController
	projectImpactAssessmentTypeController := app.ProjectModule.ProjectImpactAssessmentTypeController
	projectSecurityGateClassificationController := app.ProjectModule.ProjectSecurityGateClassificationController
	projectSecurityGateStageController := app.ProjectModule.ProjectSecurityGateStageController
	projectSecurityGateChecklistTemplateController := app.ProjectModule.ProjectSecurityGateChecklistTemplateController
	projectSecurityGateChecklistResponseStatusController := app.ProjectModule.ProjectSecurityGateChecklistResponseStatusController
	projectWorkflowLevelController := app.ProjectModule.ProjectWorkflowLevelController
	projectDecisionLogController := app.ProjectModule.ProjectDecisionLogController
	projectDecisionLogWorkflowController := app.ProjectModule.ProjectDecisionLogWorkflowController
	projectDecisionLogReviewController := app.ProjectModule.ProjectDecisionLogReviewController
	projectDecisionLogReviewCommentController := app.ProjectModule.ProjectDecisionLogReviewCommentController
	projectSupportingDocumentController := app.ProjectModule.ProjectSupportingDocumentController
	projectClosureController := app.ProjectModule.ProjectClosureController
	projectSecurityGateReadinessStatusController := app.ProjectModule.ProjectSecurityGateReadinessStatusController
	projectSecurityGateController := app.ProjectModule.ProjectSecurityGateController
	projectSecurityGateChecklistController := app.ProjectModule.ProjectSecurityGateChecklistController
	projectReportController := app.ProjectModule.ProjectReportController

	// Project Types
	projectRoutes.GET("/project-types", projectTypeController.GetAllProjectTypes)
	projectRoutes.GET("/project-types/:id", projectTypeController.GetProjectTypeDetails)
	projectRoutes.POST("/project-types", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectTypeController.CreateProjectType)
	projectRoutes.PUT("/project-types/:id", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectTypeController.UpdateProjectType)
	projectRoutes.POST("/project-types/delete", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectTypeController.DeleteProjectTypes)
	projectRoutes.PUT("/project-types/activate", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectTypeController.ActivateProjectTypes)
	projectRoutes.PUT("/project-types/deactivate", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectTypeController.DeactivateProjectTypes)

	// Project Status
	projectRoutes.GET("/project-statuses", projectStatusController.GetAllProjectStatuses)
	projectRoutes.GET("/project-statuses/:id", projectStatusController.GetProjectStatusDetails)

	// Projects
	projectRoutes.GET("/projects", projectController.GetAllProjects)
	projectRoutes.GET("/projects/export", projectController.GetProjectsExport)
	projectRoutes.GET("/projects/:projectId/export-pdf", projectReportController.ExportProjectReportToPDF)
	projectRoutes.GET("/projects/:projectId", projectController.GetProjectDetails)
	projectRoutes.POST("/projects", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectController.CreateProject)
	projectRoutes.PUT("/projects/:projectId", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectController.UpdateProject)
	projectRoutes.POST("/projects/delete", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectController.DeleteProjects)
	projectRoutes.PUT("/projects/:projectId/status-update", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectController.UpdateProjectStatus)

	// Project Milestones
	projectRoutes.GET("/project-milestones", projectMilestoneController.GetAllProjectMilestones)
	projectRoutes.GET("/projects/:projectId/milestones", projectMilestoneController.GetAllProjectMilestones)
	projectRoutes.GET("/projects/:projectId/milestones/:milestoneId", projectMilestoneController.GetProjectMilestoneDetails)
	projectRoutes.POST("/projects/:projectId/milestones", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectMilestoneController.CreateProjectMilestone)
	projectRoutes.PUT("/projects/:projectId/milestones/:milestoneId", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectMilestoneController.UpdateProjectMilestone)
	projectRoutes.PUT("/projects/:projectId/milestones/:milestoneId/status-update", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectMilestoneController.UpdateProjectMilestoneStatus)
	projectRoutes.POST("/projects/:projectId/milestones/delete", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectMilestoneController.DeleteProjectMilestones)

	// Project Change Requests
	projectRoutes.GET("/projects/:projectId/change-requests", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectChangeRequestController.GetAllProjectChangeRequests)
	projectRoutes.GET("/projects/:projectId/change-requests/summary", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectChangeRequestController.GetProjectChangeRequestSummary)
	projectRoutes.GET("/projects/:projectId/change-requests/:changeRequestId", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectChangeRequestController.GetProjectChangeRequestDetails)
	projectRoutes.POST("/projects/:projectId/change-requests", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectChangeRequestController.CreateProjectChangeRequest)
	projectRoutes.PUT("/projects/:projectId/change-requests/:changeRequestId", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectChangeRequestController.UpdateProjectChangeRequest)
	projectRoutes.POST("/projects/:projectId/change-requests/delete", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectChangeRequestController.DeleteProjectChangeRequests)

	// Project Change Request Review & Workflow
	projectRoutes.POST("/projects/:projectId/change-requests/:changeRequestId/submit", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectChangeRequestReviewController.SubmitChangeRequest)
	projectRoutes.POST("/projects/:projectId/change-requests/:changeRequestId/resubmit", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectChangeRequestReviewController.ResubmitChangeRequest)
	projectRoutes.POST("/projects/:projectId/change-requests/:changeRequestId/approve", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectChangeRequestReviewController.ApproveChangeRequest)
	projectRoutes.POST("/projects/:projectId/change-requests/:changeRequestId/revert", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectChangeRequestReviewController.RevertChangeRequest)
	projectRoutes.POST("/projects/:projectId/change-requests/:changeRequestId/reject", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectChangeRequestReviewController.RejectChangeRequest)
	projectRoutes.GET("/projects/:projectId/change-requests/:changeRequestId/workflow-list", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectChangeRequestWorkflowController.GetWorkflowList)

	// Project Change Request Review Comments
	unauthenticatedRoutes.GET("/projects/:projectId/change-requests/:changeRequestId/review-comments/:commentId/document-download", projectChangeRequestReviewCommentController.DownloadDocument)
	unauthenticatedRoutes.GET("/projects/:projectId/change-requests/:changeRequestId/review-comments/:commentId/document-preview", projectChangeRequestReviewCommentController.GetDocumentPreview)
	projectRoutes.GET("/projects/:projectId/change-requests/:changeRequestId/review-comments", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectChangeRequestReviewCommentController.GetAllComments)
	projectRoutes.GET("/projects/:projectId/change-requests/:changeRequestId/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectChangeRequestReviewCommentController.GetCommentDetails)
	projectRoutes.POST("/projects/:projectId/change-requests/:changeRequestId/review-comments", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectChangeRequestReviewCommentController.CreateComment)
	projectRoutes.PUT("/projects/:projectId/change-requests/:changeRequestId/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectChangeRequestReviewCommentController.UpdateComment)
	projectRoutes.POST("/projects/:projectId/change-requests/:changeRequestId/review-comments/delete", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectChangeRequestReviewCommentController.DeleteComments)
	projectRoutes.PUT("/projects/:projectId/change-requests/:changeRequestId/review-comments/:commentId/resolve", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectChangeRequestReviewCommentController.ResolveComment)

	// Project Decision Logs
	unauthenticatedRoutes.GET("/projects/:projectId/decision-logs/:decisionLogId/documents/:documentId/download", projectDecisionLogController.DownloadDocument)
	unauthenticatedRoutes.GET("/projects/:projectId/decision-logs/:decisionLogId/documents/:documentId/preview", projectDecisionLogController.GetDocumentPreview)
	projectRoutes.GET("/projects/:projectId/decision-logs", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectDecisionLogController.GetAllDecisionLogs)
	projectRoutes.GET("/projects/:projectId/decision-logs/summary", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectDecisionLogController.GetProjectDecisionLogSummary)
	projectRoutes.GET("/projects/:projectId/decision-logs/:decisionLogId", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectDecisionLogController.GetDecisionLogDetails)
	projectRoutes.POST("/projects/:projectId/decision-logs", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectDecisionLogController.CreateDecisionLog)
	projectRoutes.PUT("/projects/:projectId/decision-logs/:decisionLogId", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectDecisionLogController.UpdateDecisionLog)
	projectRoutes.POST("/projects/:projectId/decision-logs/delete", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectDecisionLogController.DeleteDecisionLogs)

	// Project Decision Log Review & Workflow
	projectRoutes.POST("/projects/:projectId/decision-logs/:decisionLogId/submit", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectDecisionLogReviewController.SubmitDecisionLog)
	projectRoutes.POST("/projects/:projectId/decision-logs/:decisionLogId/resubmit", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectDecisionLogReviewController.ResubmitDecisionLog)
	projectRoutes.POST("/projects/:projectId/decision-logs/:decisionLogId/approve", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectDecisionLogReviewController.ApproveDecisionLog)
	projectRoutes.POST("/projects/:projectId/decision-logs/:decisionLogId/revert", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectDecisionLogReviewController.RevertDecisionLog)
	projectRoutes.POST("/projects/:projectId/decision-logs/:decisionLogId/reject", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectDecisionLogReviewController.RejectDecisionLog)
	projectRoutes.GET("/projects/:projectId/decision-logs/:decisionLogId/workflow-list", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectDecisionLogWorkflowController.GetWorkflowList)

	// Project Decision Log Review Comments
	unauthenticatedRoutes.GET("/projects/:projectId/decision-logs/:decisionLogId/review-comments/:commentId/document-download", projectDecisionLogReviewCommentController.DownloadDocument)
	unauthenticatedRoutes.GET("/projects/:projectId/decision-logs/:decisionLogId/review-comments/:commentId/document-preview", projectDecisionLogReviewCommentController.GetDocumentPreview)
	projectRoutes.GET("/projects/:projectId/decision-logs/:decisionLogId/review-comments", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectDecisionLogReviewCommentController.GetAllComments)
	projectRoutes.GET("/projects/:projectId/decision-logs/:decisionLogId/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectDecisionLogReviewCommentController.GetCommentDetails)
	projectRoutes.POST("/projects/:projectId/decision-logs/:decisionLogId/review-comments", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectDecisionLogReviewCommentController.CreateComment)
	projectRoutes.PUT("/projects/:projectId/decision-logs/:decisionLogId/review-comments/:commentId", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectDecisionLogReviewCommentController.UpdateComment)
	projectRoutes.POST("/projects/:projectId/decision-logs/:decisionLogId/review-comments/delete", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectDecisionLogReviewCommentController.DeleteComments)
	projectRoutes.PUT("/projects/:projectId/decision-logs/:decisionLogId/review-comments/:commentId/resolve", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectDecisionLogReviewCommentController.ResolveComment)

	// Project Supporting Documents
	unauthenticatedRoutes.GET("/projects/:projectId/supporting-documents/:supportingDocumentId/documents/:documentId/download", projectSupportingDocumentController.DownloadDocument)
	unauthenticatedRoutes.GET("/projects/:projectId/supporting-documents/:supportingDocumentId/documents/:documentId/preview", projectSupportingDocumentController.GetDocumentPreview)
	projectRoutes.GET("/projects/:projectId/supporting-documents", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectSupportingDocumentController.GetAllSupportingDocuments)
	projectRoutes.GET("/projects/:projectId/supporting-documents/:supportingDocumentId", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectSupportingDocumentController.GetSupportingDocumentDetails)
	projectRoutes.POST("/projects/:projectId/supporting-documents", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectSupportingDocumentController.CreateSupportingDocument)
	projectRoutes.PUT("/projects/:projectId/supporting-documents/:supportingDocumentId", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectSupportingDocumentController.UpdateSupportingDocument)
	projectRoutes.POST("/projects/:projectId/supporting-documents/delete", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectSupportingDocumentController.DeleteSupportingDocuments)

	// Project Action Plans
	projectRoutes.GET("/project-action-plans", projectActionPlanController.GetAllProjectActionPlans)
	projectRoutes.GET("/projects/:projectId/action-plans", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectActionPlanController.GetAllProjectActionPlans)
	projectRoutes.GET("/projects/:projectId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectActionPlanController.GetProjectActionPlanDetails)
	projectRoutes.POST("/projects/:projectId/action-plans", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectActionPlanController.CreateProjectActionPlan)
	projectRoutes.PUT("/projects/:projectId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectActionPlanController.UpdateProjectActionPlan)
	projectRoutes.POST("/projects/:projectId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectActionPlanController.DeleteProjectActionPlans)

	// Project Action Plan	 Status Update
	unauthenticatedRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/status-updates/:id/document-download", projectActionPlanStatusUpdateController.DownloadProjectActionPlanStatusUpdateDocument)
	unauthenticatedRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/status-updates/:id/document-preview", projectActionPlanStatusUpdateController.GetProjectActionPlanStatusUpdateDocumentPreview)
	projectRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectActionPlanStatusUpdateController.GetAllProjectActionPlanStatusUpdates)
	projectRoutes.POST("/projects/:projectId/action-plans/:actionPlanId/status-updates", projectActionPlanStatusUpdateController.CreateProjectActionPlanStatusUpdate)
	projectRoutes.PUT("/projects/:projectId/action-plans/:actionPlanId/status-updates/:id", projectActionPlanStatusUpdateController.UpdateProjectActionPlanStatusUpdate)
	projectRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectActionPlanStatusUpdateController.GetLatestProjectActionPlanStatusUpdate)
	projectRoutes.PUT("/projects/:projectId/action-plans/:actionPlanId/approve", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectActionPlanStatusUpdateController.ApproveProjectActionPlan)
	projectRoutes.PUT("/projects/:projectId/action-plans/:actionPlanId/revert", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectActionPlanStatusUpdateController.RevertProjectActionPlan)

	// Project Action Plan Comments
	unauthenticatedRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/comments/:commentId/document-download", projectActionPlanCommentController.DownloadProjectActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/comments/:commentId/document-preview", projectActionPlanCommentController.GetProjectActionPlanCommentDocumentPreview)
	projectRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectActionPlanCommentController.GetAllComments)
	projectRoutes.GET("/projects/:projectId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectActionPlanCommentController.GetCommentDetails)
	projectRoutes.POST("/projects/:projectId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectActionPlanCommentController.CreateComment)
	projectRoutes.PUT("/projects/:projectId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectActionPlanCommentController.UpdateComment)
	projectRoutes.POST("/projects/:projectId/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectActionPlanCommentController.DeleteComment)

	projectRoutes.GET("/projects/:projectId/documents", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectAssociatedItemController.GetAllProjectDocuments)
	projectRoutes.POST("/projects/:projectId/documents", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectAssociatedItemController.CreateProjectDocument)
	projectRoutes.DELETE("/projects/:projectId/documents/:documentId", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectAssociatedItemController.DeleteProjectDocument)

	projectRoutes.GET("/projects/:projectId/controls", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectAssociatedItemController.GetAllProjectControls)
	projectRoutes.POST("/projects/:projectId/controls", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectAssociatedItemController.CreateProjectControl)
	projectRoutes.DELETE("/projects/:projectId/controls/:controlId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.DeleteProjectContol)

	projectRoutes.GET("/projects/:projectId/processes", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectAssociatedItemController.GetAllProjectProcesses)
	projectRoutes.POST("/projects/:projectId/processes", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.CreateProjectProcess)
	projectRoutes.DELETE("/projects/:projectId/processes/:processId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.DeleteProjectProcess)

	projectRoutes.GET("/projects/:projectId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectAssociatedItemController.GetAllProjectStrategicObjectives)
	projectRoutes.POST("/projects/:projectId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.CreateProjectStrategicObjective)
	projectRoutes.DELETE("/projects/:projectId/strategic-objectives/:strategicObjectiveId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.DeleteProjectStrategicObjective)

	projectRoutes.GET("/projects/:projectId/kpis", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectAssociatedItemController.GetAllProjectKPI)
	projectRoutes.POST("/projects/:projectId/kpis", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.CreateProjectKPI)
	projectRoutes.DELETE("/projects/:projectId/kpis/:kpiRegisterId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.DeleteProjectKPI)

	projectRoutes.GET("/projects/:projectId/assets", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectAssociatedItemController.GetAllProjectAssets)
	projectRoutes.POST("/projects/:projectId/assets", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.CreateProjectAsset)
	projectRoutes.DELETE("/projects/:projectId/assets/:assetId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.DeleteProjectAsset)

	projectRoutes.GET("/projects/:projectId/incidents", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectAssociatedItemController.GetAllProjectIncidents)
	projectRoutes.POST("/projects/:projectId/incidents", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.CreateProjectIncident)
	projectRoutes.DELETE("/projects/:projectId/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.DeleteProjectIncident)

	projectRoutes.GET("/projects/:projectId/risks", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectAssociatedItemController.GetAllProjectRisks)
	projectRoutes.POST("/projects/:projectId/risks", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.CreateProjectRisk)
	projectRoutes.DELETE("/projects/:projectId/risks/:riskId", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.DeleteProjectRisk)

	projectRoutes.GET("/projects/:projectId/clauses", middlewares.RolePermissionMiddleware(db, "organization_setting", "view"), projectAssociatedItemController.GetAllClauses)
	projectRoutes.PUT("/projects/:projectId/clauses", middlewares.RolePermissionMiddleware(db, "organization_setting", "modify"), projectAssociatedItemController.UpdateClauses)

	// project classifications
	projectRoutes.GET("/project/project-classifications", projectClassificationController.GetAllProjectClassifications)
	projectRoutes.GET("/project/project-classifications/:id", projectClassificationController.GetProjectClassificationDetails)
	projectRoutes.POST("/project/project-classifications", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectClassificationController.CreateProjectClassification)
	projectRoutes.PUT("/project/project-classifications/:id", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectClassificationController.UpdateProjectClassification)
	projectRoutes.POST("/project/project-classifications/delete", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectClassificationController.DeleteProjectClassifications)
	projectRoutes.PUT("/project/project-classifications/activate", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectClassificationController.ActivateProjectClassifications)
	projectRoutes.PUT("/project/project-classifications/deactivate", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectClassificationController.DeactivateProjectClassifications)

	// Project Phases
	projectRoutes.GET("/project-phases", projectPhaseController.GetAllProjectPhases)
	projectRoutes.GET("/project-phases/:id", projectPhaseController.GetProjectPhaseDetails)

	// Project Impact Assessment Types
	projectRoutes.GET("/project-impact-assessment-types", projectImpactAssessmentTypeController.GetAllProjectImpactAssessmentTypes)
	projectRoutes.GET("/project-impact-assessment-types/:id", projectImpactAssessmentTypeController.GetProjectImpactAssessmentTypeDetails)

	// Project Security Gate Classifications
	projectRoutes.GET("/project-security-gate-classifications", projectSecurityGateClassificationController.GetAllProjectSecurityGateClassifications)
	projectRoutes.GET("/project-security-gate-classifications/:id", projectSecurityGateClassificationController.GetProjectSecurityGateClassificationDetails)

	// Project Security Gate Stages
	projectRoutes.GET("/project-security-gate-stages", projectSecurityGateStageController.GetAllProjectSecurityGateStages)
	projectRoutes.GET("/project-security-gate-stages/:id", projectSecurityGateStageController.GetProjectSecurityGateStageDetails)

	// Project Security Gate Checklist Template
	projectRoutes.GET("/project-security-gate-checklist-templates", projectSecurityGateChecklistTemplateController.GetAllProjectSecurityGateChecklistTemplates)
	projectRoutes.GET("/project-security-gate-checklist-templates/:templateId", projectSecurityGateChecklistTemplateController.GetProjectSecurityGateChecklistTemplateDetails)
	projectRoutes.POST("/project-security-gate-checklist-templates/:templateId/versions", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectSecurityGateChecklistTemplateController.CreateProjectSecurityGateChecklistTemplateVersion)
	projectRoutes.PUT("/project-security-gate-checklist-templates/:templateId/versions/:versionId/publish", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectSecurityGateChecklistTemplateController.PublishProjectSecurityGateChecklistTemplateVersion)
	projectRoutes.POST("/project-security-gate-checklist-templates/:templateId/versions/:versionId/requirements", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectSecurityGateChecklistTemplateController.CreateProjectSecurityGateChecklistTemplateRequirement)
	projectRoutes.PUT("/project-security-gate-checklist-templates/:templateId/versions/:versionId/requirements/:requirementId", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectSecurityGateChecklistTemplateController.UpdateProjectSecurityGateChecklistTemplateRequirement)
	projectRoutes.DELETE("/project-security-gate-checklist-templates/:templateId/versions/:versionId/requirements/:requirementId", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectSecurityGateChecklistTemplateController.DeleteProjectSecurityGateChecklistTemplateRequirement)

	// Project Workflow Levels
	projectRoutes.POST("/project-workflows/:workflowId/levels", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectWorkflowLevelController.CreateWorkflowLevel)
	projectRoutes.DELETE("/project-workflows/:workflowId/levels/:levelId", middlewares.RolePermissionMiddleware(db, "project_setting", "modify"), projectWorkflowLevelController.DeleteWorkflowLevel)

	// Project Security Gate Checklist Response Statuses
	projectRoutes.GET("/project-security-gate-checklist-response-statuses", projectSecurityGateChecklistResponseStatusController.GetAllProjectSecurityGateChecklistResponseStatuses)
	projectRoutes.GET("/project-security-gate-checklist-response-statuses/:id", projectSecurityGateChecklistResponseStatusController.GetProjectSecurityGateChecklistResponseStatusDetails)

	// Project Closures
	projectRoutes.GET("/projects/:projectId/closures", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectClosureController.GetAllProjectClosures)
	projectRoutes.GET("/projects/:projectId/closures/:closureId", middlewares.RolePermissionMiddleware(db, "project_register", "view"), projectClosureController.GetProjectClosureDetails)
	projectRoutes.POST("/projects/:projectId/closures", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectClosureController.CreateProjectClosure)
	projectRoutes.PUT("/projects/:projectId/closures/:closureId", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectClosureController.UpdateProjectClosure)
	projectRoutes.DELETE("/projects/:projectId/closures/:closureId", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectClosureController.DeleteProjectClosure)
	projectRoutes.POST("/projects/:projectId/closures/:closureId/approve", middlewares.RolePermissionMiddleware(db, "project_register", "modify"), projectClosureController.ApproveProjectClosure)

	// Project Security Gate Readiness Statuses
	projectRoutes.GET("/project-security-gate-readiness-statuses", projectSecurityGateReadinessStatusController.GetAllProjectSecurityGateReadinessStatuses)
	projectRoutes.GET("/project-security-gate-readiness-statuses/:id", projectSecurityGateReadinessStatusController.GetProjectSecurityGateReadinessStatusDetails)

	// Project Security Gates
	projectRoutes.GET("/project-security-gates", middlewares.RolePermissionMiddleware(db, "project_security_gate", "view"), projectSecurityGateController.GetAllProjectSecurityGates)
	projectRoutes.GET("/project-security-gates/summary", middlewares.RolePermissionMiddleware(db, "project_security_gate", "view"), projectSecurityGateController.GetProjectSecurityGateSummary)
	projectRoutes.GET("/project-security-gates/:id", middlewares.RolePermissionMiddleware(db, "project_security_gate", "view"), projectSecurityGateController.GetProjectSecurityGateDetails)
	projectRoutes.POST("/project-security-gates", middlewares.RolePermissionMiddleware(db, "project_security_gate", "modify"), projectSecurityGateController.CreateProjectSecurityGate)
	projectRoutes.PUT("/project-security-gates/:id", middlewares.RolePermissionMiddleware(db, "project_security_gate", "modify"), projectSecurityGateController.UpdateProjectSecurityGate)
	projectRoutes.POST("/project-security-gates/delete", middlewares.RolePermissionMiddleware(db, "project_security_gate", "modify"), projectSecurityGateController.DeleteProjectSecurityGates)

	// Project Security Gate Checklist
	projectRoutes.GET("/project-security-gates/:id/checklists", middlewares.RolePermissionMiddleware(db, "project_security_gate", "view"), projectSecurityGateChecklistController.GetAllProjectSecurityGateChecklists)
	projectRoutes.GET("/project-security-gates/:id/checklists/summary", middlewares.RolePermissionMiddleware(db, "project_security_gate", "view"), projectSecurityGateChecklistController.GetProjectSecurityGateChecklistSummary)
	projectRoutes.GET("/project-security-gates/:id/checklists/:checklistId", middlewares.RolePermissionMiddleware(db, "project_security_gate", "view"), projectSecurityGateChecklistController.GetProjectSecurityGateChecklistDetails)
	projectRoutes.POST("/project-security-gates/:id/checklists/:checklistId/response", middlewares.RolePermissionMiddleware(db, "project_security_gate", "modify"), projectSecurityGateChecklistController.CreateProjectSecurityGateChecklistResponse)
	projectRoutes.PUT("/project-security-gates/:id/checklists/:checklistId/response", middlewares.RolePermissionMiddleware(db, "project_security_gate", "modify"), projectSecurityGateChecklistController.UpdateProjectSecurityGateChecklistResponse)

	unauthenticatedRoutes.GET("/project-security-gates/:id/checklists/:checklistId/response/documents/document-download", projectSecurityGateChecklistController.DownloadProjectSecurityGateChecklistResponseDocument)
	unauthenticatedRoutes.GET("/project-security-gates/:id/checklists/:checklistId/response/documents/document-preview", projectSecurityGateChecklistController.GetProjectSecurityGateChecklistResponseDocumentPreview)
}
