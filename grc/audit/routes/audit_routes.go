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
	auditPlanCommencementLetterController := app.AuditModule.AuditPlanCommencementLetterController
	auditFinalReportStatusController := app.AuditModule.AuditFinalReportStatusController
	auditFinalReportController := app.AuditModule.AuditFinalReportController

	//Audit Status
	auditRoutes.GET("/audit-statuses", auditStatusController.GetAllAuditStatuses)
	auditRoutes.GET("/audit-statuses/:id", auditStatusController.GetAuditStatusDetails)

	//Auditor Risk Rating
	auditRoutes.GET("/auditor-risk-ratings", auditorRiskRatingController.GetAllAuditorRiskRatings)
	auditRoutes.GET("/auditor-risk-ratings/:id", auditorRiskRatingController.GetAuditorRiskRatingDetails)

	//Audit Universe
	auditRoutes.POST("/audit-universe/mark-as-auditable-item", auditUniverseController.MarkAsAuditableItem)
	auditRoutes.POST("/audit-universe/unmark-as-auditable-item", auditUniverseController.UnmarkAsAuditableItem)

	//Auditable Item
	auditRoutes.GET("/auditable-items", auditableItemController.GetAllAuditableItems)
	auditRoutes.GET("/auditable-items/:id", auditableItemController.GetAuditableItemDetails)
	auditRoutes.POST("/auditable-items", auditableItemController.CreateAuditableItem)
	auditRoutes.PUT("/auditable-items/:id", auditableItemController.UpdateAuditableItem)
	auditRoutes.POST("/auditable-items/delete", auditableItemController.DeleteAuditableItems)

	//Audit Team
	auditRoutes.GET("/audit-teams", auditTeamController.GetAllAuditTeams)
	auditRoutes.GET("/audit-teams/:id", auditTeamController.GetAuditTeamDetails)
	auditRoutes.POST("/audit-teams", auditTeamController.CreateAuditTeam)
	auditRoutes.PUT("/audit-teams/:id", auditTeamController.UpdateAuditTeam)
	auditRoutes.POST("/audit-teams/delete", auditTeamController.DeleteAuditTeams)

	//Audit Category
	auditRoutes.GET("/audit-categories", auditCategoryController.GetAllAuditCategories)
	auditRoutes.GET("/audit-categories/:id", auditCategoryController.GetAuditCategoryDetails)
	auditRoutes.POST("/audit-categories", auditCategoryController.CreateAuditCategory)
	auditRoutes.PUT("/audit-categories/:id", auditCategoryController.UpdateAuditCategory)
	auditRoutes.POST("/audit-categories/delete", auditCategoryController.DeleteAuditCategories)
	auditRoutes.PUT("/audit-categories/activate", auditCategoryController.ActivateAuditCategories)
	auditRoutes.PUT("/audit-categories/deactivate", auditCategoryController.DeactivateAuditCategories)

	//AuditPlan
	auditRoutes.GET("/audit-plans", auditPlanController.GetAllAuditPlans)
	auditRoutes.GET("/audit-plans/:planId", auditPlanController.GetAuditPlanDetails)
	auditRoutes.POST("/audit-plans", auditPlanController.CreateAuditPlan)
	auditRoutes.PUT("/audit-plans/:planId", auditPlanController.UpdateAuditPlan)
	auditRoutes.POST("/audit-plans/delete", auditPlanController.DeleteAuditPlans)
	auditRoutes.POST("/audit-plans/:planId/publish", auditPlanController.PublishAuditPlan)

	//AuditPlan Checklist
	auditRoutes.GET("/audit-plans/:planId/checklists", auditPlanChecklistController.GetAllAuditPlanChecklists)
	auditRoutes.GET("/audit-plans/:planId/checklists/:id", auditPlanChecklistController.GetAuditPlanChecklistDetails)
	auditRoutes.POST("/audit-plans/:planId/checklists", auditPlanChecklistController.CreateAuditPlanChecklists)
	auditRoutes.PUT("/audit-plans/:planId/checklists/:id", auditPlanChecklistController.UpdateAuditPlanChecklist)
	auditRoutes.POST("/audit-plans/:planId/checklists/delete", auditPlanChecklistController.DeleteAuditPlanChecklists)

	//AuditPlan Tools
	auditRoutes.GET("/audit-plans/:planId/tools", auditPlanToolController.GetAllAuditPlanTools)
	auditRoutes.GET("/audit-plans/:planId/tools/:id", auditPlanToolController.GetAuditPlanToolDetails)
	auditRoutes.POST("/audit-plans/:planId/tools", auditPlanToolController.CreateAuditPlanTool)
	auditRoutes.PUT("/audit-plans/:planId/tools/:id", auditPlanToolController.UpdateAuditPlanTool)
	auditRoutes.POST("/audit-plans/:planId/tools/delete", auditPlanToolController.DeleteAuditPlanTools)

	//AuditPlan Auditable Item
	auditRoutes.GET("/audit-plans/:planId/auditable-items", auditPlanAuditableItemController.GetAllAuditPlanAuditableItems)
	auditRoutes.GET("/audit-plans/:planId/auditable-items/:id", auditPlanAuditableItemController.GetAllAuditPlanAuditableItemDetails)
	auditRoutes.POST("/audit-plans/:planId/auditable-items", auditPlanAuditableItemController.CreateAuditPlanAuditableItem)
	auditRoutes.POST("/audit-plans/:planId/auditable-items/delete", auditPlanAuditableItemController.DeleteAuditPlanAuditableItems)

	//Audit Activity Type
	auditRoutes.GET("/audit-activity-types", auditActivityTypeController.GetAllAuditActivityTypes)
	auditRoutes.GET("/audit-activity-types/:id", auditActivityTypeController.GetAuditActivityTypeDetails)
	auditRoutes.POST("/audit-activity-types", auditActivityTypeController.CreateAuditActivityType)
	auditRoutes.PUT("/audit-activity-types/:id", auditActivityTypeController.UpdateAuditActivityType)
	auditRoutes.POST("/audit-activity-types/delete", auditActivityTypeController.DeleteAuditActivityTypes)
	auditRoutes.PUT("/audit-activity-types/activate", auditActivityTypeController.ActivateAuditActivityTypes)
	auditRoutes.PUT("/audit-activity-types/deactivate", auditActivityTypeController.DeactivateAuditActivityTypes)

	//Audit Schedule Statuses
	auditRoutes.GET("/audit-schedule-statuses", auditScheduleStatusController.GetAllAuditScheduleStatuses)
	auditRoutes.GET("/audit-schedule-statuses/:id", auditScheduleStatusController.GetAuditScheduleStatusDetails)

	//AuditPlan Tools
	auditRoutes.GET("/audit-plans/:planId/budgets", auditPlanBudgetController.GetAllAuditPlanBudgets)
	auditRoutes.GET("/audit-plans/:planId/budgets/:id", auditPlanBudgetController.GetAuditPlanBudgetDetails)
	auditRoutes.POST("/audit-plans/:planId/budgets", auditPlanBudgetController.CreateAuditPlanBudget)
	auditRoutes.PUT("/audit-plans/:planId/budgets/:id", auditPlanBudgetController.UpdateAuditPlanBudget)
	auditRoutes.POST("/audit-plans/:planId/budgets/delete", auditPlanBudgetController.DeleteAuditPlanBudgets)

	//Audit Schedule
	auditRoutes.GET("/audit-plans/:planId/schedules", auditPlanScheduleController.GetAllAuditPlanSchedules)
	auditRoutes.GET("/audit-plans/:planId/schedules/:id", auditPlanScheduleController.GetAuditPlanScheduleDetails)
	auditRoutes.POST("/audit-plans/:planId/schedules", auditPlanScheduleController.CreateAuditPlanSchedule)
	auditRoutes.PUT("/audit-plans/:planId/schedules/:id", auditPlanScheduleController.UpdateAuditPlanSchedule)
	auditRoutes.POST("/audit-plans/:planId/schedules/delete", auditPlanScheduleController.DeleteAuditPlanSchedules)

	//Audit Engagement Type
	auditRoutes.GET("/audit-engagement-types", auditEngagementTypeController.GetAllAuditEngagementTypes)
	auditRoutes.GET("/audit-engagement-types/:id", auditEngagementTypeController.GetAuditEngagementTypeDetails)
	auditRoutes.POST("/audit-engagement-types", auditEngagementTypeController.CreateAuditEngagementType)
	auditRoutes.PUT("/audit-engagement-types/:id", auditEngagementTypeController.UpdateAuditEngagementType)
	auditRoutes.POST("/audit-engagement-types/delete", auditEngagementTypeController.DeleteAuditEngagementTypes)
	auditRoutes.PUT("/audit-engagement-types/activate", auditEngagementTypeController.ActivateAuditEngagementTypes)
	auditRoutes.PUT("/audit-engagement-types/deactivate", auditEngagementTypeController.DeactivateAuditEngagementTypes)

	//Audit Plan Stakeholder Engagement
	auditRoutes.GET("/audit-plans/:planId/stakeholder-engagements", auditPlanStakeholderEngagementController.GetAllAuditPlanStakeholderEngagements)
	auditRoutes.GET("/audit-plans/:planId/stakeholder-engagements/:id", auditPlanStakeholderEngagementController.GetAuditPlanStakeholderEngagementDetails)
	auditRoutes.POST("/audit-plans/:planId/stakeholder-engagements", auditPlanStakeholderEngagementController.CreateAuditPlanStakeholderEngagement)
	auditRoutes.PUT("/audit-plans/:planId/stakeholder-engagements/:id", auditPlanStakeholderEngagementController.UpdateAuditPlanStakeholderEngagement)
	auditRoutes.POST("/audit-plans/:planId/stakeholder-engagements/delete", auditPlanStakeholderEngagementController.DeleteAuditPlanStakeholderEngagements)

	//Audit
	auditRoutes.GET("/audits", auditController.GetAllAudits)
	auditRoutes.GET("/audits/:auditId", auditController.GetAuditDetails)
	auditRoutes.POST("/audits", auditController.CreateAudit)
	auditRoutes.PUT("/audits/:auditId", auditController.UpdateAudit)
	auditRoutes.POST("/audits/delete", auditController.DeleteAudits)

	//Audit Schedule
	auditRoutes.GET("/audits/:auditId/workpapers", auditWorkpaperController.GetAllAuditWorkpapers)
	auditRoutes.GET("/audits/:auditId/workpapers/:id", auditWorkpaperController.GetAuditWorkpaperDetails)
	auditRoutes.POST("/audits/:auditId/workpapers", auditWorkpaperController.CreateAuditWorkpaper)
	auditRoutes.PUT("/audits/:auditId/workpapers/:id", auditWorkpaperController.UpdateAuditWorkpaper)
	auditRoutes.POST("/audits/:auditId/workpapers/delete", auditWorkpaperController.DeleteAuditWorkpapers)

	unauthenticatedRoutes.GET("/audits/:auditId/workpapers/document-download", auditWorkpaperController.DownloadWorkpaperDocument)

	//Audit Meeting
	auditRoutes.GET("/audits/:auditId/meetings", auditMeetingController.GetAllAuditMeetings)
	auditRoutes.GET("/audits/:auditId/meetings/:meetingId", auditMeetingController.GetAuditMeetingDetails)
	auditRoutes.POST("/audits/:auditId/meetings", auditMeetingController.CreateAuditMeeting)
	auditRoutes.PUT("/audits/:auditId/meetings/:meetingId", auditMeetingController.UpdateAuditMeeting)
	auditRoutes.POST("/audits/:auditId/meetings/delete", auditMeetingController.DeleteAuditMeetings)

	//Audit timeline
	auditRoutes.GET("/audits/:auditId/timelines", auditTimelineController.GetAllAuditTimelines)

	//Audit Status
	auditRoutes.GET("/audit-finding-statuses", auditFindingStatusController.GetAllAuditFindingStatuses)
	auditRoutes.GET("/audit-finding-statuses/:id", auditFindingStatusController.GetAuditFindingStatusDetails)

	//Audit Finding Category
	auditRoutes.GET("/audit-finding-categories", auditFindingCategoryController.GetAllAuditFindingCategories)
	auditRoutes.GET("/audit-finding-categories/:id", auditFindingCategoryController.GetAuditFindingCategoryDetails)
	auditRoutes.PUT("/audit-finding-categories/:id", auditFindingCategoryController.UpdateAuditFindingCategory)
	auditRoutes.POST("/audit-finding-categories/delete", auditFindingCategoryController.DeleteAuditFindingCategories)

	// Audit Finding
	auditRoutes.GET("/audit-findings", auditFindingController.GetAllAuditFindings)
	auditRoutes.GET("/audit-findings/:findingId", auditFindingController.GetAuditFindingDetails)
	auditRoutes.POST("/audit-findings", auditFindingController.CreateAuditFinding)
	auditRoutes.PUT("/audit-findings/:findingId", auditFindingController.UpdateAuditFinding)
	auditRoutes.POST("/audit-findings/delete", auditFindingController.DeleteAuditFindings)
	auditRoutes.GET("audits/:auditId/audit-finding-summary", auditFindingController.GetAuditFindingCountByCategory)

	//Audit Document  Request Status
	auditRoutes.GET("/audit-document-request-statuses", auditDocumentRequestStatusController.GetAllAuditDocumentRequestStatuses)

	// Audit Document Request
	auditRoutes.GET("/audits/:auditId/document-requests", auditDocumentRequestController.GetAllAuditDocumentRequests)
	auditRoutes.GET("/audits/:auditId/document-requests/:id", auditDocumentRequestController.GetAuditDocumentRequestDetails)
	auditRoutes.POST("/audits/:auditId/document-requests", auditDocumentRequestController.CreateAuditDocumentRequest)
	auditRoutes.PUT("/audits/:auditId/document-requests/:id", auditDocumentRequestController.UpdateAuditDocumentRequest)
	auditRoutes.POST("/audits/:auditId/document-requests/delete", auditDocumentRequestController.DeleteAuditDocumentRequests)
	auditRoutes.POST("/audits/:auditId/document-requests/:id/submit", auditDocumentRequestController.SubmitAuditDocumentRequest)

	// Audit Document Request Conversations
	auditRoutes.GET("/audits/:auditId/document-requests/:id/conversations", auditDocumentRequestConversationController.GetAllAuditDocumentRequestConversations)
	auditRoutes.POST("/audits/:auditId/document-requests/:id/conversations", auditDocumentRequestConversationController.CreateAuditDocumentRequestConversation)
	unauthenticatedRoutes.GET("/audits/:auditId/document-requests/:id/conversations/:conversationId/document-download", auditDocumentRequestConversationController.DownloadDocumentRequestConversationDocument)

	auditRoutes.GET("/audit-findings/:findingId/root-cause-analysis", auditFindingRootCauseAnalysisController.GetAllAuditFindingRootCauseAnalyses)
	auditRoutes.POST("/audit-findings/:findingId/root-cause-analysis", auditFindingRootCauseAnalysisController.CreateAuditFindingRootCauseAnalysis)
	auditRoutes.PUT("/audit-findings/:findingId/root-cause-analysis/:analysisId", auditFindingRootCauseAnalysisController.UpdateAuditFindingRootCauseAnalysis)
	auditRoutes.DELETE("/audit-findings/:findingId/root-cause-analysis/:analysisId", auditFindingRootCauseAnalysisController.DeleteAuditFindingRootCauseAnalyses)

	// Audit Finding Action Plan
	auditRoutes.GET("/audit-findings/:findingId/action-plans", auditFindingActionPlanController.GetAllAuditFindingActionPlans)
	auditRoutes.GET("/audit-findings/:findingId/action-plans/:actionPlanId", auditFindingActionPlanController.GetAuditFindingActionPlanDetails)
	auditRoutes.POST("/audit-findings/:findingId/action-plans", auditFindingActionPlanController.CreateAuditFindingActionPlan)
	auditRoutes.PUT("/audit-findings/:findingId/action-plans/:actionPlanId", auditFindingActionPlanController.UpdateAuditFindingActionPlan)
	auditRoutes.POST("/audit-findings/:findingId/action-plans/delete", auditFindingActionPlanController.DeleteAuditFindingActionPlans)
	auditRoutes.GET("/audit-findings/:findingId/action-plan-summary", auditFindingActionPlanController.GetAuditFindingActionPlanCountByStatus)

	auditRoutes.GET("/audit-findings/:findingId/action-plans/:actionPlanId/status-updates", auditFindingActionPlanStatusUpdateController.GetAllAuditFindingActionPlanStatusUpdates)
	auditRoutes.POST("/audit-findings/:findingId/action-plans/:actionPlanId/status-updates", auditFindingActionPlanStatusUpdateController.CreateAuditFindingActionPlanStatusUpdate)
	auditRoutes.PUT("/audit-findings/:findingId/action-plans/:actionPlanId/status-updates/:id", auditFindingActionPlanStatusUpdateController.UpdateAuditFindingActionPlanStatusUpdate)
	auditRoutes.GET("/audit-findings/:findingId/action-plans/:actionPlanId/status-updates/latest", auditFindingActionPlanStatusUpdateController.GetLatestAuditFindingActionPlanStatusUpdate)

	//Audit Checklist
	auditRoutes.GET("/audits/:auditId/checklists", auditChecklistController.GetAllAuditChecklists)
	auditRoutes.GET("/audits/:auditId/checklists/:id", auditChecklistController.GetAuditChecklistDetails)
	auditRoutes.GET("/audits/:auditId/checklists/:id/draft-findings", auditChecklistController.GetAllAuditChecklistDraftFindings)
	auditRoutes.POST("/audits/:auditId/checklists", auditChecklistController.CreateAuditChecklist)
	auditRoutes.PUT("/audits/:auditId/checklists/:id", auditChecklistController.UpdateAuditChecklist)
	auditRoutes.PUT("/audits/:auditId/checklists/:id/answer", auditChecklistController.AuditChecklistAnswer)
	auditRoutes.POST("audits/:auditId/checklists/delete", auditChecklistController.DeleteAuditChecklists)
	auditRoutes.POST("/audits/:auditId/checklists/delete-answer", auditChecklistController.DeleteAuditChecklistAnswer)

	//Audit Schedule
	auditRoutes.GET("/audits/:auditId/schedules", auditScheduleController.GetAllAuditSchedules)
	auditRoutes.GET("/audits/:auditId/schedules/:scheduleId", auditScheduleController.GetAuditScheduleDetails)
	auditRoutes.POST("/audits/:auditId/schedules", auditScheduleController.CreateAuditSchedule)
	auditRoutes.PUT("/audits/:auditId/schedules/:scheduleId", auditScheduleController.UpdateAuditSchedule)
	auditRoutes.POST("/audits/:auditId/schedules/delete", auditScheduleController.DeleteAuditSchedules)

	auditRoutes.GET("/audits/:auditId/schedules-summary", auditScheduleController.GetScheduleSummary)

	//Audit Auditable items
	auditRoutes.GET("/audits/:auditId/auditable-items", auditController.GetAllAuditableItems)

	//Draft Audit Finding Status
	auditRoutes.GET("/audit-draft-finding-statuses", auditDraftFindingStatusController.GetAllAuditDraftFindingStatuses)
	auditRoutes.GET("/audit-draft-finding-statuses/:id", auditDraftFindingStatusController.GetAuditDraftFindingStatusDetails)

	// Draft Audit Finding
	auditRoutes.GET("audit-draft-findings", auditDraftFindingController.GetAllAuditDraftFindings)
	auditRoutes.GET("/audit-draft-findings/:findingId", auditDraftFindingController.GetAuditDraftFindingDetails)
	auditRoutes.POST("/audit-draft-findings", auditDraftFindingController.CreateAuditDraftFinding)
	auditRoutes.PUT("/audit-draft-findings/:findingId", auditDraftFindingController.UpdateAuditDraftFinding)
	auditRoutes.POST("/audit-draft-findings/delete", auditDraftFindingController.DeleteAuditDraftFindings)
	auditRoutes.GET("audits/:auditId/audit-draft-finding-summary", auditDraftFindingController.GetAuditDraftFindingCountByCategory)

	//Draft Audit Report Status
	auditRoutes.GET("/audit-draft-report-statuses", auditDraftReportStatusController.GetAllAuditDraftReportStatuses)
	auditRoutes.GET("/audit-draft-report-statuses/:id", auditDraftReportStatusController.GetAuditDraftReportStatusDetails)

	//Audit Draft Report
	auditRoutes.GET("/audits/:auditId/draft-reports", auditDraftReportController.GetAuditDraftReportDetails)
	auditRoutes.POST("/audits/:auditId/draft-reports", auditDraftReportController.CreateAuditDraftReport)
	auditRoutes.PUT("/audits/:auditId/draft-reports/:id", auditDraftReportController.UpdateAuditDraftReport)
	auditRoutes.POST("/audits/:auditId/draft-reports/delete", auditDraftReportController.DeleteAuditDraftReports)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/auditee-submit", auditDraftReportController.SubmitAuditDraftReport)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/complete", auditDraftReportController.CompleteAuditDraftReport)

	//Audit Risk Control Metrix
	auditRoutes.GET("audits/:auditId/risk-control-matrix", rcmController.GetAllRiskControls)

	// Audit Draft Finding Conversations
	auditRoutes.GET("/audit-draft-findings/:findingId/conversations", auditDraftFindingConversationController.GetAllAuditDraftFindingConversations)
	auditRoutes.GET("/audit-draft-findings/:findingId/conversations/:id", auditDraftFindingConversationController.GetAuditDraftFindingConversationDetails)
	auditRoutes.POST("/audit-draft-findings/:findingId/conversations", auditDraftFindingConversationController.CreateAuditDraftFindingConversation)
	auditRoutes.PUT("/audit-draft-findings/:findingId/conversations/:id", auditDraftFindingConversationController.UpdateAuditDraftFindingConversation)
	auditRoutes.POST("/audit-draft-findings/:findingId/conversations/delete", auditDraftFindingConversationController.DeleteAuditDraftFindingConversations)
	unauthenticatedRoutes.GET("/audit-draft-findings/:findingId/conversations/:id/document-download", auditDraftFindingConversationController.DownloadAuditDraftFindingConversationDocument)

	// Audit Draft Finding Action plans
	auditRoutes.GET("/audit-draft-findings/:findingId/action-plans", auditDraftFindingActionPlanController.GetAllAuditDraftFindingActionPlans)
	auditRoutes.GET("/audit-draft-findings/:findingId/action-plans/:id", auditDraftFindingActionPlanController.GetAuditDraftFindingActionPlanDetails)
	auditRoutes.POST("/audit-draft-findings/:findingId/action-plans", auditDraftFindingActionPlanController.CreateAuditDraftFindingActionPlan)
	auditRoutes.PUT("/audit-draft-findings/:findingId/action-plans/:id", auditDraftFindingActionPlanController.UpdateAuditDraftFindingActionPlan)
	auditRoutes.POST("/audit-draft-findings/:findingId/action-plans/delete", auditDraftFindingActionPlanController.DeleteAuditDraftFindingActionPlans)

	// Audit Schedule status update
	auditRoutes.GET("/audits/:auditId/schedules/:scheduleId/status-updates", auditScheduleStatusUpdateController.GetAllAuditScheduleStatusUpdates)
	auditRoutes.POST("/audits/:auditId/schedules/:scheduleId/status-updates", auditScheduleStatusUpdateController.CreateAuditScheduleStatusUpdate)
	auditRoutes.PUT("/audits/:auditId/schedules/:scheduleId/status-updates/:id", auditScheduleStatusUpdateController.UpdateAuditScheduleStatusUpdate)
	auditRoutes.GET("/audits/:auditId/schedules/:scheduleId/status-updates/latest", auditScheduleStatusUpdateController.GetLatestAuditScheduleStatusUpdate)

	//Audit Draft Report Review
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/submit", auditDraftReportReviewController.SubmitAuditDraftReport)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/re-submit", auditDraftReportReviewController.ReSubmitAuditDraftReport)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/approve", auditDraftReportReviewController.ApproveAuditDraftReport)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/revert", auditDraftReportReviewController.RevertAuditDraftReport)

	// Audit Draft Report Comments
	auditRoutes.GET("/audits/:auditId/draft-reports/:id/comments", auditDraftReportCommentController.GetAllAuditDraftReportComments)
	auditRoutes.GET("/audits/:auditId/draft-reports/:id/comments/:commentId", auditDraftReportCommentController.GetAuditDraftReportCommentDetails)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/comments", auditDraftReportCommentController.CreateAuditDraftReportComment)
	auditRoutes.PUT("/audits/:auditId/draft-reports/:id/comments/:commentId", auditDraftReportCommentController.UpdateAuditDraftReportComment)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/comments/delete", auditDraftReportCommentController.DeleteAuditDraftReportComments)
	auditRoutes.POST("/audits/:auditId/draft-reports/:id/comments/:commentId/status-update", auditDraftReportCommentController.StatusUpdateAuditDraftReportComment)
	unauthenticatedRoutes.GET("/audits/:auditId/draft-reports/:id/comments/:commentId/document-download", auditDraftReportCommentController.DownloadAuditDraftReportCommentDocument)
	auditRoutes.GET("/audits/:auditId/draft-reports/:id/review-history", auditDraftReportReviewController.GetAllAuditDraftReportReviews)

	auditRoutes.GET("/audit-plans/:planId/commencement-letters", auditPlanCommencementLetterController.GetAllAuditPlanCommencementLetters)
	auditRoutes.GET("/audit-plans/:planId/commencement-letters/:id", auditPlanCommencementLetterController.GetAuditPlanCommencementLetterDetails)
	auditRoutes.POST("/audit-plans/:planId/commencement-letters", auditPlanCommencementLetterController.CreateAuditPlanCommencementLetters)
	auditRoutes.PUT("/audit-plans/:planId/commencement-letters/:id", auditPlanCommencementLetterController.UpdateAuditPlanCommencementLetter)
	auditRoutes.DELETE("/audit-plans/:planId/commencement-letters/:id", auditPlanCommencementLetterController.DeleteAuditPlanCommencementLetter)
	auditRoutes.PUT("/audit-plans/:planId/commencement-letters/:id/send", auditPlanCommencementLetterController.SendAuditPlanCommencementLetter)
	
	//Audit Final Report Status
	auditRoutes.GET("/audit-final-report-statuses", auditFinalReportStatusController.GetAllAuditFinalReportStatuses)
	auditRoutes.GET("/audit-final-report-statuses/:id", auditFinalReportStatusController.GetAuditFinalReportStatusDetails)

	//Audit Final Reports
	auditRoutes.POST("/audits/:auditId/final-reports", auditFinalReportController.CreateAuditFinalReport)

}
