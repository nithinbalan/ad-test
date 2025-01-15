package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	auditRoutes := router.Group("/")

	auditRoutes.Use(middlewares.AuthMiddleware())

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	auditStatusController := app.AuditModule.AuditStatusController
	auditorRiskRatingController := app.AuditModule.AuditorRiskRatingController
	auditUniverseController := app.AuditModule.AuditUniverseController
	auditableItemController := app.AuditModule.AuditableItemController
	auditTeamController := app.AuditModule.AuditTeamController
	auditCategoryController := app.AuditModule.AuditCategoryController
	auditPlanController := app.AuditModule.AuditPlanController
	auditPlanChecklistController := app.AuditModule.AuditPlanChecklistController
	auditPlanToolController := app.AuditModule.AuditPlanToolController
	auditPlanAuditableItemController := app.AuditModule.AuditPlanAuditableItemController
	auditActivityTypeController := app.AuditModule.AuditActivityTypeController
	auditScheduleStatusController := app.AuditModule.AuditScheduleStatusController
	auditPlanBudgetController := app.AuditModule.AuditPlanBudgetController
	auditPlanScheduleController := app.AuditModule.AuditPlanScheduleController
	auditEngagementTypeController := app.AuditModule.AuditEngagementTypeController
	auditPlanStakeholderEngagementController := app.AuditModule.AuditPlanStakeholderEngagementController
	auditController := app.AuditModule.AuditController
	auditWorkpaperController := app.AuditModule.AuditWorkpaperController
	auditMeetingController := app.AuditModule.AuditMeetingController
	auditTimelineController := app.AuditModule.AuditTimelineController
	auditFindingStatusController := app.AuditModule.AuditFindingStatusController
	auditFindingCategoryController := app.AuditModule.AuditFindingCategoryController
	auditFindingController := app.AuditModule.AuditFindingController
	auditDocumentRequestStatusController := app.AuditModule.AuditDocumentRequestStatusController
	auditDocumentRequestController := app.AuditModule.AuditDocumentRequestController
	auditFindingRootCauseAnalysisController := app.AuditModule.AuditFindingRootCauseAnalysisController
	auditFindingActionPlanController := app.AuditModule.AuditFindingActionPlanController
	auditDocumentRequestConversationController := app.AuditModule.AuditDocumentRequestConversationController
	auditFindingActionPlanStatusUpdateController := app.AuditModule.AuditFindingActionPlanStatusUpdateController
	auditChecklistController := app.AuditModule.AuditChecklistController
	auditScheduleController := app.AuditModule.AuditScheduleController
	auditDraftFindingStatusController := app.AuditModule.AuditDraftFindingStatusController
	auditDraftReportStatusController := app.AuditModule.AuditDraftReportStatusController
	auditDraftFindingController := app.AuditModule.AuditDraftFindingController
	auditDraftReportController := app.AuditModule.AuditDraftReportController
	rcmController := app.AuditModule.RCMController
	auditDraftFindingConversationController := app.AuditModule.AuditDraftFindingConversationController
	auditDraftFindingActionPlanController := app.AuditModule.AuditDraftFindingActionPlanController
	auditScheduleStatusUpdateController := app.AuditModule.AuditScheduleStatusUpdateController
	auditDraftReportReviewController := app.AuditModule.AuditDraftReportReviewController
	auditDraftReportCommentController := app.AuditModule.AuditDraftReportCommentController
	auditPlanCommencementLetterStatusController := app.AuditModule.AuditPlanCommencementLetterStatusController
	auditPlanCommencementLetterController := app.AuditModule.AuditPlanCommencementLetterController
	auditFinalReportStatusController := app.AuditModule.AuditFinalReportStatusController
	auditFinalReportController := app.AuditModule.AuditFinalReportController
	auditFinalReportSectionController := app.AuditModule.AuditFinalReportSectionController
	auditReportWorkflowController := app.AuditModule.AuditReportWorkflowController
	auditReportReviewController := app.AuditModule.AuditReportReviewController

	//Audit Status
	auditRoutes.GET("/audit-statuses", auditStatusController.GetAllAuditStatuses)
	auditRoutes.GET("/audit-statuses/:id", auditStatusController.GetAuditStatusDetails)

	//Auditor Risk Rating
	auditRoutes.GET("/auditor-risk-ratings", auditorRiskRatingController.GetAllAuditorRiskRatings)
	auditRoutes.GET("/auditor-risk-ratings/:id", auditorRiskRatingController.GetAuditorRiskRatingDetails)

	//Audit Universe
	auditRoutes.POST("/audit-universe/mark-as-auditable-item", middlewares.RolePermissionMiddleware(db, "audit_universe", "modify"), auditUniverseController.MarkAsAuditableItem)
	auditRoutes.POST("/audit-universe/unmark-as-auditable-item", middlewares.RolePermissionMiddleware(db, "audit_universe", "modify"), auditUniverseController.UnmarkAsAuditableItem)

	//Auditable Item
	auditRoutes.GET("/auditable-items", middlewares.RolePermissionMiddleware(db, "audit_universe", "view"), auditableItemController.GetAllAuditableItems)
	auditRoutes.GET("/auditable-items/:id", middlewares.RolePermissionMiddleware(db, "audit_universe", "view"), auditableItemController.GetAuditableItemDetails)
	auditRoutes.POST("/auditable-items", middlewares.RolePermissionMiddleware(db, "audit_universe", "modify"), auditableItemController.CreateAuditableItem)
	auditRoutes.PUT("/auditable-items/:id", middlewares.RolePermissionMiddleware(db, "audit_universe", "modify"), auditableItemController.UpdateAuditableItem)
	auditRoutes.POST("/auditable-items/delete", middlewares.RolePermissionMiddleware(db, "audit_universe", "modify"), auditableItemController.DeleteAuditableItems)

	//Audit Team
	auditRoutes.GET("/audit-teams", auditTeamController.GetAllAuditTeams)
	auditRoutes.GET("/audit-teams/:id", auditTeamController.GetAuditTeamDetails)
	auditRoutes.POST("/audit-teams", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditTeamController.CreateAuditTeam)
	auditRoutes.PUT("/audit-teams/:id", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditTeamController.UpdateAuditTeam)
	auditRoutes.POST("/audit-teams/delete", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditTeamController.DeleteAuditTeams)

	//Audit Category
	auditRoutes.GET("/audit-categories", auditCategoryController.GetAllAuditCategories)
	auditRoutes.GET("/audit-categories/:id", auditCategoryController.GetAuditCategoryDetails)
	auditRoutes.POST("/audit-categories", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditCategoryController.CreateAuditCategory)
	auditRoutes.PUT("/audit-categories/:id", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditCategoryController.UpdateAuditCategory)
	auditRoutes.POST("/audit-categories/delete", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditCategoryController.DeleteAuditCategories)
	auditRoutes.PUT("/audit-categories/activate", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditCategoryController.ActivateAuditCategories)
	auditRoutes.PUT("/audit-categories/deactivate", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditCategoryController.DeactivateAuditCategories)

	//AuditPlan
	auditRoutes.GET("/audit-plans", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanController.GetAllAuditPlans)
	auditRoutes.GET("/audit-plans/:planId", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanController.GetAuditPlanDetails)
	auditRoutes.POST("/audit-plans", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanController.CreateAuditPlan)
	auditRoutes.PUT("/audit-plans/:planId", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanController.UpdateAuditPlan)
	auditRoutes.POST("/audit-plans/delete", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanController.DeleteAuditPlans)
	auditRoutes.POST("/audit-plans/:planId/publish", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanController.PublishAuditPlan)

	//AuditPlan Checklist
	auditRoutes.GET("/audit-plans/:planId/checklists", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanChecklistController.GetAllAuditPlanChecklists)
	auditRoutes.GET("/audit-plans/:planId/checklists/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanChecklistController.GetAuditPlanChecklistDetails)
	auditRoutes.POST("/audit-plans/:planId/checklists", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanChecklistController.CreateAuditPlanChecklists)
	auditRoutes.PUT("/audit-plans/:planId/checklists/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanChecklistController.UpdateAuditPlanChecklist)
	auditRoutes.POST("/audit-plans/:planId/checklists/delete", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanChecklistController.DeleteAuditPlanChecklists)

	//AuditPlan Tools
	auditRoutes.GET("/audit-plans/:planId/tools", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanToolController.GetAllAuditPlanTools)
	auditRoutes.GET("/audit-plans/:planId/tools/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanToolController.GetAuditPlanToolDetails)
	auditRoutes.POST("/audit-plans/:planId/tools", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanToolController.CreateAuditPlanTool)
	auditRoutes.PUT("/audit-plans/:planId/tools/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanToolController.UpdateAuditPlanTool)
	auditRoutes.POST("/audit-plans/:planId/tools/delete", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanToolController.DeleteAuditPlanTools)

	//AuditPlan Auditable Item
	auditRoutes.GET("/audit-plans/:planId/auditable-items", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanAuditableItemController.GetAllAuditPlanAuditableItems)
	auditRoutes.GET("/audit-plans/:planId/auditable-items/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanAuditableItemController.GetAllAuditPlanAuditableItemDetails)
	auditRoutes.POST("/audit-plans/:planId/auditable-items", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanAuditableItemController.CreateAuditPlanAuditableItem)
	auditRoutes.POST("/audit-plans/:planId/auditable-items/delete", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanAuditableItemController.DeleteAuditPlanAuditableItems)

	//Audit Activity Type
	auditRoutes.GET("/audit-activity-types", auditActivityTypeController.GetAllAuditActivityTypes)
	auditRoutes.GET("/audit-activity-types/:id", auditActivityTypeController.GetAuditActivityTypeDetails)
	auditRoutes.POST("/audit-activity-types", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditActivityTypeController.CreateAuditActivityType)
	auditRoutes.PUT("/audit-activity-types/:id", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditActivityTypeController.UpdateAuditActivityType)
	auditRoutes.POST("/audit-activity-types/delete", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditActivityTypeController.DeleteAuditActivityTypes)
	auditRoutes.PUT("/audit-activity-types/activate", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditActivityTypeController.ActivateAuditActivityTypes)
	auditRoutes.PUT("/audit-activity-types/deactivate", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditActivityTypeController.DeactivateAuditActivityTypes)

	//Audit Schedule Statuses
	auditRoutes.GET("/audit-schedule-statuses", auditScheduleStatusController.GetAllAuditScheduleStatuses)
	auditRoutes.GET("/audit-schedule-statuses/:id", auditScheduleStatusController.GetAuditScheduleStatusDetails)

	//AuditPlan Tools
	auditRoutes.GET("/audit-plans/:planId/budgets", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanBudgetController.GetAllAuditPlanBudgets)
	auditRoutes.GET("/audit-plans/:planId/budgets/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanBudgetController.GetAuditPlanBudgetDetails)
	auditRoutes.POST("/audit-plans/:planId/budgets", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanBudgetController.CreateAuditPlanBudget)
	auditRoutes.PUT("/audit-plans/:planId/budgets/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanBudgetController.UpdateAuditPlanBudget)
	auditRoutes.POST("/audit-plans/:planId/budgets/delete", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanBudgetController.DeleteAuditPlanBudgets)

	//Audit Schedule
	auditRoutes.GET("/audit-plans/:planId/schedules", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanScheduleController.GetAllAuditPlanSchedules)
	auditRoutes.GET("/audit-plans/:planId/schedules/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanScheduleController.GetAuditPlanScheduleDetails)
	auditRoutes.POST("/audit-plans/:planId/schedules", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanScheduleController.CreateAuditPlanSchedule)
	auditRoutes.PUT("/audit-plans/:planId/schedules/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanScheduleController.UpdateAuditPlanSchedule)
	auditRoutes.POST("/audit-plans/:planId/schedules/delete", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanScheduleController.DeleteAuditPlanSchedules)

	//Audit Engagement Type
	auditRoutes.GET("/audit-engagement-types", auditEngagementTypeController.GetAllAuditEngagementTypes)
	auditRoutes.GET("/audit-engagement-types/:id", auditEngagementTypeController.GetAuditEngagementTypeDetails)
	auditRoutes.POST("/audit-engagement-types", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditEngagementTypeController.CreateAuditEngagementType)
	auditRoutes.PUT("/audit-engagement-types/:id", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditEngagementTypeController.UpdateAuditEngagementType)
	auditRoutes.POST("/audit-engagement-types/delete", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditEngagementTypeController.DeleteAuditEngagementTypes)
	auditRoutes.PUT("/audit-engagement-types/activate", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditEngagementTypeController.ActivateAuditEngagementTypes)
	auditRoutes.PUT("/audit-engagement-types/deactivate", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditEngagementTypeController.DeactivateAuditEngagementTypes)

	//Audit Plan Stakeholder Engagement
	auditRoutes.GET("/audit-plans/:planId/stakeholder-engagements", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanStakeholderEngagementController.GetAllAuditPlanStakeholderEngagements)
	auditRoutes.GET("/audit-plans/:planId/stakeholder-engagements/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanStakeholderEngagementController.GetAuditPlanStakeholderEngagementDetails)
	auditRoutes.POST("/audit-plans/:planId/stakeholder-engagements", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanStakeholderEngagementController.CreateAuditPlanStakeholderEngagement)
	auditRoutes.PUT("/audit-plans/:planId/stakeholder-engagements/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanStakeholderEngagementController.UpdateAuditPlanStakeholderEngagement)
	auditRoutes.POST("/audit-plans/:planId/stakeholder-engagements/delete", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanStakeholderEngagementController.DeleteAuditPlanStakeholderEngagements)

	//Audit
	auditRoutes.GET("/audits", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditController.GetAllAudits)
	auditRoutes.GET("/audits/:auditId", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditController.GetAuditDetails)
	auditRoutes.POST("/audits", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditController.CreateAudit)
	auditRoutes.PUT("/audits/:auditId", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditController.UpdateAudit)
	auditRoutes.POST("/audits/delete", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditController.DeleteAudits)

	//Audit Schedule
	auditRoutes.GET("/audits/:auditId/workpapers", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditWorkpaperController.GetAllAuditWorkpapers)
	auditRoutes.GET("/audits/:auditId/workpapers/:id", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditWorkpaperController.GetAuditWorkpaperDetails)
	auditRoutes.POST("/audits/:auditId/workpapers", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditWorkpaperController.CreateAuditWorkpaper)
	auditRoutes.PUT("/audits/:auditId/workpapers/:id", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditWorkpaperController.UpdateAuditWorkpaper)
	auditRoutes.POST("/audits/:auditId/workpapers/delete", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditWorkpaperController.DeleteAuditWorkpapers)
	unauthenticatedRoutes.GET("/audits/:auditId/workpapers/document-preview", auditWorkpaperController.GetWorkpaperDocumentPreview)
	unauthenticatedRoutes.GET("/audits/:auditId/workpapers/document-download", auditWorkpaperController.DownloadWorkpaperDocument)

	//Audit Meeting
	auditRoutes.GET("/audits/:auditId/meetings", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditMeetingController.GetAllAuditMeetings)
	auditRoutes.GET("/audits/:auditId/meetings/:meetingId", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditMeetingController.GetAuditMeetingDetails)
	auditRoutes.POST("/audits/:auditId/meetings", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditMeetingController.CreateAuditMeeting)
	auditRoutes.PUT("/audits/:auditId/meetings/:meetingId", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditMeetingController.UpdateAuditMeeting)
	auditRoutes.POST("/audits/:auditId/meetings/delete", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditMeetingController.DeleteAuditMeetings)

	//Audit timeline
	auditRoutes.GET("/audits/:auditId/timelines", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditTimelineController.GetAllAuditTimelines)

	//Audit Status
	auditRoutes.GET("/audit-finding-statuses", auditFindingStatusController.GetAllAuditFindingStatuses)
	auditRoutes.GET("/audit-finding-statuses/:id", auditFindingStatusController.GetAuditFindingStatusDetails)

	//Audit Finding Category
	auditRoutes.GET("/audit-finding-categories", auditFindingCategoryController.GetAllAuditFindingCategories)
	auditRoutes.GET("/audit-finding-categories/:id", auditFindingCategoryController.GetAuditFindingCategoryDetails)
	auditRoutes.PUT("/audit-finding-categories/:id", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditFindingCategoryController.UpdateAuditFindingCategory)
	auditRoutes.POST("/audit-finding-categories/delete", middlewares.RolePermissionMiddleware(db, "audit_setting", "modify"), auditFindingCategoryController.DeleteAuditFindingCategories)

	// Audit Finding
	unauthenticatedRoutes.GET("/audit-findings/:findingId/document-preview", auditFindingController.GetAuditFindingDocumentPreview)
	unauthenticatedRoutes.GET("/audit-findings/:findingId/document-download", auditFindingController.DownloadAuditFindingDocument)
	auditRoutes.GET("/audit-findings", middlewares.RolePermissionMiddleware(db, "audit_findings", "view"), auditFindingController.GetAllAuditFindings)
	auditRoutes.GET("/audit-findings/:findingId", middlewares.RolePermissionMiddleware(db, "audit_findings", "view"), auditFindingController.GetAuditFindingDetails)
	auditRoutes.POST("/audit-findings", middlewares.RolePermissionMiddleware(db, "audit_findings", "modify"), auditFindingController.CreateAuditFinding)
	auditRoutes.PUT("/audit-findings/:findingId", middlewares.RolePermissionMiddleware(db, "audit_findings", "modify"), auditFindingController.UpdateAuditFinding)
	auditRoutes.POST("/audit-findings/delete", middlewares.RolePermissionMiddleware(db, "audit_findings", "modify"), auditFindingController.DeleteAuditFindings)
	auditRoutes.GET("audits/:auditId/audit-finding-summary", middlewares.RolePermissionMiddleware(db, "audit_findings", "view"), auditFindingController.GetAuditFindingCountByCategory)

	//Audit Document  Request Status
	auditRoutes.GET("/audit-document-request-statuses", auditDocumentRequestStatusController.GetAllAuditDocumentRequestStatuses)

	// Audit Document Request
	auditRoutes.GET("/audits/:auditId/document-requests", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDocumentRequestController.GetAllAuditDocumentRequests)
	auditRoutes.GET("/audits/:auditId/document-requests/:id", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDocumentRequestController.GetAuditDocumentRequestDetails)
	auditRoutes.POST("/audits/:auditId/document-requests", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDocumentRequestController.CreateAuditDocumentRequest)
	auditRoutes.PUT("/audits/:auditId/document-requests/:id", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDocumentRequestController.UpdateAuditDocumentRequest)
	auditRoutes.POST("/audits/:auditId/document-requests/delete", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDocumentRequestController.DeleteAuditDocumentRequests)
	auditRoutes.POST("/audits/:auditId/document-requests/:id/submit", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDocumentRequestController.SubmitAuditDocumentRequest)

	// Audit Document Request Conversations
	auditRoutes.GET("/audits/:auditId/document-requests/:id/conversations", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDocumentRequestConversationController.GetAllAuditDocumentRequestConversations)
	auditRoutes.POST("/audits/:auditId/document-requests/:id/conversations", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDocumentRequestConversationController.CreateAuditDocumentRequestConversation)
	unauthenticatedRoutes.GET("/audits/:auditId/document-requests/:id/conversations/:conversationId/document-download", auditDocumentRequestConversationController.DownloadDocumentRequestConversationDocument)
	unauthenticatedRoutes.GET("/audits/:auditId/document-requests/:id/conversations/:conversationId/document-preview", auditDocumentRequestConversationController.GetDocumentRequestConversationDocumentPreview)

	auditRoutes.GET("/audit-findings/:findingId/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "audit_findings", "view"), auditFindingRootCauseAnalysisController.GetAllAuditFindingRootCauseAnalyses)
	auditRoutes.POST("/audit-findings/:findingId/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "audit_findings", "modify"), auditFindingRootCauseAnalysisController.CreateAuditFindingRootCauseAnalysis)
	auditRoutes.PUT("/audit-findings/:findingId/root-cause-analysis/:analysisId", middlewares.RolePermissionMiddleware(db, "audit_findings", "modify"), auditFindingRootCauseAnalysisController.UpdateAuditFindingRootCauseAnalysis)
	auditRoutes.DELETE("/audit-findings/:findingId/root-cause-analysis/:analysisId", middlewares.RolePermissionMiddleware(db, "audit_findings", "modify"), auditFindingRootCauseAnalysisController.DeleteAuditFindingRootCauseAnalyses)

	// Audit Finding Action Plan
	auditRoutes.GET("/audit-findings/:findingId/action-plans", middlewares.RolePermissionMiddleware(db, "audit_findings", "view"), auditFindingActionPlanController.GetAllAuditFindingActionPlans)
	auditRoutes.GET("/audit-findings/:findingId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "audit_findings", "view"), auditFindingActionPlanController.GetAuditFindingActionPlanDetails)
	auditRoutes.POST("/audit-findings/:findingId/action-plans", middlewares.RolePermissionMiddleware(db, "audit_findings", "modify"), auditFindingActionPlanController.CreateAuditFindingActionPlan)
	auditRoutes.PUT("/audit-findings/:findingId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "audit_findings", "modify"), auditFindingActionPlanController.UpdateAuditFindingActionPlan)
	auditRoutes.POST("/audit-findings/:findingId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "audit_findings", "modify"), auditFindingActionPlanController.DeleteAuditFindingActionPlans)
	auditRoutes.GET("/audit-findings/:findingId/action-plan-summary", middlewares.RolePermissionMiddleware(db, "audit_findings", "view"), auditFindingActionPlanController.GetAuditFindingActionPlanCountByStatus)

	unauthenticatedRoutes.GET("/audit-findings/:findingId/action-plans/:actionPlanId/status-updates/:id/document-preview", auditFindingActionPlanStatusUpdateController.GetAuditFindingActionPlanStatusUpdateDocumentPreview)
	unauthenticatedRoutes.GET("/audit-findings/:findingId/action-plans/:actionPlanId/status-updates/:id/document-download", auditFindingActionPlanStatusUpdateController.DownloadAuditFindingActionPlanStatusUpdateDocument)
	auditRoutes.GET("/audit-findings/:findingId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "audit_findings", "view"), auditFindingActionPlanStatusUpdateController.GetAllAuditFindingActionPlanStatusUpdates)
	auditRoutes.POST("/audit-findings/:findingId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "audit_findings", "modify"), auditFindingActionPlanStatusUpdateController.CreateAuditFindingActionPlanStatusUpdate)
	auditRoutes.PUT("/audit-findings/:findingId/action-plans/:actionPlanId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "audit_findings", "modify"), auditFindingActionPlanStatusUpdateController.UpdateAuditFindingActionPlanStatusUpdate)
	auditRoutes.GET("/audit-findings/:findingId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "audit_findings", "view"), auditFindingActionPlanStatusUpdateController.GetLatestAuditFindingActionPlanStatusUpdate)

	//Audit Checklist
	auditRoutes.GET("/audits/:auditId/checklists", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditChecklistController.GetAllAuditChecklists)
	auditRoutes.GET("/audits/:auditId/checklists/:id", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditChecklistController.GetAuditChecklistDetails)
	auditRoutes.GET("/audits/:auditId/checklists/:id/draft-findings", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditChecklistController.GetAllAuditChecklistDraftFindings)
	auditRoutes.POST("/audits/:auditId/checklists", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditChecklistController.CreateAuditChecklist)
	auditRoutes.PUT("/audits/:auditId/checklists/:id", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditChecklistController.UpdateAuditChecklist)
	auditRoutes.PUT("/audits/:auditId/checklists/:id/answer", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditChecklistController.AuditChecklistAnswer)
	auditRoutes.POST("audits/:auditId/checklists/delete", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditChecklistController.DeleteAuditChecklists)
	auditRoutes.POST("/audits/:auditId/checklists/delete-answer", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditChecklistController.DeleteAuditChecklistAnswer)
	unauthenticatedRoutes.GET("/audits/:auditId/checklists/:id/document-download", auditChecklistController.DownloadAuditChecklistDocument)

	//Audit Schedule
	auditRoutes.GET("/audits/:auditId/schedules", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditScheduleController.GetAllAuditSchedules)
	auditRoutes.GET("/audits/:auditId/schedules/:scheduleId", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditScheduleController.GetAuditScheduleDetails)
	auditRoutes.POST("/audits/:auditId/schedules", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditScheduleController.CreateAuditSchedule)
	auditRoutes.PUT("/audits/:auditId/schedules/:scheduleId", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditScheduleController.UpdateAuditSchedule)
	auditRoutes.POST("/audits/:auditId/schedules/delete", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditScheduleController.DeleteAuditSchedules)

	auditRoutes.GET("/audits/:auditId/schedules-summary", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditScheduleController.GetScheduleSummary)

	//Audit Auditable items
	auditRoutes.GET("/audits/:auditId/auditable-items", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditController.GetAllAuditableItems)

	//Draft Audit Finding Status
	auditRoutes.GET("/audit-draft-finding-statuses", auditDraftFindingStatusController.GetAllAuditDraftFindingStatuses)
	auditRoutes.GET("/audit-draft-finding-statuses/:id", auditDraftFindingStatusController.GetAuditDraftFindingStatusDetails)

	// Draft Audit Finding
	auditRoutes.GET("audit-draft-findings", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDraftFindingController.GetAllAuditDraftFindings)
	auditRoutes.GET("/audit-draft-findings/:findingId", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDraftFindingController.GetAuditDraftFindingDetails)
	auditRoutes.POST("/audit-draft-findings", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftFindingController.CreateAuditDraftFinding)
	auditRoutes.PUT("/audit-draft-findings/:findingId", middlewares.RolePermissionMiddleware(db, "audits", "modfiy"), auditDraftFindingController.UpdateAuditDraftFinding)
	auditRoutes.POST("/audit-draft-findings/delete", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftFindingController.DeleteAuditDraftFindings)
	auditRoutes.GET("audits/:auditId/audit-draft-finding-summary", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDraftFindingController.GetAuditDraftFindingCountByCategory)
	unauthenticatedRoutes.GET("/audit-draft-findings/:findingId/document-preview", auditDraftFindingController.GetAuditDraftFindingDocumentPreview)
	unauthenticatedRoutes.GET("/audit-draft-findings/:findingId/document-download", auditDraftFindingController.DownloadAuditDraftFindingDocument)

	//Draft Audit Report Status
	auditRoutes.GET("/audit-draft-report-statuses", auditDraftReportStatusController.GetAllAuditDraftReportStatuses)
	auditRoutes.GET("/audit-draft-report-statuses/:id", auditDraftReportStatusController.GetAuditDraftReportStatusDetails)

	//Audit Draft Report
	auditRoutes.GET("/audits/:auditId/draft-reports", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDraftReportController.GetAuditDraftReportDetails)
	auditRoutes.POST("/audits/:auditId/draft-reports", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftReportController.CreateAuditDraftReport)
	auditRoutes.PUT("/audits/:auditId/draft-reports/:id", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftReportController.UpdateAuditDraftReport)
	auditRoutes.POST("/audits/:auditId/draft-reports/delete", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftReportController.DeleteAuditDraftReports)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/auditee-submit", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftReportController.SubmitAuditDraftReport)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/complete", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftReportController.CompleteAuditDraftReport)

	//Audit Risk Control Metrix
	auditRoutes.GET("audits/:auditId/risk-control-matrix", middlewares.RolePermissionMiddleware(db, "audits", "view"), rcmController.GetAllRiskControls)

	// Audit Draft Finding Conversations
	auditRoutes.GET("/audit-draft-findings/:findingId/conversations", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDraftFindingConversationController.GetAllAuditDraftFindingConversations)
	auditRoutes.GET("/audit-draft-findings/:findingId/conversations/:id", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDraftFindingConversationController.GetAuditDraftFindingConversationDetails)
	auditRoutes.POST("/audit-draft-findings/:findingId/conversations", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftFindingConversationController.CreateAuditDraftFindingConversation)
	auditRoutes.PUT("/audit-draft-findings/:findingId/conversations/:id", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftFindingConversationController.UpdateAuditDraftFindingConversation)
	auditRoutes.POST("/audit-draft-findings/:findingId/conversations/delete", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftFindingConversationController.DeleteAuditDraftFindingConversations)
	unauthenticatedRoutes.GET("/audit-draft-findings/:findingId/conversations/:id/document-download", auditDraftFindingConversationController.DownloadAuditDraftFindingConversationDocument)
	unauthenticatedRoutes.GET("/audit-draft-findings/:findingId/conversations/:id/document-preview", auditDraftFindingConversationController.GetAuditDraftFindingConversationDocumentPreview)

	// Audit Draft Finding Action plans
	auditRoutes.GET("/audit-draft-findings/:findingId/action-plans", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDraftFindingActionPlanController.GetAllAuditDraftFindingActionPlans)
	auditRoutes.GET("/audit-draft-findings/:findingId/action-plans/:id", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDraftFindingActionPlanController.GetAuditDraftFindingActionPlanDetails)
	auditRoutes.POST("/audit-draft-findings/:findingId/action-plans", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftFindingActionPlanController.CreateAuditDraftFindingActionPlan)
	auditRoutes.PUT("/audit-draft-findings/:findingId/action-plans/:id", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftFindingActionPlanController.UpdateAuditDraftFindingActionPlan)
	auditRoutes.POST("/audit-draft-findings/:findingId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftFindingActionPlanController.DeleteAuditDraftFindingActionPlans)

	// Audit Schedule status update
	unauthenticatedRoutes.GET("/audits/:auditId/schedules/:scheduleId/status-updates/:id/document-preview", auditScheduleStatusUpdateController.GetAuditScheduleStatusUpdateDocumentPreview)
	unauthenticatedRoutes.GET("/audits/:auditId/schedules/:scheduleId/status-updates/:id/document-download", auditScheduleStatusUpdateController.DownloadAuditScheduleStatusUpdateDocument)
	auditRoutes.GET("/audits/:auditId/schedules/:scheduleId/status-updates", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditScheduleStatusUpdateController.GetAllAuditScheduleStatusUpdates)
	auditRoutes.POST("/audits/:auditId/schedules/:scheduleId/status-updates", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditScheduleStatusUpdateController.CreateAuditScheduleStatusUpdate)
	auditRoutes.PUT("/audits/:auditId/schedules/:scheduleId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditScheduleStatusUpdateController.UpdateAuditScheduleStatusUpdate)
	auditRoutes.GET("/audits/:auditId/schedules/:scheduleId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditScheduleStatusUpdateController.GetLatestAuditScheduleStatusUpdate)

	//Audit Draft Report Review
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/submit", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftReportReviewController.SubmitAuditDraftReport)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/re-submit", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftReportReviewController.ReSubmitAuditDraftReport)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/approve", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftReportReviewController.ApproveAuditDraftReport)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/revert", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftReportReviewController.RevertAuditDraftReport)

	// Audit Draft Report Comments
	auditRoutes.GET("/audits/:auditId/draft-reports/:id/comments", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDraftReportCommentController.GetAllAuditDraftReportComments)
	auditRoutes.GET("/audits/:auditId/draft-reports/:id/comments/:commentId", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDraftReportCommentController.GetAuditDraftReportCommentDetails)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/comments", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftReportCommentController.CreateAuditDraftReportComment)
	auditRoutes.PUT("/audits/:auditId/draft-reports/:id/comments/:commentId", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftReportCommentController.UpdateAuditDraftReportComment)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/comments/delete", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftReportCommentController.DeleteAuditDraftReportComments)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/comments/:commentId/status-update", middlewares.RolePermissionMiddleware(db, "audits", "modify"), auditDraftReportCommentController.StatusUpdateAuditDraftReportComment)
	unauthenticatedRoutes.GET("/audits/:auditId/draft-reports/:id/comments/:commentId/document-download", auditDraftReportCommentController.DownloadAuditDraftReportCommentDocument)
	unauthenticatedRoutes.GET("/audits/:auditId/draft-reports/:id/comments/:commentId/document-preview", auditDraftReportCommentController.GetAuditDraftReportCommentDocumentPreview)
	auditRoutes.GET("/audits/:auditId/draft-reports/:id/review-history", middlewares.RolePermissionMiddleware(db, "audits", "view"), auditDraftReportReviewController.GetAllAuditDraftReportReviews)

	auditRoutes.GET("/audit-plan-commencement-letter-statuses", auditPlanCommencementLetterStatusController.GetAllStatuses)
	auditRoutes.GET("/audit-plan-commencement-letter-statuses/:id", auditPlanCommencementLetterStatusController.GetStatusDetails)

	auditRoutes.GET("/audit-plans/:planId/commencement-letters", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanCommencementLetterController.GetAllAuditPlanCommencementLetters)
	auditRoutes.GET("/audit-plans/:planId/commencement-letters/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "view"), auditPlanCommencementLetterController.GetAuditPlanCommencementLetterDetails)
	auditRoutes.POST("/audit-plans/:planId/commencement-letters", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanCommencementLetterController.CreateAuditPlanCommencementLetters)
	auditRoutes.PUT("/audit-plans/:planId/commencement-letters/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanCommencementLetterController.UpdateAuditPlanCommencementLetter)
	auditRoutes.DELETE("/audit-plans/:planId/commencement-letters/:id", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanCommencementLetterController.DeleteAuditPlanCommencementLetter)
	auditRoutes.PUT("/audit-plans/:planId/commencement-letters/:id/send", middlewares.RolePermissionMiddleware(db, "audit_plans", "modify"), auditPlanCommencementLetterController.SendAuditPlanCommencementLetter)

	//Audit Final Report Status
	auditRoutes.GET("/audit-final-report-statuses", auditFinalReportStatusController.GetAllAuditFinalReportStatuses)
	auditRoutes.GET("/audit-final-report-statuses/:id", auditFinalReportStatusController.GetAuditFinalReportStatusDetails)

	//Audit Final Reports
	auditRoutes.GET("/audits/:auditId/final-reports", auditFinalReportController.GetAuditFinalReportDetails)
	auditRoutes.POST("/audits/:auditId/final-reports", auditFinalReportController.CreateAuditFinalReport)
	auditRoutes.POST("/audits/:auditId/final-reports/delete", auditFinalReportController.DeleteAuditFinalReport)
	auditRoutes.GET("/audits/:auditId/final-reports/last-edit-summary", auditFinalReportController.GetLastEditSummary)
	auditRoutes.POST("/audits/:auditId/final-reports/publish", auditFinalReportController.PublishAuditFinalReport)

	//Audit Final Report Section
	auditRoutes.GET("/audits/:auditId/final-report-sections/:id", auditFinalReportSectionController.GetAuditFinalReportSectionDetails)
	auditRoutes.PUT("/audits/:auditId/final-report-sections/:id", auditFinalReportSectionController.UpdateAuditFinalReportSection)
	auditRoutes.DELETE("/audits/:auditId/final-report-sections/:id", auditFinalReportSectionController.DeleteAuditFinalReportSection)

	//workflows
	auditRoutes.GET("/audits/:auditId/final-reports/workflow-list", auditReportWorkflowController.GetWorkflowList)

	auditRoutes.POST("/audits/:auditId/final-reports/submit", auditReportReviewController.SubmitAuditFinalReport)
	auditRoutes.GET("/audits/:auditId/final-reports/comments", auditReportReviewController.GetComments)
	auditRoutes.POST("/audits/:auditId/final-reports/comments", auditReportReviewController.AddComment)

}
