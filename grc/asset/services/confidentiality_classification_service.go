package services

import (
	"errors"
	"jethur/grc/asset/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ConfidentialityClassificationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ConfidentialityClassification, error)
	GetByID(id uint) (*models.ConfidentialityClassification, error)
	GetByType(status string) (*models.ConfidentialityClassification, error)
}

type ConfidentialityClassificationService struct {
	repo repositories.ConfidentialityClassificationRepositoryInterface
}

func NewConfidentialityClassificationService(
	repo repositories.ConfidentialityClassificationRepositoryInterface,
) ConfidentialityClassificationServiceInterface {
	return &ConfidentialityClassificationService{
		repo: repo,
	}
}

func (ccs *ConfidentialityClassificationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ConfidentialityClassification, error) {
	return ccs.repo.GetAll(sort, sortBy, keyword, status)
}

func (ccs ConfidentialityClassificationService) GetByID(id uint) (*models.ConfidentialityClassification, error) {
	confidentialityClassification, err := ccs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return confidentialityClassification, nil
}

func (ccs *ConfidentialityClassificationService) GetByType(status string) (*models.ConfidentialityClassification, error) {
	confidentialityClassification, err := ccs.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return confidentialityClassification, nil
}
