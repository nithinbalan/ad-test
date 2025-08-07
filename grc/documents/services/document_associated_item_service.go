package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	documentDTOs "jethur/grc/documents/dtos"
	documentRepsitories "jethur/grc/documents/repositories"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
	"time"
)

// Service Interface
type DocumentAssociatedItemServiceInterface interface {
	GetAllControls(documentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	GetAllComplianceRequirements(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemComplianceRequirementListResponseDTO, error)
	GetAllProcesses(documentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	GetAllSLAs(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemSLAListResponseDTO, error)
	GetAllIncidents(documentID uint, params *documentDTOs.DocumentAssociatedItemListQueryParamsDTO) (*documentDTOs.DocumentAssociatedItemIncidentListResponseDTO, error)
	GetAllRisks(documentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetAllClauses(documentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	UpdateClauses(documentID uint, currentUserID uint, ControlClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error
	CreateDocumentRisk(documentID uint, currentUserID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteDocumentRisk(documentID uint, riskID uint) error
	CreateDocumentControl(documentID uint, currentUserID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteDocumentControl(documentID uint, controlID uint) error
	CreateDocumentProcess(documentID uint, currentUserID uint, dto *commonDTOs.AssociatedProcessMapDTO) error
	DeleteDocumentProcess(documentID uint, processID uint) error
}

// Service Implementation
type DocumentAssociatedItemService struct {
	repo          documentRepsitories.DocumentAssociatedItemRepositoryInterface
	documentSvc   DocumentServiceInterface
	validationSvc commonServices.AssociatedItemValidationServiceInterface
}

// Constructor Function
func NewDocumentAssociatedItemService(
	repo documentRepsitories.DocumentAssociatedItemRepositoryInterface,
	documentSvc DocumentServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) DocumentAssociatedItemServiceInterface {
	return &DocumentAssociatedItemService{
		repo:          repo,
		documentSvc:   documentSvc,
		validationSvc: validationSvc,
	}
}

func (dais *DocumentAssociatedItemService) GetAllControls(documentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
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
				ComplianceClauseTitleEn:              item.ComplianceClauseTitleEn,
				ComplianceClauseTitleAr:              item.ComplianceClauseTitleAr,
				ComplianceClauseNo:                   item.ComplianceClauseNo,
				ComplianceRequirementCategoryTitleEn: item.ComplianceRequirementCategoryTitleEn,
				ComplianceRequirementCategoryTitleAr: item.ComplianceRequirementCategoryTitleAr,
				ComplianceRequirementTitleEn:         item.ComplianceRequirementTitleEn,
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

func (dais *DocumentAssociatedItemService) GetAllProcesses(documentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
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

func (dais *DocumentAssociatedItemService) GetAllRisks(documentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return nil, err
	}
	return dais.repo.GetAllRisks(documentID, params)
}

func (dais *DocumentAssociatedItemService) GetAllClauses(documentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return nil, err
	}

	response, err := dais.repo.GetAllClauses(documentID, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (dais *DocumentAssociatedItemService) UpdateClauses(documentID uint, currentUserID uint, ControlClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error {

	tx := dais.repo.BeginTransaction()

	if _, err := dais.documentSvc.GetByID(documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	clauses, err := dais.validateClauseIDs(ControlClauseFormDTO.ClauseIDs, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlClauses, err := dais.repo.GetClausesByDocumentID(documentID, tx)
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

	var clausesToAdd []*models.DocumentComplianceRequirement
	now := time.Now()
	for id, clause := range newClauses {
		if !existingClauses[id] {
			clausesToAdd = append(clausesToAdd, &models.DocumentComplianceRequirement{
				DocumentID:           documentID,
				ComplianceRegisterID: clause.ComplianceRegisterID,
				ComplianceClauseID:   &clause.ID,
				CreatedByID:          currentUserID,
				CreatedAt:            now,
			})
		}
	}

	if len(clausesToAdd) > 0 {
		if err := dais.repo.CreateDocumentClauses(clausesToAdd, tx); err != nil {
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
		if err := dais.repo.DeleteDocumentClauses(documentID, deleteClauseIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (dais *DocumentAssociatedItemService) validateClauseIDs(clauseIDs []uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceClause, error) {
	seen := make(map[uint]bool)
	for _, id := range clauseIDs {
		if seen[id] {
			return nil, utility.NewValidationError(fmt.Errorf("duplicate clause ID: %d", id), "duplicate", "clause_id")
		}
		seen[id] = true
	}

	clauses, err := dais.repo.GetClausesByIDs(clauseIDs, tx...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch clauses: %w", err)
	}

	foundClauseIDs := make(map[uint]bool)
	for _, clause := range clauses {
		foundClauseIDs[clause.ID] = true
	}

	for _, id := range clauseIDs {
		if !foundClauseIDs[id] {
			return nil, utility.NewValidationError(
				fmt.Errorf("clause does not exist for ID: %d", id), "not_exist", "clause_id")
		}
	}

	return clauses, nil
}

func (dais *DocumentAssociatedItemService) CreateDocumentRisk(documentID, currentUserID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document does not exist for ID: %d", documentID), "not_exist", "document_id")
	}
	if err := dais.validationSvc.ValidateRiskMapping(documentID, dto.RiskIDs, "risk_associated_documents", "document_id"); err != nil {
		return err
	}
	tx := dais.repo.BeginTransaction()
	var mappings []*models.RiskAssociatedDocument
	currentTime := time.Now()

	for _, riskID := range dto.RiskIDs {
		mappings = append(mappings, &models.RiskAssociatedDocument{
			RiskRegisterID: riskID,
			DocumentID:     documentID,
			CreatedByID:    currentUserID,
			CreatedAt:      currentTime,
		})
	}

	if err := dais.repo.CreateDocumentRisk(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (dais *DocumentAssociatedItemService) DeleteDocumentRisk(documentID uint, riskID uint) error {
	// Validate control existence
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document does not exist for ID: %d", documentID), "not_exist", "document_id")
	}

	tx := dais.repo.BeginTransaction()

	if err := dais.repo.DeleteDocumentRisk(documentID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
func (dais *DocumentAssociatedItemService) CreateDocumentControl(documentID, currentUserID uint, dto *commonDTOs.AssociatedControlMapDTO) error {
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document does not exist for ID: %d", documentID), "not_exist", "document_id")
	}
	if err := dais.validationSvc.ValidateControlMapping(documentID, dto.ControlIDs, "control_associated_documents", "document_id"); err != nil {
		return err
	}
	tx := dais.repo.BeginTransaction()
	var mappings []*models.ControlAssociatedDocument
	currentTime := time.Now()

	for _, controlID := range dto.ControlIDs {
		mappings = append(mappings, &models.ControlAssociatedDocument{
			ControlID:   controlID,
			DocumentID:  documentID,
			CreatedByID: currentUserID,
			CreatedAt:   currentTime,
		})
	}

	if err := dais.repo.CreateDocumentControl(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (dais *DocumentAssociatedItemService) DeleteDocumentControl(documentID uint, controlID uint) error {
	// Validate control existence
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document does not exist for ID: %d", documentID), "not_exist", "document_id")
	}

	tx := dais.repo.BeginTransaction()

	if err := dais.repo.DeleteDocumentControl(documentID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (dais *DocumentAssociatedItemService) CreateDocumentProcess(documentID, currentUserID uint, dto *commonDTOs.AssociatedProcessMapDTO) error {
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document does not exist for ID: %d", documentID), "not_exist", "document_id")
	}
	if err := dais.validationSvc.ValidateProcessMapping(documentID, dto.ProcessIDs, "process_associated_documents", "document_id"); err != nil {
		return err
	}
	tx := dais.repo.BeginTransaction()
	var mappings []*models.ProcessAssociatedDocument
	currentTime := time.Now()

	for _, processID := range dto.ProcessIDs {
		mappings = append(mappings, &models.ProcessAssociatedDocument{
			ProcessID:   processID,
			DocumentID:  documentID,
			CreatedByID: currentUserID,
			CreatedAt:   currentTime,
		})
	}

	if err := dais.repo.CreateDocumentProcess(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (dais *DocumentAssociatedItemService) DeleteDocumentProcess(documentID uint, processID uint) error {
	// Validate control existence
	if _, err := dais.documentSvc.GetByID(documentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document does not exist for ID: %d", documentID), "not_exist", "document_id")
	}

	tx := dais.repo.BeginTransaction()

	if err := dais.repo.DeleteDocumentProcess(documentID, processID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
