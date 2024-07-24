package services

import (
	"errors"
	"jethur/grc/audit/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditorRiskRatingServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditorRiskRating, error)
	GetByID(id uint) (*models.AuditorRiskRating, error)
	GetByType(riskRatingType string) (*models.AuditorRiskRating, error)
}

// Service and implementations
type AuditorRiskRatingService struct {
	repo repositories.AuditorRiskRatingRepositoryInterface
}

func NewAuditorRiskRatingService(repo repositories.AuditorRiskRatingRepositoryInterface) AuditorRiskRatingServiceInterface {
	return &AuditorRiskRatingService{
		repo: repo,
	}
}

func (arrs *AuditorRiskRatingService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditorRiskRating, error) {
	return arrs.repo.GetAll(sort, sortBy, keyword, status)
}

func (arrs AuditorRiskRatingService) GetByID(id uint) (*models.AuditorRiskRating, error) {
	auditorRiskRating, err := arrs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditorRiskRating, nil
}

func (arrs *AuditorRiskRatingService) GetByType(riskRatingType string) (*models.AuditorRiskRating, error) {
	riskRatingResult, err := arrs.repo.GetByType(riskRatingType)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return riskRatingResult, nil
}
