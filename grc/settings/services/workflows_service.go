package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"log"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type WorkflowServiceInterface interface {
	GetAll(params *settingsDTOs.WorkflowListQueryParamsDTO) ([]*settingsDTOs.WorkflowListDTO, error)
	GetFullDetailsByID(id uint) (*settingsDTOs.WorkflowDetailDTO, error)
	GetByID(id uint) (*models.Workflow, error)
	Create(currentUserID uint, escalationRuleFormDTO *settingsDTOs.WorkflowFormDTO) (uint, error)
	Update(id uint, currentUserID uint, escalationRuleFormDTO *settingsDTOs.WorkflowFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	GetAllWorkflowsByType(workflowType string, moduleIdentifier string) ([]*settingsDTOs.WorkflowListDTO, error)
	GetWorkflowByDocumentType(documentTypeId uint, tx ...utility.DatabaseTransaction) (*models.Workflow, error)
	GetWorkflowLevels(workflowID uint, tx ...utility.DatabaseTransaction) (*settingsDTOs.WorkflowLevelListDTO, error)
	CreateWorkFlowLevel(workflowId uint, currentUserID uint, workflowFormDTO *settingsDTOs.WorkflowLevelFormDTO) (uint, error)
	DeleteWorkflowLevel(currentUserID, workflowLevelID uint) error
}

type WorkflowService struct {
	repo                 repositories.WorkflowRepositoryInterface
	activityLogSvc       services.ActivityLogServiceInterface
	workflowLevelTypeSvc WorkflowLevelTypeServiceInterface
}

func NewWorkflowService(
	repo repositories.WorkflowRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	workflowLevelTypeSvc WorkflowLevelTypeServiceInterface) WorkflowServiceInterface {
	return &WorkflowService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		workflowLevelTypeSvc: workflowLevelTypeSvc,
	}
}

func (ws *WorkflowService) GetAll(params *settingsDTOs.WorkflowListQueryParamsDTO) ([]*settingsDTOs.WorkflowListDTO, error) {
	return ws.repo.GetAll(params)
}

func (ws *WorkflowService) GetByID(id uint) (*models.Workflow, error) {
	escalationRule, err := ws.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return escalationRule, nil
}

func (ws *WorkflowService) GetFullDetailsByID(id uint) (*settingsDTOs.WorkflowDetailDTO, error) {
	var escalationRuleDetailDTO settingsDTOs.WorkflowDetailDTO

	escalationRule, err := ws.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(escalationRule, &escalationRuleDetailDTO)
	if err != nil {
		return nil, err
	}

	return &escalationRuleDetailDTO, nil
}

func (ws *WorkflowService) Create(currentUserID uint, dto *settingsDTOs.WorkflowFormDTO) (uint, error) {

	level := &models.Workflow{
		WorkflowContextID: dto.WorkflowContextID,
		TitleEn:           &dto.TitleEn,
		TitleAr:           &dto.TitleAr,
		Status:            dto.Status,
		CreatedBYID:       currentUserID,
	}

	id, err := ws.repo.Create(level)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ws *WorkflowService) Update(id uint, currentUserID uint, dto *settingsDTOs.WorkflowFormDTO) error {

	oldData, err := ws.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	existing, err := ws.GetByID(id)
	if err != nil {
		return err
	}

	existing.TitleEn = &dto.TitleEn
	existing.TitleAr = &dto.TitleAr
	existing.UpdatedBYID = &currentUserID

	err = ws.repo.Update(existing)
	if err != nil {
		return err
	}

	go func() {
		// Fetch updated full details after update
		newData, err := ws.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		oldFields := map[string]interface{}{
			"title_en": oldData.TitleEn,
			"title_ar": oldData.TitleAr,
		}

		newFields := map[string]interface{}{
			"title_en": newData.TitleEn,
			"title_ar": newData.TitleAr,
		}

		changedData := utility.CompareFieldChanges(oldFields, newFields)

		if len(changedData) > 0 {
			err := ws.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.Workflow{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log escalation level update", map[string]interface{}{
					"action":    models.UpdateAction,
					"table":     models.Workflow{}.TableName(),
					"record_id": id,
					"user_id":   currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (ws *WorkflowService) Delete(ids []uint, currentUserID uint) ([]uint, error) {

	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	for _, id := range ids {
		item, err := ws.repo.GetFullDetailsByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		if err := ws.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
			deletedItems = append(deletedItems, map[string]interface{}{
				"id":       id,
				"title_en": item.TitleEn,
				"title_ar": item.TitleAr,
			})
		}
	}

	if len(successfulIds) > 0 {
		go ws.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.Workflow{}.TableName(),
			successfulIds,
			deletedItems,
			currentUserID,
			nil,
		)
	}

	return failedIds, nil
}

func (ws *WorkflowService) GetAllWorkflowsByType(workflowType string, moduleIdentifier string) ([]*settingsDTOs.WorkflowListDTO, error) {
	return ws.repo.GetAllWorkflowsByType(workflowType, moduleIdentifier)
}

func (ws *WorkflowService) GetWorkflowByDocumentType(documentTypeId uint, tx ...utility.DatabaseTransaction) (*models.Workflow, error) {
	workflow, err := ws.repo.GetWorkflowByDocumentType(documentTypeId, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, documentTypeId)
		}
		return nil, err
	}
	return workflow, nil
}

func (ws *WorkflowService) GetWorkflowLevels(workflowID uint, tx ...utility.DatabaseTransaction) (*settingsDTOs.WorkflowLevelListDTO, error) {
	var workflowLevelDTO settingsDTOs.WorkflowLevelListDTO

	workflowLevels, err := ws.repo.GetWorkflowLevels(workflowID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, workflowID)
		}
		return nil, err
	}

	err = mapstructure.Decode(workflowLevels, &workflowLevelDTO)
	if err != nil {
		return nil, err
	}

	return &workflowLevelDTO, nil
}

func (ws *WorkflowService) CreateWorkFlowLevel(workflowId uint, currentUserID uint, workflowLevelFormDTO *settingsDTOs.WorkflowLevelFormDTO) (uint, error) {
	tx := ws.repo.BeginTransaction()

	if err := ws.validateWorkflowFormDTO(workflowLevelFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	maxLevel, err := ws.GetMaxLevelForWorkflow(workflowId, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var newLevel int
	if workflowLevelFormDTO.ParentLevel == nil {
		if maxLevel == 0 {
			newLevel = 1 // Set as the first level
		} else {
			tx.Rollback()
			return 0, fmt.Errorf("parent level must be specified when other levels exist")
		}
	} else {
		newLevel = int(*workflowLevelFormDTO.ParentLevel) + 1
	}

	// Adjust subsequent levels if the new level is not the next immediate level and is less than or equal to max level
	if newLevel <= maxLevel {
		err := ws.repo.IncrementSubsequentLevels(workflowId, uint(newLevel), tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	workflowLevel := &models.WorkflowLevel{
		Level:               newLevel,
		WorkflowID:          workflowId,
		WorkflowLevelTypeID: workflowLevelFormDTO.LevelTypeID,
		ApprovalType:        *workflowLevelFormDTO.ApprovalType,
		CreatedByID:         currentUserID,
	}

	id, err := ws.repo.CreateWorkflowLevel(workflowLevel, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	workflowRoleType, err := ws.workflowLevelTypeSvc.GetByID(workflowLevelFormDTO.LevelTypeID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	switch workflowRoleType.Type {
	case "individual", "multiple-users":
		// Handle individual and multiple users by looping through user IDs provided in the form DTO
		fmt.Println(workflowLevelFormDTO.UserIds)
		if len(workflowLevelFormDTO.UserIds) < 1 {
			tx.Rollback()
			return 0, fmt.Errorf("user is required for selected level type")
		}
		for _, userID := range workflowLevelFormDTO.UserIds {
			participant := &models.WorkflowLevelParticipant{
				WorkflowLevelID: id,
				UserID:          userID,
			}
			if err := ws.repo.CreateParticipant(participant, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	case "role":
		// Handle role by creating a participant with RoleID
		if workflowLevelFormDTO.RoleID == nil {
			tx.Rollback()
			return 0, fmt.Errorf("role is required for selected level type")
		}
		participant := &models.WorkflowLevelParticipant{
			WorkflowLevelID: id,
			RoleID:          workflowLevelFormDTO.RoleID,
		}
		if err := ws.repo.CreateParticipant(participant, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()
	go func() {
		workflowLevelDetail, err := ws.repo.GetWorkflowLevelByID(id)
		if err != nil {
			return
		}

		err = ws.activityLogSvc.LogActivity(
			models.CreateAction,
			models.WorkflowLevel{}.TableName(),
			id,
			workflowLevelDetail,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.WorkflowLevel{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (ws *WorkflowService) DeleteWorkflowLevel(currentUserID, workflowLevelID uint) error {
	tx := ws.repo.BeginTransaction()

	var deletedWorkflowLevel []map[string]interface{}

	// Retrieve the workflow level to get its WorkflowID and level
	workflowLevel, err := ws.repo.GetWorkflowLevelByID(workflowLevelID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := ws.repo.DeleteWorkflowLevelParticipants(workflowLevelID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := ws.repo.DeleteWorkflowLevel(workflowLevelID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Check if the deleted level is the last one
	maxLevel, err := ws.GetMaxLevelForWorkflow(workflowLevel.WorkflowID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// If the deleted level is not the last one, adjust levels of subsequent levels
	if workflowLevel.Level < maxLevel {
		if err := ws.repo.DecrementSubsequentLevels(workflowLevel.WorkflowID, workflowLevel.Level, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	deletedWorkflowLevel = append(deletedWorkflowLevel, map[string]interface{}{
		"id":          workflowLevelID,
		"workflow_id": workflowLevel.WorkflowID,
		"level":       workflowLevel.Level,
	})

	// Commit transaction before logging
	if err := tx.Commit(); err != nil {
		return err
	}

	// Log activity after successful deletion
	go func() {
		err := ws.activityLogSvc.LogActivity(
			models.DeleteAction, // Make sure this action constant is defined appropriately
			models.WorkflowLevel{}.TableName(),
			workflowLevelID,
			deletedWorkflowLevel, // You can define a struct or map that describes the change
			currentUserID,
			nil,
		)
		if err != nil {
			// Optionally handle logging error, e.g., log to a file or another system
			log.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return nil
}

func (ws *WorkflowService) validateWorkflowFormDTO(workflowFormDTO *settingsDTOs.WorkflowLevelFormDTO) error {
	if _, err := ws.workflowLevelTypeSvc.GetByID(workflowFormDTO.LevelTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("workflow role type does not exist"), "not_exist", "workflow role type")
	}
	return nil
}

// GetMaxLevelForWorkflow ensures a proper level is set for new workflow entries.
func (ws *WorkflowService) GetMaxLevelForWorkflow(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	maxLevel, err := ws.repo.GetMaxLevelByWorkflowID(workflowID, tx...)
	if err != nil {
		return 0, err
	}

	return maxLevel, nil
}
