package services

import (
	"fmt"
	"io"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"
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
	GetFullDetailsByID(clauseID uint) (*complianceDTOs.ComplianceClauseDetailDTO, error)

	GetAssociatedComplianceRequirements(complianceID uint, clauseID uint) ([]complianceDTOs.ComplianceRequirementListDTO, error)
	GetAssociatedControls(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	CreateComplianceClauseControls(complianceID, clauseID, currentUserID uint, dto *complianceDTOs.ComplianceClauseControlMapDTO) error
	DeleteComplianceClauseControl(clauseID, controlID uint) error
	GetClausesByIDs(ids []uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceClause, error)

	GetAssociatedRisks(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	CreateComplianceClauseRisk(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteComplianceClauseRisk(clauseID, riskID uint) error

	GetAssociatedDocuments(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateComplianceClauseDocument(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteComplianceClauseDocument(clauseID, documentID uint) error

	GetAssociatedIncident(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error)
	CreateComplianceClauseIncident(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error
	DeleteComplianceClauseIncident(clauseID, incidentID uint) error

	GetAssociatedStrategicObjective(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error)
	CreateComplianceClauseStrategicObjective(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error
	DeleteComplianceClausStrategicObjective(clauseID, strategicObjectiveID uint) error

	GetAssociatedAssets(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error)
	CreateComplianceClauseAssets(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedAssetMapDTO) error
	DeleteComplianceAssets(clauseID, assetID uint) error

	GetAssociatedProcesses(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	CreateComplianceClausProcesses(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedProcessMapDTO) error
	DeleteComplianceClausProcesses(clauseID, processID uint) error

	GetAssociatedKPI(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error)
	CreateComplianceClausKPI(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedKPIMapDTO) error
	DeleteComplianceClausKPI(clauseID, kpiID uint) error

	GenerateClauseTemplate() (*os.File, error)
	ValidateClauseTemplate(complianceID uint, file io.Reader) ([]map[string]interface{}, error)
	ValidateImportClause(complianceID uint, input []map[string]interface{}) ([]map[string]interface{}, error)
	validateClauses(complianceID uint, input []map[string]interface{}) ([]map[string]interface{}, error)
	ImportClauses(complianceID uint, input []map[string]interface{}, currentUserID uint) ([]uint, error)
}

type ComplianceClauseService struct {
	repo           repositories.ComplianceClauseRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	controlSvc     organizationServices.ControlServiceInterface
	validationSvc  commonServices.AssociatedItemValidationServiceInterface
}

func NewComplianceClauseService(
	repo repositories.ComplianceClauseRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) *ComplianceClauseService {
	return &ComplianceClauseService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		controlSvc:     controlSvc,
		validationSvc:  validationSvc,
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
			ID:                   clause.ID,
			ClauseNo:             clause.ClauseNo,
			ComplianceRegisterID: clause.ComplianceRegisterID,
			ParentClauseID:       clause.ParentClauseID,
			TitleEn:              clause.TitleEn,
			TitleAr:              clause.TitleAr,
			ScopeEn:              clause.ScopeEn,
			ScopeAr:              clause.ScopeAr,
			ObjectiveEn:          clause.ObjectiveEn,
			ObjectiveAr:          clause.ObjectiveAr,
			ChildClauses:         []*complianceDTOs.ComplianceClauseListDTO{},
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

func (ccs *ComplianceClauseService) GetFullDetailsByID(clauseID uint) (*complianceDTOs.ComplianceClauseDetailDTO, error) {
	var clauseDTO complianceDTOs.ComplianceClauseDetailDTO
	clause, err := ccs.repo.GetFullDetailsByID(clauseID)
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(clause, &clauseDTO)
	if err != nil {
		return nil, err
	}

	return &clauseDTO, nil
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

func (ccs *ComplianceClauseService) GetAssociatedControls(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	// Fetch associated controls through requirements from repository
	controls, err := ccs.repo.GetAssociatedControls(complianceID, clauseID, params)
	if err != nil {
		return nil, err
	}

	return controls, nil
}

func (ccs *ComplianceClauseService) CreateComplianceClauseControls(complianceID, clauseID, currentUserID uint, dto *complianceDTOs.ComplianceClauseControlMapDTO) error {
	if err := ccs.validateControlIDs(clauseID, dto.ControlIDs); err != nil {
		return err
	}

	tx := ccs.repo.BeginTransaction()

	var controlMappings []*models.ComplianceControl

	for _, controlID := range dto.ControlIDs {
		controlMappings = append(controlMappings, &models.ComplianceControl{
			ComplianceRegisterID: complianceID,
			ComplianceClauseID:   &clauseID,
			ControlID:            controlID,
			CreatedByID:          currentUserID,
		})
	}

	if err := ccs.repo.CreateComplianceControls(controlMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ccs *ComplianceClauseService) DeleteComplianceClauseControl(clauseID, controlID uint) error {
	if _, err := ccs.controlSvc.GetByID(controlID); err != nil {
		return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", controlID), "not_exist", "control_id")
	}
	tx := ccs.repo.BeginTransaction()
	if _, err := ccs.repo.GetByID(clauseID, tx); err != nil {
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", controlID), "not_exist", "clause_id")
	}
	if err := ccs.repo.DeleteComplianceControlByClauseAndControl(clauseID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (ccs *ComplianceClauseService) validateControlIDs(clauseID uint, controlIDs []uint) error {
	if len(controlIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no controls provided"), "required", "control_ids")
	}

	seen := make(map[uint]bool)
	for _, id := range controlIDs {
		alreadyMapped, err := ccs.repo.IsControlMappedToClause(clauseID, id)
		if err != nil {
			return fmt.Errorf("failed to check existing mapping for control ID %d: %v", id, err)
		}
		if alreadyMapped {
			return utility.NewValidationError(fmt.Errorf("control ID %d is already mapped to this clause", id), "duplicate", "control_ids")
		}
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", id), "duplicate", "control_ids")
		}
		if _, err := ccs.controlSvc.GetByID(id); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", id), "not_exist", "control_ids")
		}
		seen[id] = true
	}
	return nil
}
func (ccs *ComplianceClauseService) GetClausesByIDs(ids []uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceClause, error) {
	clauses, err := ccs.repo.GetClausesByIDs(ids, tx...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch compliance clauses: %v", err)
	}
	return clauses, nil
}

func (ccs *ComplianceClauseService) GetAssociatedRisks(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	// Fetch associated controls through requirements from repository
	risks, err := ccs.repo.GetAssociatedRisks(complianceID, clauseID, params)
	if err != nil {
		return nil, err
	}

	return risks, nil
}

func (ccs *ComplianceClauseService) CreateComplianceClauseRisk(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {

	if err := ccs.validationSvc.ValidateRiskMapping(clauseID, dto.RiskIDs, "risk_compliance_requirements", "compliance_clause_id"); err != nil {
		return err
	}
	tx := ccs.repo.BeginTransaction()
	_, err := ccs.repo.GetByID(clauseID, tx)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}

	var riskMappings []*models.RiskComplianceRequirement

	for _, riskID := range dto.RiskIDs {
		riskMappings = append(riskMappings, &models.RiskComplianceRequirement{
			ComplianceRegisterID: complianceID,
			ComplianceClauseID:   &clauseID,
			RiskRegisterID:       riskID,
			CreatedByID:          currentUserID,
		})
	}

	if err := ccs.repo.CreateComplianceClauseRisk(riskMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ccs *ComplianceClauseService) DeleteComplianceClauseRisk(clauseID, riskID uint) error {
	tx := ccs.repo.BeginTransaction()
	if _, err := ccs.repo.GetByID(clauseID, tx); err != nil {
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}
	if err := ccs.repo.DeleteComplianceClauseRisk(clauseID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (ccs *ComplianceClauseService) GetAssociatedDocuments(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	// Fetch associated controls through requirements from repository
	documents, err := ccs.repo.GetAssociatedDocuments(complianceID, clauseID, params)
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (ccs *ComplianceClauseService) CreateComplianceClauseDocument(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {
	if err := ccs.validationSvc.ValidateDocumentMapping(clauseID, dto.DocumentIDs, "document_compliance_requirements", "compliance_clause_id"); err != nil {
		return err
	}
	tx := ccs.repo.BeginTransaction()
	_, err := ccs.repo.GetByID(clauseID, tx)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}
	var documentMappings []*models.DocumentComplianceRequirement

	for _, documentID := range dto.DocumentIDs {
		documentMappings = append(documentMappings, &models.DocumentComplianceRequirement{
			ComplianceRegisterID: complianceID,
			ComplianceClauseID:   &clauseID,
			DocumentID:           documentID,
			CreatedByID:          currentUserID,
		})
	}

	if err := ccs.repo.CreateComplianceClauseDocument(documentMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ccs *ComplianceClauseService) DeleteComplianceClauseDocument(clauseID, documentId uint) error {
	tx := ccs.repo.BeginTransaction()
	if _, err := ccs.repo.GetByID(clauseID, tx); err != nil {
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}
	if err := ccs.repo.DeleteComplianceClauseDocument(clauseID, documentId, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (ccs *ComplianceClauseService) GetAssociatedIncident(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error) {
	// Fetch associated controls through requirements from repository
	incidents, err := ccs.repo.GetAssociatedIncident(complianceID, clauseID, params)
	if err != nil {
		return nil, err
	}

	return incidents, nil
}

func (ccs *ComplianceClauseService) CreateComplianceClauseIncident(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error {
	if err := ccs.validationSvc.ValidateIncidentMapping(clauseID, dto.IncidentIDs, "incident_compliance_requirements", "compliance_clause_id"); err != nil {
		return err
	}
	tx := ccs.repo.BeginTransaction()
	_, err := ccs.repo.GetByID(clauseID, tx)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}
	var incidentMappings []*models.IncidentComplianceRequirement

	for _, incidentID := range dto.IncidentIDs {
		incidentMappings = append(incidentMappings, &models.IncidentComplianceRequirement{
			ComplianceRegisterID: complianceID,
			ComplianceClauseID:   clauseID,
			IncidentRegisterID:   incidentID,
			CreatedByID:          currentUserID,
		})
	}

	if err := ccs.repo.CreateComplianceClauseIncident(incidentMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ccs *ComplianceClauseService) DeleteComplianceClauseIncident(clauseID, incidentId uint) error {
	tx := ccs.repo.BeginTransaction()
	if _, err := ccs.repo.GetByID(clauseID, tx); err != nil {
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}
	if err := ccs.repo.DeleteComplianceClauseIncident(clauseID, incidentId, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (ccs *ComplianceClauseService) GetAssociatedStrategicObjective(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error) {
	// Fetch associated controls through requirements from repository
	strategicObjectives, err := ccs.repo.GetAssociatedStrategicObjective(complianceID, clauseID, params)
	if err != nil {
		return nil, err
	}

	return strategicObjectives, nil
}

func (ccs *ComplianceClauseService) CreateComplianceClauseStrategicObjective(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error {
	if err := ccs.validationSvc.ValidateStrtegicObjectiveMapping(clauseID, dto.StrategicObjectiveIDs, "compliance_requirement_strategic_objectives", "compliance_clause_id"); err != nil {
		return err
	}
	tx := ccs.repo.BeginTransaction()
	_, err := ccs.repo.GetByID(clauseID, tx)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}
	var strategicObjectiveMappings []*models.ComplianceRequirementStrategicObjective

	for _, strtegicObjectiveID := range dto.StrategicObjectiveIDs {
		strategicObjectiveMappings = append(strategicObjectiveMappings, &models.ComplianceRequirementStrategicObjective{
			ComplianceRegisterID: complianceID,
			ComplianceClauseID:   clauseID,
			StrategicObjectiveID: strtegicObjectiveID,
			CreatedByID:          currentUserID,
		})
	}

	if err := ccs.repo.CreateComplianceClauseStrategicObjective(strategicObjectiveMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ccs *ComplianceClauseService) DeleteComplianceClausStrategicObjective(clauseID, strategicObjectiveId uint) error {
	tx := ccs.repo.BeginTransaction()
	if _, err := ccs.repo.GetByID(clauseID, tx); err != nil {
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}
	if err := ccs.repo.DeleteComplianceClausStrategicObjective(clauseID, strategicObjectiveId, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (ccs *ComplianceClauseService) GetAssociatedAssets(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error) {
	// Fetch associated controls through requirements from repository
	strategicObjectives, err := ccs.repo.GetAssociatedAssets(complianceID, clauseID, params)
	if err != nil {
		return nil, err
	}

	return strategicObjectives, nil
}

func (ccs *ComplianceClauseService) CreateComplianceClauseAssets(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedAssetMapDTO) error {
	if err := ccs.validationSvc.ValidateAssetMapping(clauseID, dto.AssetIDs, "asset_compliance_requirements", "compliance_clause_id"); err != nil {
		return err
	}
	tx := ccs.repo.BeginTransaction()
	_, err := ccs.repo.GetByID(clauseID, tx)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}
	var assetMappings []*models.AssetComplianceRequirement

	for _, assetID := range dto.AssetIDs {
		assetMappings = append(assetMappings, &models.AssetComplianceRequirement{
			ComplianceRegisterID: complianceID,
			ComplianceClauseID:   &clauseID,
			AssetRegisterID:      assetID,
			CreatedBYID:          currentUserID,
		})
	}

	if err := ccs.repo.CreateComplianceClauseAssets(assetMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ccs *ComplianceClauseService) DeleteComplianceAssets(clauseID, assetID uint) error {
	tx := ccs.repo.BeginTransaction()
	if _, err := ccs.repo.GetByID(clauseID, tx); err != nil {
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}
	if err := ccs.repo.DeleteComplianceAssets(clauseID, assetID, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (ccs *ComplianceClauseService) GetAssociatedProcesses(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	// Fetch associated controls through requirements from repository
	strategicObjectives, err := ccs.repo.GetAssociatedProcesses(complianceID, clauseID, params)
	if err != nil {
		return nil, err
	}

	return strategicObjectives, nil
}

func (ccs *ComplianceClauseService) CreateComplianceClausProcesses(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedProcessMapDTO) error {
	if err := ccs.validationSvc.ValidateProcessMapping(clauseID, dto.ProcessIDs, "process_compliance_requirements", "compliance_clause_id"); err != nil {
		return err
	}
	tx := ccs.repo.BeginTransaction()
	_, err := ccs.repo.GetByID(clauseID, tx)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}
	var processMappings []*models.ProcessComplianceRequirement

	for _, processID := range dto.ProcessIDs {
		processMappings = append(processMappings, &models.ProcessComplianceRequirement{
			ComplianceRegisterID: complianceID,
			ComplianceClauseID:   clauseID,
			ProcessID:            processID,
			CreatedByID:          currentUserID,
		})
	}

	if err := ccs.repo.CreateComplianceClausProcesses(processMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ccs *ComplianceClauseService) DeleteComplianceClausProcesses(clauseID, processID uint) error {
	tx := ccs.repo.BeginTransaction()
	if _, err := ccs.repo.GetByID(clauseID, tx); err != nil {
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}
	if err := ccs.repo.DeleteComplianceClausProcesses(clauseID, processID, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (ccs *ComplianceClauseService) GetAssociatedKPI(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error) {

	kpis, err := ccs.repo.GetAssociatedKPI(complianceID, clauseID, params)
	if err != nil {
		return nil, err
	}

	return kpis, nil
}

func (ccs *ComplianceClauseService) CreateComplianceClausKPI(complianceID, clauseID, currentUserID uint, dto *commonDTOs.AssociatedKPIMapDTO) error {
	if err := ccs.validationSvc.ValidateKPIMapping(clauseID, dto.KpiIDs, "kpi_compliances", "compliance_clause_id"); err != nil {
		return err
	}
	tx := ccs.repo.BeginTransaction()
	_, err := ccs.repo.GetByID(clauseID, tx)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}
	var mappings []*models.KPICompliance

	for _, kpiID := range dto.KpiIDs {
		mappings = append(mappings, &models.KPICompliance{
			ComplianceRegisterID: complianceID,
			ComplianceClauseID:   &clauseID,
			KPIRegisterID:        kpiID,
			CreatedByID:          currentUserID,
		})
	}

	if err := ccs.repo.CreateComplianceClausKPI(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ccs *ComplianceClauseService) DeleteComplianceClausKPI(clauseID, kpiID uint) error {
	tx := ccs.repo.BeginTransaction()
	if _, err := ccs.repo.GetByID(clauseID, tx); err != nil {
		return utility.NewValidationError(fmt.Errorf("clause does not exist for ID: %d", clauseID), "not_exist", "clause_id")
	}
	if err := ccs.repo.DeleteComplianceClausKPI(clauseID, kpiID, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (css *ComplianceClauseService) GenerateClauseTemplate() (*os.File, error) {
	templateHeaders := []string{"ClauseNo", "TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "ScopeEn", "ScopeAr", "ObjectiveEn", "ObjectiveAr"}

	templateData := [][]interface{}{}

	// Generate file
	file, err := utility.GenerateExcelFileWithSheets(map[string][][]interface{}{
		"Template": append([][]interface{}{utility.ConvertToInterfaceSlice(templateHeaders)}, templateData...),
	})
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (css *ComplianceClauseService) ValidateClauseTemplate(complianceID uint, file io.Reader) ([]map[string]interface{}, error) {
	rows, err := utility.ReadExcelFile(file)
	if err != nil {
		return nil, err
	}

	inputHeaders := []string{"ClauseNo", "TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "ScopeEn", "ScopeAr", "ObjectiveEn", "ObjectiveAr"}
	snakeHeaders := []string{"clause_no", "title_en", "title_ar", "description_en", "description_ar", "scope_en", "scope_ar", "objective_en", "objective_ar"}

	if !utility.ValidateHeaders(inputHeaders, rows[0]) {
		return nil, utility.NewValidationError(fmt.Errorf("invalid headers"), "invalid_header", "Requirement Template")
	}

	var input []map[string]interface{}
	for _, row := range rows[1:] {
		record := make(map[string]interface{})
		for i, h := range snakeHeaders {
			if i < len(row) {
				record[h] = row[i]
			} else {
				record[h] = ""
			}
		}
		input = append(input, record)
	}

	return css.validateClauses(complianceID, input)
}

func (css *ComplianceClauseService) ValidateImportClause(complianceID uint, input []map[string]interface{}) ([]map[string]interface{}, error) {
	return css.validateClauses(complianceID, input)
}

func (css *ComplianceClauseService) validateClauses(complianceID uint, input []map[string]interface{}) ([]map[string]interface{}, error) {
	headers := []string{"clause_no", "title_en", "title_ar", "description_en", "description_ar", "scope_en", "scope_ar", "objective_en", "objective_ar"}

	// Clauses
	clauses, err := css.GetAll(complianceID)
	if err != nil {
		return nil, fmt.Errorf("error fetching clauses: %v", err)
	}
	clauseMap := make(map[string]uint)
	for _, cl := range clauses {
		key := strings.TrimSpace(cl.ClauseNo)
		clauseMap[key] = cl.ID
	}

	var validated []map[string]interface{}

	for _, record := range input {
		rowErrors := make(map[string]interface{})
		processed := make(map[string]interface{})

		// Copy original input
		for k, v := range record {
			processed[k] = v
		}

		// Clause
		if val, ok := record["clause_no"].(string); ok && val != "" {
			if _, exists := clauseMap[strings.TrimSpace(val)]; exists {
				rowErrors["clause_no"] = fmt.Sprintf("ClauseNo already exist: %s", val)
			} else {
				processed["clause_no"] = val
			}
		}

		// TitleEn
		if val, ok := record["title_en"].(string); !ok || val == "" {
			rowErrors["title_en"] = "Field is required"
		}

		// TitleAr
		if val, ok := record["title_ar"].(string); !ok || val == "" {
			rowErrors["title_ar"] = "Field is required"
		}

		// DescriptionEn
		if val, ok := record["description_en"].(string); ok {
			trimmed := strings.TrimSpace(val)
			if trimmed == "" {
				processed["description_en"] = nil
			} else {
				processed["description_en"] = trimmed
			}
		}

		// DescriptionAr
		if val, ok := record["description_ar"].(string); ok {
			trimmed := strings.TrimSpace(val)
			if trimmed == "" {
				processed["description_ar"] = nil
			} else {
				processed["description_ar"] = trimmed
			}
		}

		// // ScopeEn
		if val, ok := record["scope_en"].(string); ok {
			trimmed := strings.TrimSpace(val)
			if trimmed == "" {
				processed["scope_en"] = nil
			} else {
				processed["scope_en"] = trimmed
			}
		}

		// // ScopeAr
		if val, ok := record["scope_ar"].(string); ok {
			trimmed := strings.TrimSpace(val)
			if trimmed == "" {
				processed["scope_ar"] = nil
			} else {
				processed["scope_ar"] = trimmed
			}
		}

		// ObjectiveEn
		if val, ok := record["objective_en"].(string); ok {
			trimmed := strings.TrimSpace(val)
			if trimmed == "" {
				processed["objective_en"] = nil
			} else {
				processed["objective_en"] = trimmed
			}
		}

		// ObjectiveAr
		if val, ok := record["objective_ar"].(string); ok {
			trimmed := strings.TrimSpace(val)
			if trimmed == "" {
				processed["objective_ar"] = nil
			} else {
				processed["objective_ar"] = trimmed
			}
		}

		if len(rowErrors) == 0 {
			processed["errors"] = nil
		} else {
			for _, h := range headers {
				if _, exists := rowErrors[h]; !exists {
					rowErrors[h] = ""
				}
			}
			processed["errors"] = rowErrors
		}

		validated = append(validated, processed)
	}
	return validated, nil
}

func (css *ComplianceClauseService) ImportClauses(complianceID uint, input []map[string]interface{}, currentUserID uint) ([]uint, error) {
	validated, err := css.validateClauses(complianceID, input)
	if err != nil {
		return nil, err
	}

	for _, r := range validated {
		if errs, ok := r["errors"].(map[string]interface{}); ok && len(errs) > 0 {
			return nil, utility.NewValidationError(fmt.Errorf("invalid data"), "invalid_data", "Invalid Data")
		}
	}

	var createdIDs []uint
	for _, record := range validated {

		dto := &complianceDTOs.ComplianceClauseFormDTO{
			ClauseNo: record["clause_no"].(string),
			TitleEn:  record["title_en"].(string),
		}

		if v, ok := record["title_ar"].(string); ok {
			dto.TitleAr = &v
		}
		if v, ok := record["description_en"].(string); ok {
			dto.DescriptionEn = &v
		}
		if v, ok := record["description_ar"].(string); ok {
			dto.DescriptionAr = &v
		}
		if v, ok := record["scope_en"].(string); ok {
			dto.ScopeEn = &v
		}
		if v, ok := record["scope_ar"].(string); ok {
			dto.ScopeAr = &v
		}
		if v, ok := record["objective_en"].(string); ok {
			dto.ObjectiveEn = &v
		}
		if v, ok := record["objective_ar"].(string); ok {
			dto.ObjectiveAr = &v
		}

		id, err := css.Create(complianceID, currentUserID, dto)
		if err != nil {
			return nil, err
		}
		createdIDs = append(createdIDs, id)
	}
	return createdIDs, nil
}
