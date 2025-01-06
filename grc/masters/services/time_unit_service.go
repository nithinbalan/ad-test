package services

import (
	"errors"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type TimeUnitServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TimeUnit, error)
	GetByID(id uint) (*models.TimeUnit, error)
}

// Service and implementations
type TimeUnitService struct {
	repo repositories.TimeUnitRepositoryInterface
}

func NewTimeUnitService(repo repositories.TimeUnitRepositoryInterface) TimeUnitServiceInterface {
	return &TimeUnitService{
		repo: repo,
	}
}

func (ics *TimeUnitService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TimeUnit, error) {
	return ics.repo.GetAll(sort, sortBy, keyword, status)
}

func (ics TimeUnitService) GetByID(id uint) (*models.TimeUnit, error) {
	timeUnit, err := ics.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return timeUnit, nil
}
