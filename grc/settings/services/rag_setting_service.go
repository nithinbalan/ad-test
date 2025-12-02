package services

import (
	"jethur/grc/models"
	"jethur/grc/settings/repositories"
)

type RAGSettingServiceInterface interface {
	GetAllRAGSetting() ([]*models.RagSetting, error)
}

// Service and implementations
type RAGSettingService struct {
	repo repositories.RAGSettingRepositoryInterface
}

func NewRAGSettingService(repo repositories.RAGSettingRepositoryInterface) RAGSettingServiceInterface {
	return &RAGSettingService{
		repo: repo,
	}
}

func (rss *RAGSettingService) GetAllRAGSetting() ([]*models.RagSetting, error) {
	return rss.repo.GetAllRAGSettings()
}
