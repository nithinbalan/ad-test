package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	complianceServices "jethur/grc/compliance/services"
	"jethur/grc/models"
	"jethur/grc/strategy/repositories"
	"jethur/utility"
	"time"
)

type StrategicObjectiveAssociatedItemServiceInterface interface {
	GetAllClauses(strategicObjectiveID uint, params *commonDTOs.AssociatedItemClauseListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	UpdateClauses(strategicObjectiveID uint, currentUserID uint, ControlClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error
}

// Service and implementations
type StrategicObjectiveAssociatedItemService struct {
	repo                  repositories.StrategicObjectiveAssociatedItemRepositoryInterface
	strategicObjectiveSvc StrategicObjectiveServiceInterface
	clauseSvc             complianceServices.ComplianceClauseServiceInterface
}

func NewstrategicObjectiveAssociatedItemService(
	repo repositories.StrategicObjectiveAssociatedItemRepositoryInterface,
	strategicObjectiveSvc StrategicObjectiveServiceInterface,
	clauseSvc complianceServices.ComplianceClauseServiceInterface,

) StrategicObjectiveAssociatedItemServiceInterface {
	return &StrategicObjectiveAssociatedItemService{
		repo:                  repo,
		strategicObjectiveSvc: strategicObjectiveSvc,
		clauseSvc:             clauseSvc,
	}
}

func (rais *StrategicObjectiveAssociatedItemService) GetAllClauses(strategicObjectiveID uint, params *commonDTOs.AssociatedItemClauseListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
	if _, err := rais.strategicObjectiveSvc.GetByID(strategicObjectiveID); err != nil {
		return nil, err
	}

	response, err := rais.repo.GetAllClauses(strategicObjectiveID, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (rais *StrategicObjectiveAssociatedItemService) UpdateClauses(strategicObjectiveID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error {

	tx := rais.repo.BeginTransaction()

	if _, err := rais.strategicObjectiveSvc.GetByID(strategicObjectiveID); err != nil {
		return err
	}

	clauses, err := rais.validateClauseIDs(ClauseFormDTO.ClauseIDs, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlClauses, err := rais.repo.GetClausesByStrategicObjectiveID(strategicObjectiveID, tx)
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

	var clausesToAdd []*models.ComplianceRequirementStrategicObjective
	now := time.Now()
	for id, clause := range newClauses {
		if !existingClauses[id] {
			clausesToAdd = append(clausesToAdd, &models.ComplianceRequirementStrategicObjective{
				StrategicObjectiveID: strategicObjectiveID,
				ComplianceRegisterID: clause.ComplianceRegisterID,
				ComplianceClauseID:   clause.ID,
				CreatedByID:          currentUserID,
				CreatedAt:            now,
			})
		}
	}

	if len(clausesToAdd) > 0 {
		if err := rais.repo.CreateStrategicObjectiveClauses(clausesToAdd, tx); err != nil {
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
		if err := rais.repo.DeleteStrategicObjectiveClauses(strategicObjectiveID, deleteClauseIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (rais *StrategicObjectiveAssociatedItemService) validateClauseIDs(clauseIDs []uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceClause, error) {
	seen := make(map[uint]bool)
	for _, id := range clauseIDs {
		if seen[id] {
			return nil, utility.NewValidationError(fmt.Errorf("duplicate clause ID: %d", id), "duplicate", "clause_id")
		}
		seen[id] = true
	}

	clauses, err := rais.clauseSvc.GetClausesByIDs(clauseIDs, tx...)
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
