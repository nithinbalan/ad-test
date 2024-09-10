package services

import (
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	documentRepsitories "jethur/grc/documents/repositories"
	organizationDTOs "jethur/grc/organization/dtos"
)

// Service Interface
type DocumentAssociatedItemServiceInterface interface {
	GetAllControls(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemControlListResponseDTO, error)
	GetAllComplianceRequirements(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemComplianceRequirementListResponseDTO, error)
	GetAllProcesses(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemProcessListResponseDTO, error)
	GetAllSLAs(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemSLAListResponseDTO, error)
	GetAllIncidents(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemIncidentListResponseDTO, error)
	GetAllRisks(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemRiskListResponseDTO, error)
}

// Service Implementation
type DocumentAssociatedItemService struct {
	repo        documentRepsitories.DocumentAssociatedItemRepositoryInterface
	documentSvc DocumentServiceInterface
}

// Constructor Function
func NewDocumentAssociatedItemService(
	repo documentRepsitories.DocumentAssociatedItemRepositoryInterface,
	documentSvc DocumentServiceInterface,
) DocumentAssociatedItemServiceInterface {
	return &DocumentAssociatedItemService{
		repo:        repo,
		documentSvc: documentSvc,
	}
}

func (dais *DocumentAssociatedItemService) GetAllControls(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemControlListResponseDTO, error) {
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return nil, err
	}
	return dais.repo.GetAllControls(documentID, params)
}

func (dais *DocumentAssociatedItemService) GetAllComplianceRequirements(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemComplianceRequirementListResponseDTO, error) {
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return nil, err
	}

	complianceRequirementList, totalCount, err := dais.repo.GetAllComplianceRequirements(documentID, params)
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
	groupedComplianceRequirementList := make([]*organizationDTOs.ComplianceRegisterDTO, 0)
	for _, register := range groupedRequirements {
		groupedComplianceRequirementList = append(groupedComplianceRequirementList, register)
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &documentDTOs.DocumentAssociatedItemComplianceRequirementListResponseDTO{
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

func (dais *DocumentAssociatedItemService) GetAllProcesses(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemProcessListResponseDTO, error) {
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return nil, err
	}
	return dais.repo.GetAllProcesses(documentID, params)
}

func (dais *DocumentAssociatedItemService) GetAllSLAs(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemSLAListResponseDTO, error) {
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return nil, err
	}
	return dais.repo.GetAllSLAs(documentID, params)
}

func (dais *DocumentAssociatedItemService) GetAllIncidents(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemIncidentListResponseDTO, error) {
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return nil, err
	}
	return dais.repo.GetAllIncidents(documentID, params)
}

func (dais *DocumentAssociatedItemService) GetAllRisks(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemRiskListResponseDTO, error) {
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return nil, err
	}
	return dais.repo.GetAllRisks(documentID, params)
}
