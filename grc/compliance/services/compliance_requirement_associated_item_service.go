package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
)

type ComplianceRequirementAssociatedItemServiceInterface interface {
	GetAllRisks(complianceID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementRiskListDTO, error)
	GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementIncidentListDTO, error)
}

// Service and implementations
type ComplianceRequirementAssociatedItemService struct {
	repo                  complianceRepositories.ComplianceRequirementAssociatedItemRepositoryInterface
	complianceRegisterSvc ComplianceRegisterServiceInterface
}

func NewComplianceRequirementAssociatedItemService(
	repo complianceRepositories.ComplianceRequirementAssociatedItemRepositoryInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,

) ComplianceRequirementAssociatedItemServiceInterface {
	return &ComplianceRequirementAssociatedItemService{
		repo:                  repo,
		complianceRegisterSvc: complianceRegisterSvc,
	}
}

func (crais *ComplianceRequirementAssociatedItemService) GetAllRisks(complianceID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementRiskListDTO, error) {
	if _, err := crais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return crais.repo.GetAllRisks(complianceID, params)
}

func (crais *ComplianceRequirementAssociatedItemService) GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementIncidentListDTO, error) {
	if _, err := crais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return crais.repo.GetAllIncidents(complianceID, params)
}
