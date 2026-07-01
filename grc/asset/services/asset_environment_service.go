package services

import (
	"errors"
	"jethur/grc/asset/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetEnvironmentServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetEnvironment, error)
	GetByID(id uint) (*models.AssetEnvironment, error)
	GetByType(status string) (*models.AssetEnvironment, error)
}

type AssetEnvironmentService struct {
	repo repositories.AssetEnvironmentRepositoryInterface
}

func NewAssetEnvironmentService(
	repo repositories.AssetEnvironmentRepositoryInterface,
) AssetEnvironmentServiceInterface {
	return &AssetEnvironmentService{
		repo: repo,
	}
}

func (aes *AssetEnvironmentService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetEnvironment, error) {
	return aes.repo.GetAll(sort, sortBy, keyword, status)
}

func (aes AssetEnvironmentService) GetByID(id uint) (*models.AssetEnvironment, error) {
	assetEnvironment, err := aes.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return assetEnvironment, nil
}

func (aes *AssetEnvironmentService) GetByType(status string) (*models.AssetEnvironment, error) {
	assetEnvironment, err := aes.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return assetEnvironment, nil
}