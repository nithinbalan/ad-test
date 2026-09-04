package services

import (
	"jethur/grc/models"
	"jethur/grc/project/repositories"
	"jethur/utility"
)

type ProjectSecurityGateReadinessStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateReadinessStatus, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateReadinessStatus, error)
	GetByType(readinessStatusType string, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateReadinessStatus, error)
}

type ProjectSecurityGateReadinessStatusService struct {
	repo repositories.ProjectSecurityGateReadinessStatusRepositoryInterface
}

func NewProjectSecurityGateReadinessStatusService(repo repositories.ProjectSecurityGateReadinessStatusRepositoryInterface) ProjectSecurityGateReadinessStatusServiceInterface {
	return &ProjectSecurityGateReadinessStatusService{
		repo: repo,
	}
}

func (psgrss *ProjectSecurityGateReadinessStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateReadinessStatus, error) {
	return psgrss.repo.GetAll(sort, sortBy, keyword, status)
}

func (psgrss *ProjectSecurityGateReadinessStatusService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateReadinessStatus, error) {
	return psgrss.repo.GetByID(id, tx...)
}

func (psgrss *ProjectSecurityGateReadinessStatusService) GetByType(readinessStatusType string, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateReadinessStatus, error) {
	return psgrss.repo.GetByType(readinessStatusType, tx...)
}
