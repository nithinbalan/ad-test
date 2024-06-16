package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
)

type ComplianceAssociatedItemServiceInterface interface {
	GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceControlListDTO, error)
	GetAllRisks(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRiskListDTO, error)
	GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceIncidentListDTO, error)
}

// Service and implementations
type ComplianceAssociatedItemService struct {
	repo                  complianceRepositories.ComplianceAssociatedItemRepositoryInterface
	complianceRegisterSvc ComplianceRegisterServiceInterface
}

func NewComplianceAssociatedItemService(
	repo complianceRepositories.ComplianceAssociatedItemRepositoryInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
) ComplianceAssociatedItemServiceInterface {
	return &ComplianceAssociatedItemService{
		repo:                  repo,
		complianceRegisterSvc: complianceRegisterSvc,
	}
}

func (cais *ComplianceAssociatedItemService) GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceControlListDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllControls(complianceID, params)
}

func (cais *ComplianceAssociatedItemService) GetAllRisks(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRiskListDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllRisks(complianceID, params)
}

func (cais *ComplianceAssociatedItemService) GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceIncidentListDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllIncidents(complianceID, params)
}
