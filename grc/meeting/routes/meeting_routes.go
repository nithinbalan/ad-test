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
	meetingRoutes.GET("/meetings", meetingController.GetAllMeetings)
	meetingRoutes.GET("/meetings/:id", meetingController.GetMeetingDetails)
	meetingRoutes.POST("/meetings", meetingController.CreateMeeting)
	meetingRoutes.PUT("/meetings/:id", meetingController.UpdateMeeting)
	meetingRoutes.POST("/meetings/delete", meetingController.DeleteMeetings)

	//Meeting action plan
	meetingRoutes.GET("/meeting-action-plans", meetingActionPlanController.GetAllMeetingActionPlans)
	meetingRoutes.GET("/meetings/:id/action-plans", meetingActionPlanController.GetAllMeetingActionPlans)
	meetingRoutes.GET("/meetings/:id/action-plans/:meetingActionPlanId", meetingActionPlanController.GetMeetingActionPlanDetails)
	meetingRoutes.POST("/meetings/:id/action-plans", meetingActionPlanController.CreateMeetingActionPlan)
	meetingRoutes.PUT("/meetings/:id/action-plans/:meetingActionPlanId", meetingActionPlanController.UpdateMeetingActionPlan)
	meetingRoutes.POST("/meetings/:id/action-plans/delete", meetingActionPlanController.DeleteMeetingActionPlans)

	//Meeting action plan status updates
	meetingRoutes.GET("/meetings/:id/action-plans/:meetingActionPlanId/status-updates", meetingActionPlanStatusUpdateController.GetAllMeetingActionPlanStatusUpdates)
	meetingRoutes.GET("/meetings/:id/action-plans/:meetingActionPlanId/status-updates/:meetingActionPlanUpdateId", meetingActionPlanStatusUpdateController.GetMeetingActionPlanStatusUpdateDetails)
	meetingRoutes.POST("/meetings/:id/action-plans/:meetingActionPlanId/status-updates", meetingActionPlanStatusUpdateController.CreateMeetingActionPlanStatusUpdate)
	meetingRoutes.PUT("/meetings/:id/action-plans/:meetingActionPlanId/status-updates/:meetingActionPlanUpdateId", meetingActionPlanStatusUpdateController.UpdateMeetingActionPlanStatusUpdate)

}
