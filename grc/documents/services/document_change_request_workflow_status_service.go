package services

import (
	"errors"
	"jethur/grc/documents/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentChangeRequestWorkflowStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestWorkflowStatus, error)
	GetByID(id uint) (*models.DocumentChangeRequestWorkflowStatus, error)
	GetByType(status string) (*models.DocumentChangeRequestWorkflowStatus, error)
}

// Service and implementations
type DocumentChangeRequestWorkflowStatusService struct {
	repo repositories.DocumentChangeRequestWorkflowStatusRepositoryInterface
}

func NewDocumentChangeRequestWorkflowStatusService(repo repositories.DocumentChangeRequestWorkflowStatusRepositoryInterface) DocumentChangeRequestWorkflowStatusServiceInterface {
	return &DocumentChangeRequestWorkflowStatusService{
		repo: repo,
	}
}

func (dws *DocumentChangeRequestWorkflowStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestWorkflowStatus, error) {
	return dws.repo.GetAll(sort, sortBy, keyword, status)
}

func (dws *DocumentChangeRequestWorkflowStatusService) GetByID(id uint) (*models.DocumentChangeRequestWorkflowStatus, error) {
	documentWorkflowStatus, err := dws.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return documentWorkflowStatus, nil
}

func (dws *DocumentChangeRequestWorkflowStatusService) GetByType(status string) (*models.DocumentChangeRequestWorkflowStatus, error) {
	documentWorkflowStatus, err := dws.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return documentWorkflowStatus, nil
}
