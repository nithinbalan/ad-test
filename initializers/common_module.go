package initializers

import (
	commonRepositories "jethur/grc/common/repositories"
	commonServices "jethur/grc/common/services"

	"gorm.io/gorm"
)

type CommonModule struct {
	NotificationService          commonServices.NotificationServiceInterface
	FileRepo                     commonRepositories.FileRepositoryInterface
	RiskMappingValidationService commonServices.RiskMappingServiceInterface
}

func NewCommonModule(db *gorm.DB, settingsModule *SettingsModule) *CommonModule {

	notificationRepo := commonRepositories.NewNotificationRepository(db)
	riskValidationRepo := commonRepositories.NewRiskMappingRepository(db)
	notificationSvc := commonServices.NewNotificationService(notificationRepo, settingsModule.ModuleSvc, settingsModule.EmailTemplateSettingsSvc)
	riskMappingValidationSvc := commonServices.NewMappingService(riskValidationRepo)

	fileRepo := commonRepositories.NewFileRepository()

	return &CommonModule{
		NotificationService:          notificationSvc,
		FileRepo:                     fileRepo,
		RiskMappingValidationService: riskMappingValidationSvc,
	}
}
