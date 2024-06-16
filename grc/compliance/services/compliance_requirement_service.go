package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceRequirementServiceInterface interface {
	GetAll(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error)
	GetByID(complianceID uint, id uint) (*models.ComplianceRequirement, error)
	GetRequirementsByComplianceID(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO, error)
	GetByComplianceRequirementID(id uint) (*models.ComplianceRequirement, error)
	GetFullDetailsByID(complianceID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceRequirementDetailDTO, error)
	Create(complianceID uint, currentUserID uint, ComplianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) (uint, error)
	Update(complianceID uint, id uint, currentUserID uint, compliance *complianceDTOs.ComplianceRequirementFormDTO) error
	Delete(complianceID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(complianceID uint, id uint, currentUserID uint, complianceStatusID uint, tx utility.DatabaseTransaction) error
}

// Service and implementations
type ComplianceRequirementService struct {
	repo                             complianceRepositories.ComplianceRequirementRepositoryInterface
	activityLogSvc                   logServices.ActivityLogServiceInterface
	complianceRegisterSvc            ComplianceRegisterServiceInterface
	complianceRequirementCategorySvc ComplianceRequirementCategoryServiceInterface
	complianceStatusSvc              ComplianceStatusServiceInterface
	severityLevelSvc                 masterServices.SeverityLevelServiceInterface
}

func NewComplianceRequirementService(
	repo complianceRepositories.ComplianceRequirementRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	complianceRequirementCategorySvc ComplianceRequirementCategoryServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
) ComplianceRequirementServiceInterface {
	return &ComplianceRequirementService{
		repo:                             repo,
		activityLogSvc:                   activityLogSvc,
		complianceRegisterSvc:            complianceRegisterSvc,
		complianceRequirementCategorySvc: complianceRequirementCategorySvc,
		complianceStatusSvc:              complianceStatusSvc,
		severityLevelSvc:                 severityLevelSvc,
	}
}

func (crs *ComplianceRequirementService) GetAll(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error) {
	return crs.repo.GetAll(params)
}

func (crs *ComplianceRequirementService) GetRequirementsByComplianceID(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO, error) {
	return crs.repo.GetRequirementsByComplianceID(complianceID, params)
}

func (crs ComplianceRequirementService) GetByID(complianceID uint, id uint) (*models.ComplianceRequirement, error) {
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	complianceRequirement, err := crs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceRequirement, nil
}

func (crs ComplianceRequirementService) GetByComplianceRequirementID(id uint) (*models.ComplianceRequirement, error) {

	complianceRequirement, err := crs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceRequirement, nil
}

func (crs ComplianceRequirementService) GetFullDetailsByID(complianceID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceRequirementDetailDTO, error) {
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	var complianceRequirementDTO complianceDTOs.ComplianceRequirementDetailDTO

	complianceRequirement, err := crs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to DTO
	err = mapstructure.Decode(complianceRequirement, &complianceRequirementDTO)
	if err != nil {
		return nil, err
	}

	if complianceRequirement.Deadline != nil {
		deadline := complianceRequirement.Deadline.Format("2006-01-02")
		complianceRequirementDTO.DeadlineString = &deadline
	}

	// Convert compliance requirement status update to DTO
	var complianceRequirementStatusUpdateDTO complianceDTOs.ComplianceRequirementStatusUpdateDTO
	if len(complianceRequirement.ComplianceRequirementStatusUpdates) > 0 {
		err = mapstructure.Decode(complianceRequirement.ComplianceRequirementStatusUpdates[0], &complianceRequirementStatusUpdateDTO)
		if err != nil {
			return nil, err
		}
		complianceRequirementDTO.ComplianceRequirementStatusUpdate = &complianceRequirementStatusUpdateDTO
	}

	return &complianceRequirementDTO, nil
}

func (crs *ComplianceRequirementService) Create(complianceID uint, currentUserID uint, complianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) (uint, error) {
	tx := crs.repo.BeginTransaction()

	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return 0, err
	}

	if err := crs.validateComplianceRequirementFormDTO(complianceRequirementFormDTO); err != nil {
		return 0, err
	}

	dates, err := crs.parseComplianceRequirementDates(complianceRequirementFormDTO)
	if err != nil {
		return 0, err
	}

	// Create or retrieve the compliance clause
	complianceClause := &models.ComplianceClause{
		ComplianceRegisterID: complianceID,
		ClauseNo:             complianceRequirementFormDTO.ComplianceClauseNo,
		Title:                complianceRequirementFormDTO.ComplianceClauseTitle,
	}

	clause, err := crs.repo.GetOrCreateComplianceClause(complianceClause, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	complianceRequirement := &models.ComplianceRequirement{
		ComplianceRequirementCategoryID: complianceRequirementFormDTO.ComplianceRequirementCategoryID,
		ComplianceStatusID:              complianceRequirementFormDTO.ComplianceStatusID,
		ResponsibleUserID:               complianceRequirementFormDTO.ResponsibleUserID,
		SeverityLevelID:                 complianceRequirementFormDTO.SeverityLevelID,
		Title:                           complianceRequirementFormDTO.Title,
		Deadline:                        dates.Deadline,
		CreatedByID:                     currentUserID,
		ComplianceClauseID:              clause.ID,
	}

	id, err := crs.repo.Create(complianceRequirement, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	complianceStatusID := complianceRequirementFormDTO.ComplianceStatusID
	complianceRequirementUpdate := &models.ComplianceRequirementStatusUpdate{
		ComplianceRequirementID: id,
		ComplianceStatusID:      complianceStatusID,
		Comment:                 "Initial compliance status",
		CreatedByID:             currentUserID,
	}

	_, err = crs.repo.CreateComplianceRequirementStatusUpdate(complianceRequirementUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Add Controls
	var controls []*models.ComplianceControl
	for _, controlID := range complianceRequirementFormDTO.ControlIDs {
		control := &models.ComplianceControl{
			ComplianceRegisterID:    complianceID,
			ComplianceRequirementID: id,
			ControlID:               *controlID,
			CreatedByID:             currentUserID,
		}
		controls = append(controls, control)
	}

	if len(controls) > 0 {
		if err := crs.repo.CreateComplianceRequirementControls(controls, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Add Strategic Objectives
	var strategicObjectives []*models.ComplianceRequirementStrategicObjective
	for _, strategicObjectiveID := range complianceRequirementFormDTO.StrategicObjectiveIDs {
		strategicObjective := &models.ComplianceRequirementStrategicObjective{
			ComplianceRequirementID: id,
			StrategicObjectiveID:    *strategicObjectiveID,
			CreatedByID:             currentUserID,
		}
		strategicObjectives = append(strategicObjectives, strategicObjective)
	}

	if len(strategicObjectives) > 0 {
		if err := crs.repo.CreateComplianceRequirementStrategicObjectives(strategicObjectives, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := crs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceRequirement{}.TableName(),
			id,
			complianceRequirementFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceRequirement{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (crs *ComplianceRequirementService) Update(complianceID uint, id uint, currentUserID uint, complianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) error {
	tx := crs.repo.BeginTransaction()

	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return err
	}

	existingComplianceRequirement, err := crs.GetByID(complianceID, id)

	if err != nil {
		return err
	}

	if err := crs.validateComplianceRequirementFormDTO(complianceRequirementFormDTO); err != nil {
		return err
	}

	dates, err := crs.parseComplianceRequirementDates(complianceRequirementFormDTO)
	if err != nil {
		return err
	}

	complianceClause := &models.ComplianceClause{
		ComplianceRegisterID: complianceID,
		ClauseNo:             complianceRequirementFormDTO.ComplianceClauseNo,
		Title:                complianceRequirementFormDTO.ComplianceClauseTitle,
	}

	clause, err := crs.repo.GetOrCreateComplianceClause(complianceClause)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update Compliance Requirement
	existingComplianceRequirement.ComplianceRequirementCategoryID = complianceRequirementFormDTO.ComplianceRequirementCategoryID
	existingComplianceRequirement.ComplianceStatusID = complianceRequirementFormDTO.ComplianceStatusID
	existingComplianceRequirement.ResponsibleUserID = complianceRequirementFormDTO.ResponsibleUserID
	existingComplianceRequirement.SeverityLevelID = complianceRequirementFormDTO.SeverityLevelID
	existingComplianceRequirement.Title = complianceRequirementFormDTO.Title
	existingComplianceRequirement.Deadline = dates.Deadline
	existingComplianceRequirement.UpdatedByID = &currentUserID
	existingComplianceRequirement.ComplianceClauseID = clause.ID

	err = crs.repo.Update(existingComplianceRequirement, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update Controls
	complianceRequirement, err := crs.GetFullDetailsByID(complianceID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlIDs := make(map[uint]bool)
	for _, ctrl := range complianceRequirement.Controls {
		existingControlIDs[ctrl.ControlID] = true
	}

	// Collect new control IDs
	newControlIDs := make(map[uint]bool)
	for _, ctrl := range complianceRequirementFormDTO.ControlIDs {
		newControlIDs[*ctrl] = true
	}

	// Add new controls
	var controlsToAdd []*models.ComplianceControl
	for _, newControlID := range complianceRequirementFormDTO.ControlIDs {
		if _, exists := existingControlIDs[*newControlID]; !exists {
			// Create a new compliance control record
			control := &models.ComplianceControl{
				ComplianceRegisterID:    complianceID,
				ComplianceRequirementID: id,
				ControlID:               *newControlID,
				CreatedByID:             currentUserID,
			}
			controlsToAdd = append(controlsToAdd, control)
		}
	}

	// Save the new controls
	if len(controlsToAdd) > 0 {
		if err := crs.repo.CreateComplianceRequirementControls(controlsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove controls not in the new list
	var deleteControlIDs []*uint

	for _, existingControl := range complianceRequirement.Controls {
		if _, exists := newControlIDs[existingControl.ControlID]; !exists {
			// Append control to delete list
			deleteControlIDs = append(deleteControlIDs, existingControl.ID)
		}
	}

	// Delete controls not in the new list
	if len(deleteControlIDs) > 0 {
		if err := crs.repo.DeleteComplianceRequirementControls(deleteControlIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update Strategic Objectives
	existingStrategicObjectiveIDs := make(map[uint]bool)
	for _, obj := range complianceRequirement.StrategicObjectives {
		existingStrategicObjectiveIDs[obj.StrategicObjectiveID] = true
	}

	// Collect new strategic objective IDs
	newStrategicObjectiveIDs := make(map[uint]bool)
	for _, obj := range complianceRequirementFormDTO.StrategicObjectiveIDs {
		newStrategicObjectiveIDs[*obj] = true
	}

	// Add new strategic objectives
	var strategicObjectivesToAdd []*models.ComplianceRequirementStrategicObjective
	for _, newStrategicObjectiveID := range complianceRequirementFormDTO.StrategicObjectiveIDs {
		if _, exists := existingStrategicObjectiveIDs[*newStrategicObjectiveID]; !exists {
			// Create a new compliance strategic objective record
			strategicObjective := &models.ComplianceRequirementStrategicObjective{
				ComplianceRequirementID: id,
				StrategicObjectiveID:    *newStrategicObjectiveID,
				CreatedByID:             currentUserID,
			}
			strategicObjectivesToAdd = append(strategicObjectivesToAdd, strategicObjective)
		}
	}

	// Save the new strategic objectives
	if len(strategicObjectivesToAdd) > 0 {
		if err := crs.repo.CreateComplianceRequirementStrategicObjectives(strategicObjectivesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove strategic objectives not in the new list
	var deleteStrategicObjectiveIDs []*uint

	for _, existingStrategicObjective := range complianceRequirement.StrategicObjectives {
		if _, exists := newStrategicObjectiveIDs[existingStrategicObjective.StrategicObjectiveID]; !exists {
			// Append strategic objective to delete list
			deleteStrategicObjectiveIDs = append(deleteStrategicObjectiveIDs, existingStrategicObjective.ID)
		}
	}

	// Delete strategic objectives not in the new list
	if len(deleteStrategicObjectiveIDs) > 0 {
		if err := crs.repo.DeleteComplianceRequirementStrategicObjectives(deleteStrategicObjectiveIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (crs *ComplianceRequirementService) Delete(complianceID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := crs.repo.BeginTransaction()
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint

	// Track clause IDs for deletion
	clauseIDsToDelete := make(map[uint]bool)

	for _, id := range ids {
		requirement, err := crs.GetFullDetailsByID(complianceID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Check if controls exist
		if len(requirement.Controls) > 0 {
			var deleteControlIds []*uint
			for _, ctrl := range requirement.Controls {
				deleteControlIds = append(deleteControlIds, ctrl.ID)
			}
			if err := crs.repo.DeleteComplianceRequirementControls(deleteControlIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Check if strategic objectives exist
		if len(requirement.StrategicObjectives) > 0 {
			var deleteStrategicObjectiveIds []*uint
			for _, obj := range requirement.StrategicObjectives {
				deleteStrategicObjectiveIds = append(deleteStrategicObjectiveIds, obj.ID)
			}
			if err := crs.repo.DeleteComplianceRequirementStrategicObjectives(deleteStrategicObjectiveIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if err := crs.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
			// Track the clause IDs
			clauseIDsToDelete[uint(requirement.ComplianceClause.ID)] = true
		}
	}

	// Delete clauses if all their requirements are deleted
	for clauseID := range clauseIDsToDelete {
		count, err := crs.repo.CountRequirementsByClauseID(clauseID, tx)
		if err != nil {
			tx.Rollback()
			fmt.Printf("Error counting requirements for clause %d: %v\n", clauseID, err)
		} else if count == 0 {
			if err := crs.repo.DeleteClause(clauseID, tx); err != nil {
				tx.Rollback()
				// Log error, but proceed with other deletions
				fmt.Printf("Error deleting clause %d: %v\n", clauseID, err)
			}
		}
	}

	// Log successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ComplianceRequirement{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// Log error, but proceed as it might be recoverable later
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	tx.Commit()

	return failedIds, nil
}

func (crs *ComplianceRequirementService) StatusUpdate(complianceID uint, id uint, currentUserID uint, complianceStatusID uint, tx utility.DatabaseTransaction) error {
	existingComplianceRequirement, err := crs.GetByID(complianceID, id)

	if err != nil {
		return err
	}

	existingComplianceRequirement.ComplianceStatusID = complianceStatusID
	existingComplianceRequirement.UpdatedByID = &currentUserID

	err = crs.repo.Update(existingComplianceRequirement, tx)

	if err != nil {
		return err
	}

	return nil
}

func (crs *ComplianceRequirementService) validateComplianceRequirementFormDTO(ComplianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) error {
	if _, err := crs.complianceRequirementCategorySvc.GetByID(ComplianceRequirementFormDTO.ComplianceRequirementCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance requirement category does not exist"), "not_exist", "compliance requirement category")
	}

	if _, err := crs.complianceStatusSvc.GetByID(ComplianceRequirementFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	if _, err := crs.severityLevelSvc.GetByID(ComplianceRequirementFormDTO.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	return nil
}

func (crs *ComplianceRequirementService) parseComplianceRequirementDates(dto *complianceDTOs.ComplianceRequirementFormDTO) (*complianceDTOs.ComplianceRequirementParsedDates, error) {
	parsed := &complianceDTOs.ComplianceRequirementParsedDates{}

	var err error
	if dto.Deadline != "" {
		parsed.Deadline, err = utility.ParseDate(dto.Deadline)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid deadline"), "invalid_date", "deadline")
		}
	}

	return parsed, nil
}
