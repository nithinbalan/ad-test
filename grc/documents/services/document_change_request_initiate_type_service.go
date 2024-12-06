package services

import (
	"errors"
	"jethur/grc/documents/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentChangeRequestInitiateTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestInitiateType, error)
	GetByID(id uint) (*models.DocumentChangeRequestInitiateType, error)
	GetByType(changeType string) (*models.DocumentChangeRequestInitiateType, error)
}

// Service and implementations
type DocumentChangeRequestInitiateTypeService struct {
	repo repositories.DocumentChangeRequestInitiateTypeRepositoryInterface
}

func NewDocumentChangeRequestInitiateTypeService(repo repositories.DocumentChangeRequestInitiateTypeRepositoryInterface) DocumentChangeRequestInitiateTypeServiceInterface {
	return &DocumentChangeRequestInitiateTypeService{
		repo: repo,
	}
}

func (dcrits *DocumentChangeRequestInitiateTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestInitiateType, error) {
	return dcrits.repo.GetAll(sort, sortBy, keyword, status)
}

func (dcrits *DocumentChangeRequestInitiateTypeService) GetByID(id uint) (*models.DocumentChangeRequestInitiateType, error) {
	changeRequestInitiateType, err := dcrits.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return changeRequestInitiateType, nil
}

func (dcrits *DocumentChangeRequestInitiateTypeService) GetByType(changeType string) (*models.DocumentChangeRequestInitiateType, error) {
	documentChangeRequestInitiateType, err := dcrits.repo.GetByType(changeType)

	if err != nil {
		return nil, err
	}

	return documentChangeRequestInitiateType, nil
}
