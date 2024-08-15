package services

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
)

type RiskAssociatedItemServiceInterface interface {
	GetAllControls(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemControlListResponseDTO, error)
	GetAllComplianceRequirements(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemComplianceRequirementListResponseDTO, error)
	GetAllProcesses(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemProcessListResponseDTO, error)
	GetAllStrategicObjectives(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemStrategicObjectiveListResponseDTO, error)
	GetAllSLAs(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemSLAListResponseDTO, error)
	GetAllIncidents(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemIncidentListResponseDTO, error)
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

	complianceRequirementList, totalCount, err := rais.repo.GetAllComplianceRequirements(riskID, params)
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

	response := &riskDTOs.RiskAssociatedItemComplianceRequirementListResponseDTO{
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

func (rais *RiskAssociatedItemService) GetAllIncidents(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemIncidentListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllIncidents(riskID, params)
}
