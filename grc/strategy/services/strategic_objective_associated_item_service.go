package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceServices "jethur/grc/compliance/services"
	"jethur/grc/models"
	"jethur/grc/strategy/repositories"
	"jethur/utility"
	"time"
)

type StrategicObjectiveAssociatedItemServiceInterface interface {
	GetAllClauses(strategicObjectiveID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	UpdateClauses(strategicObjectiveID uint, currentUserID uint, ControlClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error
	GetAllRisks(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	CreateStrategicObjectiveRisk(strategicObjectiveID uint, currentUserID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteStrategicObjectiveRisk(strategicObjectiveID uint, riskID uint) error
	GetAllStrategicObjectiveDocuments(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateStrategicObjectiveDocuments(strategicObjectiveID uint, currentUserID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteStrategicObjectiveDocument(strategicObjectiveID uint, documentID uint) error
}

// Service and implementations
type StrategicObjectiveAssociatedItemService struct {
	repo                  repositories.StrategicObjectiveAssociatedItemRepositoryInterface
	strategicObjectiveSvc StrategicObjectiveServiceInterface
	clauseSvc             complianceServices.ComplianceClauseServiceInterface
	validationSvc         commonServices.AssociatedItemValidationServiceInterface
}

func NewstrategicObjectiveAssociatedItemService(
	repo repositories.StrategicObjectiveAssociatedItemRepositoryInterface,
	strategicObjectiveSvc StrategicObjectiveServiceInterface,
	clauseSvc complianceServices.ComplianceClauseServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,

) StrategicObjectiveAssociatedItemServiceInterface {
	return &StrategicObjectiveAssociatedItemService{
		repo:                  repo,
		strategicObjectiveSvc: strategicObjectiveSvc,
		clauseSvc:             clauseSvc,
		validationSvc:         validationSvc,
	}
}

func (rais *StrategicObjectiveAssociatedItemService) GetAllClauses(strategicObjectiveID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
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
func (rais *StrategicObjectiveAssociatedItemService) GetAllRisks(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	if _, err := rais.strategicObjectiveSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllRisks(processID, params)
}

func (rais *StrategicObjectiveAssociatedItemService) CreateStrategicObjectiveRisk(strategicObjectiveID, currentUserID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {
	if _, err := rais.strategicObjectiveSvc.GetByID(strategicObjectiveID); err != nil {
		return utility.NewValidationError(fmt.Errorf("strategic objective does not exist for ID: %d", strategicObjectiveID), "not_exist", "strategic_objective_id")
	}
	if err := rais.validationSvc.ValidateRiskMapping(strategicObjectiveID, dto.RiskIDs, "risk_strategic_objectives", "strategic_objective_id"); err != nil {
		return err
	}
	tx := rais.repo.BeginTransaction()
	var mappings []*models.RiskStrategicObjective
	currentTime := time.Now()

	for _, riskID := range dto.RiskIDs {
		mappings = append(mappings, &models.RiskStrategicObjective{
			RiskRegisterID:       riskID,
			StrategicObjectiveID: strategicObjectiveID,
			CreatedByID:          currentUserID,
			CreatedAt:            currentTime,
		})
	}

	if err := rais.repo.CreateStrategicObjectiveRisk(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (rais *StrategicObjectiveAssociatedItemService) DeleteStrategicObjectiveRisk(strategicObjectiveID uint, riskID uint) error {
	// Validate control existence
	if _, err := rais.strategicObjectiveSvc.GetByID(strategicObjectiveID); err != nil {
		return utility.NewValidationError(fmt.Errorf("strategic objective does not exist for ID: %d", strategicObjectiveID), "not_exist", "strategic_objective_id")
	}

	tx := rais.repo.BeginTransaction()

	if err := rais.repo.DeleteStrategicObjectiveRisk(strategicObjectiveID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (rais *StrategicObjectiveAssociatedItemService) GetAllStrategicObjectiveDocuments(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := rais.strategicObjectiveSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllStrategicObjectiveDocuments(processID, params)
}

func (rais *StrategicObjectiveAssociatedItemService) CreateStrategicObjectiveDocuments(strategicObjectiveID, currentUserID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {
	if _, err := rais.strategicObjectiveSvc.GetByID(strategicObjectiveID); err != nil {
		return utility.NewValidationError(fmt.Errorf("strategic objective does not exist for ID: %d", strategicObjectiveID), "not_exist", "strategic_objective_id")
	}
	if err := rais.validationSvc.ValidateDocumentMapping(strategicObjectiveID, dto.DocumentIDs, "strategic_objective_associated_documents", "strategic_objective_id"); err != nil {
		return err
	}
	tx := rais.repo.BeginTransaction()
	var mappings []*models.StrategicObjectiveAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.StrategicObjectiveAssociatedDocument{
			DocumentID:           documentID,
			StrategicObjectiveID: strategicObjectiveID,
			CreatedByID:          currentUserID,
			CreatedAt:            currentTime,
		})
	}

	if err := rais.repo.CreateStrategicObjectiveDocuments(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (rais *StrategicObjectiveAssociatedItemService) DeleteStrategicObjectiveDocument(strategicObjectiveID uint, documentID uint) error {
	// Validate control existence
	if _, err := rais.strategicObjectiveSvc.GetByID(strategicObjectiveID); err != nil {
		return utility.NewValidationError(fmt.Errorf("strategic objective does not exist for ID: %d", strategicObjectiveID), "not_exist", "strategic_objective_id")
	}

	tx := rais.repo.BeginTransaction()

	if err := rais.repo.DeleteStrategicObjectiveDocument(strategicObjectiveID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
