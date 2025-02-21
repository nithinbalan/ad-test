package services

import (
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// ComplianceClauseServiceInterface defines the methods for the service layer
type ComplianceClauseServiceInterface interface {
	GetAll(complianceID uint) ([]*complianceDTOs.ComplianceClauseListDTO, error)
	GetByID(clauseID uint) (*complianceDTOs.ComplianceClauseDetailDTO, error)
	Create(complianceID uint, currentUserID uint, formDTO *complianceDTOs.ComplianceClauseFormDTO) (uint, error)
	Update(complianceID uint, clauseID uint, currentUserID uint, formDTO *complianceDTOs.ComplianceClauseFormDTO) error
	Delete(id uint) error

	GetAssociatedComplianceRequirements(complianceID uint, clauseID uint) ([]complianceDTOs.ComplianceRequirementListDTO, error)
}

type ComplianceClauseService struct {
	repo           repositories.ComplianceClauseRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewComplianceClauseService(
	repo repositories.ComplianceClauseRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) *ComplianceClauseService {
	return &ComplianceClauseService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ccs *ComplianceClauseService) GetAll(complianceID uint) ([]*complianceDTOs.ComplianceClauseListDTO, error) {
	// Fetch clauses with preloaded child clauses
	clauses, err := ccs.repo.GetAll(complianceID)
	if err != nil {
		return nil, err
	}

	// If no clauses are found, return an empty slice
	if len(clauses) == 0 {
		return []*complianceDTOs.ComplianceClauseListDTO{}, nil
	}

	// Build a proper hierarchy by filtering out child clauses from the root level
	var rootClauses []*complianceDTOs.ComplianceClauseListDTO
	clauseMap := make(map[uint]*complianceDTOs.ComplianceClauseListDTO)

	// Convert all clauses into DTOs and map them by ID
	for _, clause := range clauses {
		dto := &complianceDTOs.ComplianceClauseListDTO{
			ID:             clause.ID,
			ClauseNo:       clause.ClauseNo,
			ParentClauseID: clause.ParentClauseID,
			TitleEn:        clause.TitleEn,
			TitleAr:        clause.TitleAr,
			ScopeEn:        clause.ScopeEn,
			ScopeAr:        clause.ScopeAr,
			ObjectiveEn:    clause.ObjectiveEn,
			ObjectiveAr:    clause.ObjectiveAr,
			ChildClauses:   []*complianceDTOs.ComplianceClauseListDTO{},
		}
		clauseMap[clause.ID] = dto
	}

	// Build hierarchy by linking children to their parent
	for _, clause := range clauses {
		if clause.ParentClauseID != nil {
			parent := clauseMap[*clause.ParentClauseID]
			parent.ChildClauses = append(parent.ChildClauses, clauseMap[clause.ID])
		} else {
			rootClauses = append(rootClauses, clauseMap[clause.ID])
		}
	}

	return rootClauses, nil
}

// Directly map child clauses (recursive)
func mapChildClauses(childClauses []*models.ComplianceClause) []*complianceDTOs.ComplianceClauseListDTO {
	var childDTOs []*complianceDTOs.ComplianceClauseListDTO
	for _, child := range childClauses {
		childDTOs = append(childDTOs, &complianceDTOs.ComplianceClauseListDTO{
			ID:             child.ID,
			ClauseNo:       child.ClauseNo,
			ParentClauseID: child.ParentClauseID,
			TitleEn:        child.TitleEn,
			TitleAr:        child.TitleAr,
			ScopeEn:        child.ScopeEn,
			ScopeAr:        child.ScopeAr,
			ObjectiveEn:    child.ObjectiveEn,
			ObjectiveAr:    child.ObjectiveAr,
			ChildClauses:   mapChildClauses(child.ChildClauses), // Recursive mapping
		})
	}
	return childDTOs
}

func (ccs *ComplianceClauseService) GetByID(clauseID uint) (*complianceDTOs.ComplianceClauseDetailDTO, error) {
	clause, err := ccs.repo.GetByID(clauseID)
	if err != nil {
		return nil, fmt.Errorf("compliance clause not found: %v", err)
	}

	var clauseDetailsDTO complianceDTOs.ComplianceClauseDetailDTO

	err = mapstructure.Decode(clause, &clauseDetailsDTO)
	if err != nil {
		return nil, err
	}

	return &clauseDetailsDTO, nil
}

func (ccs *ComplianceClauseService) Create(complianceID uint, currentUserID uint, formDTO *complianceDTOs.ComplianceClauseFormDTO) (uint, error) {
	tx := ccs.repo.BeginTransaction()

	var clauseNo string

	// If PreviousClauseID is provided, validate its relationship with ParentClauseID
	if formDTO.AddAsFirst != nil && *formDTO.AddAsFirst {
		if formDTO.ParentClauseID != nil {
			// Adding as first under a parent
			parentClause, err := ccs.repo.GetByID(*formDTO.ParentClauseID, tx)
			if err != nil {
				tx.Rollback()
				return 0, utility.NewValidationError(fmt.Errorf("invalid parent_clause_id: %v", err), "invalid", "parent_clause_id")
			}
			clauseNo = fmt.Sprintf("%s.1", parentClause.ClauseNo)
		} else {
			// Adding as first at root level
			clauseNo = "1"
		}

		// Check if first clause already exists
		existingClause, _ := ccs.repo.GetByClauseNo(clauseNo, complianceID, tx)
		if existingClause != nil {
			tx.Rollback()
			return 0, utility.NewValidationError(fmt.Errorf("first clause already exists"), "duplicate", "clause_no")
		}
	} else if formDTO.PreviousClauseID != nil {
		// Fetch previous clause details
		previousClause, err := ccs.repo.GetByID(*formDTO.PreviousClauseID, tx)
		if err != nil {
			tx.Rollback()
			return 0, utility.NewValidationError(fmt.Errorf("invalid previous_clause_id: %v", err), "invalid", "previous_clause_id")
		}
		if previousClause == nil {
			tx.Rollback()
			return 0, utility.NewValidationError(fmt.Errorf("previous clause does not exist"), "not_exist", "previous clause")
		}

		// **CASE 1: PreviousClauseID with ParentClauseID**
		if formDTO.ParentClauseID != nil {
			// Fetch parent clause details
			parentClause, err := ccs.repo.GetByID(*formDTO.ParentClauseID, tx)
			if err != nil {
				tx.Rollback()
				return 0, utility.NewValidationError(fmt.Errorf("invalid parent_clause_id: %v", err), "invalid", "parent_clause_id")
			}
			if parentClause == nil {
				tx.Rollback()
				return 0, utility.NewValidationError(fmt.Errorf("parent clause with ID %d not found", *formDTO.ParentClauseID), "not_exist", "parent_clause_id")
			}

			// Ensure PreviousClauseID is a direct child of ParentClauseID
			if previousClause.ParentClauseID == nil || *previousClause.ParentClauseID != *formDTO.ParentClauseID {
				tx.Rollback()
				//return 0, fmt.Errorf("previous clause %s is not a first-level child of parent clause %s", previousClause.ClauseNo, parentClause.ClauseNo)
				return 0, utility.NewValidationError(fmt.Errorf("previous clause %s is not a first-level child of parent clause %s", previousClause.ClauseNo, parentClause.ClauseNo), "invalid_clause_number", "clause_sequence")
			}
		}

		// **CASE 2: PreviousClauseID without ParentClauseID (Top-Level Validation)**
		if formDTO.ParentClauseID == nil {
			if previousClause.ParentClauseID != nil {
				tx.Rollback()
				return 0, utility.NewValidationError(fmt.Errorf("previous clause %s is not a top-level clause", previousClause.ClauseNo), "invalid_clause_number", "clause_sequence")
			}
		}

		// Generate the next sequence clause number
		baseNumberParts := strings.Split(previousClause.ClauseNo, ".")
		lastPartIndex := len(baseNumberParts) - 1
		lastPartNumber, _ := strconv.Atoi(baseNumberParts[lastPartIndex]) // Convert last part to integer
		nextSequence := lastPartNumber + 1                                // Get next clause number

		// Generate the new clause number (e.g., "3.2")
		baseNumberParts[lastPartIndex] = strconv.Itoa(nextSequence)
		clauseNo = strings.Join(baseNumberParts, ".")

		// Validate if clause number already exists
		existingClause, _ := ccs.repo.GetByClauseNo(clauseNo, complianceID, tx)
		if existingClause != nil {
			tx.Rollback()
			return 0, utility.NewValidationError(fmt.Errorf("clause number '%s' already exists", clauseNo), "duplicate", "clause_no")
		}
	} else if formDTO.ParentClauseID != nil {
		// Validate parent clause
		parentClause, err := ccs.repo.GetByID(*formDTO.ParentClauseID, tx)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("invalid parent_clause_id: %v", err)
		}
		if parentClause == nil {
			tx.Rollback()
			return 0, fmt.Errorf("parent clause with ID %d not found", *formDTO.ParentClauseID)
		}

		// Generate clause number based on parent
		nextSequence, err := ccs.repo.GetNextClauseSequence(*formDTO.ParentClauseID, complianceID, tx)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to generate clause number: %v", err)
		}
		clauseNo = fmt.Sprintf("%s.%d", parentClause.ClauseNo, nextSequence)
	} else {
		// Generate top-level clause number
		nextSequence, err := ccs.repo.GetNextTopLevelClauseSequence(complianceID, tx)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to generate top-level clause number: %v", err)
		}
		clauseNo = fmt.Sprintf("%d", nextSequence)
	}

	// Validate unique clause number
	existingClause, _ := ccs.repo.GetByClauseNo(clauseNo, complianceID, tx)
	if existingClause != nil {
		tx.Rollback()
		return 0, fmt.Errorf("clause number '%s' already exists", clauseNo)
	}

	// Create ComplianceClause model
	complianceClause := &models.ComplianceClause{
		ComplianceRegisterID: complianceID,
		ParentClauseID:       formDTO.ParentClauseID,
		ClauseNo:             clauseNo,
		TitleEn:              formDTO.TitleEn,
		TitleAr:              formDTO.TitleAr,
		DescriptionEn:        formDTO.DescriptionEn,
		DescriptionAr:        formDTO.DescriptionAr,
		ScopeEn:              formDTO.ScopeEn,
		ScopeAr:              formDTO.ScopeAr,
		ObjectiveEn:          formDTO.ObjectiveEn,
		ObjectiveAr:          formDTO.ObjectiveAr,
		CreatedBYID:          currentUserID,
	}

	// Save to database
	id, err := ccs.repo.Create(complianceClause, tx)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create compliance clause: %v", err)
	}

	tx.Commit()

	go func() {
		err := ccs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceClause{}.TableName(),
			id,
			formDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceClause{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log control activity", fields)
		}
	}()
	return id, nil
}

func (ccs *ComplianceClauseService) Update(complianceID uint, clauseID uint, currentUserID uint, formDTO *complianceDTOs.ComplianceClauseFormDTO) error {
	tx := ccs.repo.BeginTransaction()

	// Fetch existing clause
	clause, err := ccs.repo.GetByID(clauseID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update fields
	clause.TitleEn = formDTO.TitleEn
	clause.TitleAr = formDTO.TitleAr
	clause.DescriptionEn = formDTO.DescriptionEn
	clause.DescriptionAr = formDTO.DescriptionAr
	clause.ScopeEn = formDTO.ScopeEn
	clause.ScopeAr = formDTO.ScopeAr
	clause.ObjectiveEn = formDTO.ObjectiveEn
	clause.ObjectiveAr = formDTO.ObjectiveAr
	clause.UpdatedBYID = &currentUserID

	// Save updates
	err = ccs.repo.Update(clause, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (ccs *ComplianceClauseService) Delete(clauseID uint) error {
	tx := ccs.repo.BeginTransaction()

	// Recursively delete the clause and its children
	err := ccs.repo.Delete(clauseID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (ccs *ComplianceClauseService) GetAssociatedComplianceRequirements(complianceID uint, clauseID uint) ([]complianceDTOs.ComplianceRequirementListDTO, error) {
	// Fetch associated compliance requirements from repository
	requirements, err := ccs.repo.GetAssociatedComplianceRequirements(complianceID, clauseID)
	if err != nil {
		return nil, err
	}

	return requirements, nil
}
