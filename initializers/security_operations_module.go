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
}

func NewSecurityOperationsModule(
	db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
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

	return &SecurityOperationsModule{
		CSLogDataEntityTypeController:         csLogDataEntityTypeController,
		CSLogDataReceptionMechanismController: csLogDataReceptionMechanismController,
		CSLogCommunicationChannelController:   csLogCommunicationChannelController,
	}
}
