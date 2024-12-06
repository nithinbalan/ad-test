package services

import (
	"errors"
	"jethur/grc/documents/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentChangeRequestInitiateStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestInitiateStatus, error)
	GetByID(id uint) (*models.DocumentChangeRequestInitiateStatus, error)
	GetByType(status string) (*models.DocumentChangeRequestInitiateStatus, error)
}

// Service and implementations
type DocumentChangeRequestInitiateStatusService struct {
	repo repositories.DocumentChangeRequestInitiateStatusRepositoryInterface
}

func NewDocumentChangeRequestInitiateStatusService(repo repositories.DocumentChangeRequestInitiateStatusRepositoryInterface) DocumentChangeRequestInitiateStatusServiceInterface {
	return &DocumentChangeRequestInitiateStatusService{
		repo: repo,
	}
}

func (dcriss *DocumentChangeRequestInitiateStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestInitiateStatus, error) {
	return dcriss.repo.GetAll(sort, sortBy, keyword, status)
}

func (dcriss *DocumentChangeRequestInitiateStatusService) GetByID(id uint) (*models.DocumentChangeRequestInitiateStatus, error) {
	documentChangeRequestInitiateStatus, err := dcriss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return documentChangeRequestInitiateStatus, nil
}

func (dcriss *DocumentChangeRequestInitiateStatusService) GetByType(status string) (*models.DocumentChangeRequestInitiateStatus, error) {
	documentChangeRequestInitiateStatus, err := dcriss.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return documentChangeRequestInitiateStatus, nil
}
