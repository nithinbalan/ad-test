package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	masterRoutes := router.Group("/")

	masterRoutes.Use(middlewares.AuthMiddleware())

	meetingFrequencyController := app.MeetingModule.MeetingFrequencyController

	masterRoutes.GET("/meeting-frequencies", meetingFrequencyController.GetAllMeetingFrequencies)
	masterRoutes.GET("/meeting-frequencies/:id", meetingFrequencyController.GetMeetingFrequencyDetails)

}
