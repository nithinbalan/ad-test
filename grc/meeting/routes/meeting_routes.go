package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	meetingRoutes := router.Group("/")

	meetingRoutes.Use(middlewares.AuthMiddleware())

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	meetingFrequencyController := app.MeetingModule.MeetingFrequencyController
	meetingStatusController := app.MeetingModule.MeetingStatusController
	meetingTypeController := app.MeetingModule.MeetingTypeController
	meetingController := app.MeetingModule.MeetingController
	meetingActionPlanController := app.MeetingModule.MeetingActionPlanController
	meetingActionPlanStatusUpdateController := app.MeetingModule.MeetingActionPlanStatusUpdateController
	meetingAgendaController := app.MeetingModule.MeetingAgendaController
	meetingMinuteController := app.MeetingModule.MeetingMinuteController
	meetingAssociatedItemController := app.MeetingModule.MeetingAssociatedItemController
	meetingMinuteStatusController := app.MeetingModule.MeetingMinuteStatusController

	meetingRoutes.GET("/meeting-frequencies", meetingFrequencyController.GetAllMeetingFrequencies)
	meetingRoutes.GET("/meeting-frequencies/:id", meetingFrequencyController.GetMeetingFrequencyDetails)

	// Meeting Status
	meetingRoutes.GET("/meeting-statuses", meetingStatusController.GetAllMeetingStatuses)
	meetingRoutes.GET("/meeting-statuses/:id", meetingStatusController.GetMeetingStatusDetails)

	// Meeting Type
	meetingRoutes.GET("/meeting-types", meetingTypeController.GetAllMeetingTypes)
	meetingRoutes.GET("/meeting-types/:id", meetingTypeController.GetMeetingTypeDetails)
	meetingRoutes.POST("/meeting-types", middlewares.RolePermissionMiddleware(db, "meeting_setting", "modify"), meetingTypeController.CreateMeetingType)
	meetingRoutes.PUT("/meeting-types/:id", middlewares.RolePermissionMiddleware(db, "meeting_setting", "modify"), meetingTypeController.UpdateMeetingType)
	meetingRoutes.POST("/meeting-types/delete", middlewares.RolePermissionMiddleware(db, "meeting_setting", "modify"), meetingTypeController.DeleteMeetingTypes)
	meetingRoutes.PUT("/meeting-types/activate", middlewares.RolePermissionMiddleware(db, "meeting_setting", "modify"), meetingTypeController.ActivateMeetingTypes)
	meetingRoutes.PUT("/meeting-types/deactivate", middlewares.RolePermissionMiddleware(db, "meeting_setting", "modify"), meetingTypeController.DeactivateMeetingTypes)

	// Meetings
	meetingRoutes.GET("/meetings/export", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingController.GetMeetingExport)
	meetingRoutes.GET("/meetings", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingController.GetAllMeetings)
	meetingRoutes.GET("/meetings/:id", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingController.GetMeetingDetails)
	meetingRoutes.POST("/meetings", middlewares.RolePermissionMiddleware(db, "meeting", "modify"), meetingController.CreateMeeting)
	meetingRoutes.PUT("/meetings/:id", middlewares.RolePermissionMiddleware(db, "meeting", "modify"), meetingController.UpdateMeeting)
	meetingRoutes.POST("/meetings/delete", middlewares.RolePermissionMiddleware(db, "meeting", "modify"), meetingController.DeleteMeetings)
	meetingRoutes.GET("/meetings/:id/meeting-summary", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingController.GetAllCounts)
	meetingRoutes.GET("/meetings/:id/effectiveness", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingController.GetMeetingEffectiveness)

	unauthenticatedRoutes.GET("/meetings/:id/document-download", meetingController.DownloadMeetingDocument)
	// Meeting action plan
	meetingRoutes.GET("/meeting-action-plans/export", middlewares.RolePermissionMiddleware(db, "meeting_action_plan", "view"), meetingActionPlanController.GetMeetingActionPlanExport)
	meetingRoutes.GET("/meeting-action-plans", middlewares.RolePermissionMiddleware(db, "meeting_action_plan", "view"), meetingActionPlanController.GetAllMeetingActionPlans)
	meetingRoutes.GET("/meetings/:id/action-plans", middlewares.RolePermissionMiddleware(db, "meeting_action_plan", "view"), meetingActionPlanController.GetAllMeetingActionPlans)
	meetingRoutes.GET("/meetings/:id/action-plans/:meetingActionPlanId", middlewares.RolePermissionMiddleware(db, "meeting_action_plan", "view"), meetingActionPlanController.GetMeetingActionPlanDetails)
	meetingRoutes.POST("/meetings/:id/action-plans", middlewares.RolePermissionMiddleware(db, "meeting_action_plan", "modify"), meetingActionPlanController.CreateMeetingActionPlan)
	meetingRoutes.PUT("/meetings/:id/action-plans/:meetingActionPlanId", middlewares.RolePermissionMiddleware(db, "meeting_action_plan", "modify"), meetingActionPlanController.UpdateMeetingActionPlan)
	meetingRoutes.POST("/meetings/:id/action-plans/delete", middlewares.RolePermissionMiddleware(db, "meeting_action_plan", "modify"), meetingActionPlanController.DeleteMeetingActionPlans)

	meetingRoutes.GET("/meetings/:id/action-plan-summary", middlewares.RolePermissionMiddleware(db, "meeting_action_plan", "view"), meetingActionPlanController.GetMeetingActionPlanSummary)

	// Meeting action plan status updates
	meetingRoutes.GET("/meetings/:id/action-plans/:meetingActionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "meeting_action_plan", "view"), meetingActionPlanStatusUpdateController.GetAllMeetingActionPlanStatusUpdates)
	meetingRoutes.GET("/meetings/:id/action-plans/:meetingActionPlanId/status-updates/:meetingActionPlanUpdateId", middlewares.RolePermissionMiddleware(db, "meeting_action_plan", "view"), meetingActionPlanStatusUpdateController.GetMeetingActionPlanStatusUpdateDetails)
	meetingRoutes.POST("/meetings/:id/action-plans/:meetingActionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "meeting_action_plan", "modify"), meetingActionPlanStatusUpdateController.CreateMeetingActionPlanStatusUpdate)
	meetingRoutes.PUT("/meetings/:id/action-plans/:meetingActionPlanId/status-updates/:meetingActionPlanUpdateId", middlewares.RolePermissionMiddleware(db, "meeting_action_plan", "modify"), meetingActionPlanStatusUpdateController.UpdateMeetingActionPlanStatusUpdate)
	meetingRoutes.GET("/meetings/:id/action-plans/:meetingActionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "meeting_action_plan", "view"), meetingActionPlanStatusUpdateController.GetLatestMeetingActionPlanStatusUpdate)
	// Preview and download
	unauthenticatedRoutes.GET("/meetings/:id/action-plans/:meetingActionPlanId/status-updates/:meetingActionPlanUpdateId/document-preview", meetingActionPlanStatusUpdateController.GetMeetingActionPlanDocumentPreview)
	unauthenticatedRoutes.GET("/meetings/:id/action-plans/:meetingActionPlanId/status-updates/:meetingActionPlanUpdateId/document-download", meetingActionPlanStatusUpdateController.DownloadMeetingActionPlanDocument)

	// Meeting Agendas
	meetingRoutes.GET("/meetings/:id/agendas", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingAgendaController.GetAllMeetingAgendas)
	meetingRoutes.GET("/meetings/:id/agendas/:agendaID", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingAgendaController.GetMeetingAgendaDetails)
	meetingRoutes.POST("/meetings/:id/agendas", middlewares.RolePermissionMiddleware(db, "meeting", "modify"), meetingAgendaController.CreateMeetingAgenda)
	meetingRoutes.PUT("/meetings/:id/agendas/:agendaID", middlewares.RolePermissionMiddleware(db, "meeting", "modify"), meetingAgendaController.UpdateMeetingAgenda)
	meetingRoutes.POST("/meetings/:id/agendas/delete", middlewares.RolePermissionMiddleware(db, "meeting", "modify"), meetingAgendaController.DeleteMeetingAgendas)

	// Meeting Minutes
	meetingRoutes.GET("/meetings/:id/meeting-minutes", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingMinuteController.GetAllMeetingMinutes)
	meetingRoutes.GET("/meetings/:id/meeting-minutes/:meetingMinuteID", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingMinuteController.GetMeetingMinuteDetails)
	meetingRoutes.POST("/meetings/:id/meeting-minutes", middlewares.RolePermissionMiddleware(db, "meeting", "modify"), meetingMinuteController.CreateMeetingMinute)
	meetingRoutes.PUT("/meetings/:id/meeting-minutes/:meetingMinuteID", middlewares.RolePermissionMiddleware(db, "meeting", "modify"), meetingMinuteController.UpdateMeetingMinute)
	meetingRoutes.POST("/meetings/:id/meeting-minutes/delete", middlewares.RolePermissionMiddleware(db, "meeting", "modify"), meetingMinuteController.DeleteMeetingMinutes)

	// Meeting Associated Items
	meetingRoutes.GET("/meetings/:id/controls", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingAssociatedItemController.GetAllMeetingControls)
	meetingRoutes.GET("/meetings/:id/compliance-requirements", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingAssociatedItemController.GetAllMeetingComplianceRequirements)
	meetingRoutes.GET("/meetings/:id/risks", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingAssociatedItemController.GetAllMeetingRisks)
	meetingRoutes.GET("/meetings/:id/processes", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingAssociatedItemController.GetAllMeetingProcesses)
	meetingRoutes.GET("/meetings/:id/incidents", middlewares.RolePermissionMiddleware(db, "meeting", "view"), meetingAssociatedItemController.GetAllMeetingIncidents)

	// Meeting Minute Status
	meetingRoutes.GET("/meeting-minute-statuses", meetingMinuteStatusController.GetAllMeetingMinuteStatuses)
	meetingRoutes.GET("/meeting-minute-statuses/:id", meetingMinuteStatusController.GetMeetingMinuteStatusDetails)
}
