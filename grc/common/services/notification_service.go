package services

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/common/repositories"
	"jethur/grc/models"
	"jethur/grc/settings/services"
	"jethur/utility"
)

type NotificationServiceInterface interface {
	SendNotification(notificationDTO *commonDTOs.DBNotificationDTO, emailNotificationDTO *commonDTOs.EmailNotificationDTO, notificationMasterIdentifier string, tx ...utility.DatabaseTransaction) error
}

type NotificationService struct {
	repo      repositories.NotificationRepositoryInteface
	moduleSvc services.ModuleServiceInterface
}

func NewNotificationService(repo repositories.NotificationRepositoryInteface,
	moduleSvc services.ModuleServiceInterface) NotificationServiceInterface {
	return &NotificationService{
		repo:      repo,
		moduleSvc: moduleSvc,
	}
}

func (ns *NotificationService) SendNotification(notificationDTO *commonDTOs.DBNotificationDTO, emailNotificationDTO *commonDTOs.EmailNotificationDTO, notificationMasterIdentifier string, tx ...utility.DatabaseTransaction) error {
	// Check if notification is enabled for the user
	isEnabled, err := ns.IsNotificationEnabled(notificationDTO.UserID, notificationMasterIdentifier)
	if err != nil {
		return err
	}

	if isEnabled {
		go func() {
			err := utility.SendTemplatedEmail(
				emailNotificationDTO.EmailTo,
				emailNotificationDTO.Subject,
				emailNotificationDTO.TemplateRelativePath,
				emailNotificationDTO.EmailData,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to send email", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}()
	}

	// DB notification
	moduleID, err := ns.moduleSvc.GetModuleIDByIdentifier(notificationDTO.ModuleIdentifier)
	if err != nil {
		return err
	}

	userNotification := &models.UserNotification{
		UserID:               notificationDTO.UserID,
		ModuleID:             moduleID,
		NotificationType:     notificationDTO.NotificationType,
		ReferenceNumber:      notificationDTO.ReferenceNumber,
		HeadingEn:            notificationDTO.HeadingEn,
		HeadingAr:            notificationDTO.HeadingAr,
		TitleEn:              notificationDTO.TitleEn,
		TitleAr:              notificationDTO.TitleAr,
		URLParam:             notificationDTO.URLParam,
		QueryParam:           notificationDTO.QueryParam,
		IsSystemNotification: notificationDTO.IsSystemNotification,
		CreatedByID:          notificationDTO.CreatedByID,
	}

	err = ns.repo.Create(userNotification)
	if err != nil {
		return err
	}

	return nil
}

func (ns *NotificationService) IsNotificationEnabled(userID uint, notificationMasterIdentifier string) (bool, error) {
	notificationMaster, err := ns.repo.GetNotificationMasterByIdentifier(notificationMasterIdentifier)
	if err != nil {
		return false, err
	}

	setting, err := ns.repo.GetUserNotificationSetting(userID, notificationMaster.ID)
	if err != nil {
		return false, err
	}

	return setting != nil && setting.IsEnable, nil
}

// func (ns *NotificationService) SendNotification(notificationDTO *commonDTOs.NotificationDTO, notificationMasterIdentifier string, tx ...utility.DatabaseTransaction) error {
// 	// Check if notification is enabled for the user
// 	isEnabled, err := ns.IsNotificationEnabled(notificationDTO.UserID, notificationMasterIdentifier)
// 	if err != nil {
// 		return err
// 	}

// 	if isEnabled {
// 		//@todo email notification
// 	}

// 	moduleID, err := ns.moduleSvc.GetModuleIDByIdentifier(notificationDTO.ModuleIdentifier)
// 	if err != nil {
// 		return err
// 	}

// 	userNotification := &models.UserNotification{
// 		UserID:      notificationDTO.UserID,
// 		ModuleID:    moduleID,
// 		Icon:        notificationDTO.Icon,
// 		MessageEn:   notificationDTO.MessageEn,
// 		MessageAr:   notificationDTO.MessageAr,
// 		RedirectURL: notificationDTO.RedirectURL,
// 		CreatedByID:   notificationDTO.CreatedBy,
// 	}

// 	err = ns.repo.Create(userNotification)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (ns *NotificationService) IsNotificationEnabled(userID uint, notificationMasterIdentifier string) (bool, error) {
// 	notificationMaster, err := ns.repo.GetNotificationMasterByIdentifier(notificationMasterIdentifier)
// 	if err != nil {
// 		return false, err
// 	}

// 	setting, err := ns.repo.GetUserNotificationSetting(userID, notificationMaster.ID)
// 	if err != nil {
// 		return false, err
// 	}

// 	return setting != nil && setting.IsEnable, nil
// }
