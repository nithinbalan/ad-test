package services

import (
	"fmt"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	"jethur/utility"
	"time"
)

type ControlAssociatedItemServiceInterface interface {
	GetAllProcesses(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemProcessListResponseDTO, error)
	GetAllComplianceRequirements(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) ([]organizationDTOs.ComplianceRegisterDTO, error)
	GetAllRisks(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemRiskListResponseDTO, error)
	GetAllIncidents(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemIncidentListResponseDTO, error)
	GetAllAssociatedCounts(controlID uint) (*organizationDTOs.ControlAssociatedItemsCountResponseDTO, error)
	GetAllClauses(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemClauseListResponseDTO, error)
	UpdateClauses(controlID uint, currentUserID uint, formDTO *organizationDTOs.ControlClauseFormDTO) (uint, error)
}

type ControlAssociatedItemService struct {
	repo                           organizationRepositories.ControlAssociatedItemRepositoryInterface
	controlSvc                     ControlServiceInterface
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface
	controlEfficiencyStatusService ControlEfficiencyStatusServiceInterface
}

func NewControlAssociatedItemService(
	repo organizationRepositories.ControlAssociatedItemRepositoryInterface,
	controlSvc ControlServiceInterface,
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusService ControlEfficiencyStatusServiceInterface,
) ControlAssociatedItemServiceInterface {
	return &ControlAssociatedItemService{
		repo:                           repo,
		controlSvc:                     controlSvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusService: controlEfficiencyStatusService,
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

func (cais *ControlAssociatedItemService) GetAllRisks(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemRiskListResponseDTO, error) {
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

func (cais *ControlAssociatedItemService) GetAllClauses(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemClauseListResponseDTO, error) {
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	response, err := cais.repo.GetAllClauses(controlID, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (cais *ControlAssociatedItemService) UpdateClauses(controlID uint, currentUserID uint, ControlClauseFormDTO *organizationDTOs.ControlClauseFormDTO) (uint, error) {

	tx := cais.repo.BeginTransaction()

	if _, err := cais.controlSvc.GetByID(controlID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	clauses, err := cais.repo.GetClausesByIDs(ControlClauseFormDTO.ClauseIDs, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(clauses) != len(ControlClauseFormDTO.ClauseIDs) {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("one or more clause_ids are invalid"), "not_exist", "clause_id")
	}

	existingControlClauses, err := cais.repo.GetClausesByControlID(controlID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
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
			return 0, err
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
			return 0, err
		}
	}

	tx.Commit()

	return uint(len(clausesToAdd)), nil
}
