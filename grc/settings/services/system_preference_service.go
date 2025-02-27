package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"time"

	"jethur/utility"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"
)

type SystemPreferenceServiceInterface interface {
	GetSystemPreference() (*dtos.SystemPreferenceDTO, error)
	Update(id uint, currentUserID uint, systemPreferenceFormDTO *dtos.SystemPreferenceFormDTO) error
	UpdateLicenseKey(id uint, systemPreferenceFormDTO *dtos.LicenseKeyUpdateFormDTO) error

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

func (sps SystemPreferenceService) GetSystemPreference() (*dtos.SystemPreferenceDTO, error) {
	var systemPreferenceDetailDTO dtos.SystemPreferenceDTO

	systemPreference, err := sps.repo.GetSystemPreference()

	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(systemPreference, &systemPreferenceDetailDTO)
	if err != nil {
		return nil, err
	}

	return &systemPreferenceDetailDTO, nil
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

	var SystemPreferenceDTO dtos.SystemPreferenceDTO

	// Map the result to your DTO
	err = mapstructure.Decode(systemPreference, &SystemPreferenceDTO)
	if err != nil {
		return nil, err
	}

	// Combine both responses into a single struct
	combinedResponse := dtos.OrganizationPreferenceResponseDTO{
		SystemPreference: SystemPreferenceDTO,
		Integrations:     integrations,
	}

	return &combinedResponse, nil
}

func (sps *SystemPreferenceService) UpdateLicenseKey(id uint, systemPreferenceFormDTO *dtos.LicenseKeyUpdateFormDTO) error {
	// Retrieve existing system preference by ID
	existingSystemPreference, err := sps.repo.GetSystemPreference()
	if err != nil {
		return err
	}
	if err := sps.validateLicenseKey(*systemPreferenceFormDTO.LicenseKey); err != nil {
		return err
	}

	// Update fields with values from DTO
	existingSystemPreference.LicenseKey = systemPreferenceFormDTO.LicenseKey

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
			1,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.SystemPreference{}.TableName(),
				"recordID":  id,
				"userID":    1,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (sps *SystemPreferenceService) validateLicenseKey(licenseKey string) error {
	// Parse the JWT token
	token, err := jwt.Parse(licenseKey, func(token *jwt.Token) (interface{}, error) {
		return []byte("f0a3e3b61e7082367820944efeb9d81c7e242ec4d22dc6fa941aaf76d333f158"), nil
	})

	if err != nil {
		return errors.New("invalid license key")
	}

	// Validate token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Validate Expiry Date (expiry_date)
		expiryDateStr, ok := claims["expiry_date"].(string)
		if !ok || expiryDateStr == "" {
			return errors.New("missing expiry date in license key")
		}

		// Convert expiry_date string to time.Time format
		expiryDate, err := time.Parse("2006-01-02", expiryDateStr) // Ensure format matches the stored format
		if err != nil {
			return errors.New("invalid expiry date format, expected YYYY-MM-DD")
		}

		// Check if expiry_date is in the future
		if expiryDate.Before(time.Now()) {
			return errors.New("license key has expired")
		}

		// Validate API URL (if needed)

		if apiURL, ok := claims["api_url"].(string); ok {
			systemAPIURL := utility.GetApiURL() // Change this to your actual API URL
			if apiURL != systemAPIURL {
				return errors.New("license key does not match system API URL")
			}
		}

		return nil
	}

	return errors.New("invalid license key format")
}
