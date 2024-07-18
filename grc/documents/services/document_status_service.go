package services

import (
	"errors"
	"jethur/grc/documents/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentStatus, error)
	GetByID(id uint) (*models.DocumentStatus, error)
	GetByType(status string) (*models.DocumentStatus, error)
}

// Service and implementations
type DocumentStatusService struct {
	repo repositories.DocumentStatusRepositoryInterface
}

func NewDocumentStatusService(repo repositories.DocumentStatusRepositoryInterface) DocumentStatusServiceInterface {
	return &DocumentStatusService{
		repo: repo,
	}
}

func (dss *DocumentStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentStatus, error) {
	return dss.repo.GetAll(sort, sortBy, keyword, status)
}

func (dss DocumentStatusService) GetByID(id uint) (*models.DocumentStatus, error) {
	documentStatus, err := dss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return documentStatus, nil
}

func (dss *DocumentStatusService) GetByType(status string) (*models.DocumentStatus, error) {
	documentStatus, err := dss.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return documentStatus, nil
}
