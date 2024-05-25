package services

import (
	"errors"
	"jethur/grc/incident/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentPostReviewEffectivenessServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentPostReviewEffectiveness, error)
	GetByID(id uint) (*models.IncidentPostReviewEffectiveness, error)
}

// Service and implementations
type IncidentPostReviewEffectivenessService struct {
	repo repositories.IncidentPostReviewEffectivenessRepositoryInterface
}

func NewIncidentPostReviewEffectivenessService(repo repositories.IncidentPostReviewEffectivenessRepositoryInterface) IncidentPostReviewEffectivenessServiceInterface {
	return &IncidentPostReviewEffectivenessService{
		repo: repo,
	}
}

func (ipres *IncidentPostReviewEffectivenessService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentPostReviewEffectiveness, error) {
	return ipres.repo.GetAll(sort, sortBy, keyword, status)
}

func (ipres IncidentPostReviewEffectivenessService) GetByID(id uint) (*models.IncidentPostReviewEffectiveness, error) {
	incidentPostReviewEffectiveness, err := ipres.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentPostReviewEffectiveness, nil
}
