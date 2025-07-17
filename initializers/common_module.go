package initializers

import (
	commonRepositories "jethur/grc/common/repositories"
	commonServices "jethur/grc/common/services"

	"gorm.io/gorm"
)

type CommonModule struct {
	NotificationService             commonServices.NotificationServiceInterface
	FileRepo                        commonRepositories.FileRepositoryInterface
	AssociatedItemValidationService commonServices.AssociatedItemValidationServiceInterface
}

func NewCommonModule(db *gorm.DB, settingsModule *SettingsModule) *CommonModule {

	notificationRepo := commonRepositories.NewNotificationRepository(db)
	notificationSvc := commonServices.NewNotificationService(notificationRepo, settingsModule.ModuleSvc, settingsModule.EmailTemplateSettingsSvc)
	associatedItemValidationRepo := commonRepositories.NewRiskMappingRepository(db)
	associatedItemValidationSvc := commonServices.NewAssociatedItemValidationService(associatedItemValidationRepo)

	fileRepo := commonRepositories.NewFileRepository()

	return &CommonModule{
		NotificationService:             notificationSvc,
		FileRepo:                        fileRepo,
		AssociatedItemValidationService: associatedItemValidationSvc,
	}
}
