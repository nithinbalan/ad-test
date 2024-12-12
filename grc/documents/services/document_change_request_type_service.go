package services

import (
	"errors"
	"jethur/grc/documents/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentChangeRequestTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestType, error)
	GetByID(id uint) (*models.DocumentChangeRequestType, error)
	GetByType(changeType string) (*models.DocumentChangeRequestType, error)
}

// Service and implementations
type DocumentChangeRequestTypeService struct {
	repo repositories.DocumentChangeRequestTypeRepositoryInterface
}

func NewDocumentChangeRequestTypeService(repo repositories.DocumentChangeRequestTypeRepositoryInterface) DocumentChangeRequestTypeServiceInterface {
	return &DocumentChangeRequestTypeService{
		repo: repo,
	}
}

func (dcrts *DocumentChangeRequestTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestType, error) {
	return dcrts.repo.GetAll(sort, sortBy, keyword, status)
}

func (dcrts *DocumentChangeRequestTypeService) GetByID(id uint) (*models.DocumentChangeRequestType, error) {
	changeRequestType, err := dcrts.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return changeRequestType, nil
}

func (dcrts *DocumentChangeRequestTypeService) GetByType(changeType string) (*models.DocumentChangeRequestType, error) {
	documentChangeRequestType, err := dcrts.repo.GetByType(changeType)

	if err != nil {
		return nil, err
	}

	return documentChangeRequestType, nil
}
