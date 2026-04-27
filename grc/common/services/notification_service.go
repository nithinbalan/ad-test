package services

import (
	"encoding/json"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"time"

	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	masterServices "jethur/grc/masters/services"

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
	emailStatusSvc           masterServices.EmailStatusServiceInterface
}

func NewNotificationService(repo repositories.NotificationRepositoryInteface,
	moduleSvc services.ModuleServiceInterface,
	emailTemplateSettingsSvc services.EmailTemplateSettingsServiceInterface,
	emailStatusSvc masterServices.EmailStatusServiceInterface,
) NotificationServiceInterface {
	return &NotificationService{
		repo:                     repo,
		moduleSvc:                moduleSvc,
		emailTemplateSettingsSvc: emailTemplateSettingsSvc,
		emailStatusSvc:           emailStatusSvc,
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
		var createdByID *uint
		if notificationDTO != nil {
			createdByID = notificationDTO.CreatedByID
		}
		_, err := ns.EnqueueEmail(emailNotificationDTO, notificationMasterIdentifier, createdByID)
		if err != nil {
			return err
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

// EnqueueEmail renders the HTML body from the template, then inserts one row into
// emails and one row per address into email_recipients (status = "pending").
// For each address:
//   - If the address belongs to a system user → internal (is_external=false, user_id set).
//     Notification settings are checked; the address is skipped if notifications are disabled.
//   - If the address is not in the system → external (is_external=true).
//     External recipients are always added (no notification settings to check).
//
// No SMTP call is made here; the worker handles sending.
func (ns *NotificationService) EnqueueEmail(dto *commonDTOs.EmailNotificationDTO, notificationMasterIdentifier string, createdByID *uint) (uint, error) {
	if dto == nil {
		return 0, fmt.Errorf("email notification DTO is nil")
	}
	if len(dto.EmailTo) == 0 {
		return 0, fmt.Errorf("at least one TO recipient is required")
	}

	templateSettings, err := ns.emailTemplateSettingsSvc.GetEmailTemplateSettings()
	if err != nil {
		return 0, fmt.Errorf("failed to load email template settings: %w", err)
	}

	settingsDTO := &commonDTOs.EmailTemplateAndSocialDTO{
		SupportText:   templateSettings.SupportText,
		ClosingText:   templateSettings.ClosingText,
		FooterAddress: templateSettings.FooterAddress,
		FooterEmail:   templateSettings.FooterEmail,
	}

	body, err := utility.RenderEmailBody(dto.TemplateRelativePath, dto.EmailData, settingsDTO)
	if err != nil {
		return 0, fmt.Errorf("failed to render email body: %w", err)
	}

	pendingStatus, err := ns.emailStatusSvc.GetByType("pending")
	if err != nil {
		return 0, fmt.Errorf("failed to get pending email status: %w", err)
	}

	email := &models.Email{
		Subject:       dto.Subject,
		Body:          body,
		EmailStatusID: pendingStatus.ID,
		MaxAttempts:   5,
		CreatedByID:   createdByID,
	}

	var recipients []*models.EmailRecipient

	// buildToRecipients: internal users must have IsNotification=true and notification-master enabled.
	// External users are always added.
	buildToRecipients := func(addresses []string) {
		for _, addr := range utility.DedupeEmails(utility.SanitizeEmails(addresses)) {
			user, _ := ns.repo.GetUserByEmail(addr)

			if user != nil {
				// Check global notification flag
				userSetting, err := ns.repo.GetUserSettingByUserID(user.ID)
				if err != nil {
					utility.LogError("Failed to get user setting", map[string]interface{}{"user_id": user.ID, "error": err.Error()})
					continue
				}
				if userSetting == nil || !userSetting.IsNotification {
					utility.LogInfo("Skipping TO recipient: global notification flag disabled", map[string]interface{}{"email": addr})
					continue
				}

				// Check notification-master setting (skipped when ForceSend is true)
				if !dto.ForceSend {
					enabled, err := ns.IsNotificationEnabled(user.ID, notificationMasterIdentifier)
					if err != nil || !enabled {
						utility.LogInfo("Skipping TO recipient: notification disabled for this module", map[string]interface{}{"email": addr})
						continue
					}
				}
				recipients = append(recipients, &models.EmailRecipient{
					RecipientType: "to",
					EmailAddress:  addr,
					IsExternal:    false,
					UserID:        &user.ID,
				})
			} else {
				recipients = append(recipients, &models.EmailRecipient{
					RecipientType: "to",
					EmailAddress:  addr,
					IsExternal:    true,
				})
			}
		}
	}

	// buildCCRecipients: CC is explicitly set by the sender — always stored regardless of
	// notification settings. Internal users get user_id set; external users get is_external=true.
	buildCCRecipients := func(addresses []string) {
		for _, addr := range utility.DedupeEmails(utility.SanitizeEmails(addresses)) {
			user, _ := ns.repo.GetUserByEmail(addr)

			if user != nil {
				recipients = append(recipients, &models.EmailRecipient{
					RecipientType: "cc",
					EmailAddress:  addr,
					IsExternal:    false,
					UserID:        &user.ID,
				})
			} else {
				recipients = append(recipients, &models.EmailRecipient{
					RecipientType: "cc",
					EmailAddress:  addr,
					IsExternal:    true,
				})
			}
		}
	}

	buildToRecipients(dto.EmailTo)
	buildCCRecipients(dto.CC)

	// Must have at least one TO recipient to enqueue — CC alone is not valid.
	hasTo := false
	for _, r := range recipients {
		if r.RecipientType == "to" {
			hasTo = true
			break
		}
	}
	if !hasTo {
		utility.LogInfo("EnqueueEmail: no eligible TO recipients, skipping insert", map[string]interface{}{
			"subject": dto.Subject,
		})
		return 0, nil
	}

	// Step 1: insert emails + recipients to get email.ID
	if err := ns.repo.EnqueueEmail(email, recipients, nil); err != nil {
		return 0, fmt.Errorf("failed to enqueue email: %w", err)
	}

	// Step 2: stage each attachment file under emails/{email_id}/ then insert rows.
	// Path: STORAGE_PATH/emails/{email_id}/{filename}
	// Worker resolves full path as: STORAGE_PATH + "/" + storage_key
	storagePath := config.GetEnv("STORAGE_PATH", "")
	var attachments []*models.EmailAttachment

	for i, att := range dto.Attachments {
		if att == "" {
			continue
		}

		subDir := fmt.Sprintf("emails/%d", email.ID)
		destDir := filepath.Join(storagePath, subDir)
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			utility.LogError("Failed to create attachment dir", map[string]interface{}{"error": err.Error(), "path": destDir})
			continue
		}

		fileName := filepath.Base(att)
		destPath := filepath.Join(destDir, fileName)

		// Always copy so the original file survives subsequent recipient iterations.
		// The caller is responsible for deleting the original after all emails are queued
		// (see cleanupPaths goroutine in the calling service).
		if copyErr := utility.CopyFile(att, destPath); copyErr != nil {
			utility.LogError("Failed to stage attachment", map[string]interface{}{"file": att, "error": copyErr.Error()})
			continue
		}

		storageKey := filepath.Join(subDir, fileName)
		contentType := mime.TypeByExtension(filepath.Ext(fileName))

		fi, _ := os.Stat(destPath)
		var sizeBytes *int64
		if fi != nil {
			sz := fi.Size()
			sizeBytes = &sz
		}

		attachments = append(attachments, &models.EmailAttachment{
			EmailID:     email.ID,
			StorageKey:  storageKey,
			FileName:    fileName,
			ContentType: &contentType,
			SizeBytes:   sizeBytes,
			SortOrder:   i,
		})
	}

	if len(attachments) > 0 {
		if err := ns.repo.InsertEmailAttachments(attachments); err != nil {
			return 0, fmt.Errorf("failed to insert email attachments: %w", err)
		}
	}

	utility.LogInfo("Email enqueued", map[string]interface{}{
		"email_id":   email.ID,
		"subject":    email.Subject,
		"recipients": len(recipients),
	})

	return email.ID, nil
}

func (ns *NotificationService) GetSocialMediaLinks() (*commonDTOs.SocialMediaLinksDTO, error) {
	socialLinks, err := ns.repo.GetSocialMediaLinks()
	if err != nil {
		return nil, err
	}
	return socialLinks, nil
}
