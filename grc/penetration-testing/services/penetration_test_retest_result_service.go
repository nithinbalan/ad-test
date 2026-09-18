package services

import (
	"jethur/grc/models"
	"jethur/grc/penetration-testing/repositories"
)

type PenetrationTestRetestResultServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestRetestResult, error)
	GetByID(id uint) (*models.PenetrationTestRetestResult, error)
}

type PenetrationTestRetestResultService struct {
	repo repositories.PenetrationTestRetestResultRepositoryInterface
}

func NewPenetrationTestRetestResultService(repo repositories.PenetrationTestRetestResultRepositoryInterface) PenetrationTestRetestResultServiceInterface {
	return &PenetrationTestRetestResultService{repo: repo}
}

func (ptrrs *PenetrationTestRetestResultService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestRetestResult, error) {
	return ptrrs.repo.GetAll(sort, sortBy, keyword, status)
}

func (ptrrs *PenetrationTestRetestResultService) GetByID(id uint) (*models.PenetrationTestRetestResult, error) {
	return ptrrs.repo.GetByID(id)
}
