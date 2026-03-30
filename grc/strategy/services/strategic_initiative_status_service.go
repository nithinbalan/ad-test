package services

import (
	"jethur/grc/models"
	"jethur/grc/strategy/repositories"
)

type StrategicInitiativeStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StrategicInitiativeStatus, error)
	GetByID(id uint) (*models.StrategicInitiativeStatus, error)
}

type StrategicInitiativeStatusService struct {
	repo repositories.StrategicInitiativeStatusRepositoryInterface
}

func NewStrategicInitiativeStatusService(repo repositories.StrategicInitiativeStatusRepositoryInterface) StrategicInitiativeStatusServiceInterface {
	return &StrategicInitiativeStatusService{
		repo: repo,
	}
}

func (s *StrategicInitiativeStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StrategicInitiativeStatus, error) {
	return s.repo.List(sort, sortBy, keyword, status)
}

func (s StrategicInitiativeStatusService) GetByID(id uint) (*models.StrategicInitiativeStatus, error) {
	return s.repo.GetByID(id)
}
