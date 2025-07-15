package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
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
	GetAllClauses(riskID uint, params *commonDTOs.AssociatedItemClauseListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	UpdateClauses(riskID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error
}

// Service and implementations
type RiskAssociatedItemService struct {
	repo            riskRepositories.RiskAssociatedItemRepositoryInterface
	riskRegisterSvc RiskRegisterServiceInterface
	clauseSvc       complianceServices.ComplianceClauseServiceInterface
}

func NewRiskAssociatedItemService(
	repo riskRepositories.RiskAssociatedItemRepositoryInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	clauseSvc complianceServices.ComplianceClauseServiceInterface,

) RiskAssociatedItemServiceInterface {
	return &RiskAssociatedItemService{
		repo:            repo,
		riskRegisterSvc: riskRegisterSvc,
		clauseSvc:       clauseSvc,
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

func (rais *RiskAssociatedItemService) GetAllClauses(riskID uint, params *commonDTOs.AssociatedItemClauseListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
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
