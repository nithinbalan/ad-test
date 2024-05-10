package services

import (
	"jethur/grc/models"
	"jethur/grc/settings/repositories"
)

type TimezoneServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.Timezone, error)
	GetByID(id uint) (*models.Timezone, error)
}

// Service and implementations
type TimezoneService struct {
	repo repositories.TimezoneRepositoryInterface
}

func NewTimezoneService(repo repositories.TimezoneRepositoryInterface) TimezoneServiceInterface {
	return &TimezoneService{
		repo: repo,
	}
}

func (ts *TimezoneService) GetAll(sort string, sortBy string, keyword string) ([]*models.Timezone, error) {
	timezones, err := ts.repo.GetAll(sort, sortBy, keyword)
	if err != nil {
		return nil, err
	}

	return timezones, nil
}

func (ts TimezoneService) GetByID(id uint) (*models.Timezone, error) {
	return ts.repo.GetByID(id)
}
