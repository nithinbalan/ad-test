package services

import (
	"jethur/grc/common/repositories"
	"jethur/grc/models"
	"jethur/grc/settings/services"
	"jethur/utility"
)

type NotificationServiceInterface interface {
	SendNotification(userID uint, moduleIdentifier string, icon models.Icon, messageEn string, messageAr string, redirectURL string, currentUserID uint, tx ...utility.DatabaseTransaction) error
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

func (ns *NotificationService) SendNotification(userID uint, moduleIdentifier string, icon models.Icon, messageEn string, messageAr string, redirectURL string, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	moduleID, err := ns.moduleSvc.GetModuleIDByIdentifier(moduleIdentifier)
	if err != nil {
		return err
	}

	userNotification := &models.UserNotification{
		UserID:      userID,
		ModuleID:    moduleID,
		Icon:        icon,
		MessageEn:   messageEn,
		MessageAr:   messageAr,
		RedirectURL: redirectURL,
		CreatedBy:   currentUserID,
	}

	err = ns.repo.Create(userNotification)
	if err != nil {
		return err
	}

	return nil

}
