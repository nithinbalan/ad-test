package services

import (
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
)

type DistributionStatusServiceInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.DistributionStatus, error)
	GetByID(id uint) (*models.DistributionStatus, error)
}

type DistributionStatusService struct {
	repo repositories.DistributionStatusRepositoryInterface
}

func NewDistributionStatusService(repo repositories.DistributionStatusRepositoryInterface) DistributionStatusServiceInterface {
	return &DistributionStatusService{repo: repo}
}

func (s *DistributionStatusService) GetAll(sort, sortBy, keyword, status string) ([]*models.DistributionStatus, error) {
	return s.repo.GetAll(sort, sortBy, keyword, status)
}

func (s *DistributionStatusService) GetByID(id uint) (*models.DistributionStatus, error) {
	return s.repo.GetByID(id)
}
