package initializers

import (
	"jethur/grc/meeting/controllers"
	"jethur/grc/meeting/repositories"
	"jethur/grc/meeting/services"

	"gorm.io/gorm"
)

type MeetingModule struct {
	MeetingFrequencyController *controllers.MeetingFrequencyController

	MeetingFrequencySvc services.MeetingFrequencyServiceInterface
}

func NewMeetingModule(db *gorm.DB, sharedModule *SharedModule) *MeetingModule {
	//Meeting Frequency
	meetingFrequencyRepo := repositories.NewMeetingFrequencyRepository(db)
	meetingFrequencySvc := services.NewMeetingFrequencyService(meetingFrequencyRepo)
	meetingFrequencyController := controllers.NewMeetingFrequencyController(meetingFrequencySvc)

	return &MeetingModule{
		MeetingFrequencyController: meetingFrequencyController,

		MeetingFrequencySvc: meetingFrequencySvc,
	}
}
