package services

import (
	"fmt"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"
)

type AISettingServiceInterface interface {
	GetAll() ([]*models.AIModule, error)
	GetByModuleID(moduleID int) ([]dtos.AIModuleFeatureListDTO, error)
	EnableAISetting(currentUserID uint) error
	DisableAISetting(currentUserID uint) error
}

type AISettingService struct {
	repo                repositories.AISettingRepositoryInterface
	systemPreferenceSvc SystemPreferenceServiceInterface
}

func NewAISettingService(repo repositories.AISettingRepositoryInterface, systemPreferenceSvc SystemPreferenceServiceInterface) AISettingServiceInterface {
	return &AISettingService{
		repo:                repo,
		systemPreferenceSvc: systemPreferenceSvc,
	}
}

func (amc *AISettingService) EnableAISetting(currentUserID uint) error {
	existingSystemPreference, err := amc.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return err
	}
	if existingSystemPreference.LicenseKeyDetails != nil {
		// Extract allowed modules from the license key
		claims, err := utility.DecryptLicenseKey(*existingSystemPreference.LicenseKeyDetails.LicenseKey)
		if err != nil {
			return err
		}

		// Extract allowed addons (can be a string or array)
		rawAddons, ok := claims["enabled_addons"]
		if !ok {
			return fmt.Errorf("invalid or missing enabled addons in license key")
		}

		allowedAddons, ok := rawAddons.([]string)
		if !ok {
			return utility.NewValidationError(fmt.Errorf("invalid enabled_addons format"), "invalid", "enabled ai")
		}

		// Check if the addon is allowed in the license
		isAllowed := false
		for _, allowedAddon := range allowedAddons {
			if allowedAddon == "jethur_ai" {
				isAllowed = true
				amc.repo.UpdateSystemPreferenceAI(true, currentUserID)
			}
		}

		if !isAllowed {
			return utility.NewValidationError(fmt.Errorf("invalid ai"), "ai_not_allowed", "Jethur AI")
		}
	} else {
		amc.repo.UpdateSystemPreferenceAI(true, currentUserID)
	}

	return err
}

func (amc *AISettingService) DisableAISetting(currentUserID uint) error {
	return amc.repo.UpdateSystemPreferenceAI(false, currentUserID)
}

func (amc *AISettingService) GetAll() ([]*models.AIModule, error) {
	return amc.repo.GetAll()
}

func (amc *AISettingService) GetByModuleID(moduleID int) ([]dtos.AIModuleFeatureListDTO, error) {
	return amc.repo.GetFeatureListByModuleID(moduleID)
}
