package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/strategy/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type SwotQuestionGroupServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SwotQuestionGroup, error)
	GetByID(id uint) (*models.SwotQuestionGroup, error)
	Create(currentUserID uint, swotQuestionGroup *models.SwotQuestionGroup) (*models.SwotQuestionGroup, error)
	Update(id uint, currentUserID uint, swotQuestionGroup *models.SwotQuestionGroup) (*models.SwotQuestionGroup, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type SwotQuestionGroupService struct {
	repo           repositories.SwotQuestionGroupRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewSwotQuestionGroupService(repo repositories.SwotQuestionGroupRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) SwotQuestionGroupServiceInterface {
	return &SwotQuestionGroupService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (sqgs *SwotQuestionGroupService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SwotQuestionGroup, error) {
	return sqgs.repo.GetAll(sort, sortBy, keyword, status)
}

func (sqgs SwotQuestionGroupService) GetByID(id uint) (*models.SwotQuestionGroup, error) {
	swotQuestionGroup, err := sqgs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return swotQuestionGroup, nil
}

func (sqgs *SwotQuestionGroupService) Create(currentUserID uint, swotQuestionGroup *models.SwotQuestionGroup) (*models.SwotQuestionGroup, error) {
	if exists, err := sqgs.repo.ExistsByTitle(swotQuestionGroup.TitleEn, swotQuestionGroup.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	existingGroups, err := sqgs.repo.GetAll("ASC", "[order]", "", "active")
	if err != nil {
		return nil, err
	}

	newOrder := swotQuestionGroup.Order
	for _, group := range existingGroups {
		if group.Order >= newOrder {
			group.Order++ // Increment order for existing records
			if _, err := sqgs.repo.Update(group.ID, group); err != nil {
				return nil, err
			}
		}
	}

	createdSwotQuestionGroup, err := sqgs.repo.Create(swotQuestionGroup)
	if err != nil {
		return nil, err
	}

	go func() {
		err := sqgs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SwotQuestionGroup{}.TableName(),
			createdSwotQuestionGroup.ID,
			createdSwotQuestionGroup,
			currentUserID,
			nil,
		)

		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdSwotQuestionGroup, nil
}

func (sqgs *SwotQuestionGroupService) Update(id uint, currentUserID uint, updatedSwotQuestionGroup *models.SwotQuestionGroup) (*models.SwotQuestionGroup, error) {
	existingSwotQuestionGroup, err := sqgs.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedSwotQuestionGroup.TitleEn != existingSwotQuestionGroup.TitleEn || updatedSwotQuestionGroup.TitleAr != existingSwotQuestionGroup.TitleAr {
		if exists, err := sqgs.repo.ExistsByTitleExcludingID(id, updatedSwotQuestionGroup.TitleEn, updatedSwotQuestionGroup.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	swotGroups, err := sqgs.repo.GetAll("ASC", "[order]", "", "active")
	if err != nil {
		return nil, err
	}

	if updatedSwotQuestionGroup.Order != existingSwotQuestionGroup.Order {
		updatedOrder := updatedSwotQuestionGroup.Order
		oldOrder := existingSwotQuestionGroup.Order

		for _, group := range swotGroups {
			if group.ID == id {
				continue
			}
			if oldOrder < updatedOrder && group.Order > oldOrder && group.Order <= updatedOrder {
				group.Order--
				_, err = sqgs.repo.Update(group.ID, group)
			} else if oldOrder > updatedOrder && group.Order >= updatedOrder && group.Order < oldOrder {
				group.Order++
				_, err = sqgs.repo.Update(group.ID, group)
			}

			if err != nil {
				return nil, err
			}
		}
	}

	result, err := sqgs.repo.Update(id, updatedSwotQuestionGroup)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "order", "Status"}
	changes, err := utility.CompareFields(existingSwotQuestionGroup, updatedSwotQuestionGroup, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := sqgs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.SwotQuestionGroup{}.TableName(),
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

func (sqgs *SwotQuestionGroupService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	// Fetch all records in order before deletion
	swotGroups, err := sqgs.repo.GetAll("ASC", "[order]", "", "active")
	if err != nil {
		return nil, err
	}

	for _, id := range ids {

		// Get existing record
		swotQuestionGroup, err := sqgs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			failedIds = append(failedIds, id)
			continue
		}

		err = sqgs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		// Adjust ordering for remaining records
		for _, group := range swotGroups {
			if group.ID != id && group.Order > swotQuestionGroup.Order {
				group.Order--
				_, err := sqgs.repo.Update(group.ID, group)
				if err != nil {
					failedIds = append(failedIds, id) // If order update fails, mark as failed
					continue
				}
			}
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := sqgs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.SwotQuestionGroup{}.TableName(),
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

func (sqgs *SwotQuestionGroupService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := sqgs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := sqgs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.SwotQuestionGroup{}.TableName(),
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

func (sqgs *SwotQuestionGroupService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := sqgs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := sqgs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.SwotQuestionGroup{}.TableName(),
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
