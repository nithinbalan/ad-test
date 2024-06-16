package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
)

type ComplianceRequirementAssociatedItemServiceInterface interface {
	GetAllControls(complianceID uint, requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementControlListDTO, error)
	GetAllStrategicObjectives(complianceID uint, requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementStrategicObjectiveListDTO, error)
	GetAllRisks(complianceID uint, requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementRiskListDTO, error)
	GetAllIncidents(complianceID uint, requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementIncidentListDTO, error)
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

func (crais *ComplianceRequirementAssociatedItemService) GetAllControls(complianceID uint, requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementControlListDTO, error) {
	if _, err := crais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return crais.repo.GetAllControls(complianceID, requirementID, params)
}

func (crais *ComplianceRequirementAssociatedItemService) GetAllStrategicObjectives(complianceID uint, requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementStrategicObjectiveListDTO, error) {
	if _, err := crais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return crais.repo.GetAllStrategicObjectives(complianceID, requirementID, params)
}

func (crais *ComplianceRequirementAssociatedItemService) GetAllRisks(complianceID uint, requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementRiskListDTO, error) {
	if _, err := crais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return crais.repo.GetAllRisks(complianceID, requirementID, params)
}

func (crais *ComplianceRequirementAssociatedItemService) GetAllIncidents(complianceID uint, requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementIncidentListDTO, error) {
	if _, err := crais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return crais.repo.GetAllIncidents(complianceID, requirementID, params)
}
