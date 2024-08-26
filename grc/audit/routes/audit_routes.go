package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	auditRoutes := router.Group("/")

	auditRoutes.Use(middlewares.AuthMiddleware())

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
	auditRoutes.GET("/audit-findings/:auditId", auditFindingController.GetAuditFindingDetails)
	auditRoutes.POST("/audit-findings", auditFindingController.CreateAuditFinding)
	auditRoutes.PUT("/audit-findings/:auditId", auditFindingController.UpdateAuditFinding)
	auditRoutes.POST("/audit-findings/delete", auditFindingController.DeleteAuditFindings)

}
