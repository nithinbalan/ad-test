package services

import (
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type SourceServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Source, error)
	GetByID(id uint) (*models.Source, error)
	GetByType(typeStr string) (*models.Source, error)
}

// Service and implementations
type SourceService struct {
	repo repositories.SourceRepositoryInterface
}

func NewSourceService(repo repositories.SourceRepositoryInterface) SourceServiceInterface {
	return &SourceService{
		repo: repo,
	}
}

func (ss *SourceService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Source, error) {
	return ss.repo.GetAll(sort, sortBy, keyword, status)
}

func (ss SourceService) GetByID(id uint) (*models.Source, error) {
	return ss.repo.GetByID(id)
}

func (ss *SourceService) GetByType(typeStr string) (*models.Source, error) {
	// Fetch the action plan typeStr by its typeStr
	typeStrResult, err := ss.repo.GetByType(typeStr)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return typeStrResult, nil
}
