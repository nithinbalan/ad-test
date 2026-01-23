package services

import (
	"jethur/grc/compliance/repositories"
	"jethur/grc/models"
)

// Interface
type ImportanceClassificationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ImportanceClassification, error)
	GetByID(id uint) (*models.ImportanceClassification, error)
	GetByType(status string) (*models.ImportanceClassification, error)
}

// Service implementation
type ImportanceClassificationService struct {
	repo repositories.ImportanceClassificationRepositoryInterface
}

func NewImportanceClassificationService(
	repo repositories.ImportanceClassificationRepositoryInterface,
) ImportanceClassificationServiceInterface {
	return &ImportanceClassificationService{
		repo: repo,
	}
}

func (ics *ImportanceClassificationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ImportanceClassification, error) {
	return ics.repo.GetAll(sort, sortBy, keyword, status)
}

func (ics *ImportanceClassificationService) GetByID(id uint) (*models.ImportanceClassification, error) {
	return ics.repo.GetByID(id)
}

func (ics *ImportanceClassificationService) GetByType(status string) (*models.ImportanceClassification, error) {
	classification, err := ics.repo.GetByType(status)
	if err != nil {
		return nil, err
	}
	return classification, nil
}
