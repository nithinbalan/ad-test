package services

import (
	"errors"
	"jethur/grc/documents/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentWorkflowStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentWorkflowStatus, error)
	GetByID(id uint) (*models.DocumentWorkflowStatus, error)
	GetByType(status string) (*models.DocumentWorkflowStatus, error)
}

// Service and implementations
type DocumentWorkflowStatusService struct {
	repo repositories.DocumentWorkflowStatusRepositoryInterface
}

func NewDocumentWorkflowStatusService(repo repositories.DocumentWorkflowStatusRepositoryInterface) DocumentWorkflowStatusServiceInterface {
	return &DocumentWorkflowStatusService{
		repo: repo,
	}
}

func (dws *DocumentWorkflowStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentWorkflowStatus, error) {
	return dws.repo.GetAll(sort, sortBy, keyword, status)
}

func (dws *DocumentWorkflowStatusService) GetByID(id uint) (*models.DocumentWorkflowStatus, error) {
	documentWorkflowStatus, err := dws.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return documentWorkflowStatus, nil
}

func (dws *DocumentWorkflowStatusService) GetByType(status string) (*models.DocumentWorkflowStatus, error) {
	documentWorkflowStatus, err := dws.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return documentWorkflowStatus, nil
}
