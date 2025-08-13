package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceServices "jethur/grc/compliance/services"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	"jethur/grc/models"
	"jethur/utility"
	"time"
)

type IncidentAssociatedItemServiceInterface interface {
	GetAllControls(incidentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	GetAllComplianceRequirements(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemComplianceRequirementListResponseDTO, error)
	GetAllRisks(incidentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetAllProcesses(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemProcessListResponseDTO, error)
	GetAllSLAs(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemSLAListResponseDTO, error)
	GetAllDocuments(incidentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	GetAllAssets(incidentID uint, params *incidentDTOs.IncidentAssociatedItemListQueryParamsDTO) (*incidentDTOs.IncidentAssociatedItemAssetListResponseDTO, error)
	GetAllClauses(incidentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	UpdateClauses(incidentID uint, currentUserID uint, ControlClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error
	CreateIncidentRisk(incidentID uint, currentUserID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteIncidentRisk(incidentID uint, riskID uint) error
	CreateIncidentDocument(incidentID uint, currentUserID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteIncidentDocument(incidentID uint, documentID uint) error
	CreateIncidentControl(incidentID uint, currentUserID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteIncidentControl(incidentID uint, controlID uint) error
}

// Service and implementations
type IncidentAssociatedItemService struct {
	repo                incidentRepositories.IncidentAssociatedItemRepositoryInterface
	incidentRegisterSvc IncidentRegisterServiceInterface
	clauseSvc           complianceServices.ComplianceClauseServiceInterface
	validationSvc       commonServices.AssociatedItemValidationServiceInterface
}

func NewIncidentAssociatedItemService(
	repo incidentRepositories.IncidentAssociatedItemRepositoryInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface,
	clauseSvc complianceServices.ComplianceClauseServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,

) IncidentAssociatedItemServiceInterface {
	return &IncidentAssociatedItemService{
		repo:                repo,
		incidentRegisterSvc: incidentRegisterSvc,
		clauseSvc:           clauseSvc,
		validationSvc:       validationSvc,
	}
}

func (iais *IncidentAssociatedItemService) GetAllControls(incidentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
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

func (iais *IncidentAssociatedItemService) GetAllRisks(incidentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
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

func (iais *IncidentAssociatedItemService) GetAllDocuments(incidentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
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

func (iais *IncidentAssociatedItemService) GetAllClauses(incidentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
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

	var clausesToAdd []*models.IncidentComplianceRequirement
	now := time.Now()
	for id, clause := range newClauses {
		if !existingClauses[id] {
			clausesToAdd = append(clausesToAdd, &models.IncidentComplianceRequirement{
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

func (iais *IncidentAssociatedItemService) CreateIncidentRisk(incidentID, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident does not exist for ID: %d", incidentID), "not_exist", "incident_id")
	}
	if err := iais.validationSvc.ValidateRiskMapping(incidentID, dto.RiskIDs, "incident_risks", "incident_register_id"); err != nil {
		return err
	}
	tx := iais.repo.BeginTransaction()
	var mappings []*models.IncidentRisk
	currentTime := time.Now()

	for _, riskID := range dto.RiskIDs {
		mappings = append(mappings, &models.IncidentRisk{
			RiskRegisterID:     riskID,
			IncidentRegisterID: incidentID,
			CreatedByID:        userID,
			CreatedAt:          currentTime,
		})
	}

	if err := iais.repo.CreateIncidentRisk(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (iais *IncidentAssociatedItemService) DeleteIncidentRisk(incidentID uint, riskID uint) error {
	// Validate control existence
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident does not exist for ID: %d", incidentID), "not_exist", "incident_id")
	}

	tx := iais.repo.BeginTransaction()

	if err := iais.repo.DeleteIncidentRisk(incidentID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (iais *IncidentAssociatedItemService) CreateIncidentDocument(incidentID, currentUserID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident does not exist for ID: %d", incidentID), "not_exist", "incident_id")
	}
	if err := iais.validationSvc.ValidateDocumentMapping(incidentID, dto.DocumentIDs, "incident_associated_documents", "incident_register_id"); err != nil {
		return err
	}
	tx := iais.repo.BeginTransaction()
	var mappings []*models.IncidentAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.IncidentAssociatedDocument{
			DocumentID:         documentID,
			IncidentRegisterID: incidentID,
			CreatedByID:        currentUserID,
			CreatedAt:          currentTime,
		})
	}

	if err := iais.repo.CreateIncidentDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (iais *IncidentAssociatedItemService) DeleteIncidentDocument(incidentID uint, documentID uint) error {
	// Validate control existence
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident does not exist for ID: %d", incidentID), "not_exist", "incident_id")
	}

	tx := iais.repo.BeginTransaction()

	if err := iais.repo.DeleteIncidentDocument(incidentID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
func (iais *IncidentAssociatedItemService) CreateIncidentControl(incidentID, currentUserID uint, dto *commonDTOs.AssociatedControlMapDTO) error {
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident does not exist for ID: %d", incidentID), "not_exist", "incident_id")
	}
	if err := iais.validationSvc.ValidateControlMapping(incidentID, dto.ControlIDs, "incident_controls", "incident_register_id"); err != nil {
		return err
	}
	tx := iais.repo.BeginTransaction()
	var mappings []*models.IncidentControl
	currentTime := time.Now()

	for _, controlID := range dto.ControlIDs {
		mappings = append(mappings, &models.IncidentControl{
			ControlID:          controlID,
			IncidentRegisterID: incidentID,
			CreatedByID:        currentUserID,
			CreatedAt:          currentTime,
		})
	}

	if err := iais.repo.CreateIncidentControl(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (iais *IncidentAssociatedItemService) DeleteIncidentControl(incidentID uint, controlID uint) error {
	// Validate control existence
	if _, err := iais.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident does not exist for ID: %d", incidentID), "not_exist", "incident_id")
	}

	tx := iais.repo.BeginTransaction()

	if err := iais.repo.DeleteIncidentControl(incidentID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
