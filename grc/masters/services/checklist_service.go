package services

import (
	"fmt"
	activityLogModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/grc/masters/models"
	"jethur/grc/masters/repositories"
	"jethur/utility"
	"log"
)

type ChecklistServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Checklist, error)
	GetByID(id uint) (*models.Checklist, error)
	Create(currentUserID uint, checklist *models.Checklist) (*models.Checklist, error)
	Update(id uint, currentUserID uint, checklist *models.Checklist) (*models.Checklist, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ChecklistService struct {
	repo           repositories.ChecklistRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewChecklistService(repo repositories.ChecklistRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ChecklistServiceInterface {
	return &ChecklistService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (cs *ChecklistService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Checklist, error) {
	return cs.repo.GetAll(sort, sortBy, keyword, status)
}

func (cs ChecklistService) GetByID(id uint) (*models.Checklist, error) {
	return cs.repo.GetByID(id)
}

func (cs *ChecklistService) Create(currentUserID uint, checklist *models.Checklist) (*models.Checklist, error) {

	if err := cs.validateChecklistForm(checklist, nil); err != nil {
		return nil, err
	}

	createdChecklist, err := cs.repo.Create(checklist)
	if err != nil {
		return nil, err
	}

	go func() {
		err := cs.activityLogSvc.LogActivity(
			activityLogModels.CreateAction,
			models.Checklist{}.TableName(),
			checklist.ID,
			createdChecklist,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdChecklist, nil
}

func (cs *ChecklistService) Update(id uint, currentUserID uint, updatedChecklist *models.Checklist) (*models.Checklist, error) {
	existingChecklist, err := cs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if err := cs.validateChecklistForm(updatedChecklist, &id); err != nil {
		return nil, err
	}

	result, err := cs.repo.Update(id, updatedChecklist)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingChecklist, updatedChecklist, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := cs.activityLogSvc.LogActivity(
			activityLogModels.UpdateAction,
			models.Checklist{}.TableName(),
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

func (cs *ChecklistService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := cs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = cs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := cs.activityLogSvc.LogBulkActivities(
				activityLogModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.Checklist{}.TableName(),
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

func (cs *ChecklistService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := cs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := cs.activityLogSvc.LogBulkActivities(
				activityLogModels.ActivateAction,
				models.Checklist{}.TableName(),
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

func (cs *ChecklistService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := cs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := cs.activityLogSvc.LogBulkActivities(
				activityLogModels.DeactivateAction,
				models.Checklist{}.TableName(),
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

func (cs *ChecklistService) validateChecklistForm(checklistForm *models.Checklist, id *uint) error {

	if id == nil {
		if exists, err := cs.repo.ExistsByTitle(checklistForm.TitleEn, checklistForm.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {
		if exists, err := cs.repo.ExistsByTitleExcludingID(*id, checklistForm.TitleEn, checklistForm.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	return nil
}
