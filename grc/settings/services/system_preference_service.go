package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"time"

	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type SystemPreferenceServiceInterface interface {
	GetSystemPreference() (*dtos.SystemPreferenceDTO, error)
	GetByID(id uint) (*models.SystemPreference, error)
	Update(id uint, currentUserID uint, systemPreferenceFormDTO *dtos.SystemPreferenceFormDTO) error
	UpdateLicenseKey(systemPreferenceFormDTO *dtos.LicenseKeyUpdateFormDTO) error
	UpdateBalanceUserCount(systemPreference *models.SystemPreference, tx utility.DatabaseTransaction) error

	GetOrganizationPreference(params *dtos.IntegrationListQueryParamsDTO) (*dtos.OrganizationPreferenceResponseDTO, error)
}

// Service and implementations
type SystemPreferenceService struct {
	repo           repositories.SystemPreferenceRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
	currencySvc    masterServices.CurrencyServiceInterface
	integrationSvc IntegrationServiceInterface
	moduleSvc      ModuleServiceInterface
}

func NewSystemPreferenceService(repo repositories.SystemPreferenceRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface, currencySvc masterServices.CurrencyServiceInterface, integrationSvc IntegrationServiceInterface, moduleSvc ModuleServiceInterface) SystemPreferenceServiceInterface {
	return &SystemPreferenceService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		currencySvc:    currencySvc,
		integrationSvc: integrationSvc,
		moduleSvc:      moduleSvc,
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

	// Populate License Key as an object
	if systemPreference.LicenseKey != nil {
		claims, err := utility.DecryptJWTToken(*systemPreference.LicenseKey)
		if err != nil {
			return nil, err
		}

		expiryDate := systemPreference.LicenseKeyExpiryDate.Format("2006-01-02")

		licenseKey := *systemPreference.LicenseKey
		licenseDTO := &dtos.LicenseKeyDTO{
			LicenseKey:                 &licenseKey,
			LicenseKeyExpiryDateString: &expiryDate,
			BalanceLicensedUserCount:   systemPreference.BalanceLicensedUserCount,
			TotalLicensedUserCount:     systemPreference.TotalLicensedUserCount,
			EnabledModules:             claims["enabled_modules"].([]string),
		}

		systemPreferenceDetailDTO.LicenseKeyDetails = licenseDTO
	}

	return &systemPreferenceDetailDTO, nil
}

func (sps SystemPreferenceService) GetByID(id uint) (*models.SystemPreference, error) {
	systemPreference, err := sps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return systemPreference, nil
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

func (sps *SystemPreferenceService) UpdateLicenseKey(systemPreferenceFormDTO *dtos.LicenseKeyUpdateFormDTO) error {
	// Retrieve existing system preference by ID
	existingSystemPreference, err := sps.repo.GetSystemPreference()
	if err != nil {
		return err
	}

	if err := sps.validateLicenseKey(*systemPreferenceFormDTO.LicenseKey); err != nil {
		return err
	}

	claims, err := utility.DecryptJWTToken(*systemPreferenceFormDTO.LicenseKey)
	if err != nil {
		return err
	}

	// Extract expiry_date from claims
	expiryDateStr, ok := claims["expiry_date"].(string) // Check if it's a string
	if !ok {
		return utility.NewValidationError(fmt.Errorf("invalid token"), "invalid", "expiry date")
	}

	// Convert expiryDateStr to time.Time
	expiryDate, err := time.Parse("2006-01-02", expiryDateStr)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid_date", "expiry date")
	}

	// Get max_users from JWT claims
	maxUsersFloat, ok := claims["max_users"].(float64) // JWT stores numbers as float64
	if !ok {
		return utility.NewValidationError(fmt.Errorf("invalid token"), "invalid", "max users")
	}
	maxUsers := int(maxUsersFloat)

	activeUsersCount, err := sps.repo.GetActiveUserCount()
	if err != nil {
		return err
	}

	balanceLicensedUserCount := maxUsers - int(activeUsersCount)
	if balanceLicensedUserCount < 0 {
		balanceLicensedUserCount = 0
	}

	// Update fields with values from DTO
	existingSystemPreference.LicenseKey = systemPreferenceFormDTO.LicenseKey
	existingSystemPreference.LicenseKeyExpiryDate = &expiryDate
	existingSystemPreference.BalanceLicensedUserCount = &balanceLicensedUserCount
	existingSystemPreference.TotalLicensedUserCount = &maxUsers

	// Update system preference in the repository
	err = sps.repo.Update(existingSystemPreference)
	if err != nil {
		return err
	}

	rawModules, ok := claims["enabled_modules"]
	if !ok {
		return utility.NewValidationError(fmt.Errorf("invalid token"), "invalid", "enabled modules")
	}

	var enabledModules []string
	if str, ok := rawModules.(string); ok {
		if err := json.Unmarshal([]byte(str), &enabledModules); err != nil {
			return utility.NewValidationError(fmt.Errorf("failed to parse enabled_modules"), "invalid", "enabled modules")
		}
	} else {
		return utility.NewValidationError(fmt.Errorf("invalid enabled_modules format"), "invalid", "enabled modules")
	}

	// Enable modules
	for _, moduleIdentifier := range enabledModules {
		moduleID, err := sps.moduleSvc.GetModuleIDByIdentifier(moduleIdentifier)
		if err != nil {
			return err
		}

		// Ensure EnableModule returns an error if it fails
		if _, err := sps.moduleSvc.EnableModule(moduleID); err != nil {
			return err
		}
	}

	go func() {
		err := sps.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.SystemPreference{}.TableName(),
			existingSystemPreference.ID,
			systemPreferenceFormDTO,
			1,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.SystemPreference{}.TableName(),
				"recordID":  existingSystemPreference.ID,
				"userID":    1,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (sps *SystemPreferenceService) UpdateBalanceUserCount(systemPreference *models.SystemPreference, tx utility.DatabaseTransaction) error {

	err := sps.repo.Update(systemPreference)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (sps *SystemPreferenceService) validateLicenseKey(licenseKey string) error {
	claims, err := utility.DecryptJWTToken(licenseKey)
	if err != nil {
		return err
	}

	// Validate Expiry Date (expiry_date)
	expiryDateStr, ok := claims["expiry_date"].(string)
	if !ok || expiryDateStr == "" {
		return utility.NewValidationError(fmt.Errorf("invalid token"), "not_exist", "expiry date")
	}

	// Convert expiry_date string to time.Time format
	expiryDate, err := time.Parse("2006-01-02", expiryDateStr) // Ensure format matches the stored format
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid_date", "expiry date")
	}

	// Check if expiry_date is in the future
	if expiryDate.Before(time.Now()) {
		return utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid", "expiry date")
	}

	// Validate API URL (if needed)

	if apiURL, ok := claims["api_url"].(string); ok {
		systemAPIURL := utility.GetApiURL() // Change this to your actual API URL
		if apiURL != systemAPIURL {
			return utility.NewValidationError(fmt.Errorf("invalid api key"), "invalid", "api url")
		}
	}

	return nil
}
