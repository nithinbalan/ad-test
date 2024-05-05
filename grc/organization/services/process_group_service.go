package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type ProcessGroupServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProcessGroup, error)
	GetByID(id uint) (*models.ProcessGroup, error)
	Create(currentUserID uint, processGroup *models.ProcessGroup) (*models.ProcessGroup, error)
	Update(id uint, currentUserID uint, processGroup *models.ProcessGroup) (*models.ProcessGroup, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ProcessGroupService struct {
	repo           repositories.ProcessGroupRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewProcessGroupService(repo repositories.ProcessGroupRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ProcessGroupServiceInterface {
	return &ProcessGroupService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (pgs *ProcessGroupService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProcessGroup, error) {
	return pgs.repo.GetAll(sort, sortBy, keyword, status)
}

func (pgs *ProcessGroupService) GetByID(id uint) (*models.ProcessGroup, error) {
	processGroup, err := pgs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return processGroup, nil
}

func (pgs *ProcessGroupService) Create(currentUserID uint, processGroup *models.ProcessGroup) (*models.ProcessGroup, error) {
	if exists, err := pgs.repo.ExistsByTitle(processGroup.TitleEn, processGroup.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdProcessGroup, err := pgs.repo.Create(processGroup)
	if err != nil {
		return nil, err
	}

	go func() {
		err := pgs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProcessGroup{}.TableName(),
			createdProcessGroup.ID,
			createdProcessGroup,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdProcessGroup, nil
}

func (pgs *ProcessGroupService) Update(id uint, currentUserID uint, updatedProcessGroup *models.ProcessGroup) (*models.ProcessGroup, error) {
	existingProcessGroup, err := pgs.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedProcessGroup.TitleEn != existingProcessGroup.TitleEn || updatedProcessGroup.TitleAr != existingProcessGroup.TitleAr {
		if exists, err := pgs.repo.ExistsByTitleExcludingID(id, updatedProcessGroup.TitleEn, updatedProcessGroup.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := pgs.repo.Update(id, updatedProcessGroup)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingProcessGroup, updatedProcessGroup, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := pgs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ProcessGroup{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (pgs *ProcessGroupService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := pgs.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = pgs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := pgs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ProcessGroup{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (pgs *ProcessGroupService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := pgs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pgs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ProcessGroup{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (pgs *ProcessGroupService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := pgs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pgs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ProcessGroup{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}
