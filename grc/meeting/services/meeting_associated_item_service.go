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
	GetAllRisks(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetAllProcesses(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	GetAllIncidents(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error)
	GetAllMeetingDocuments(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateMeetingDocument(meetingID uint, currentUserID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteMeetingDocument(meetingID uint, documentID uint) error
	CreateMeetingControl(meetingID uint, currentUserID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteMeetingControl(meetingID uint, controlID uint) error
	CreateMeetingProcess(meetingID uint, currentUserID uint, dto *commonDTOs.AssociatedProcessMapDTO) error
	DeleteMeetingProcess(meetingID uint, processID uint) error
	GetAllMeetingStrategicObjectives(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error)
	CreateMeetingStrategicObjective(meetingID uint, currentUserID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error
	DeleteMeetingStrategicObjective(meetingID uint, strategicObjectiveID uint) error
	CreateMeetingRisk(meetingID uint, currentUserID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteMeetingRisk(meetingID uint, riskID uint) error
	CreateMeetingIncident(meetingID uint, currentUserID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error
	DeleteMeetingIncident(meetingID uint, incidentID uint) error
	GetAllMeetingAssets(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error)
	CreateMeetingAsset(meetingID uint, currentUserID uint, dto *commonDTOs.AssociatedAssetMapDTO) error
	DeleteMeetingAsset(meetingID uint, assetID uint) error
	GetAllMeetingKPI(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error)
	CreateMeetingKPI(meetingID uint, currentUserID uint, dto *commonDTOs.AssociatedKPIMapDTO) error
	DeleteMeetingKPI(meetingID uint, kpiID uint) error
	GetAllClauses(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	UpdateClauses(meetingID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error
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

func (mais *MeetingAssociatedItemService) GetAllRisks(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
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

func (mais *MeetingAssociatedItemService) GetAllIncidents(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error) {
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

func (mais *MeetingAssociatedItemService) GetAllMeetingStrategicObjectives(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mais.repo.GetAllMeetingStrategicObjectives(meetingID, params)
}

func (mais *MeetingAssociatedItemService) CreateMeetingStrategicObjective(meetingID, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error {

	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting does not exist for ID: %d", meetingID), "not_exist", "meeting_id")
	}

	if err := mais.validationSvc.ValidateStrtegicObjectiveMapping(meetingID, dto.StrategicObjectiveIDs, "meeting_strategic_objectives", "meeting_id"); err != nil {
		return err
	}

	tx := mais.repo.BeginTransaction()
	var mappings []*models.MeetingStrategicObjective
	currentTime := time.Now()

	for _, StrategicObjectiveID := range dto.StrategicObjectiveIDs {
		mappings = append(mappings, &models.MeetingStrategicObjective{
			StrategicObjectiveID: StrategicObjectiveID,
			MeetingID:            meetingID,
			CreatedBYID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := mais.repo.CreateMeetingStrategicObjective(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) DeleteMeetingStrategicObjective(meetingID uint, strategicObjectiveID uint) error {
	// Validate control existence
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting does not exist for ID: %d", meetingID), "not_exist", "meeting_id")
	}

	tx := mais.repo.BeginTransaction()

	if err := mais.repo.DeleteMeetingStrategicObjective(meetingID, strategicObjectiveID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) CreateMeetingRisk(meetingID, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {

	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting does not exist for ID: %d", meetingID), "not_exist", "meeting_id")
	}

	if err := mais.validationSvc.ValidateRiskMapping(meetingID, dto.RiskIDs, "meeting_risks", "meeting_id"); err != nil {
		return err
	}

	tx := mais.repo.BeginTransaction()
	var mappings []*models.MeetingRisk
	currentTime := time.Now()

	for _, riskID := range dto.RiskIDs {
		mappings = append(mappings, &models.MeetingRisk{
			RiskRegisterID: riskID,
			MeetingID:      meetingID,
			CreatedBYID:    userID,
			CreatedAt:      currentTime,
		})
	}

	if err := mais.repo.CreateMeetingRisk(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) DeleteMeetingRisk(meetingID uint, riskID uint) error {
	// Validate control existence
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting does not exist for ID: %d", meetingID), "not_exist", "meeting_id")
	}

	tx := mais.repo.BeginTransaction()

	if err := mais.repo.DeleteMeetingRisk(meetingID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) CreateMeetingIncident(meetingID, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error {

	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting does not exist for ID: %d", meetingID), "not_exist", "meeting_id")
	}

	if err := mais.validationSvc.ValidateIncidentMapping(meetingID, dto.IncidentIDs, "meeting_incidents", "meeting_id"); err != nil {
		return err
	}

	tx := mais.repo.BeginTransaction()
	var mappings []*models.MeetingIncident
	currentTime := time.Now()

	for _, incidentID := range dto.IncidentIDs {
		mappings = append(mappings, &models.MeetingIncident{
			IncidentRegisterID: incidentID,
			MeetingID:          meetingID,
			CreatedBYID:        userID,
			CreatedAt:          currentTime,
		})
	}

	if err := mais.repo.CreateMeetingIncident(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) DeleteMeetingIncident(meetingID uint, incidentID uint) error {
	// Validate control existence
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting does not exist for ID: %d", meetingID), "not_exist", "meeting_id")
	}

	tx := mais.repo.BeginTransaction()

	if err := mais.repo.DeleteMeetingIncident(meetingID, incidentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) GetAllMeetingAssets(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mais.repo.GetAllMeetingAssets(meetingID, params)
}

func (mais *MeetingAssociatedItemService) CreateMeetingAsset(meetingID, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error {

	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return err
	}

	if err := mais.validationSvc.ValidateAssetMapping(meetingID, dto.AssetIDs, "meeting_assets", "meeting_id"); err != nil {
		return err
	}

	tx := mais.repo.BeginTransaction()
	var mappings []*models.MeetingAsset
	currentTime := time.Now()

	for _, assetID := range dto.AssetIDs {
		mappings = append(mappings, &models.MeetingAsset{
			AssetRegisterID: assetID,
			MeetingID:       meetingID,
			CreatedBYID:     userID,
			CreatedAt:       currentTime,
		})
	}

	if err := mais.repo.CreateMeetingAsset(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) DeleteMeetingAsset(meetingID uint, assetID uint) error {
	// Validate control existence
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return err
	}

	tx := mais.repo.BeginTransaction()

	if err := mais.repo.DeleteMeetingAsset(meetingID, assetID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) GetAllMeetingKPI(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mais.repo.GetAllMeetingKPI(meetingID, params)
}

func (mais *MeetingAssociatedItemService) CreateMeetingKPI(meetingID, userID uint, dto *commonDTOs.AssociatedKPIMapDTO) error {

	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return err
	}

	if err := mais.validationSvc.ValidateKPIMapping(meetingID, dto.KpiIDs, "meeting_kpis", "meeting_id"); err != nil {
		return err
	}

	tx := mais.repo.BeginTransaction()
	var mappings []*models.MeetingKpi
	currentTime := time.Now()

	for _, kpiID := range dto.KpiIDs {
		mappings = append(mappings, &models.MeetingKpi{
			KPIRegisterID: kpiID,
			MeetingID:     meetingID,
			CreatedByID:   userID,
			CreatedAt:     currentTime,
		})
	}

	if err := mais.repo.CreateMeetingKPI(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) DeleteMeetingKPI(meetingID uint, kpiID uint) error {
	// Validate control existence
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return err
	}

	tx := mais.repo.BeginTransaction()

	if err := mais.repo.DeleteMeetingKPI(meetingID, kpiID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (mais *MeetingAssociatedItemService) GetAllClauses(meetingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}

	response, err := mais.repo.GetAllClauses(meetingID, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (mais *MeetingAssociatedItemService) UpdateClauses(meetingID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error {

	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return err
	}
	tx := mais.repo.BeginTransaction()
	clauses, err := mais.validationSvc.ValidateClauseIDs(ClauseFormDTO.ClauseIDs)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlClauses, err := mais.repo.GetClausesByMeetingId(meetingID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingClauses := make(map[uint]bool) // clauseID => exists
	for _, clauseID := range existingControlClauses {
		existingClauses[clauseID] = true
	}

	newClauses := make(map[uint]*models.ComplianceClause)
	for _, clause := range clauses {
		newClauses[clause.ID] = clause
	}

	var clausesToAdd []*models.MeetingComplianceRequirement
	now := time.Now()
	for id, clause := range newClauses {
		if !existingClauses[id] {
			clausesToAdd = append(clausesToAdd, &models.MeetingComplianceRequirement{
				MeetingID:            meetingID,
				ComplianceRegisterID: clause.ComplianceRegisterID,
				ComplianceClauseID:   &clause.ID,
				CreatedBYID:          currentUserID,
				CreatedAt:            now,
			})
		}
	}

	if len(clausesToAdd) > 0 {
		if err := mais.repo.CreateMeetingClauses(clausesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var deleteClauseIDs []uint
	for existingID := range existingClauses {
		if _, found := newClauses[existingID]; !found {
			deleteClauseIDs = append(deleteClauseIDs, existingID)
		}
	}

	if len(deleteClauseIDs) > 0 {
		if err := mais.repo.DeleteMeetingClauses(meetingID, deleteClauseIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}
