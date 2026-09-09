package initializers

import (
	"jethur/grc/security-operations/controllers"
	"jethur/grc/security-operations/repositories"
	"jethur/grc/security-operations/services"

	"gorm.io/gorm"
)

type SecurityOperationsModule struct {
	CSLogDataEntityTypeController         *controllers.CSLogDataEntityTypeController
	CSLogDataReceptionMechanismController *controllers.CSLogDataReceptionMechanismController
	CSLogCommunicationChannelController   *controllers.CSLogCommunicationChannelController
	CSLogSourceDirectoryController        *controllers.CSLogSourceDirectoryController
	CSLogSocRoomAccessLogController       *controllers.CSLogSocRoomAccessLogController

	CSLogSourceDirectorySvc  services.CSLogSourceDirectoryServiceInterface
	CSLogSocRoomAccessLogSvc services.CSLogSocRoomAccessLogServiceInterface
}

func NewSecurityOperationsModule(
	db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	organizationModule *OrganizationModule,
	commonModule *CommonModule,
) *SecurityOperationsModule {
	csLogDataEntityTypeRepo := repositories.NewCSLogDataEntityTypeRepository(db)
	csLogDataEntityTypeSvc := services.NewCSLogDataEntityTypeService(csLogDataEntityTypeRepo, sharedModule.ActivityLogService)
	csLogDataEntityTypeController := controllers.NewCSLogDataEntityTypeController(csLogDataEntityTypeSvc)

	csLogDataReceptionMechanismRepo := repositories.NewCSLogDataReceptionMechanismRepository(db)
	csLogDataReceptionMechanismSvc := services.NewCSLogDataReceptionMechanismService(csLogDataReceptionMechanismRepo, sharedModule.ActivityLogService)
	csLogDataReceptionMechanismController := controllers.NewCSLogDataReceptionMechanismController(csLogDataReceptionMechanismSvc)

	csLogCommunicationChannelRepo := repositories.NewCSLogCommunicationChannelRepository(db)
	csLogCommunicationChannelSvc := services.NewCSLogCommunicationChannelService(csLogCommunicationChannelRepo, sharedModule.ActivityLogService)
	csLogCommunicationChannelController := controllers.NewCSLogCommunicationChannelController(csLogCommunicationChannelSvc)

	csLogSourceDirectoryRepo := repositories.NewCSLogSourceDirectoryRepository(db)
	csLogSourceDirectorySvc := services.NewCSLogSourceDirectoryService(
		csLogSourceDirectoryRepo,
		sharedModule.ActivityLogService,
		csLogDataEntityTypeSvc,
		csLogDataReceptionMechanismSvc,
		csLogCommunicationChannelSvc,
		settingsModule.OrganizationUnitSvc,
		organizationModule.UserSettingSvc,
	)
	csLogSourceDirectoryController := controllers.NewCSLogSourceDirectoryController(csLogSourceDirectorySvc)

	csLogSocRoomAccessLogRepo := repositories.NewCSLogSocRoomAccessLogRepository(db)
	csLogSocRoomAccessLogSvc := services.NewCSLogSocRoomAccessLogService(
		csLogSocRoomAccessLogRepo,
		sharedModule.ActivityLogService,
		settingsModule.OrganizationUnitSvc,
		organizationModule.UserSvc,
		organizationModule.UserSettingSvc,
		commonModule.NotificationService,
	)
	csLogSocRoomAccessLogController := controllers.NewCSLogSocRoomAccessLogController(csLogSocRoomAccessLogSvc)

	return &SecurityOperationsModule{
		CSLogDataEntityTypeController:         csLogDataEntityTypeController,
		CSLogDataReceptionMechanismController: csLogDataReceptionMechanismController,
		CSLogCommunicationChannelController:   csLogCommunicationChannelController,
		CSLogSourceDirectoryController:        csLogSourceDirectoryController,
		CSLogSocRoomAccessLogController:       csLogSocRoomAccessLogController,
		CSLogSourceDirectorySvc:               csLogSourceDirectorySvc,
		CSLogSocRoomAccessLogSvc:              csLogSocRoomAccessLogSvc,
	}
}
