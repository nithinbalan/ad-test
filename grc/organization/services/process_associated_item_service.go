package services

import (
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
)

type ProcessAssociatedItemServiceInterface interface {
	GetAllControls(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (map[string]map[string][]*organizationDTOs.ProcessControlListDTO, error)
	GetAllRisks(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessRiskListDTO, error)
	GetAllIncidents(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessIncidentListDTO, error)
}

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

func (pais *ProcessAssociatedItemService) GetAllControls(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (map[string]map[string][]*organizationDTOs.ProcessControlListDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllControls(processID, params)
}

func (pais *ProcessAssociatedItemService) GetAllRisks(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessRiskListDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllRisks(processID, params)
}

func (pais *ProcessAssociatedItemService) GetAllIncidents(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) ([]*organizationDTOs.ProcessIncidentListDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllIncidents(processID, params)
}
