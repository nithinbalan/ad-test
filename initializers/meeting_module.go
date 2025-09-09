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
	MeetingAgendaController                 *controllers.MeetingAgendaController
	MeetingMinuteController                 *controllers.MeetingMinuteController
	MeetingAssociatedItemController         *controllers.MeetingAssociatedItemController
	MeetingMinuteStatusController           *controllers.MeetingMinuteStatusController
	MeetingActionPlanCommentController      *controllers.MeetingActionPlanCommentController

	MeetingFrequencySvc    services.MeetingFrequencyServiceInterface
	MeetingStatusSvc       services.MeetingStatusServiceInterface
	MeetingTypeSvc         services.MeetingTypeServiceInterface
	MeetingSvc             services.MeetingServiceInterface
	MeetingMinuteStatusSvc services.MeetingMinuteStatusServiceInterface
	MeetingActionPlanSvc   services.MeetingActionPlanServiceInterface
}

func NewMeetingModule(db *gorm.DB, sharedModule *SharedModule, organizationModule *OrganizationModule, masterModule *MasterModule, commonModule *CommonModule, complianceModule *ComplianceModule, riskModule *RiskModule, incidentModule *IncidentModule, taskModule *TaskModule) *MeetingModule {
	// Meeting Frequency
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
		meetingStatusSvc,
		organizationModule.ControlSvc,
		complianceModule.ComplianceRequirementSvc,
		riskModule.RiskRegisterSvc,
		organizationModule.ProcessSvc,
		incidentModule.IncidentRegisterSvc,
		organizationModule.UserSettingSvc, commonModule.NotificationService)
	meetingController := controllers.NewMeetingController(meetingSvc)

	meetingAgendaRepo := repositories.NewMeetingAgendaRepository(db)
	meetingAgendaSvc := services.NewMeetingAgendaService(meetingAgendaRepo, sharedModule.ActivityLogService, meetingSvc)
	meetingAgendaController := controllers.NewMeetingAgendaController(meetingAgendaSvc)

	meetingMinuteStatusRepo := repositories.NewMeetingMinuteStatusRepository(db)
	meetingMinuteStatusSvc := services.NewMeetingMinuteStatusService(meetingMinuteStatusRepo)
	meetingMinuteStatusController := controllers.NewMeetingMinuteStatusController(meetingMinuteStatusSvc)

	meetingMinuteRepo := repositories.NewMeetingMinuteRepository(db)
	meetingMinuteSvc := services.NewMeetingMinuteService(meetingMinuteRepo, sharedModule.ActivityLogService, meetingSvc, meetingAgendaSvc, meetingMinuteStatusSvc)
	meetingMinuteController := controllers.NewMeetingMinuteController(meetingMinuteSvc)

	meetingActionPlanRepo := repositories.NewMeetingActionPlanRepository(db)
	meetingActionPlanSvc := services.NewMeetingActionPlanService(meetingActionPlanRepo, sharedModule.ActivityLogService, meetingSvc, masterModule.SeverityLevelSvc, masterModule.ActionPlanStatusSvc, organizationModule.UserSvc, commonModule.NotificationService, meetingMinuteSvc, meetingMinuteStatusSvc, organizationModule.UserSettingSvc, taskModule.TaskService, commonModule.DelegateeUserService)
	meetingActionPlanController := controllers.NewMeetingActionPlanController(meetingActionPlanSvc)

	meetingActionPlanStatusUpdateRepo := repositories.NewMeetingActionPlanStatusUpdateRepository(db)
	meetingActionPlanStatusUpdateSvc := services.NewMeetingActionPlanStatusUpdateService(meetingActionPlanStatusUpdateRepo, sharedModule.ActivityLogService, meetingActionPlanSvc, masterModule.ActionPlanStatusSvc, meetingMinuteSvc, meetingMinuteStatusSvc)
	meetingActionPlanStatusUpdateController := controllers.NewMeetingActionPlanStatusUpdateController(meetingActionPlanStatusUpdateSvc)

	meetingAssociatedItemRepo := repositories.NewMeetingAssociatedItemRepository(db)
	meetingAssociatedItemSvc := services.NewMeetingAssociatedItemService(meetingAssociatedItemRepo, meetingSvc, commonModule.AssociatedItemValidationService)
	meetingAssociatedItemController := controllers.NewMeetingAssociatedItemController(meetingAssociatedItemSvc)

	meetingActionPlanCommentRepository := repositories.NewMeetingActionPlanCommentRepository(db)
	meetingActionPlanCommentSvc := services.NewMeetingActionPlanCommentService(meetingActionPlanCommentRepository, organizationModule.UserSvc, sharedModule.ActivityLogService, meetingActionPlanSvc, commonModule.NotificationService)
	meetingActionPlanCommentController := controllers.NewMeetingActionPlanCommentController(meetingActionPlanCommentSvc)

	return &MeetingModule{
		MeetingFrequencyController:              meetingFrequencyController,
		MeetingStatusController:                 meetingStatusController,
		MeetingTypeController:                   meetingTypeController,
		MeetingController:                       meetingController,
		MeetingActionPlanController:             meetingActionPlanController,
		MeetingActionPlanStatusUpdateController: meetingActionPlanStatusUpdateController,
		MeetingAgendaController:                 meetingAgendaController,
		MeetingMinuteController:                 meetingMinuteController,
		MeetingAssociatedItemController:         meetingAssociatedItemController,
		MeetingMinuteStatusController:           meetingMinuteStatusController,
		MeetingActionPlanCommentController:      meetingActionPlanCommentController,

		MeetingFrequencySvc:    meetingFrequencySvc,
		MeetingStatusSvc:       meetingStatusSvc,
		MeetingTypeSvc:         meetingTypeSvc,
		MeetingSvc:             meetingSvc,
		MeetingMinuteStatusSvc: meetingMinuteStatusSvc,
		MeetingActionPlanSvc:   meetingActionPlanSvc,
	}
}
