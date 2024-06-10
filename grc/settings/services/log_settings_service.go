package services

import (
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"

	"jethur/utility"
)

type LogSettingsServiceInterface interface {
	GetLogSettings() (*models.LogSettings, error)
	Update(currentUserID uint, logSettingsFormDTO *dtos.LogSettingsFormDTO) error
}

type LogSettingsService struct {
	repo           repositories.LogSettingsRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewLogSettingsService(repo repositories.LogSettingsRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) LogSettingsServiceInterface {
	return &LogSettingsService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (lss *LogSettingsService) GetLogSettings() (*models.LogSettings, error) {
	return lss.repo.GetLogSettings()
}

func (lss *LogSettingsService) Update(currentUserID uint, logSettingsFormDTO *dtos.LogSettingsFormDTO) error {
	existingLogSettings, err := lss.repo.GetLogSettings()
	if err != nil {
		return err
	}

	// Update fields with values from DTO
	existingLogSettings.AuditLogRetentionDays = logSettingsFormDTO.AuditLogRetentionDays
	existingLogSettings.EmailLogRetentionDays = logSettingsFormDTO.EmailLogRetentionDays
	existingLogSettings.SystemLogRetentionDays = logSettingsFormDTO.SystemLogRetentionDays
	existingLogSettings.UpdatedBy = &currentUserID

	// Update log settings in the repository
	err = lss.repo.Update(existingLogSettings)
	if err != nil {
		return err
	}

	go func() {
		err := lss.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.LogSettings{}.TableName(),
			existingLogSettings.ID,
			logSettingsFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.LogSettings{}.TableName(),
				"recordID":  existingLogSettings.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
