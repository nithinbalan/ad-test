package services

import (
	"jethur/grc/models"
	"jethur/grc/penetration-testing/repositories"
)

type PenetrationTestStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestStatus, error)
	GetByID(id uint) (*models.PenetrationTestStatus, error)
	GetByType(testType string) (*models.PenetrationTestStatus, error)
}

type PenetrationTestStatusService struct {
	repo repositories.PenetrationTestStatusRepositoryInterface
}

func NewPenetrationTestStatusService(repo repositories.PenetrationTestStatusRepositoryInterface) PenetrationTestStatusServiceInterface {
	return &PenetrationTestStatusService{repo: repo}
}

func (ptss *PenetrationTestStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestStatus, error) {
	return ptss.repo.GetAll(sort, sortBy, keyword, status)
}

func (ptss *PenetrationTestStatusService) GetByID(id uint) (*models.PenetrationTestStatus, error) {
	return ptss.repo.GetByID(id)
}

func (ptss *PenetrationTestStatusService) GetByType(testType string) (*models.PenetrationTestStatus, error) {
	return ptss.repo.GetByType(testType)
}
