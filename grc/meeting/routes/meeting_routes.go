package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	meetingRoutes := router.Group("/")

	meetingRoutes.Use(middlewares.AuthMiddleware())

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
	meetingRoutes.POST("/meeting-types", meetingTypeController.CreateMeetingType)
	meetingRoutes.PUT("/meeting-types/:id", meetingTypeController.UpdateMeetingType)
	meetingRoutes.POST("/meeting-types/delete", meetingTypeController.DeleteMeetingTypes)
	meetingRoutes.PUT("/meeting-types/activate", meetingTypeController.ActivateMeetingTypes)
	meetingRoutes.PUT("/meeting-types/deactivate", meetingTypeController.DeactivateMeetingTypes)

	// Meetings
	meetingRoutes.GET("/meetings/export", meetingController.GetMeetingExport)
	meetingRoutes.GET("/meetings", meetingController.GetAllMeetings)
	meetingRoutes.GET("/meetings/:id", meetingController.GetMeetingDetails)
	meetingRoutes.POST("/meetings", meetingController.CreateMeeting)
	meetingRoutes.PUT("/meetings/:id", meetingController.UpdateMeeting)
	meetingRoutes.POST("/meetings/delete", meetingController.DeleteMeetings)
	meetingRoutes.GET("/meetings/:id/meeting-summary", meetingController.GetAllCounts)

	//Meeting action plan
	meetingRoutes.GET("/meeting-action-plans/export", meetingActionPlanController.GetMeetingActionPlanExport)
	meetingRoutes.GET("/meeting-action-plans", meetingActionPlanController.GetAllMeetingActionPlans)
	meetingRoutes.GET("/meetings/:id/action-plans", meetingActionPlanController.GetAllMeetingActionPlans)
	meetingRoutes.GET("/meetings/:id/action-plans/:meetingActionPlanId", meetingActionPlanController.GetMeetingActionPlanDetails)
	meetingRoutes.POST("/meetings/:id/action-plans", meetingActionPlanController.CreateMeetingActionPlan)
	meetingRoutes.PUT("/meetings/:id/action-plans/:meetingActionPlanId", meetingActionPlanController.UpdateMeetingActionPlan)
	meetingRoutes.POST("/meetings/:id/action-plans/delete", meetingActionPlanController.DeleteMeetingActionPlans)

	meetingRoutes.GET("/meetings/:id/action-plan-summary", meetingActionPlanController.GetMeetingActionPlanSummary)

	//Meeting action plan status updates
	meetingRoutes.GET("/meetings/:id/action-plans/:meetingActionPlanId/status-updates", meetingActionPlanStatusUpdateController.GetAllMeetingActionPlanStatusUpdates)
	meetingRoutes.GET("/meetings/:id/action-plans/:meetingActionPlanId/status-updates/:meetingActionPlanUpdateId", meetingActionPlanStatusUpdateController.GetMeetingActionPlanStatusUpdateDetails)
	meetingRoutes.POST("/meetings/:id/action-plans/:meetingActionPlanId/status-updates", meetingActionPlanStatusUpdateController.CreateMeetingActionPlanStatusUpdate)
	meetingRoutes.PUT("/meetings/:id/action-plans/:meetingActionPlanId/status-updates/:meetingActionPlanUpdateId", meetingActionPlanStatusUpdateController.UpdateMeetingActionPlanStatusUpdate)
	meetingRoutes.GET("/meetings/:id/action-plans/:meetingActionPlanId/status-updates/latest", meetingActionPlanStatusUpdateController.GetLatestMeetingActionPlanStatusUpdate)

	//Meeting Agendas
	meetingRoutes.GET("/meetings/:id/agendas", meetingAgendaController.GetAllMeetingAgendas)
	meetingRoutes.GET("/meetings/:id/agendas/:agendaID", meetingAgendaController.GetMeetingAgendaDetails)
	meetingRoutes.POST("/meetings/:id/agendas", meetingAgendaController.CreateMeetingAgenda)
	meetingRoutes.PUT("/meetings/:id/agendas/:agendaID", meetingAgendaController.UpdateMeetingAgenda)
	meetingRoutes.POST("/meetings/:id/agendas/delete", meetingAgendaController.DeleteMeetingAgendas)

	//Meeting Minutes
	meetingRoutes.GET("/meetings/:id/meeting-minutes", meetingMinuteController.GetAllMeetingMinutes)
	meetingRoutes.GET("/meetings/:id/meeting-minutes/:meetingMinuteID", meetingMinuteController.GetMeetingMinuteDetails)
	meetingRoutes.POST("/meetings/:id/meeting-minutes", meetingMinuteController.CreateMeetingMinute)
	meetingRoutes.PUT("/meetings/:id/meeting-minutes/:meetingMinuteID", meetingMinuteController.UpdateMeetingMinute)
	meetingRoutes.POST("/meetings/:id/meeting-minutes/delete", meetingMinuteController.DeleteMeetingMinutes)

	//Meeting Associated Items
	meetingRoutes.GET("/meetings/:id/controls", meetingAssociatedItemController.GetAllMeetingControls)
	meetingRoutes.GET("/meetings/:id/compliance-requirements", meetingAssociatedItemController.GetAllMeetingComplianceRequirements)
	meetingRoutes.GET("/meetings/:id/risks", meetingAssociatedItemController.GetAllMeetingRisks)
	meetingRoutes.GET("/meetings/:id/processes", meetingAssociatedItemController.GetAllMeetingProcesses)
	meetingRoutes.GET("/meetings/:id/incidents", meetingAssociatedItemController.GetAllMeetingIncidents)

	// Meeting Minute Status
	meetingRoutes.GET("/meeting-minute-statuses", meetingMinuteStatusController.GetAllMeetingMinuteStatuses)
	meetingRoutes.GET("/meeting-minute-statuses/:id", meetingMinuteStatusController.GetMeetingMinuteStatusDetails)

}
