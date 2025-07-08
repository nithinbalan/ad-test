package services

import (
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	"time"

	// complianceDTOs "jethur/grc/compliance/dtos"
	// complianceRepo "jethur/grc/compliance/repositories"
	"fmt"
	"jethur/utility"
	// "github.com/mitchellh/mapstructure"
)

type ControlAssociatedItemServiceInterface interface {
	GetAllProcesses(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemProcessListResponseDTO, error)
	GetAllComplianceRequirements(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) ([]organizationDTOs.ComplianceRegisterDTO, error)
	GetAllRisks(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemRiskListResponseDTO, error)
	GetAllIncidents(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemIncidentListResponseDTO, error)
	GetAllAssociatedCounts(controlID uint) (*organizationDTOs.ControlAssociatedItemsCountResponseDTO, error)
	GetAllAssociatedClauses(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemClauseListResponseDTO, error)
	UpdateAssociatedClause(controlID uint, currentUserID uint, formDTO *organizationDTOs.ControlClauseFormDTO) (uint, error)
}

type ControlAssociatedItemService struct {
	repo                           organizationRepositories.ControlAssociatedItemRepositoryInterface
	controlSvc                     ControlServiceInterface
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface
	controlEfficiencyStatusService ControlEfficiencyStatusServiceInterface
	// complianceRepo  complianceRepo.ComplianceClauseRepositoryInterface
}

func NewControlAssociatedItemService(
	repo organizationRepositories.ControlAssociatedItemRepositoryInterface,
	controlSvc ControlServiceInterface,
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusService ControlEfficiencyStatusServiceInterface,
	// complianceRepo complianceRepo.ComplianceClauseRepositoryInterface,
) ControlAssociatedItemServiceInterface {
	return &ControlAssociatedItemService{
		repo:                           repo,
		controlSvc:                     controlSvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusService: controlEfficiencyStatusService,
		// complianceRepo:                 complianceRepo,
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

func (cais *ControlAssociatedItemService) GetAllAssociatedClauses(controlID uint, params *organizationDTOs.ControlAssociatedItemListQueryParamsDTO) (*organizationDTOs.ControlAssociatedItemClauseListResponseDTO, error) {
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}
	response, err := cais.repo.GetAllAssociatedClauses(controlID, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (cais *ControlAssociatedItemService) UpdateAssociatedClause(controlID uint, currentUserID uint, ControlClauseFormDTO *organizationDTOs.ControlClauseFormDTO) (uint, error) {
	tx := cais.repo.BeginTransaction()
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		tx.Rollback()
		return 0, err
	}
	controlClauses, err := cais.buildControlClausesFromDTO(controlID, currentUserID, ControlClauseFormDTO)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	insertedCount, err := cais.repo.BulkInsertClauses(controlClauses, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return insertedCount, nil
}

func (cais *ControlAssociatedItemService) buildControlClausesFromDTO(controlID uint, currentUserID uint, ControlClauseFormDTO *organizationDTOs.ControlClauseFormDTO) ([]*models.ControlClause, error) {
	if len(ControlClauseFormDTO.ClauseIDs) == 0 {
		return nil, utility.NewValidationError(
			fmt.Errorf("clause_ids must not be empty"),
			"required",
			"clause_ids",
		)
	}
	var controlClauses []*models.ControlClause
	now := time.Now()
	//  Fix: Fetch all clauses in a single query:TODO
	for _, clauseID := range ControlClauseFormDTO.ClauseIDs {
		clause, err := cais.repo.GetClauseByID(clauseID)
		if err != nil {
			return nil, utility.NewValidationError(
				fmt.Errorf("clause ID %d does not exist", clauseID),
				"not_exist",
				"clause_id",
			)
		}
		controlClauses = append(controlClauses, &models.ControlClause{
			ControlID:            controlID,
			ComplianceRegisterID: clause.ComplianceRegisterID,
			ComplianceClauseID:   clause.ID,
			CreatedByID:          currentUserID,
			CreatedAt:            now,
		})
	}
	return controlClauses, nil
}
