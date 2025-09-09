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

type ProcessAssociatedItemServiceInterface interface {
	GetAllControls(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) ([]map[string]interface{}, error)
	GetAllControlsWithPagination(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	GetAllRisks(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetAllIncidents(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error)
	GetAllDocuments(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	GetAllAssociatedCounts(processID uint) (*organizationDTOs.ProcessAssociatedItemsCountResponseDTO, error)
	GetAllAssets(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemAssetListResponseDTO, error)
	GetAllClauses(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	UpdateClauses(processID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error
	CreateProcessRisk(processID uint, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteProcessRisk(processID uint, riskID uint) error
	CreateProcessDocument(processID uint, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteProcessDocument(processID uint, riskID uint) error
	CreateProcessControl(processID uint, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteProcessControl(processID uint, controlID uint) error
	GetAllProcessKPI(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error)
	CreateProcessKPI(processID uint, userID uint, dto *commonDTOs.AssociatedKPIMapDTO) error
	DeleteProcessKPI(processID uint, kpiID uint) error
	CreateProcessAsset(processID uint, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error
	DeleteProcessAsset(processID uint, assetID uint) error
	CreateProcessIncident(processID uint, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error
	DeleteProcessIncident(processID uint, incidentID uint) error
	GetAllProcessStrategicObjectives(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error)
	CreateProcessStrategicObjective(processID uint, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error
	DeleteProcessStrategicObjective(processID uint, strategicObjectiveID uint) error
}

type ProcessAssociatedItemService struct {
	repo                           organizationRepositories.ProcessAssociatedItemRepositoryInterface
	processSvc                     ProcessServiceInterface
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface
	controlEfficiencyStatusService ControlEfficiencyStatusServiceInterface
	validationSvc                  commonServices.AssociatedItemValidationServiceInterface
}

func NewProcessAssociatedItemService(
	repo organizationRepositories.ProcessAssociatedItemRepositoryInterface,
	processSvc ProcessServiceInterface,
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusService ControlEfficiencyStatusServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) ProcessAssociatedItemServiceInterface {
	return &ProcessAssociatedItemService{
		repo:                           repo,
		processSvc:                     processSvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusService: controlEfficiencyStatusService,
		validationSvc:                  validationSvc,
	}
}

func (pais *ProcessAssociatedItemService) GetAllControls(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) ([]map[string]interface{}, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}

	controlList, err := pais.repo.GetAllControls(processID, params)
	if err != nil {
		return nil, err
	}

	if len(controlList) == 0 {
		return []map[string]interface{}{}, nil
	}

	implementationStatuses, err := pais.controlImplementationStatusSvc.GetAll("", "", "", "")
	if err != nil {
		return nil, err
	}
	implementationStatusMap := make(map[string]models.ControlImplementationStatus)
	for _, status := range implementationStatuses {
		implementationStatusMap[status.Type] = *status
	}

	efficiencyStatuses, err := pais.controlEfficiencyStatusService.GetAll("", "", "", "")
	if err != nil {
		return nil, err
	}
	efficiencyStatusMap := make(map[string]models.ControlEfficiencyStatus)
	for _, status := range efficiencyStatuses {
		efficiencyStatusMap[status.Type] = *status
	}

	groupedControls := make(map[string]map[string][]*organizationDTOs.ProcessControlListDTO)
	for _, control := range controlList {
		implementationStatusType := control.ControlImplementationStatusType
		efficiencyStatusType := control.ControlEfficiencyStatusType

		if _, exists := groupedControls[implementationStatusType]; !exists {
			groupedControls[implementationStatusType] = make(map[string][]*organizationDTOs.ProcessControlListDTO)
		}
		if _, exists := groupedControls[implementationStatusType][efficiencyStatusType]; !exists {
			groupedControls[implementationStatusType][efficiencyStatusType] = []*organizationDTOs.ProcessControlListDTO{}
		}
		groupedControls[implementationStatusType][efficiencyStatusType] = append(groupedControls[implementationStatusType][efficiencyStatusType], control)
	}

	var response []map[string]interface{}
	for implStatusType, effMap := range groupedControls {
		implItem := map[string]interface{}{
			"status":   implStatusType,
			"title_en": "",
			"title_ar": "",
			"count":    0, // Count of implementation status items
			"items":    []map[string]interface{}{},
		}

		if status, exists := implementationStatusMap[implStatusType]; exists {
			implItem["title_en"] = status.TitleEn
			implItem["title_ar"] = status.TitleAr
		}

		for effStatusType, controls := range effMap {
			effItem := map[string]interface{}{
				"efficiency": effStatusType,
				"title_en":   "",
				"title_ar":   "",
				"count":      len(controls), // Count of efficiency status items
				"items":      []map[string]interface{}{},
			}

			if status, exists := efficiencyStatusMap[effStatusType]; exists {
				effItem["title_en"] = status.TitleEn
				effItem["title_ar"] = status.TitleAr
			}

			for _, control := range controls {
				controlItem := map[string]interface{}{
					"id":                                       control.ID,
					"process_id":                               control.ProcessID,
					"control_id":                               control.ControlID,
					"reference_number":                         control.ControlReferenceNumber,
					"code":                                     control.ControlCode,
					"last_test_date":                           control.ControlLastTestDate,
					"control_category_title_en":                control.ControlCategoryTitleEn,
					"control_category_title_ar":                control.ControlCategoryTitleAr,
					"control_type_title_en":                    control.ControlTypeTitleEn,
					"control_type_title_ar":                    control.ControlTypeTitleAr,
					"title_en":                                 control.ControlTitleEn,
					"title_ar":                                 control.ControlTitleAr,
					"description":                              control.ControlDescription,
					"control_efficiency_status_title_en":       control.ControlEfficiencyStatusTitleEn,
					"control_efficiency_status_title_ar":       control.ControlEfficiencyStatusTitleAr,
					"control_efficiency_status_type":           control.ControlEfficiencyStatusType,
					"control_efficiency_status_color_code":     control.ControlEfficiencyStatusColorCode,
					"control_implementation_status_title_en":   control.ControlImplementationStatusTitleEn,
					"control_implementation_status_title_ar":   control.ControlImplementationStatusTitleAr,
					"control_implementation_status_type":       control.ControlImplementationStatusType,
					"control_implementation_status_color_code": control.ControlImplementationStatusColorCode,
					"owner_first_name_en":                      control.OwnerFirstNameEn,
					"owner_first_name_ar":                      control.OwnerFirstNameAr,
					"owner_last_name_en":                       control.OwnerLastNameEn,
					"owner_last_name_ar":                       control.OwnerLastNameAr,
					"owner_initial":                            control.OwnerInitial,
					"owner_color":                              control.OwnerColor,
					"owner_email":                              control.OwnerEmail,
					"owner_profile_image_path":                 control.OwnerProfileImagePath,
					"objective":                                control.Objective,
					"test_responsible_user_first_name_en":      control.TestResponsibleUserFirstNameEn,
					"test_responsible_user_first_name_ar":      control.TestResponsibleUserFirstNameAr,
					"test_responsible_user_last_name_en":       control.TestResponsibleUserLastNameEn,
					"test_responsible_user_last_name_ar":       control.TestResponsibleUserLastNameAr,
					"test_responsible_user_initial":            control.TestResponsibleUserInitial,
					"test_responsible_user_color":              control.TestResponsibleUserColor,
					"test_responsible_user_email":              control.TestResponsibleUserEmail,
					"test_responsible_user_profile_image_path": control.TestResponsibleUserProfileImagePath,
					"test_responsible_user_id":                 control.TestResponsibleUserID,
					"implemented_at":                           control.ImplementedAt,
					"last_review_date":                         control.LastReviewDate,
					"next_review_date":                         control.NextReviewDate,
					"last_review_note":                         control.LastReviewNote,
					"test_frequency_title_en":                  control.TestFrequencyTitleEn,
					"test_frequency_title_ar":                  control.TestFrequencyTitleAr,
					"next_test_date":                           control.NextTestDate,
					"risk_rating_title_en":                     control.RiskRatingTitleEn,
					"risk_rating_title_ar":                     control.RiskRatingTitleAr,
					"risk_rating_color_code":                   control.RiskRatingColorCode,
					"risk_rating_value":                        control.RiskRatingValue,
					"lifecycle_status_title_en":                control.LifecycleStatusTitleEn,
					"lifecycle_status_title_ar":                control.LifecycleStatusTitleAr,
					"lifecycle_status_color_code":              control.LifecycleStatusColorCode,
					"review_frequency_title_en":                control.ReviewFrequencyTitleEn,
					"review_frequency_title_ar":                control.ReviewFrequencyTitleAr,
				}
				effItem["items"] = append(effItem["items"].([]map[string]interface{}), controlItem)
			}

			implItem["items"] = append(implItem["items"].([]map[string]interface{}), effItem)
			implItem["count"] = implItem["count"].(int) + effItem["count"].(int)
		}

		response = append(response, implItem)
	}

	return response, nil
}

func (pais *ProcessAssociatedItemService) GetAllControlsWithPagination(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllControlsWithPagination(processID, params)
}

func (pais *ProcessAssociatedItemService) GetAllRisks(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllRisks(processID, params)
}

func (pais *ProcessAssociatedItemService) GetAllIncidents(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllIncidents(processID, params)
}

func (pais *ProcessAssociatedItemService) GetAllDocuments(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllDocuments(processID, params)
}

func (pais *ProcessAssociatedItemService) GetAllAssets(processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemAssetListResponseDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllAssets(processID, params)
}

func (pais *ProcessAssociatedItemService) GetAllAssociatedCounts(processID uint) (*organizationDTOs.ProcessAssociatedItemsCountResponseDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	response, err := pais.repo.GetAllAssociatedCounts(processID)
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (pais *ProcessAssociatedItemService) GetAllClauses(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}

	response, err := pais.repo.GetAllClauses(processID, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (pais *ProcessAssociatedItemService) UpdateClauses(processID uint, currentUserID uint, ClauseFormDTO *commonDTOs.AssociatedClauseFormDTO) error {

	tx := pais.repo.BeginTransaction()

	if _, err := pais.processSvc.GetByID(processID, tx); err != nil {
		tx.Rollback()
		return err
	}

	clauses, err := pais.validateClauseIDs(ClauseFormDTO.ClauseIDs, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlClauses, err := pais.repo.GetClausesByProcessID(processID, tx)
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

	var clausesToAdd []*models.ProcessComplianceRequirement
	now := time.Now()
	for id, clause := range newClauses {
		if !existingClauses[id] {
			clausesToAdd = append(clausesToAdd, &models.ProcessComplianceRequirement{
				ProcessID:            processID,
				ComplianceRegisterID: clause.ComplianceRegisterID,
				ComplianceClauseID:   clause.ID,
				CreatedByID:          currentUserID,
				CreatedAt:            now,
			})
		}
	}

	if len(clausesToAdd) > 0 {
		if err := pais.repo.CreateProcessClauses(clausesToAdd, tx); err != nil {
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
		if err := pais.repo.DeleteProcessClauses(processID, deleteClauseIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (pais *ProcessAssociatedItemService) validateClauseIDs(clauseIDs []uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceClause, error) {
	seen := make(map[uint]bool)
	for _, id := range clauseIDs {
		if seen[id] {
			return nil, utility.NewValidationError(fmt.Errorf("duplicate clause ID: %d", id), "duplicate", "clause_id")
		}
		seen[id] = true
	}

	clauses, err := pais.repo.GetClausesByIDs(clauseIDs, tx...)
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

func (pais *ProcessAssociatedItemService) CreateProcessRisk(processID, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", processID), "not_exist", "process_id")
	}
	if err := pais.validationSvc.ValidateRiskMapping(processID, dto.RiskIDs, "risk_processes", "process_id"); err != nil {
		return err
	}
	tx := pais.repo.BeginTransaction()
	var mappings []*models.RiskProcess
	currentTime := time.Now()

	for _, riskID := range dto.RiskIDs {
		mappings = append(mappings, &models.RiskProcess{
			RiskRegisterID: riskID,
			ProcessID:      processID,
			CreatedByID:    userID,
			CreatedAt:      currentTime,
		})
	}

	if err := pais.repo.CreateProcessRisk(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProcessAssociatedItemService) DeleteProcessRisk(processID uint, riskID uint) error {
	// Validate control existence
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", processID), "not_exist", "process_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProcessRisk(processID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *ProcessAssociatedItemService) CreateProcessDocument(processID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", processID), "not_exist", "process_id")
	}
	if err := pais.validationSvc.ValidateDocumentMapping(processID, dto.DocumentIDs, "process_associated_documents", "process_id"); err != nil {
		return err
	}
	tx := pais.repo.BeginTransaction()
	var mappings []*models.ProcessAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.ProcessAssociatedDocument{
			DocumentID:  documentID,
			ProcessID:   processID,
			CreatedByID: userID,
			CreatedAt:   currentTime,
		})
	}

	if err := pais.repo.CreateProcessDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProcessAssociatedItemService) DeleteProcessDocument(processID uint, documentID uint) error {
	// Validate control existence
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", processID), "not_exist", "process_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProcessDocument(processID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *ProcessAssociatedItemService) CreateProcessControl(processID, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", processID), "not_exist", "process_id")
	}
	if err := pais.validationSvc.ValidateControlMapping(processID, dto.ControlIDs, "process_controls", "process_id"); err != nil {
		return err
	}
	tx := pais.repo.BeginTransaction()
	var mappings []*models.ProcessControl
	currentTime := time.Now()

	for _, controlID := range dto.ControlIDs {
		mappings = append(mappings, &models.ProcessControl{
			ControlID:   controlID,
			ProcessID:   processID,
			CreatedByID: userID,
			CreatedAt:   currentTime,
		})
	}

	if err := pais.repo.CreateProcessControl(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProcessAssociatedItemService) DeleteProcessControl(processID uint, controlID uint) error {
	// Validate control existence
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", processID), "not_exist", "process_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProcessControl(processID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *ProcessAssociatedItemService) GetAllProcessKPI(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}

	response, err := pais.repo.GetAllProcessKPI(processID, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (pais *ProcessAssociatedItemService) CreateProcessKPI(processID, userID uint, dto *commonDTOs.AssociatedKPIMapDTO) error {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return err
	}
	if err := pais.validationSvc.ValidateKPIMapping(processID, dto.KpiIDs, "kpi_processes", "process_id"); err != nil {
		return err
	}
	tx := pais.repo.BeginTransaction()
	var mappings []*models.KPIProcess
	currentTime := time.Now()

	for _, KpiID := range dto.KpiIDs {
		mappings = append(mappings, &models.KPIProcess{
			KPIRegisterID: KpiID,
			ProcessID:     processID,
			CreatedByID:   userID,
			CreatedAt:     currentTime,
		})
	}

	if err := pais.repo.CreateProcessKPI(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProcessAssociatedItemService) DeleteProcessKPI(processID uint, kpiID uint) error {
	// Validate control existence
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProcessKPI(processID, kpiID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *ProcessAssociatedItemService) CreateProcessAsset(processID, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return err
	}
	if err := pais.validationSvc.ValidateAssetMapping(processID, dto.AssetIDs, "process_assets", "process_id"); err != nil {
		return err
	}
	tx := pais.repo.BeginTransaction()
	var mappings []*models.ProcessAsset
	currentTime := time.Now()

	for _, AssetID := range dto.AssetIDs {
		mappings = append(mappings, &models.ProcessAsset{
			AssetRegisterID: AssetID,
			ProcessID:       processID,
			CreatedByID:     userID,
			CreatedAt:       currentTime,
		})
	}

	if err := pais.repo.CreateProcessAsset(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProcessAssociatedItemService) DeleteProcessAsset(processID uint, assetID uint) error {
	// Validate control existence
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProcessAsset(processID, assetID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *ProcessAssociatedItemService) CreateProcessIncident(processID, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return err
	}
	if err := pais.validationSvc.ValidateIncidentMapping(processID, dto.IncidentIDs, "incident_processes", "process_id"); err != nil {
		return err
	}
	tx := pais.repo.BeginTransaction()
	var mappings []*models.IncidentProcess
	currentTime := time.Now()

	for _, IncidentID := range dto.IncidentIDs {
		mappings = append(mappings, &models.IncidentProcess{
			IncidentRegisterID: IncidentID,
			ProcessID:          processID,
			CreatedByID:        userID,
			CreatedAt:          currentTime,
		})
	}

	if err := pais.repo.CreateProcessIncident(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProcessAssociatedItemService) DeleteProcessIncident(processID uint, incidentID uint) error {
	// Validate control existence
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProcessIncident(processID, incidentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *ProcessAssociatedItemService) GetAllProcessStrategicObjectives(processID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error) {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllProcessStrategicObjectives(processID, params)
}

func (pais *ProcessAssociatedItemService) CreateProcessStrategicObjective(processID, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error {
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return err
	}
	if err := pais.validationSvc.ValidateStrtegicObjectiveMapping(processID, dto.StrategicObjectiveIDs, "strategic_objective_processes", "process_id"); err != nil {
		return err
	}
	tx := pais.repo.BeginTransaction()
	var mappings []*models.StrategicObjectiveProcess
	currentTime := time.Now()

	for _, strategicObjectiveID := range dto.StrategicObjectiveIDs {
		mappings = append(mappings, &models.StrategicObjectiveProcess{
			StrategicObjectiveID: strategicObjectiveID,
			ProcessID:            processID,
			CreatedByID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := pais.repo.CreateProcessStrategicObjective(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProcessAssociatedItemService) DeleteProcessStrategicObjective(processID uint, strategicObjectiveID uint) error {
	// Validate control existence
	if _, err := pais.processSvc.GetByID(processID); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProcessStrategicObjective(processID, strategicObjectiveID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
