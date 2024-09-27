package services

import (
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"

	"jethur/utility"
)

type EmailTemplateSettingsServiceInterface interface {
	GetEmailTemplateSettings() (*models.EmailTemplateSettings, error)
	Update(id uint, currentUserID uint, emailTemplateSettingsFormDTO *dtos.EmailTemplateSettingsFormDTO) error
}

type EmailTemplateSettingsService struct {
	repo           repositories.EmailTemplateSettingsRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewEmailTemplateSettingsService(repo repositories.EmailTemplateSettingsRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) EmailTemplateSettingsServiceInterface {
	return &EmailTemplateSettingsService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ets *EmailTemplateSettingsService) GetEmailTemplateSettings() (*models.EmailTemplateSettings, error) {
	return ets.repo.GetEmailTemplateSettings()
}

func (ets *EmailTemplateSettingsService) Update(id uint, currentUserID uint, emailTemplateSettingsFormDTO *dtos.EmailTemplateSettingsFormDTO) error {
	// Retrieve existing email template settings by ID
	existingSettings, err := ets.repo.GetEmailTemplateSettings()
	if err != nil {
		return err
	}

	// Update fields with values from DTO
	existingSettings.HeaderImageURL = emailTemplateSettingsFormDTO.HeaderImageURL
	existingSettings.SupportText = emailTemplateSettingsFormDTO.SupportText
	existingSettings.ClosingText = emailTemplateSettingsFormDTO.ClosingText
	existingSettings.FooterLogoURL = emailTemplateSettingsFormDTO.FooterLogoURL
	existingSettings.FooterAddress = emailTemplateSettingsFormDTO.FooterAddress
	existingSettings.FooterEmail = emailTemplateSettingsFormDTO.FooterEmail
	existingSettings.UpdatedByID = &currentUserID

	// Update email template settings in the repository
	err = ets.repo.Update(existingSettings)
	if err != nil {
		return err
	}

	go func() {
		err := ets.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.EmailTemplateSettings{}.TableName(),
			id,
			emailTemplateSettingsFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.EmailTemplateSettings{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
