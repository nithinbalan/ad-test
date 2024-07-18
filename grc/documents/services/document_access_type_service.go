package services

import (
	"errors"
	"jethur/grc/documents/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentAccessTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentAccessType, error)
	GetByID(id uint) (*models.DocumentAccessType, error)
	GetByType(accessType string) (*models.DocumentAccessType, error)
}

// Service and implementations
type DocumentAccessTypeService struct {
	repo repositories.DocumentAccessTypeRepositoryInterface
}

func NewDocumentAccessTypeService(repo repositories.DocumentAccessTypeRepositoryInterface) DocumentAccessTypeServiceInterface {
	return &DocumentAccessTypeService{
		repo: repo,
	}
}

func (das *DocumentAccessTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentAccessType, error) {
	return das.repo.GetAll(sort, sortBy, keyword, status)
}

func (das DocumentAccessTypeService) GetByID(id uint) (*models.DocumentAccessType, error) {
	accessType, err := das.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return accessType, nil
}

func (das *DocumentAccessTypeService) GetByType(accessType string) (*models.DocumentAccessType, error) {
	documentAccessType, err := das.repo.GetByType(accessType)

	if err != nil {
		return nil, err
	}

	return documentAccessType, nil
}
