package initializers

import (
	"jethur/grc/meeting/controllers"
	"jethur/grc/meeting/repositories"
	"jethur/grc/meeting/services"

	"gorm.io/gorm"
)

type MeetingModule struct {
	MeetingFrequencyController *controllers.MeetingFrequencyController
	MeetingStatusController    *controllers.MeetingStatusController
	MeetingTypeController      *controllers.MeetingTypeController

	MeetingFrequencySvc services.MeetingFrequencyServiceInterface
	MeetingStatusSvc    services.MeetingStatusServiceInterface
	MeetingTypeSvc      services.MeetingTypeServiceInterface
}

func NewMeetingModule(db *gorm.DB, sharedModule *SharedModule) *MeetingModule {
	//Meeting Frequency
	meetingFrequencyRepo := repositories.NewMeetingFrequencyRepository(db)
	meetingFrequencySvc := services.NewMeetingFrequencyService(meetingFrequencyRepo)
	meetingFrequencyController := controllers.NewMeetingFrequencyController(meetingFrequencySvc)

	meetingStatusRepo := repositories.NewMeetingStatusRepository(db)
	meetingStatusSvc := services.NewMeetingStatusService(meetingStatusRepo)
	meetingStatusController := controllers.NewMeetingStatusController(meetingStatusSvc)

	meetingTypeRepo := repositories.NewMeetingTypeRepository(db)
	meetingTypeSvc := services.NewMeetingTypeService(meetingTypeRepo, sharedModule.ActivityLogService)
	meetingTypeController := controllers.NewMeetingTypeController(meetingTypeSvc)

	return &MeetingModule{
		MeetingFrequencyController: meetingFrequencyController,
		MeetingStatusController:    meetingStatusController,
		MeetingTypeController:      meetingTypeController,

		MeetingFrequencySvc: meetingFrequencySvc,
		MeetingStatusSvc:    meetingStatusSvc,
		MeetingTypeSvc:      meetingTypeSvc,
	}
}
