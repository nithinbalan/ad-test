package services

import (
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
)

type AISettingServiceInterface interface {
	GetAll() ([]*models.AIModule, error)
	GetByModuleID(moduleID int) ([]dtos.AIModuleFeatureListDTO, error)
	EnableAISetting(currentUserID uint) error
	DisableAISetting(currentUserID uint) error
}

type AISettingService struct {
	repo repositories.AISettingRepositoryInterface
}

func NewAISettingService(repo repositories.AISettingRepositoryInterface) AISettingServiceInterface {
	return &AISettingService{
		repo: repo,
	}
}

func (amc *AISettingService) EnableAISetting(currentUserID uint) error {
	return amc.repo.UpdateSystemPreferenceAI(true, currentUserID)
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
