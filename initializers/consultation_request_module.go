package initializers

import (
	"jethur/grc/consultation-request/controllers"
	"jethur/grc/consultation-request/repositories"
	"jethur/grc/consultation-request/services"

	"gorm.io/gorm"
)

type ConsultationRequestModule struct {
	ConsultationRequestStatusController *controllers.ConsultationRequestStatusController
	ConsultationRequestController       *controllers.ConsultationRequestController
	// Add more controllers if needed
}

func NewConsultationRequestModule(db *gorm.DB, sharedModule *SharedModule, masterModule *MasterModule, organizationModule *OrganizationModule, taskModule *TaskModule, settingsModule *SettingsModule, commonModule *CommonModule, documentModule *DocumentModule) *ConsultationRequestModule {
	consultationRequestStatusRepo := repositories.NewConsultationRequestStatusRepository(db)
	consultationRequestStatusSvc := services.NewConsultationRequestStatusService(consultationRequestStatusRepo)
	consultationRequestStatusController := controllers.NewConsultationRequestStatusController(consultationRequestStatusSvc)

	consultationRequestRepo := repositories.NewConsultationRequestRepository(db)
	consultationRequestSvc := services.NewConsultationRequestService(consultationRequestRepo, organizationModule.UserSvc, sharedModule.ActivityLogService, settingsModule.OrganizationUnitSvc, consultationRequestStatusSvc)
	consultationRequestController := controllers.NewConsultationRequestController(consultationRequestSvc)
	return &ConsultationRequestModule{
		ConsultationRequestController:       consultationRequestController,
		ConsultationRequestStatusController: consultationRequestStatusController,
	}
}
