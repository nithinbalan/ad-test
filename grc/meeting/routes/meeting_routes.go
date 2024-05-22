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

}
