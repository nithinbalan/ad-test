package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
)

type ComplianceAssociatedItemServiceInterface interface {
	GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemControlListResponseDTO, error)
	GetAllRisks(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemRiskListResponseDTO, error)
	GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemIncidentListResponseDTO, error)
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

func (cais *ComplianceAssociatedItemService) GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemControlListResponseDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllControls(complianceID, params)
}

func (cais *ComplianceAssociatedItemService) GetAllRisks(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemRiskListResponseDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllRisks(complianceID, params)
}

func (cais *ComplianceAssociatedItemService) GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemIncidentListResponseDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllIncidents(complianceID, params)
}
