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
	DelegateeUserService            commonServices.DelegateeUserserviceInterface
	FileService                     commonServices.FileServiceInterface
}

func NewCommonModule(db *gorm.DB, settingsModule *SettingsModule) *CommonModule {

	notificationRepo := commonRepositories.NewNotificationRepository(db)
	notificationSvc := commonServices.NewNotificationService(notificationRepo, settingsModule.ModuleSvc, settingsModule.EmailTemplateSettingsSvc)
	associatedItemValidationRepo := commonRepositories.NewRiskMappingRepository(db)
	associatedItemValidationSvc := commonServices.NewAssociatedItemValidationService(associatedItemValidationRepo)
	DelegateeUserRepo := commonRepositories.NewDelegateeUserRepository(db)
	DelegateeUserSvc := commonServices.NewDelegateeUserservice(DelegateeUserRepo, settingsModule.ModuleSvc)

	fileRepo := commonRepositories.NewFileRepository()
	fileSvc := commonServices.NewFileService(fileRepo)
	return &CommonModule{
		NotificationService:             notificationSvc,
		FileRepo:                        fileRepo,
		AssociatedItemValidationService: associatedItemValidationSvc,
		DelegateeUserService:            DelegateeUserSvc,
		FileService:                     fileSvc,
	}
}
