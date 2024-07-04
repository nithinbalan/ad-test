package services

import (
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
)

type RiskAssociatedItemServiceInterface interface {
	GetAllControls(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemControlListResponseDTO, error)
	GetAllComplianceRequirements(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemComplianceRequirementListResponseDTO, error)
	GetAllProcesses(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemProcessListResponseDTO, error)
	GetAllStrategicObjectives(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemStrategicObjectiveListResponseDTO, error)
	GetAllSLAs(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemSLAListResponseDTO, error)
}

// Service and implementations
type RiskAssociatedItemService struct {
	repo            riskRepositories.RiskAssociatedItemRepositoryInterface
	riskRegisterSvc RiskRegisterServiceInterface
}

func NewRiskAssociatedItemService(
	repo riskRepositories.RiskAssociatedItemRepositoryInterface,
	riskRegisterSvc RiskRegisterServiceInterface,

) RiskAssociatedItemServiceInterface {
	return &RiskAssociatedItemService{
		repo:            repo,
		riskRegisterSvc: riskRegisterSvc,
	}
}

func (rais *RiskAssociatedItemService) GetAllControls(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemControlListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllControls(riskID, params)
}

func (rais *RiskAssociatedItemService) GetAllComplianceRequirements(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemComplianceRequirementListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllComplianceRequirements(riskID, params)
}

func (rais *RiskAssociatedItemService) GetAllProcesses(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemProcessListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllProcesses(riskID, params)
}

func (rais *RiskAssociatedItemService) GetAllStrategicObjectives(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemStrategicObjectiveListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllStrategicObjectives(riskID, params)
}

func (rais *RiskAssociatedItemService) GetAllSLAs(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemSLAListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllSLAs(riskID, params)
}