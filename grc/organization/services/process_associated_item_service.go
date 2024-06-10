package services

import (
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
)

type ProcessAssociatedItemServiceInterface interface {
	GetAllControls(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessControlListDTO, error)
}

// Service and implementations
type ProcessAssociatedItemService struct {
	repo       organizationRepositories.ProcessAssociatedItemRepositoryInterface
	processSvc ProcessServiceInterface
}

func NewProcessAssociatedItemService(
	repo organizationRepositories.ProcessAssociatedItemRepositoryInterface,
	processSvc ProcessServiceInterface,

) ProcessAssociatedItemServiceInterface {
	return &ProcessAssociatedItemService{
		repo:       repo,
		processSvc: processSvc,
	}
}

func (pais *ProcessAssociatedItemService) GetAllControls(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessControlListDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllControls(params)
}
