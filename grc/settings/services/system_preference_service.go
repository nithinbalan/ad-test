package services

import (
	"fmt"
	logModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/models"
	"jethur/grc/settings/repositories"
	"jethur/utility"
)

type SystemPreferenceServiceInterface interface {
	GetByID(id uint) (*models.SystemPreference, error)
	Update(id uint, currentUserID uint, systemPreferenceFormDTO *dtos.SystemPreferenceFormDTO) error
}

// Service and implementations
type SystemPreferenceService struct {
	repo           repositories.SystemPreferenceRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
	currencySvc    CurrencyServiceInterface
}

func NewSystemPreferenceService(repo repositories.SystemPreferenceRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface, currencySvc CurrencyServiceInterface) SystemPreferenceServiceInterface {
	return &SystemPreferenceService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		currencySvc:    currencySvc,
	}
}

func (sps SystemPreferenceService) GetByID(id uint) (*models.SystemPreference, error) {
	return sps.repo.GetByID(id)
}

func (sps *SystemPreferenceService) Update(id uint, currentUserID uint, systemPreferenceFormDTO *dtos.SystemPreferenceFormDTO) error {
	// Retrieve existing system preference by ID
	existingSystemPreference, err := sps.repo.GetByID(id)
	if err != nil {
		return err
	}
	if err := sps.validateSystemPreferenceFormDTO(systemPreferenceFormDTO); err != nil {
		return err
	}

	// Update fields with values from DTO
	existingSystemPreference.CurrencyID = systemPreferenceFormDTO.CurrencyID
	existingSystemPreference.DateFormat = systemPreferenceFormDTO.DateFormat
	existingSystemPreference.TimeFormat = systemPreferenceFormDTO.TimeFormat
	existingSystemPreference.DocUploadSize = systemPreferenceFormDTO.DocUploadSize
	existingSystemPreference.LogoUploadSize = systemPreferenceFormDTO.LogoUploadSize
	existingSystemPreference.IdleTimeout = systemPreferenceFormDTO.IdleTimeout
	existingSystemPreference.IsMaintanceModeOn = systemPreferenceFormDTO.IsMaintanceModeOn
	existingSystemPreference.UpdatedBy = &currentUserID

	// Update system preference in the repository
	err = sps.repo.Update(existingSystemPreference)
	if err != nil {
		return err
	}

	go func() {
		err := sps.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.SystemPreference{}.TableName(),
			id,
			systemPreferenceFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.UpdateAction,
				"tableName": models.SystemPreference{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (sps *SystemPreferenceService) validateSystemPreferenceFormDTO(systemPreferenceFormDTO *dtos.SystemPreferenceFormDTO) error {
	if _, err := sps.currencySvc.GetByID(systemPreferenceFormDTO.CurrencyID); err != nil {
		return utility.NewValidationError(fmt.Errorf("currency does not exist"), "not_exist", "currency")
	}
	return nil
}
