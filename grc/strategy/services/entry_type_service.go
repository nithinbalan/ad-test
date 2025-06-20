package services

import (
	"jethur/grc/models"
	"jethur/grc/strategy/repositories"
)

type EntryTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.EntryType, error)
	GetByID(id uint) (*models.EntryType, error)
	GetByType(status string) (*models.EntryType, error) // Renamed method

}

// Service and implementations
type EntryTypeService struct {
	repo repositories.EntryTypeRepositoryInterface
}

func NewEntryTypeService(repo repositories.EntryTypeRepositoryInterface) EntryTypeServiceInterface {
	return &EntryTypeService{
		repo: repo,
	}
}

func (ets *EntryTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.EntryType, error) {
	return ets.repo.GetAll(sort, sortBy, keyword, status)
}

func (ets EntryTypeService) GetByID(id uint) (*models.EntryType, error) {
	return ets.repo.GetByID(id)
}

func (ets *EntryTypeService) GetByType(status string) (*models.EntryType, error) {
	// Fetch the kpi status by its status
	statusResult, err := ets.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
