package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/strategy/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type SwotPestelCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SwotPestelCategory, error)
	GetByID(id uint) (*models.SwotPestelCategory, error)
	GetByType(status string) (*models.SwotPestelCategory, error) // Renamed method
}

// Service and implementations
type SwotPestelCategoryService struct {
	repo repositories.SwotPestelCategoryRepositoryInterface
}

func NewSwotPestelCategoryService(repo repositories.SwotPestelCategoryRepositoryInterface) SwotPestelCategoryServiceInterface {
	return &SwotPestelCategoryService{
		repo: repo,
	}
}

func (cess *SwotPestelCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SwotPestelCategory, error) {
	return cess.repo.GetAll(sort, sortBy, keyword, status)
}

func (cess SwotPestelCategoryService) GetByID(id uint) (*models.SwotPestelCategory, error) {
	swotPestelCategory, err := cess.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return swotPestelCategory, nil
}

func (cess *SwotPestelCategoryService) GetByType(status string) (*models.SwotPestelCategory, error) {
	swotPestelCategory, err := cess.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return swotPestelCategory, nil
}
