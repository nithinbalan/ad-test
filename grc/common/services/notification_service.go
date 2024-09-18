package services

import (
	"encoding/json"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/common/repositories"
	"jethur/grc/models"
	"jethur/grc/settings/services"
	"jethur/utility"
	"time"
)

type NotificationServiceInterface interface {
	SendNotification(notificationDTO *commonDTOs.DBNotificationDTO, emailNotificationDTO *commonDTOs.EmailNotificationDTO, notificationMasterIdentifier string, tx ...utility.DatabaseTransaction) error
	LogEmail(recipientEmail string, subject string, templatePath string, emailData interface{}, status string, errorMessage string, sentAt *time.Time, createdByID uint) error
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
	}

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
		Entity:               notificationDTO.Entity,
		Action:               notificationDTO.Action,
		URLParam:             notificationDTO.URLParam,
		QueryParam:           notificationDTO.QueryParam,
		IsSystemNotification: notificationDTO.IsSystemNotification,
		CreatedByID:          notificationDTO.CreatedByID,
	}

	if err := ns.repo.Create(userNotification); err != nil {
		utility.LogError("Failed to create DB notification", map[string]interface{}{
			"error": err.Error(),
		})
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

func (ns *NotificationService) LogEmail(recipientEmail string, subject string, templatePath string, emailData interface{}, status string, errorMessage string, sentAt *time.Time, createdByID uint) error {
	// Assert that emailData is of type map[string]interface{}
	emailDataMap, ok := emailData.(map[string]interface{})
	if !ok {
		utility.LogError("Failed to assert email data to map[string]interface{}", map[string]interface{}{
			"error": "Type assertion failed",
		})
		return fmt.Errorf("invalid email data format")
	}

	// Serialize email data to store in log
	emailDataBytes, err := json.Marshal(emailDataMap)
	if err != nil {
		utility.LogError("Failed to serialize email data", map[string]interface{}{"error": err.Error()})
		return err
	}

	emailLog := &models.EmailLog{
		RecipientEmail: recipientEmail,
		Subject:        subject,
		TemplatePath:   templatePath,
		EmailData:      string(emailDataBytes),
		SendStatus:     status,
		ErrorMessage:   errorMessage,
		SentAt:         sentAt,
		CreatedByID:    createdByID,
		CreatedAt:      time.Now(),
	}

	if err := ns.repo.CreateEmailLog(emailLog); err != nil {
		utility.LogError("Failed to create email log", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	return nil
}
