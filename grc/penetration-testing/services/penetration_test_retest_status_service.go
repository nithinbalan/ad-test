package services

import (
	"jethur/grc/models"
	"jethur/grc/penetration-testing/repositories"
)

type PenetrationTestRetestStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestRetestStatus, error)
	GetByID(id uint) (*models.PenetrationTestRetestStatus, error)
	GetByType(statusType string) (*models.PenetrationTestRetestStatus, error)
}

type PenetrationTestRetestStatusService struct {
	repo repositories.PenetrationTestRetestStatusRepositoryInterface
}

func NewPenetrationTestRetestStatusService(repo repositories.PenetrationTestRetestStatusRepositoryInterface) PenetrationTestRetestStatusServiceInterface {
	return &PenetrationTestRetestStatusService{repo: repo}
}

func (ptrss *PenetrationTestRetestStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestRetestStatus, error) {
	return ptrss.repo.GetAll(sort, sortBy, keyword, status)
}

func (ptrss *PenetrationTestRetestStatusService) GetByID(id uint) (*models.PenetrationTestRetestStatus, error) {
	return ptrss.repo.GetByID(id)
}

func (ptrss *PenetrationTestRetestStatusService) GetByType(statusType string) (*models.PenetrationTestRetestStatus, error) {
	return ptrss.repo.GetByType(statusType)
}
