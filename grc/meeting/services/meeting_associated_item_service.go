package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	meetingDTOs "jethur/grc/meeting/dtos"
	meetingRepositories "jethur/grc/meeting/repositories"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
	"time"
)

type MeetingAssociatedItemServiceInterface interface {
	GetAllControls(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	GetAllComplianceRequirements(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) (*meetingDTOs.MeetingAssociatedItemComplianceRequirementListResponseDTO, error)
	GetAllRisks(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) (*meetingDTOs.MeetingAssociatedItemRiskListResponseDTO, error)
	GetAllProcesses(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	GetAllIncidents(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) (*meetingDTOs.MeetingAssociatedItemIncidentListResponseDTO, error)
	GetAllMeetingDocuments(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateMeetingDocument(meetingID uint, currentUserID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteMeetingDocument(meetingID uint, documentID uint) error
	CreateMeetingControl(meetingID uint, currentUserID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteMeetingControl(meetingID uint, controlIDID uint) error
	CreateMeetingProcess(meetingID uint, currentUserID uint, dto *commonDTOs.AssociatedProcessMapDTO) error
	DeleteMeetingProcess(meetingID uint, controlIDID uint) error
}

// Service and implementations
type MeetingAssociatedItemService struct {
	repo          meetingRepositories.MeetingAssociatedItemRepositoryInterface
	meetingSvc    MeetingServiceInterface
	validationSvc commonServices.AssociatedItemValidationServiceInterface
}

func NewMeetingAssociatedItemService(
	repo meetingRepositories.MeetingAssociatedItemRepositoryInterface,
	meetingSvc MeetingServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,

) MeetingAssociatedItemServiceInterface {
	return &MeetingAssociatedItemService{
		repo:          repo,
		meetingSvc:    meetingSvc,
		validationSvc: validationSvc,
	}
}

func (mais *MeetingAssociatedItemService) GetAllControls(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
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
				ComplianceClauseTitleEn:              item.ComplianceClauseTitleEn,
				ComplianceClauseTitleAr:              item.ComplianceClauseTitleAr,
				ComplianceClauseNo:                   item.ComplianceClauseNo,
				ComplianceRequirementCategoryTitleEn: item.ComplianceRequirementCategoryTitleEn,
				ComplianceRequirementCategoryTitleAr: item.ComplianceRequirementCategoryTitleAr,
				ComplianceRequirementTitleEn:         item.ComplianceRequirementTitleEn,
				ComplianceRequirementTitleAr:         item.ComplianceRequirementTitleAr,
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

func (mais *MeetingAssociatedItemService) GetAllProcesses(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mais.repo.GetAllProcesses(meetingID, params)
}

func (mais *MeetingAssociatedItemService) GetAllIncidents(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) (*meetingDTOs.MeetingAssociatedItemIncidentListResponseDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mais.repo.GetAllIncidents(meetingID, params)
}

func (mais *MeetingAssociatedItemService) GetAllMeetingDocuments(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mais.repo.GetAllMeetingDocuments(meetingID, params)
}

func (mais *MeetingAssociatedItemService) CreateMeetingDocument(meetingID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {

	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting does not exist for ID: %d", meetingID), "not_exist", "meeting_id")
	}

	if err := mais.validationSvc.ValidateDocumentMapping(meetingID, dto.DocumentIDs, "meeting_associated_documents", "meeting_id"); err != nil {
		return err
	}

	tx := mais.repo.BeginTransaction()
	var mappings []*models.MeetingAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.MeetingAssociatedDocument{
			DocumentID:  documentID,
			MeetingID:   meetingID,
			CreatedByID: userID,
			CreatedAt:   currentTime,
		})
	}

	if err := mais.repo.CreateMeetingDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) DeleteMeetingDocument(meetingID uint, documentID uint) error {
	// Validate control existence
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting does not exist for ID: %d", meetingID), "not_exist", "meeting_id")
	}

	tx := mais.repo.BeginTransaction()

	if err := mais.repo.DeleteMeetingDocument(meetingID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) CreateMeetingControl(meetingID, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error {

	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting does not exist for ID: %d", meetingID), "not_exist", "meeting_id")
	}

	if err := mais.validationSvc.ValidateControlMapping(meetingID, dto.ControlIDs, "meeting_controls", "meeting_id"); err != nil {
		return err
	}

	tx := mais.repo.BeginTransaction()
	var mappings []*models.MeetingControl
	currentTime := time.Now()

	for _, controlID := range dto.ControlIDs {
		mappings = append(mappings, &models.MeetingControl{
			ControlID:   controlID,
			MeetingID:   meetingID,
			CreatedBYID: userID,
			CreatedAt:   currentTime,
		})
	}

	if err := mais.repo.CreateMeetingControl(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) DeleteMeetingControl(meetingID uint, controlID uint) error {
	// Validate control existence
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting does not exist for ID: %d", meetingID), "not_exist", "meeting_id")
	}

	tx := mais.repo.BeginTransaction()

	if err := mais.repo.DeleteMeetingControl(meetingID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) CreateMeetingProcess(meetingID, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error {

	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting does not exist for ID: %d", meetingID), "not_exist", "meeting_id")
	}

	if err := mais.validationSvc.ValidateProcessMapping(meetingID, dto.ProcessIDs, "meeting_processes", "meeting_id"); err != nil {
		return err
	}

	tx := mais.repo.BeginTransaction()
	var mappings []*models.MeetingProcess
	currentTime := time.Now()

	for _, processID := range dto.ProcessIDs {
		mappings = append(mappings, &models.MeetingProcess{
			ProcessID:   processID,
			MeetingID:   meetingID,
			CreatedBYID: userID,
			CreatedAt:   currentTime,
		})
	}

	if err := mais.repo.CreateMeetingProcess(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) DeleteMeetingProcess(meetingID uint, processID uint) error {
	// Validate control existence
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting does not exist for ID: %d", meetingID), "not_exist", "meeting_id")
	}

	tx := mais.repo.BeginTransaction()

	if err := mais.repo.DeleteMeetingProcess(meetingID, processID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
