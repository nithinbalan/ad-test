package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/strategy/repositories"
	"jethur/utility"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type SwotPestelQuestionGroupServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, isSwot string, isPestel string) ([]*models.SwotPestelQuestionGroup, error)
	GetByID(id uint) (*models.SwotPestelQuestionGroup, error)
	Create(currentUserID uint, swotPestelQuestionGroup *models.SwotPestelQuestionGroup) (*models.SwotPestelQuestionGroup, error)
	Update(id uint, currentUserID uint, swotPestelQuestionGroup *models.SwotPestelQuestionGroup) (*models.SwotPestelQuestionGroup, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type SwotPestelQuestionGroupService struct {
	repo           repositories.SwotPestelQuestionGroupRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewSwotPestelQuestionGroupService(repo repositories.SwotPestelQuestionGroupRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) SwotPestelQuestionGroupServiceInterface {
	return &SwotPestelQuestionGroupService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (sqgs *SwotPestelQuestionGroupService) GetAll(sort string, sortBy string, keyword string, status string, isSwot string, isPestel string) ([]*models.SwotPestelQuestionGroup, error) {
	return sqgs.repo.GetAll(sort, sortBy, keyword, status, isSwot, isPestel)
}

func (sqgs SwotPestelQuestionGroupService) GetByID(id uint) (*models.SwotPestelQuestionGroup, error) {
	swotPestelQuestionGroup, err := sqgs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return swotPestelQuestionGroup, nil
}

func (sqgs *SwotPestelQuestionGroupService) Create(currentUserID uint, swotPestelQuestionGroup *models.SwotPestelQuestionGroup) (*models.SwotPestelQuestionGroup, error) {
	if exists, err := sqgs.repo.ExistsByTitle(swotPestelQuestionGroup.TitleEn, swotPestelQuestionGroup.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	existingGroups, err := sqgs.repo.GetAll("ASC", "[order]", "", "active", strconv.FormatBool(swotPestelQuestionGroup.IsSwot), strconv.FormatBool(swotPestelQuestionGroup.IsPestel))
	if err != nil {
		return nil, err
	}

	newOrder := swotPestelQuestionGroup.Order
	for _, group := range existingGroups {
		if group.Order >= newOrder {
			group.Order++ // Increment order for existing records
			if _, err := sqgs.repo.Update(group.ID, group); err != nil {
				return nil, err
			}
		}
	}

	createdSwotPestelQuestionGroup, err := sqgs.repo.Create(swotPestelQuestionGroup)
	if err != nil {
		return nil, err
	}

	go func() {
		err := sqgs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SwotPestelQuestionGroup{}.TableName(),
			createdSwotPestelQuestionGroup.ID,
			createdSwotPestelQuestionGroup,
			currentUserID,
			nil,
		)

		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdSwotPestelQuestionGroup, nil
}

func (sqgs *SwotPestelQuestionGroupService) Update(id uint, currentUserID uint, updatedSwotPestelQuestionGroup *models.SwotPestelQuestionGroup) (*models.SwotPestelQuestionGroup, error) {
	existingSwotPestelQuestionGroup, err := sqgs.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedSwotPestelQuestionGroup.TitleEn != existingSwotPestelQuestionGroup.TitleEn || updatedSwotPestelQuestionGroup.TitleAr != existingSwotPestelQuestionGroup.TitleAr {
		if exists, err := sqgs.repo.ExistsByTitleExcludingID(id, updatedSwotPestelQuestionGroup.TitleEn, updatedSwotPestelQuestionGroup.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := sqgs.repo.Update(id, updatedSwotPestelQuestionGroup)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "order", "Status"}
	changes, err := utility.CompareFields(existingSwotPestelQuestionGroup, updatedSwotPestelQuestionGroup, fieldsToCompare)

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
			models.SwotPestelQuestionGroup{}.TableName(),
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

func (sqgs *SwotPestelQuestionGroupService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {

		// Get existing record
		swotPestelQuestionGroup, err := sqgs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			failedIds = append(failedIds, id)
			continue
		}

		// Fetch all records in order before deletion
		swotGroups, err := sqgs.repo.GetAll("ASC", "[order]", "", "active", strconv.FormatBool(swotPestelQuestionGroup.IsSwot), strconv.FormatBool(swotPestelQuestionGroup.IsPestel))
		if err != nil {
			return nil, err
		}

		err = sqgs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		// Adjust ordering for remaining records
		for _, group := range swotGroups {
			if group.ID != id && group.Order > swotPestelQuestionGroup.Order {
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
				models.SwotPestelQuestionGroup{}.TableName(),
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

func (sqgs *SwotPestelQuestionGroupService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := sqgs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := sqgs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.SwotPestelQuestionGroup{}.TableName(),
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

func (sqgs *SwotPestelQuestionGroupService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := sqgs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := sqgs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.SwotPestelQuestionGroup{}.TableName(),
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
