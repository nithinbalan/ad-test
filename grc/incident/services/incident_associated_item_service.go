package services

import (
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
)

type IncidentAssociatedItemServiceInterface interface {
	GetAllControls(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemControlListResponseDTO, error)
	GetAllComplianceRequirements(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemComplianceRequirementListResponseDTO, error)
	GetAllRisks(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemRiskListResponseDTO, error)
	GetAllProcesses(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemProcessListResponseDTO, error)
	GetAllSLAs(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemSLAListResponseDTO, error)
}

// Service and implementations
type IncidentAssociatedItemService struct {
	repo                incidentRepositories.IncidentAssociatedItemRepositoryInterface
	incidentRegisterSvc IncidentRegisterServiceInterface
}

func NewIncidentAssociatedItemService(
	repo incidentRepositories.IncidentAssociatedItemRepositoryInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface,

) IncidentAssociatedItemServiceInterface {
	return &IncidentAssociatedItemService{
		repo:                repo,
		incidentRegisterSvc: incidentRegisterSvc,
	}
}

func (iais *IncidentAssociatedItemService) GetAllControls(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemControlListResponseDTO, error) {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return iais.repo.GetAllControls(incidentID, params)
}

func (iais *IncidentAssociatedItemService) GetAllComplianceRequirements(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemComplianceRequirementListResponseDTO, error) {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return iais.repo.GetAllComplianceRequirements(incidentID, params)
}

func (iais *IncidentAssociatedItemService) GetAllRisks(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemRiskListResponseDTO, error) {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return iais.repo.GetAllRisks(incidentID, params)
}

func (iais *IncidentAssociatedItemService) GetAllProcesses(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemProcessListResponseDTO, error) {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return iais.repo.GetAllProcesses(incidentID, params)
}

func (iais *IncidentAssociatedItemService) GetAllSLAs(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemSLAListResponseDTO, error) {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return iais.repo.GetAllSLAs(incidentID, params)
}
