package services

import (
	organizationDTOs "jethur/grc/organization/dtos"

	organizationRepositories "jethur/grc/organization/repositories"
)

type ControlAssociatedItemServiceInterface interface {
	GetAllProcesses(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemProcessListResponseDTO, error)
	GetAllComplianceRequirements(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) ([]organizationDTOs.ComplianceRegisterDTO, error)
	GetAllRisks(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemRiskListResponseDTO, error)
	GetAllIncidents(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemIncidentListResponseDTO, error)
	GetAllAssociatedCounts(controlID uint) (*organizationDTOs.ControlAssociatedItemsCountResponseDTO, error)
}

type ControlAssociatedItemService struct {
	repo                           organizationRepositories.ControlAssociatedItemRepositoryInterface
	controlSvc                     ControlServiceInterface
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface
	controlEfficiencyStatusService ControlEfficiencyStatusServiceInterface
}

func NewControlAssociatedItemService(
	repo organizationRepositories.ControlAssociatedItemRepositoryInterface,
	controlSvc ControlServiceInterface,
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusService ControlEfficiencyStatusServiceInterface,
) ControlAssociatedItemServiceInterface {
	return &ControlAssociatedItemService{
		repo:                           repo,
		controlSvc:                     controlSvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusService: controlEfficiencyStatusService,
	}
}

func (cais *ControlAssociatedItemService) GetAllProcesses(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemProcessListResponseDTO, error) {
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllProcesses(controlID, params)
}

func (cais *ControlAssociatedItemService) GetAllComplianceRequirements(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) ([]organizationDTOs.ComplianceRegisterDTO, error) {
	// Validate the control ID
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	// Fetch all compliance requirements
	complianceRequirementList, err := cais.repo.GetAllComplianceRequirements(controlID, params)
	if err != nil {
		return nil, err
	}

	// Group by compliance_register_id
	groupedRequirements := make(map[uint]*organizationDTOs.ComplianceRegisterDTO)

	for _, item := range complianceRequirementList {
		registerID := item.ComplianceRegisterID

		// Check if the compliance register is already added to the group
		if _, exists := groupedRequirements[registerID]; !exists {
			groupedRequirements[registerID] = &organizationDTOs.ComplianceRegisterDTO{
				ComplianceRegisterID:      item.ComplianceRegisterID,
				ComplianceRegisterTitleEn: item.ComplianceRegisterTitleEn,
				ComplianceRegisterTitleAr: item.ComplianceRegisterTitleAr,
				Requirements:              []*organizationDTOs.ComplianceRequirementMinimalDTO{},
			}
		}

		if item.ComplianceRequirementID == 0 {
			continue
		}
		// Add the compliance requirement to the corresponding register
		groupedRequirements[registerID].Requirements = append(groupedRequirements[registerID].Requirements, &organizationDTOs.ComplianceRequirementMinimalDTO{
			ComplianceRequirementID:              item.ComplianceRequirementID,
			ComplianceClauseTitle:                item.ComplianceClauseTitle,
			ComplianceClauseNo:                   item.ComplianceClauseNo,
			ComplianceRequirementCategoryTitleEn: item.ComplianceRequirementCategoryTitleEn,
			ComplianceRequirementCategoryTitleAr: item.ComplianceRequirementCategoryTitleAr,
			ComplianceRequirementTitle:           item.ComplianceRequirementTitle,
			ComplianceStatusTitleEn:              item.ComplianceStatusTitleEn,
			ComplianceStatusTitleAr:              item.ComplianceStatusTitleAr,
			ComplianceStatusColorCode:            item.ComplianceStatusColorCode,
		})
	}

	// Convert the map to a slice
	var complianceRegisters []organizationDTOs.ComplianceRegisterDTO
	for _, register := range groupedRequirements {
		complianceRegisters = append(complianceRegisters, *register)
	}

	return complianceRegisters, nil
}

func (cais *ControlAssociatedItemService) GetAllRisks(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemRiskListResponseDTO, error) {
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllRisks(controlID, params)
}

func (cais *ControlAssociatedItemService) GetAllIncidents(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemIncidentListResponseDTO, error) {
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllIncidents(controlID, params)
}

func (cais *ControlAssociatedItemService) GetAllAssociatedCounts(controlID uint) (*organizationDTOs.ControlAssociatedItemsCountResponseDTO, error) {
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}
	response, err := cais.repo.GetAllAssociatedCounts(controlID)
	if err != nil {
		return nil, err
	}

	return response, nil
}
