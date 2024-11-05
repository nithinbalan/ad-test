package services

import (
	"jethur/grc/asset/repositories"
)

type AssetCriticalityMatrixServiceInterface interface {
	GetFirstMatrixID() (uint, error)
}

type AssetCriticalityMatrixService struct {
	repo repositories.AssetCriticalityMatrixRepositoryInterface
}

func NewAssetCriticalityMatrixService(
	repo repositories.AssetCriticalityMatrixRepositoryInterface,
) AssetCriticalityMatrixServiceInterface {
	return &AssetCriticalityMatrixService{
		repo: repo,
	}
}

func (s *AssetCriticalityMatrixService) GetFirstMatrixID() (uint, error) {
	return s.repo.GetFirstMatrixID()
}
