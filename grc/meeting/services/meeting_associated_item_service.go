package services

import (
	commonDTOs "jethur/grc/common/dtos"
	meetingDTOs "jethur/grc/meeting/dtos"
	meetingRepositories "jethur/grc/meeting/repositories"
	organizationDTOs "jethur/grc/organization/dtos"
)

type MeetingAssociatedItemServiceInterface interface {
	GetAllControls(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) (*meetingDTOs.MeetingAssociatedItemControlListResponseDTO, error)
	GetAllComplianceRequirements(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) (*meetingDTOs.MeetingAssociatedItemComplianceRequirementListResponseDTO, error)
	GetAllRisks(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) (*meetingDTOs.MeetingAssociatedItemRiskListResponseDTO, error)
	GetAllProcesses(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) (*meetingDTOs.MeetingAssociatedItemProcessListResponseDTO, error)
}

// Service and implementations
type MeetingAssociatedItemService struct {
	repo       meetingRepositories.MeetingAssociatedItemRepositoryInterface
	meetingSvc MeetingServiceInterface
}

func NewMeetingAssociatedItemService(
	repo meetingRepositories.MeetingAssociatedItemRepositoryInterface,
	meetingSvc MeetingServiceInterface,

) MeetingAssociatedItemServiceInterface {
	return &MeetingAssociatedItemService{
		repo:       repo,
		meetingSvc: meetingSvc,
	}
}

func (mais *MeetingAssociatedItemService) GetAllControls(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) (*meetingDTOs.MeetingAssociatedItemControlListResponseDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mais.repo.GetAllControls(meetingID, params)
}

func (mais *MeetingAssociatedItemService) GetAllComplianceRequirements(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) (*meetingDTOs.MeetingAssociatedItemComplianceRequirementListResponseDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}

	complianceRequirementList, totalCount, err := mais.repo.GetAllComplianceRequirements(meetingID, params)
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

	response := &meetingDTOs.MeetingAssociatedItemComplianceRequirementListResponseDTO{
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

func (mais *MeetingAssociatedItemService) GetAllRisks(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) (*meetingDTOs.MeetingAssociatedItemRiskListResponseDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mais.repo.GetAllRisks(meetingID, params)
}

func (mais *MeetingAssociatedItemService) GetAllProcesses(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) (*meetingDTOs.MeetingAssociatedItemProcessListResponseDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mais.repo.GetAllProcesses(meetingID, params)
}
