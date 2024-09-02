package services

import (
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
	"jethur/utility"
)

type SeverityLevelServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SeverityLevel, error)
	GetByID(id uint) (*models.SeverityLevel, error)
	GetMaxSeverityWeight(tx ...utility.DatabaseTransaction) (int, error)
}

// Service and implementations
type SeverityLevelService struct {
	repo repositories.SeverityLevelRepositoryInterface
}

func NewSeverityLevelService(repo repositories.SeverityLevelRepositoryInterface) SeverityLevelServiceInterface {
	return &SeverityLevelService{
		repo: repo,
	}
}

func (sls *SeverityLevelService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SeverityLevel, error) {
	return sls.repo.GetAll(sort, sortBy, keyword, status)
}

func (sls SeverityLevelService) GetByID(id uint) (*models.SeverityLevel, error) {
	return sls.repo.GetByID(id)
}

func (sls SeverityLevelService) GetMaxSeverityWeight(tx ...utility.DatabaseTransaction) (int, error) {
	return sls.repo.GetMaxSeverityWeight(tx...)
}
