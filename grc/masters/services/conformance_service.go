package services

import (
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type ConformanceServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Conformance, error)
	GetByID(id uint) (*models.Conformance, error)
	GetByType(status string) (*models.Conformance, error)
}

// Service and implementations
type ConformanceService struct {
	repo repositories.ConformanceRepositoryInterface
}

func NewConformanceService(repo repositories.ConformanceRepositoryInterface) ConformanceServiceInterface {
	return &ConformanceService{
		repo: repo,
	}
}

func (cs *ConformanceService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Conformance, error) {
	return cs.repo.GetAll(sort, sortBy, keyword, status)
}

func (cs ConformanceService) GetByID(id uint) (*models.Conformance, error) {
	return cs.repo.GetByID(id)
}

func (cs *ConformanceService) GetByType(conformance string) (*models.Conformance, error) {
	// Fetch the conformance by its conformance
	conformanceResult, err := cs.repo.GetByType(conformance)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return conformanceResult, nil
}
