package services

import (
	"errors"
	"jethur/grc/compliance/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ControlEffectivenessRatingServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.ControlEffectivenessRating, error)
	GetByID(id uint) (*models.ControlEffectivenessRating, error)
	GetByType(typeFilter string) (*models.ControlEffectivenessRating, error)
}

// Service and implementations
type ControlEffectivenessRatingService struct {
	repo repositories.ControlEffectivenessRatingRepositoryInterface
}

func NewControlEffectivenessRatingService(repo repositories.ControlEffectivenessRatingRepositoryInterface) ControlEffectivenessRatingServiceInterface {
	return &ControlEffectivenessRatingService{
		repo: repo,
	}
}

func (cers *ControlEffectivenessRatingService) GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.ControlEffectivenessRating, error) {
	return cers.repo.GetAll(sort, sortBy, keyword, status, typeFilter)
}

func (cers ControlEffectivenessRatingService) GetByID(id uint) (*models.ControlEffectivenessRating, error) {
	controlEffectivenessRating, err := cers.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return controlEffectivenessRating, nil

}

func (cers ControlEffectivenessRatingService) GetByType(typeFilter string) (*models.ControlEffectivenessRating, error) {
	controlEffectivenessRating, err := cers.repo.GetByType(typeFilter)
	if err != nil {
		return nil, err
	}

	return controlEffectivenessRating, nil
}
