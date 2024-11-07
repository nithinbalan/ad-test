package services

import (
	"fmt"
	"jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"

	"jethur/utility"
)

type SystemPreferenceServiceInterface interface {
	GetSystemPreference() (*models.SystemPreference, error)
	Update(id uint, currentUserID uint, systemPreferenceFormDTO *dtos.SystemPreferenceFormDTO) error

	GetOrganizationPreference(params *dtos.IntegrationListQueryParamsDTO) (*dtos.OrganizationPreferenceResponseDTO, error)
}

// Service and implementations
type SystemPreferenceService struct {
	repo           repositories.SystemPreferenceRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
	currencySvc    masterServices.CurrencyServiceInterface
	integrationSvc IntegrationServiceInterface
}

func NewSystemPreferenceService(repo repositories.SystemPreferenceRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface, currencySvc masterServices.CurrencyServiceInterface, integrationSvc IntegrationServiceInterface) SystemPreferenceServiceInterface {
	return &SystemPreferenceService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		currencySvc:    currencySvc,
		integrationSvc: integrationSvc,
	}
}

func (sps SystemPreferenceService) GetSystemPreference() (*models.SystemPreference, error) {
	return sps.repo.GetSystemPreference()
}

func (sps *SystemPreferenceService) Update(id uint, currentUserID uint, systemPreferenceFormDTO *dtos.SystemPreferenceFormDTO) error {
	// Retrieve existing system preference by ID
	existingSystemPreference, err := sps.repo.GetSystemPreference()
	if err != nil {
		return err
	}
	if err := sps.validateSystemPreferenceFormDTO(systemPreferenceFormDTO); err != nil {
		return err
	}

	// Update fields with values from DTO
	existingSystemPreference.CurrencyID = systemPreferenceFormDTO.CurrencyID
	existingSystemPreference.TimezoneID = systemPreferenceFormDTO.TimezoneID
	existingSystemPreference.DateFormat = systemPreferenceFormDTO.DateFormat
	existingSystemPreference.TimeFormat = systemPreferenceFormDTO.TimeFormat
	existingSystemPreference.DocUploadSize = systemPreferenceFormDTO.DocUploadSize
	existingSystemPreference.LogoUploadSize = systemPreferenceFormDTO.LogoUploadSize
	existingSystemPreference.IdleTimeout = systemPreferenceFormDTO.IdleTimeout
	existingSystemPreference.IsMaintenanceModeOn = systemPreferenceFormDTO.IsMaintenanceModeOn
	existingSystemPreference.UpdatedByID = &currentUserID

	// Update system preference in the repository
	err = sps.repo.Update(existingSystemPreference)
	if err != nil {
		return err
	}

	go func() {
		err := sps.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.SystemPreference{}.TableName(),
			id,
			systemPreferenceFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
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

func (sps SystemPreferenceService) GetOrganizationPreference(params *dtos.IntegrationListQueryParamsDTO) (*dtos.OrganizationPreferenceResponseDTO, error) {
	systemPreference, err := sps.repo.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	// Call GetAll to get integrations list
	integrations, err := sps.integrationSvc.GetAll(params)
	if err != nil {
		return nil, err
	}

	// Combine both responses into a single struct
	combinedResponse := dtos.OrganizationPreferenceResponseDTO{
		SystemPreference: systemPreference,
		Integrations:     integrations,
	}

	return &combinedResponse, nil
}
