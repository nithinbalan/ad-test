package services

import (
	"encoding/json"
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/common/repositories"
	"jethur/grc/models"
	"jethur/grc/settings/services"
	"jethur/utility"
)

type NotificationServiceInterface interface {
	SendNotification(notificationDTO *commonDTOs.DBNotificationDTO, emailNotificationDTO *commonDTOs.EmailNotificationDTO, notificationMasterIdentifier string, tx ...utility.DatabaseTransaction) error
	LogEmail(recipientEmail string, subject string, templatePath string, emailData interface{}, status string, errorMessage string, sentAt *time.Time, createdByID uint) error

	GetSocialMediaLinks() (*commonDTOs.SocialMediaLinksDTO, error)
}

type NotificationService struct {
	repo                     repositories.NotificationRepositoryInteface
	moduleSvc                services.ModuleServiceInterface
	emailTemplateSettingsSvc services.EmailTemplateSettingsServiceInterface
}

func NewNotificationService(repo repositories.NotificationRepositoryInteface,
	moduleSvc services.ModuleServiceInterface,
	emailTemplateSettingsSvc services.EmailTemplateSettingsServiceInterface,
) NotificationServiceInterface {
	return &NotificationService{
		repo:                     repo,
		moduleSvc:                moduleSvc,
		emailTemplateSettingsSvc: emailTemplateSettingsSvc,
	}
}

func (ns *NotificationService) SendNotification(notificationDTO *commonDTOs.DBNotificationDTO, emailNotificationDTO *commonDTOs.EmailNotificationDTO, notificationMasterIdentifier string, tx ...utility.DatabaseTransaction) error {
	// **Ensure at least one of the parameters is provided**
	if notificationDTO == nil && emailNotificationDTO == nil {
		return fmt.Errorf("both notificationDTO and emailNotificationDTO are nil, cannot send notification")
	}

	var notifications []*models.UserNotification

	// **Process DB Notification if provided**
	if notificationDTO != nil {
		if len(notificationDTO.UserIDs) == 0 {
			return fmt.Errorf("no user IDs provided for notification")
		}

		for _, userID := range notificationDTO.UserIDs {
			// **Check if notification is enabled for the user**
			isEnabled, err := ns.IsNotificationEnabled(userID, notificationMasterIdentifier)
			if err != nil {
				return err
			}

			if !isEnabled {
				continue
			}

			moduleID, err := ns.moduleSvc.GetModuleIDByIdentifier(notificationDTO.ModuleIdentifier)
			if err != nil {
				return err
			}

			notifications = append(notifications, &models.UserNotification{
				UserID:               userID,
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
			})
		}

		if len(notifications) > 0 {
			if err := ns.repo.Create(notifications); err != nil {
				utility.LogError("Failed to create bulk DB notifications", map[string]interface{}{
					"error": err.Error(),
				})
				return err
			}
		}
	}

	// **Process Email Notification if provided**
	if emailNotificationDTO != nil {
		emailTemplateSettings, err := ns.emailTemplateSettingsSvc.GetEmailTemplateSettings()
		if err != nil {
			fmt.Printf("failed to get email template settings: %v", err)
			return err
		}

		emailTemplateAndSocialDTO := &commonDTOs.EmailTemplateAndSocialDTO{
			HeaderImageURL: emailTemplateSettings.HeaderImageURL,
			SupportText:    emailTemplateSettings.SupportText,
			ClosingText:    emailTemplateSettings.ClosingText,
			FooterLogoURL:  emailTemplateSettings.FooterLogoURL,
			FooterAddress:  emailTemplateSettings.FooterAddress,
			FooterEmail:    emailTemplateSettings.FooterEmail,
		}

		err = utility.SendTemplatedEmail(
			emailNotificationDTO.EmailTo,
			emailNotificationDTO.Subject,
			emailNotificationDTO.TemplateRelativePath,
			emailNotificationDTO.EmailData,
			emailTemplateAndSocialDTO,
			nil,
			emailNotificationDTO.CC,
		)
		if err != nil {
			utility.LogError("Failed to send email", map[string]interface{}{
				"error": err.Error(),
			})
		}
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

func (ns *NotificationService) GetSocialMediaLinks() (*commonDTOs.SocialMediaLinksDTO, error) {
	socialLinks, err := ns.repo.GetSocialMediaLinks()
	if err != nil {
		return nil, err
	}
	return socialLinks, nil
}
