package services

import (
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type CountryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.Country, error)
	GetByID(id uint) (*models.Country, error)
}

// Service and implementations
type CountryService struct {
	repo repositories.CountryRepositoryInterface
}

func NewCountryService(repo repositories.CountryRepositoryInterface) CountryServiceInterface {
	return &CountryService{
		repo: repo,
	}
}

func (cs *CountryService) GetAll(sort string, sortBy string, keyword string) ([]*models.Country, error) {
	countries, err := cs.repo.GetAll(sort, sortBy, keyword)
	if err != nil {
		return nil, err
	}

	return countries, nil
}

func (cs CountryService) GetByID(id uint) (*models.Country, error) {
	return cs.repo.GetByID(id)
}
