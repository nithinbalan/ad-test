package services

import (
	"errors"
	"jethur/grc/documents/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentChangeRequestStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestStatus, error)
	GetByID(id uint) (*models.DocumentChangeRequestStatus, error)
	GetByType(status string) (*models.DocumentChangeRequestStatus, error)
}

// Service and implementations
type DocumentChangeRequestStatusService struct {
	repo repositories.DocumentChangeRequestStatusRepositoryInterface
}

func NewDocumentChangeRequestStatusService(repo repositories.DocumentChangeRequestStatusRepositoryInterface) DocumentChangeRequestStatusServiceInterface {
	return &DocumentChangeRequestStatusService{
		repo: repo,
	}
}

func (dcrss *DocumentChangeRequestStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestStatus, error) {
	return dcrss.repo.GetAll(sort, sortBy, keyword, status)
}

func (dcrss DocumentChangeRequestStatusService) GetByID(id uint) (*models.DocumentChangeRequestStatus, error) {
	documentChangeRequestStatus, err := dcrss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return documentChangeRequestStatus, nil
}

func (dcrss *DocumentChangeRequestStatusService) GetByType(status string) (*models.DocumentChangeRequestStatus, error) {
	documentChangeRequestStatus, err := dcrss.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return documentChangeRequestStatus, nil
}
