package services

import (
	"errors"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type TimelinessServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Timeliness, error)
	GetByID(id uint) (*models.Timeliness, error)
	GetByType(status string) (*models.Timeliness, error)
}

// Service and implementations
type TimelinessService struct {
	repo repositories.TimelinessRepositoryInterface
}

func NewTimelinessService(repo repositories.TimelinessRepositoryInterface) TimelinessServiceInterface {
	return &TimelinessService{
		repo: repo,
	}
}

func (tess *TimelinessService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Timeliness, error) {
	return tess.repo.GetAll(sort, sortBy, keyword, status)
}

func (tess TimelinessService) GetByID(id uint) (*models.Timeliness, error) {
	timeliness, err := tess.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return timeliness, nil
}

func (tess *TimelinessService) GetByType(status string) (*models.Timeliness, error) {
	timeliness, err := tess.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return timeliness, nil
}
