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
	GetAllWorkflowsByType(workflowType string) ([]*settingsDTOs.WorkflowListDTO, error)
	GetWorkflowLevels(workflowID uint) (*settingsDTOs.WorkflowLevelListDTO, error)
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

func (ws *WorkflowService) GetAllWorkflowsByType(workflowType string) ([]*settingsDTOs.WorkflowListDTO, error) {
	return ws.repo.GetAllWorkflowsByType(workflowType)
}

func (ws *WorkflowService) GetWorkflowLevels(workflowID uint) (*settingsDTOs.WorkflowLevelListDTO, error) {
	var workflowLevelDTO settingsDTOs.WorkflowLevelListDTO

	workflowLevels, err := ws.repo.GetWorkflowLevels(workflowID)
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

	// levelsDTO := &settingsDTOs.WorkflowLevelListDTO{
	// 	WorkflowLevels: make([]settingsDTOs.WorkflowLevelDTO, len(workflowLevels)),
	// }

	// for i, level := range workflowLevels {
	// 	levelDTO := settingsDTOs.WorkflowLevelDTO{
	// 		ID:           level.ID,
	// 		Level:        level.Level,
	// 		WorkflowID:   level.WorkflowID,
	// 		ApprovalType: level.ApprovalType,
	// 		CreatedByID:  level.CreatedByID,
	// 		UpdatedByID:  level.UpdatedByID,
	// 		CreatedAt:    level.CreatedAt,
	// 		UpdatedAt:    level.UpdatedAt,
	// 		WorkflowLevelType: settingsDTOs.WorkflowLevelTypeDTO{
	// 			ID:      level.WorkflowLevelType.ID,
	// 			TitleEn: level.WorkflowLevelType.TitleEn,
	// 			TitleAr: level.WorkflowLevelType.TitleAr,
	// 			Type:    level.WorkflowLevelType.Type,
	// 		},
	// 	}

	// 	// Assuming the participants are loaded as part of the WorkflowLevel model
	// 	for _, participant := range level.Participants {
	// 		participantDTO := settingsDTOs.WorkflowLevelParticipantDTO{
	// 			ID: participant.ID,
	// 		}

	// 		if participant.User != nil {
	// 			participantDTO.User = &settingsDTOs.UserMinimalDTO{
	// 				ID:               int(participant.User.ID),
	// 				FirstNameEN:      participant.User.FirstNameEn,
	// 				LastNameEN:       participant.User.LastNameEn,
	// 				FirstNameAR:      participant.User.FirstNameAr,
	// 				LastNameAR:       participant.User.LastNameAr,
	// 				ProfileImagePath: participant.User.ProfileImagePath,
	// 				Email:            participant.User.Email,
	// 				Initial:          participant.User.Initial,
	// 				Color:            participant.User.Color,
	// 			}
	// 		}

	// 		if participant.Role != nil {
	// 			participantDTO.Role = &settingsDTOs.RoleDTO{
	// 				ID:      participant.Role.ID,
	// 				TitleEn: participant.Role.TitleEn,
	// 				TitleAr: participant.Role.TitleAr,
	// 				Status:  participant.Role.Status,
	// 			}
	// 		}

	// 		levelDTO.Participants = append(levelDTO.Participants, participantDTO)
	// 	}

	// 	levelsDTO.WorkflowLevels[i] = levelDTO
	// }

	// return levelsDTO, nil
}

func (ws *WorkflowService) CreateWorkFlowLevel(workflowId uint, currentUserID uint, workflowLevelFormDTO *settingsDTOs.WorkflowLevelFormDTO) (uint, error) {
	tx := ws.repo.BeginTransaction()

	if err := ws.validateWorkflowFormDTO(workflowLevelFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	maxLevel, err := ws.GetMaxLevelForWorkflow(workflowLevelFormDTO.WorkflowID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var newLevel int
	if workflowLevelFormDTO.ParantLevel == nil {
		if maxLevel == 0 {
			newLevel = 1 // Set as the first level
		} else {
			tx.Rollback()
			return 0, fmt.Errorf("parent level must be specified when other levels exist")
		}
	} else {
		newLevel = int(*workflowLevelFormDTO.ParantLevel) + 1
	}

	// Adjust subsequent levels if the new level is not the next immediate level and is less than or equal to max level
	if newLevel <= maxLevel {
		err := ws.repo.IncrementSubsequentLevels(workflowLevelFormDTO.WorkflowID, uint(newLevel), tx)
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
		if len(workflowLevelFormDTO.UserIDs) > 0 {
			tx.Rollback()
			return 0, fmt.Errorf("user is required for selected level type")
		}
		for _, userID := range workflowLevelFormDTO.UserIDs {
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
		err := ws.activityLogSvc.LogActivity(
			models.CreateAction,
			models.WorkflowLevel{}.TableName(),
			id,
			workflowLevelFormDTO,
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
			nil, // You can define a struct or map that describes the change
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

	// Check if no levels exist and start from 1
	// if maxLevel == 0 {
	// 	maxLevel = 1 // Set the first level explicitly if no levels exist
	// } else {
	// 	maxLevel += 1 // Increment to set the next level
	// }

	return maxLevel, nil
}
