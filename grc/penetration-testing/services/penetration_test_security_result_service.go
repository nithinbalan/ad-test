package services

import (
	"jethur/grc/models"
	"jethur/grc/penetration-testing/repositories"
)

type PenetrationTestSecurityResultServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestSecurityResult, error)
	GetByID(id uint) (*models.PenetrationTestSecurityResult, error)
}

type PenetrationTestSecurityResultService struct {
	repo repositories.PenetrationTestSecurityResultRepositoryInterface
}

func NewPenetrationTestSecurityResultService(repo repositories.PenetrationTestSecurityResultRepositoryInterface) PenetrationTestSecurityResultServiceInterface {
	return &PenetrationTestSecurityResultService{repo: repo}
}

func (ptss *PenetrationTestSecurityResultService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestSecurityResult, error) {
	return ptss.repo.GetAll(sort, sortBy, keyword, status)
}

func (ptss *PenetrationTestSecurityResultService) GetByID(id uint) (*models.PenetrationTestSecurityResult, error) {
	return ptss.repo.GetByID(id)
}
