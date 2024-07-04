package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
)

type ComplianceRequirementAssociatedItemServiceInterface interface {
	GetAllControls(requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementAssociatedItemControlListResponseDTO, error)
	GetAllStrategicObjectives(requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementAssociatedItemStrategicObjectiveListResponseDTO, error)
	GetAllRisks(requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementAssociatedItemRiskListResponseDTO, error)
	GetAllIncidents(requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementAssociatedItemIncidentListResponseDTO, error)
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

func (crais *ComplianceRequirementAssociatedItemService) GetAllControls(requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementAssociatedItemControlListResponseDTO, error) {
	return crais.repo.GetAllControls(requirementID, params)
}

func (crais *ComplianceRequirementAssociatedItemService) GetAllStrategicObjectives(requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementAssociatedItemStrategicObjectiveListResponseDTO, error) {
	return crais.repo.GetAllStrategicObjectives(requirementID, params)
}

func (crais *ComplianceRequirementAssociatedItemService) GetAllRisks(requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementAssociatedItemRiskListResponseDTO, error) {
	return crais.repo.GetAllRisks(requirementID, params)
}

func (crais *ComplianceRequirementAssociatedItemService) GetAllIncidents(requirementID uint, params *complianceDTOs.ComplianceRequirementAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementAssociatedItemIncidentListResponseDTO, error) {
	return crais.repo.GetAllIncidents(requirementID, params)
}
