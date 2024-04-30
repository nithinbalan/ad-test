package services

import (
	"jethur/grc/models"
	"jethur/grc/settings/repositories"
)

type CurrencyServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.Currency, error)
	GetByID(id uint) (*models.Currency, error)
}

// Service and implementations
type CurrencyService struct {
	repo repositories.CurrencyRepositoryInterface
}

func NewCurrencyService(repo repositories.CurrencyRepositoryInterface) CurrencyServiceInterface {
	return &CurrencyService{
		repo: repo,
	}
}

func (cs *CurrencyService) GetAll(sort string, sortBy string, keyword string) ([]*models.Currency, error) {
	currencies, err := cs.repo.GetAll(sort, sortBy, keyword)
	if err != nil {
		return nil, err
	}

	return currencies, nil
}

func (cs CurrencyService) GetByID(id uint) (*models.Currency, error) {
	return cs.repo.GetByID(id)
}
