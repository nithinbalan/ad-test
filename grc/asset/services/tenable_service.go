package services

import (
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
)

type TenableServiceInterface interface {
	GetAll(params *dtos.TenableAssetListQueryParamsDTO) (*dtos.TenableAssetListResponseDTO, error)
}

type TenableService struct {
	repo repositories.TenableRepositoryInterface
}

func NewTenableService(repo repositories.TenableRepositoryInterface) TenableServiceInterface {
	return &TenableService{repo: repo}
}

func (s *TenableService) GetAll(params *dtos.TenableAssetListQueryParamsDTO) (*dtos.TenableAssetListResponseDTO, error) {
	return s.repo.GetAllAssets(params)
}
