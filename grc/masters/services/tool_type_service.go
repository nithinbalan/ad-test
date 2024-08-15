package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
	"jethur/utility"
	"log"
)

type ToolTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ToolType, error)
	GetByID(id uint) (*models.ToolType, error)
	Create(currentUserID uint, toolType *models.ToolType) (*models.ToolType, error)
	Update(id uint, currentUserID uint, toolType *models.ToolType) (*models.ToolType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ToolTypeService struct {
	repo           repositories.ToolTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewToolTypeService(repo repositories.ToolTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ToolTypeServiceInterface {
	return &ToolTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (tts *ToolTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ToolType, error) {
	return tts.repo.GetAll(sort, sortBy, keyword, status)
}

func (tts ToolTypeService) GetByID(id uint) (*models.ToolType, error) {
	return tts.repo.GetByID(id)
}

func (tts *ToolTypeService) Create(currentUserID uint, toolType *models.ToolType) (*models.ToolType, error) {

	if err := tts.validateToolTypeForm(toolType, nil); err != nil {
		return nil, err
	}

	createdToolType, err := tts.repo.Create(toolType)
	if err != nil {
		return nil, err
	}

	go func() {
		err := tts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ToolType{}.TableName(),
			toolType.ID,
			createdToolType,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdToolType, nil
}

func (tts *ToolTypeService) Update(id uint, currentUserID uint, updatedToolType *models.ToolType) (*models.ToolType, error) {
	existingToolType, err := tts.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if err := tts.validateToolTypeForm(updatedToolType, &id); err != nil {
		return nil, err
	}

	result, err := tts.repo.Update(id, updatedToolType)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "status"}
	changes, err := utility.CompareFields(existingToolType, updatedToolType, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := tts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ToolType{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()
	return result, err
}

func (tts *ToolTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := tts.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = tts.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := tts.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ToolType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (tts *ToolTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := tts.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := tts.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ToolType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (tts *ToolTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := tts.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := tts.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ToolType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (tts *ToolTypeService) validateToolTypeForm(toolTypeForm *models.ToolType, id *uint) error {

	if id == nil {
		if exists, err := tts.repo.ExistsByTitle(toolTypeForm.TitleEn, toolTypeForm.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {
		if exists, err := tts.repo.ExistsByTitleExcludingID(*id, toolTypeForm.TitleEn, toolTypeForm.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	return nil
}
