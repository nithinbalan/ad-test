package services

import (
	"jethur/grc/models"
	"jethur/grc/project/repositories"
	"jethur/utility"
)

type ProjectSecurityGateChecklistResponseStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateChecklistResponseStatus, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistResponseStatus, error)
	GetByType(responseStatusType string, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistResponseStatus, error)
}

type ProjectSecurityGateChecklistResponseStatusService struct {
	repo repositories.ProjectSecurityGateChecklistResponseStatusRepositoryInterface
}

func NewProjectSecurityGateChecklistResponseStatusService(repo repositories.ProjectSecurityGateChecklistResponseStatusRepositoryInterface) ProjectSecurityGateChecklistResponseStatusServiceInterface {
	return &ProjectSecurityGateChecklistResponseStatusService{
		repo: repo,
	}
}

func (psgcrss *ProjectSecurityGateChecklistResponseStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateChecklistResponseStatus, error) {
	return psgcrss.repo.GetAll(sort, sortBy, keyword, status)
}

func (psgcrss *ProjectSecurityGateChecklistResponseStatusService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistResponseStatus, error) {
	return psgcrss.repo.GetByID(id, tx...)
}

func (psgcrss *ProjectSecurityGateChecklistResponseStatusService) GetByType(responseStatusType string, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistResponseStatus, error) {
	return psgcrss.repo.GetByType(responseStatusType, tx...)
}
