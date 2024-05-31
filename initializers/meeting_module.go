package initializers

import (
	"jethur/grc/meeting/controllers"
	"jethur/grc/meeting/repositories"
	"jethur/grc/meeting/services"

	"gorm.io/gorm"
)

type MeetingModule struct {
	MeetingFrequencyController              *controllers.MeetingFrequencyController
	MeetingStatusController                 *controllers.MeetingStatusController
	MeetingTypeController                   *controllers.MeetingTypeController
	MeetingController                       *controllers.MeetingController
	MeetingActionPlanController             *controllers.MeetingActionPlanController
	MeetingActionPlanStatusUpdateController *controllers.MeetingActionPlanStatusUpdateController

	MeetingFrequencySvc services.MeetingFrequencyServiceInterface
	MeetingStatusSvc    services.MeetingStatusServiceInterface
	MeetingTypeSvc      services.MeetingTypeServiceInterface
	MeetingSvc          services.MeetingServiceInterface
}

func NewMeetingModule(db *gorm.DB, sharedModule *SharedModule, organizationModule *OrganizationModule, masterModule *MasterModule, commonModule *CommonModule) *MeetingModule {
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

	meetingRepo := repositories.NewMeetingRepository(db)
	meetingSvc := services.NewMeetingService(
		meetingRepo,
		sharedModule.ActivityLogService,
		meetingTypeSvc,
		meetingFrequencySvc,
		organizationModule.UserSvc,
		meetingStatusSvc)
	meetingController := controllers.NewMeetingController(meetingSvc)

	meetingActionPlanRepo := repositories.NewMeetingActionPlanRepository(db)
	meetingActionPlanSvc := services.NewMeetingActionPlanService(meetingActionPlanRepo, sharedModule.ActivityLogService, meetingSvc, masterModule.SeverityLevelSvc, masterModule.ActionPlanStatusSvc, organizationModule.UserSvc, commonModule.NotificationService)
	meetingActionPlanController := controllers.NewMeetingActionPlanController(meetingActionPlanSvc)

	meetingActionPlanStatusUpdateRepo := repositories.NewMeetingActionPlanStatusUpdateRepository(db)
	meetingActionPlanStatusUpdateSvc := services.NewMeetingActionPlanStatusUpdateService(meetingActionPlanStatusUpdateRepo, sharedModule.ActivityLogService, meetingActionPlanSvc, masterModule.ActionPlanStatusSvc)
	meetingActionPlanStatusUpdateController := controllers.NewMeetingActionPlanStatusUpdateController(meetingActionPlanStatusUpdateSvc)

	return &MeetingModule{
		MeetingFrequencyController:              meetingFrequencyController,
		MeetingStatusController:                 meetingStatusController,
		MeetingTypeController:                   meetingTypeController,
		MeetingController:                       meetingController,
		MeetingActionPlanController:             meetingActionPlanController,
		MeetingActionPlanStatusUpdateController: meetingActionPlanStatusUpdateController,

		MeetingFrequencySvc: meetingFrequencySvc,
		MeetingStatusSvc:    meetingStatusSvc,
		MeetingTypeSvc:      meetingTypeSvc,
		MeetingSvc:          meetingSvc,
	}
}
