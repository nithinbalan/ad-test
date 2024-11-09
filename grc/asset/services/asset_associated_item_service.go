package services

import (
	assetDTOs "jethur/grc/asset/dtos"
	assetRepsitories "jethur/grc/asset/repositories"
)

// Service Interface
type AssetAssociatedItemServiceInterface interface {
	GetAllControls(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemControlListResponseDTO, error)
	GetAllIncidents(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemIncidentListResponseDTO, error)
}

// Service Implementation
type AssetAssociatedItemService struct {
	repo             assetRepsitories.AssetAssociatedItemRepositoryInterface
	assetRegisterSvc AssetRegisterServiceInterface
}

// Constructor Function
func NewAssetAssociatedItemService(
	repo assetRepsitories.AssetAssociatedItemRepositoryInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
) AssetAssociatedItemServiceInterface {
	return &AssetAssociatedItemService{
		repo:             repo,
		assetRegisterSvc: assetRegisterSvc,
	}
}

func (aais *AssetAssociatedItemService) GetAllControls(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemControlListResponseDTO, error) {
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return aais.repo.GetAllControls(assetID, params)
}

func (aais *AssetAssociatedItemService) GetAllIncidents(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemIncidentListResponseDTO, error) {
	if _, err := aais.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return aais.repo.GetAllIncidents(assetID, params)
}
