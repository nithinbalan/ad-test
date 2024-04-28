package services

import (
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	userDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"

	"gorm.io/gorm"
)

type UserSettingServiceInterface interface {
	GetNotificationSettings(userID uint) ([]userDTOs.NotificationSettingListResponseDTO, error)
	UpdateNotificationSetting(userID uint, notificationMasterID uint, enable bool) error
}

type UserSettingService struct {
	repo           organizationRepositories.UserSettingRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewUserSettingService(
	repo organizationRepositories.UserSettingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface) UserSettingServiceInterface {
	return &UserSettingService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (uss *UserSettingService) GetNotificationSettings(userID uint) ([]userDTOs.NotificationSettingListResponseDTO, error) {
	modules, err := uss.repo.GetNotificationSetting(userID)
	if err != nil {
		return nil, err
	}

	var response []userDTOs.NotificationSettingListResponseDTO
	for _, module := range modules {
		var notificationDTOs []userDTOs.NotificationListDTO
		// Check if NotificationMasters is nil or not before ranging over it
		if module.NotificationMasters != nil {
			for _, notification := range *module.NotificationMasters {
				enableStatus := false
				// Check if UserNotificationSetting is preloaded and has a value
				if notification.UserNotificationSettings != nil { // Assuming it might be a slice, handle accordingly
					for _, setting := range *notification.UserNotificationSettings {
						if setting.UserID == userID {
							enableStatus = setting.IsEnable
							break // Assuming only one setting per user per notification is possible
						}
					}
				}

				notificationDTO := userDTOs.NotificationListDTO{
					ID:         notification.ID,
					ModulesID:  notification.ModuleID,
					Identifier: notification.Identifier,
					TitleEn:    notification.TitleEn,
					TitleAr:    notification.TitleAr,
					IsEnable:   enableStatus,
				}
				notificationDTOs = append(notificationDTOs, notificationDTO)
			}
		}

		moduleDTO := userDTOs.NotificationSettingListResponseDTO{
			ModuleID:      module.ID,
			ModuleTitleEn: module.TitleEn,
			ModuleTitleAr: module.TitleAr,
			Notifications: &notificationDTOs,
		}
		response = append(response, moduleDTO)
	}

	return response, nil
}

func (uss *UserSettingService) UpdateNotificationSetting(userID uint, notificationMasterID uint, enable bool) error {
	notificationSetting, err := uss.repo.FindNotificationSetting(userID, notificationMasterID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if notificationSetting.ID == 0 {
		// No setting found, create a new one
		createNotificationSetting := &models.UserNotificationSetting{
			UserID:               userID,
			NotificationMasterID: notificationMasterID,
			IsEnable:             enable,
		}
		return uss.repo.CreateNotificationSetting(createNotificationSetting)
	}

	// Existing setting found, update it
	notificationSetting.IsEnable = enable
	return uss.repo.UpdateNotificationSetting(notificationSetting)
}
