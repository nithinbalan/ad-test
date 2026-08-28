package services

import (
	"jethur/grc/models"
	"jethur/grc/penetration-testing/repositories"
)

type PenetrationTestExecutionStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestExecutionStatus, error)
}
type PenetrationTestExecutionStatusService struct {
	repo repositories.PenetrationTestExecutionStatusRepositoryInterface
}

func NewPenetrationTestExecutionStatusService(repo repositories.PenetrationTestExecutionStatusRepositoryInterface) PenetrationTestExecutionStatusServiceInterface {
	return &PenetrationTestExecutionStatusService{repo: repo}
}
func (ptess *PenetrationTestExecutionStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestExecutionStatus, error) {
	return ptess.repo.GetAll(sort, sortBy, keyword, status)
}
