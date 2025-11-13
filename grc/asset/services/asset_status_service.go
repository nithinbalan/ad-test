package services

import (
	"errors"
	"jethur/grc/asset/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetStatus, error)
	GetByID(id uint) (*models.AssetStatus, error)
	GetByType(status string) (*models.AssetStatus, error)
}

// Service and implementations
type AssetStatusService struct {
	repo repositories.AssetStatusRepositoryInterface
}

func NewAssetStatusService(repo repositories.AssetStatusRepositoryInterface) AssetStatusServiceInterface {
	return &AssetStatusService{
		repo: repo,
	}
}

func (ass *AssetStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetStatus, error) {
	return ass.repo.GetAll(sort, sortBy, keyword, status)
}

func (ass AssetStatusService) GetByID(id uint) (*models.AssetStatus, error) {
	assetStatus, err := ass.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return assetStatus, nil
}

func (ass *AssetStatusService) GetByType(status string) (*models.AssetStatus, error) {
	assetStatus, err := ass.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return assetStatus, nil
}
