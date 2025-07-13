package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	complianceServices "jethur/grc/compliance/services"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	"jethur/grc/models"
	"jethur/utility"
	"time"
)

type IncidentAssociatedItemServiceInterface interface {
	GetAllControls(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemControlListResponseDTO, error)
	GetAllComplianceRequirements(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemComplianceRequirementListResponseDTO, error)
	GetAllRisks(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemRiskListResponseDTO, error)
	GetAllProcesses(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemProcessListResponseDTO, error)
	GetAllSLAs(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemSLAListResponseDTO, error)
	GetAllDocuments(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemDocumentListResponseDTO, error)
	GetAllAssets(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemAssetListResponseDTO, error)
	GetAllClauses(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	UpdateClauses(incidentID uint, currentUserID uint, ControlClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error
}

// Service and implementations
type IncidentAssociatedItemService struct {
	repo                incidentRepositories.IncidentAssociatedItemRepositoryInterface
	incidentRegisterSvc IncidentRegisterServiceInterface
	clauseSvc           complianceServices.ComplianceClauseServiceInterface
}

func NewIncidentAssociatedItemService(
	repo incidentRepositories.IncidentAssociatedItemRepositoryInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface,
	clauseSvc complianceServices.ComplianceClauseServiceInterface,

) IncidentAssociatedItemServiceInterface {
	return &IncidentAssociatedItemService{
		repo:                repo,
		incidentRegisterSvc: incidentRegisterSvc,
		clauseSvc:           clauseSvc,
	}
}

func (iais *IncidentAssociatedItemService) GetAllControls(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemControlListResponseDTO, error) {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return iais.repo.GetAllControls(incidentID, params)
}

func (iais *IncidentAssociatedItemService) GetAllComplianceRequirements(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemComplianceRequirementListResponseDTO, error) {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return iais.repo.GetAllComplianceRequirements(incidentID, params)
}

func (iais *IncidentAssociatedItemService) GetAllRisks(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemRiskListResponseDTO, error) {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return iais.repo.GetAllRisks(incidentID, params)
}

func (iais *IncidentAssociatedItemService) GetAllProcesses(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemProcessListResponseDTO, error) {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return iais.repo.GetAllProcesses(incidentID, params)
}

func (iais *IncidentAssociatedItemService) GetAllSLAs(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemSLAListResponseDTO, error) {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return iais.repo.GetAllSLAs(incidentID, params)
}

func (iais *IncidentAssociatedItemService) GetAllDocuments(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemDocumentListResponseDTO, error) {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return iais.repo.GetAllDocuments(incidentID, params)
}

func (iais *IncidentAssociatedItemService) GetAllAssets(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemAssetListResponseDTO, error) {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return iais.repo.GetAllAssets(incidentID, params)
}

func (iais *IncidentAssociatedItemService) GetAllClauses(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	response, err := iais.repo.GetAllClauses(incidentID, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (iais *IncidentAssociatedItemService) UpdateClauses(incidentID uint, currentUserID uint, ControlClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error {

	tx := iais.repo.BeginTransaction()

	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return err
	}

	clauses, err := iais.validateClauseIDs(ControlClauseFormDTO.ClauseIDs, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlClauses, err := iais.repo.GetClausesByIncidentID(incidentID, tx)
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

	var clausesToAdd []*models.IncidentClause
	now := time.Now()
	for id, clause := range newClauses {
		if !existingClauses[id] {
			clausesToAdd = append(clausesToAdd, &models.IncidentClause{
				IncidentRegisterID:   incidentID,
				ComplianceRegisterID: clause.ComplianceRegisterID,
				ComplianceClauseID:   clause.ID,
				CreatedByID:          currentUserID,
				CreatedAt:            now,
			})
		}
	}

	if len(clausesToAdd) > 0 {
		if err := iais.repo.CreateIncidentClauses(clausesToAdd, tx); err != nil {
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
		if err := iais.repo.DeleteIncidentClauses(incidentID, deleteClauseIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (iais *IncidentAssociatedItemService) validateClauseIDs(clauseIDs []uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceClause, error) {
	seen := make(map[uint]bool)
	for _, id := range clauseIDs {
		if seen[id] {
			return nil, utility.NewValidationError(fmt.Errorf("duplicate clause ID: %d", id), "duplicate", "clause_id")
		}
		seen[id] = true
	}

	clauses, err := iais.clauseSvc.GetClausesByIDs(clauseIDs, tx...)
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
