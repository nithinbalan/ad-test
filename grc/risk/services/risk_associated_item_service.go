package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceServices "jethur/grc/compliance/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"time"
)

type RiskAssociatedItemServiceInterface interface {
	GetAllControls(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemControlListResponseDTO, error)
	GetAllComplianceRequirements(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemComplianceRequirementListResponseDTO, error)
	GetAllProcesses(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemProcessListResponseDTO, error)
	GetAllStrategicObjectives(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemStrategicObjectiveListResponseDTO, error)
	GetAllSLAs(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemSLAListResponseDTO, error)
	GetAllIncidents(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemIncidentListResponseDTO, error)
	GetAllDocuments(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemDocumentListResponseDTO, error)
	GetAllAssets(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemAssetListResponseDTO, error)
	GetAllThreats(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemThreatListResponseDTO, error)
	GetAllClauses(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	UpdateClauses(riskID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error
	CreatRiskStrategicObjective(riskID uint, currentUserID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error
	DeletetRiskStrategicObjective(riskID uint, strategicObjectiveID uint) error
	CreatRiskKPI(riskID uint, userID uint, dto *commonDTOs.AssociatedKPIMapDTO) error
	DeletetRiskKPI(riskID uint, kpiID uint) error
}

// Service and implementations
type RiskAssociatedItemService struct {
	repo            riskRepositories.RiskAssociatedItemRepositoryInterface
	riskRegisterSvc RiskRegisterServiceInterface
	clauseSvc       complianceServices.ComplianceClauseServiceInterface
	validationSvc   commonServices.AssociatedItemValidationServiceInterface
}

func NewRiskAssociatedItemService(
	repo riskRepositories.RiskAssociatedItemRepositoryInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	clauseSvc complianceServices.ComplianceClauseServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,

) RiskAssociatedItemServiceInterface {
	return &RiskAssociatedItemService{
		repo:            repo,
		riskRegisterSvc: riskRegisterSvc,
		clauseSvc:       clauseSvc,
		validationSvc:   validationSvc,
	}
}

func (rais *RiskAssociatedItemService) GetAllControls(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemControlListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllControls(riskID, params)
}

func (rais *RiskAssociatedItemService) GetAllComplianceRequirements(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemComplianceRequirementListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	complianceRequirementList, totalCount, err := rais.repo.GetAllComplianceRequirements(riskID, params)
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

	response := &riskDTOs.RiskAssociatedItemComplianceRequirementListResponseDTO{
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

func (rais *RiskAssociatedItemService) GetAllProcesses(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemProcessListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllProcesses(riskID, params)
}

func (rais *RiskAssociatedItemService) GetAllStrategicObjectives(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemStrategicObjectiveListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllStrategicObjectives(riskID, params)
}

func (rais *RiskAssociatedItemService) GetAllSLAs(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemSLAListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllSLAs(riskID, params)
}

func (rais *RiskAssociatedItemService) GetAllIncidents(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemIncidentListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllIncidents(riskID, params)
}

func (rais *RiskAssociatedItemService) GetAllDocuments(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemDocumentListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllDocuments(riskID, params)
}

func (rais *RiskAssociatedItemService) GetAllAssets(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemAssetListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllAssets(riskID, params)
}

func (rais *RiskAssociatedItemService) GetAllThreats(riskID uint, params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) (*riskDTOs.RiskAssociatedItemThreatListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rais.repo.GetAllThreats(riskID, params)
}

func (rais *RiskAssociatedItemService) GetAllClauses(riskID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	response, err := rais.repo.GetAllClauses(riskID, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (rais *RiskAssociatedItemService) UpdateClauses(riskID uint, currentUserID uint, RiskClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error {

	tx := rais.repo.BeginTransaction()

	if _, err := rais.riskRegisterSvc.GetByID(riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	clauses, err := rais.validateClauseIDs(RiskClauseFormDTO.ClauseIDs, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlClauses, err := rais.repo.GetClausesByRiskID(riskID, tx)
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

	var clausesToAdd []*models.RiskComplianceRequirement
	now := time.Now()
	for id, clause := range newClauses {
		if !existingClauses[id] {
			clausesToAdd = append(clausesToAdd, &models.RiskComplianceRequirement{
				RiskRegisterID:       riskID,
				ComplianceRegisterID: clause.ComplianceRegisterID,
				ComplianceClauseID:   clause.ID,
				CreatedByID:          currentUserID,
				CreatedAt:            now,
			})
		}
	}

	if len(clausesToAdd) > 0 {
		if err := rais.repo.CreateRiskClauses(clausesToAdd, tx); err != nil {
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
		if err := rais.repo.DeleteRiskClauses(riskID, deleteClauseIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (rais *RiskAssociatedItemService) validateClauseIDs(clauseIDs []uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceClause, error) {
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

func (rais *RiskAssociatedItemService) CreatRiskStrategicObjective(riskID, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk does not exist for ID: %d", riskID), "not_exist", "risk_id")
	}
	if err := rais.validationSvc.ValidateRiskMapping(riskID, dto.StrategicObjectiveIDs, "risk_strategic_objectives", "strategic_objective_id"); err != nil {
		return err
	}
	tx := rais.repo.BeginTransaction()
	var mappings []*models.RiskStrategicObjective
	currentTime := time.Now()

	for _, strategicObID := range dto.StrategicObjectiveIDs {
		mappings = append(mappings, &models.RiskStrategicObjective{
			RiskRegisterID:       riskID,
			StrategicObjectiveID: strategicObID,
			CreatedByID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := rais.repo.CreatRiskStrategicObjective(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (rais *RiskAssociatedItemService) DeletetRiskStrategicObjective(strategicObjectiveID uint, riskID uint) error {
	// Validate control existence
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk does not exist for ID: %d", riskID), "not_exist", "risk_id")
	}

	tx := rais.repo.BeginTransaction()

	if err := rais.repo.DeletetRiskStrategicObjective(strategicObjectiveID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (rais *RiskAssociatedItemService) CreatRiskKPI(riskID, userID uint, dto *commonDTOs.AssociatedKPIMapDTO) error {
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk does not exist for ID: %d", riskID), "not_exist", "risk_id")
	}
	if err := rais.validationSvc.ValidateRiskMapping(riskID, dto.KpiIDs, "kpi_risks", "kpi_register_id"); err != nil {
		return err
	}
	tx := rais.repo.BeginTransaction()
	var mappings []*models.KPIRisk
	currentTime := time.Now()

	for _, kpiID := range dto.KpiIDs {
		mappings = append(mappings, &models.KPIRisk{
			RiskRegisterID: riskID,
			KPIRegisterID:  kpiID,
			CreatedByID:    userID,
			CreatedAt:      currentTime,
		})
	}

	if err := rais.repo.CreatRiskKPI(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (rais *RiskAssociatedItemService) DeletetRiskKPI(kpiID uint, riskID uint) error {
	// Validate control existence
	if _, err := rais.riskRegisterSvc.GetByID(riskID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk does not exist for ID: %d", riskID), "not_exist", "risk_id")
	}

	tx := rais.repo.BeginTransaction()

	if err := rais.repo.DeletetRiskKPI(kpiID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
