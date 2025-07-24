package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	"jethur/utility"
	"time"
)

type ControlAssociatedItemServiceInterface interface {
	GetAllProcesses(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemProcessListResponseDTO, error)
	GetAllComplianceRequirements(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) ([]organizationDTOs.ComplianceRegisterDTO, error)
	GetAllRisks(controlID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetAllIncidents(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemIncidentListResponseDTO, error)
	GetAllAssociatedCounts(controlID uint) (*organizationDTOs.ControlAssociatedItemsCountResponseDTO, error)
	GetAllClauses(controlID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	UpdateClauses(controlID uint, currentUserID uint, formDTO *organizationDTOs.ControlClauseFormDTO) error
	CreateControlRisk(controlID uint, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteControlRisk(controlID uint, riskID uint) error
}

type ControlAssociatedItemService struct {
	repo                           organizationRepositories.ControlAssociatedItemRepositoryInterface
	controlSvc                     ControlServiceInterface
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface
	controlEfficiencyStatusService ControlEfficiencyStatusServiceInterface
	validationSvc                  commonServices.AssociatedItemValidationServiceInterface
}

func NewControlAssociatedItemService(
	repo organizationRepositories.ControlAssociatedItemRepositoryInterface,
	controlSvc ControlServiceInterface,
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusService ControlEfficiencyStatusServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) ControlAssociatedItemServiceInterface {
	return &ControlAssociatedItemService{
		repo:                           repo,
		controlSvc:                     controlSvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusService: controlEfficiencyStatusService,
		validationSvc:                  validationSvc,
	}
}

func (cais *ControlAssociatedItemService) GetAllProcesses(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemProcessListResponseDTO, error) {
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllProcesses(controlID, params)
}

func (cais *ControlAssociatedItemService) GetAllComplianceRequirements(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) ([]organizationDTOs.ComplianceRegisterDTO, error) {
	// Validate the control ID
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	// Fetch all compliance requirements
	complianceRequirementList, err := cais.repo.GetAllComplianceRequirements(controlID, params)
	if err != nil {
		return nil, err
	}

	// Group by compliance_register_id
	groupedRequirements := make(map[uint]*organizationDTOs.ComplianceRegisterDTO)

	for _, item := range complianceRequirementList {
		registerID := item.ComplianceRegisterID

		// Check if the compliance register is already added to the group
		if _, exists := groupedRequirements[registerID]; !exists {
			groupedRequirements[registerID] = &organizationDTOs.ComplianceRegisterDTO{
				ComplianceRegisterID:      item.ComplianceRegisterID,
				ComplianceRegisterTitleEn: item.ComplianceRegisterTitleEn,
				ComplianceRegisterTitleAr: item.ComplianceRegisterTitleAr,
				Requirements:              []*organizationDTOs.ComplianceRequirementMinimalDTO{},
			}
		}

		if item.ComplianceRequirementID == 0 {
			continue
		}
		// Add the compliance requirement to the corresponding register
		groupedRequirements[registerID].Requirements = append(groupedRequirements[registerID].Requirements, &organizationDTOs.ComplianceRequirementMinimalDTO{
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
		})
	}

	// Convert the map to a slice
	var complianceRegisters []organizationDTOs.ComplianceRegisterDTO
	for _, register := range groupedRequirements {
		complianceRegisters = append(complianceRegisters, *register)
	}

	return complianceRegisters, nil
}

func (cais *ControlAssociatedItemService) GetAllRisks(controlID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllRisks(controlID, params)
}

func (cais *ControlAssociatedItemService) GetAllIncidents(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemIncidentListResponseDTO, error) {
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllIncidents(controlID, params)
}

func (cais *ControlAssociatedItemService) GetAllAssociatedCounts(controlID uint) (*organizationDTOs.ControlAssociatedItemsCountResponseDTO, error) {
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}
	response, err := cais.repo.GetAllAssociatedCounts(controlID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (cais *ControlAssociatedItemService) GetAllClauses(controlID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	response, err := cais.repo.GetAllClauses(controlID, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (cais *ControlAssociatedItemService) UpdateClauses(controlID uint, currentUserID uint, ControlClauseFormDTO *organizationDTOs.ControlClauseFormDTO) error {

	tx := cais.repo.BeginTransaction()

	if _, err := cais.controlSvc.GetByID(controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	clauses, err := cais.validateClauseIDs(ControlClauseFormDTO.ClauseIDs, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlClauses, err := cais.repo.GetClausesByControlID(controlID, tx)
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

	var clausesToAdd []*models.ComplianceControl
	now := time.Now()
	for id, clause := range newClauses {
		if !existingClauses[id] {
			clausesToAdd = append(clausesToAdd, &models.ComplianceControl{
				ControlID:            controlID,
				ComplianceRegisterID: clause.ComplianceRegisterID,
				ComplianceClauseID:   &clause.ID,
				CreatedByID:          currentUserID,
				CreatedAt:            now,
			})
		}
	}

	if len(clausesToAdd) > 0 {
		if err := cais.repo.CreateControlClauses(clausesToAdd, tx); err != nil {
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
		if err := cais.repo.DeleteControlClauses(controlID, deleteClauseIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (cais *ControlAssociatedItemService) validateClauseIDs(clauseIDs []uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceClause, error) {
	seen := make(map[uint]bool)
	for _, id := range clauseIDs {
		if seen[id] {
			return nil, utility.NewValidationError(fmt.Errorf("duplicate clause ID: %d", id), "duplicate", "clause_id")
		}
		seen[id] = true
	}

	clauses, err := cais.repo.GetClausesByIDs(clauseIDs, tx...)
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

func (cais *ControlAssociatedItemService) CreateControlRisk(controlID, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", controlID), "not_exist", "control_id")
	}
	if err := cais.validationSvc.ValidateRiskMapping(controlID, dto.RiskIDs, "risk_controls", "control_id"); err != nil {
		return err
	}
	tx := cais.repo.BeginTransaction()
	var mappings []*models.RiskControl
	currentTime := time.Now()

	for _, riskID := range dto.RiskIDs {
		mappings = append(mappings, &models.RiskControl{
			RiskRegisterID: riskID,
			ControlID:      controlID,
			CreatedByID:    userID,
			CreatedAt:      currentTime,
		})
	}

	if err := cais.repo.CreateControlRisk(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (cais *ControlAssociatedItemService) DeleteControlRisk(controlID uint, riskID uint) error {
	// Validate control existence
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", controlID), "not_exist", "control_id")
	}

	tx := cais.repo.BeginTransaction()

	if err := cais.repo.DeleteControlRisk(controlID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// func (cais *ControlAssociatedItemService) validateRiskIDs(controlID uint, riskIDs []uint) error {
// 	if len(riskIDs) == 0 {
// 		return utility.NewValidationError(fmt.Errorf("no risks provided"), "required", "risk_ids")
// 	}
// 	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
// 		return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", controlID), "not_exist", "control_id")
// 	}
// 	seen := make(map[uint]bool)
// 	for _, id := range riskIDs {
// 		alreadyMapped, err := cais.repo.IsRiskMappedToControl(controlID, id)
// 		if err != nil {
// 			return fmt.Errorf("failed to check existing mapping: %v", err)
// 		}
// 		if alreadyMapped {
// 			return utility.NewValidationError(fmt.Errorf("risk ID %d already mapped to this control", id), "duplicate", "risk_ids")
// 		}
// 		if seen[id] {
// 			return utility.NewValidationError(fmt.Errorf("duplicate risk ID: %d", id), "duplicate", "risk_ids")
// 		}
// 		if _, err := cais.controlSvc.GetByID(id); err != nil {
// 			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", id), "not_exist", "control_ids")
// 		}
// 		seen[id] = true
// 	}
// 	return nil
// }
