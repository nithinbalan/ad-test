package initializers

import (
	commonRepositories "jethur/grc/common/repositories"
	commonServices "jethur/grc/common/services"

	"gorm.io/gorm"
)

type CommonModule struct {
	NotificationService commonServices.NotificationServiceInterface
}

func NewCommonModule(db *gorm.DB, settingsModule *SettingsModule) *CommonModule {

	notificationRepo := commonRepositories.NewNotificationRepository(db)
	notificationSvc := commonServices.NewNotificationService(notificationRepo, settingsModule.ModuleSvc)

	return &CommonModule{
		NotificationService: notificationSvc,
	}
}
