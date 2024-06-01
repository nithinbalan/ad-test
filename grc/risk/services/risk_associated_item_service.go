package services

import (
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
)

type RiskAssociatedItemServiceInterface interface {
	GetAllControls(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) ([]*riskDTOs.RiskControlListDTO, error)
	GetAllComplianceRequirements(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) ([]*riskDTOs.RiskComplianceRequirementListDTO, error)
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

func (rais *RiskAssociatedItemService) GetAllControls(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) ([]*riskDTOs.RiskControlListDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllControls(params)
}

func (rais *RiskAssociatedItemService) GetAllComplianceRequirements(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) ([]*riskDTOs.RiskComplianceRequirementListDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllComplianceRequirements(params)
}
