package services

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	organizationRepositories "jethur/grc/organization/repositories"
)

type ControlAssociatedItemServiceInterface interface {
	GetAllProcesses(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemProcessListResponseDTO, error)
	GetAllComplianceRequirements(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemComplianceRequirementListResponseDTO, error)
	GetAllRisks(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemRiskListResponseDTO, error)
	GetAllIncidents(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemIncidentListResponseDTO, error)
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

func (cais *ControlAssociatedItemService) GetAllComplianceRequirements(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemComplianceRequirementListResponseDTO, error) {
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	complianceRequirementList, totalCount, err := cais.repo.GetAllComplianceRequirements(controlID, params)
	if err != nil {
		return nil, err
	}

	// Group by compliance register
	groupedRequirements := make(map[uint]*organizationDTOs.ComplianceRegisterDTO)
	for _, item := range complianceRequirementList {
		registerID := item.ComplianceRegisterID
		if _, exists := groupedRequirements[registerID]; !exists {
			groupedRequirements[registerID] = &organizationDTOs.ComplianceRegisterDTO{
				ComplianceRegisterID:      item.ComplianceRegisterID,
				ComplianceRegisterTitleEn: item.ComplianceRegisterTitleEn,
				ComplianceRegisterTitleAr: item.ComplianceRegisterTitleAr,
				Requirements:              []*organizationDTOs.ComplianceRequirementMinimalDTO{},
			}
		}
		if item.ComplianceRequirementID != 0 {
			requirement := &organizationDTOs.ComplianceRequirementMinimalDTO{
				ComplianceRequirementID:              item.ComplianceRequirementID,
				ComplianceClauseTitle:                item.ComplianceClauseTitle,
				ComplianceClauseNo:                   item.ComplianceClauseNo,
				ComplianceRequirementCategoryTitleEn: item.ComplianceRequirementCategoryTitleEn,
				ComplianceRequirementCategoryTitleAr: item.ComplianceRequirementCategoryTitleAr,
				ComplianceRequirementTitle:           item.ComplianceRequirementTitle,
				ComplianceStatusTitleEn:              item.ComplianceStatusTitleEn,
				ComplianceStatusTitleAr:              item.ComplianceStatusTitleAr,
				ComplianceStatusColorCode:            item.ComplianceStatusColorCode,
			}
			groupedRequirements[registerID].Requirements = append(groupedRequirements[registerID].Requirements, requirement)
		}
	}

	// Convert map to slice
	var groupedComplianceRequirementList []*organizationDTOs.ComplianceRegisterDTO
	for _, register := range groupedRequirements {
		groupedComplianceRequirementList = append(groupedComplianceRequirementList, register)
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &organizationDTOs.ControlAssociatedItemComplianceRequirementListResponseDTO{
		Items: groupedComplianceRequirementList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
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
