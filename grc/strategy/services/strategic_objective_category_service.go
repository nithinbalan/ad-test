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

type StrategicObjectiveCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StrategicObjectiveCategory, error)
	GetByID(id uint) (*models.StrategicObjectiveCategory, error)
	Create(currentUserID uint, category *models.StrategicObjectiveCategory) (*models.StrategicObjectiveCategory, error)
	Update(id uint, currentUserID uint, category *models.StrategicObjectiveCategory) (*models.StrategicObjectiveCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type StrategicObjectiveCategoryService struct {
	repo           repositories.StrategicObjectiveCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewStrategicObjectiveCategoryService(repo repositories.StrategicObjectiveCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) StrategicObjectiveCategoryServiceInterface {
	return &StrategicObjectiveCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (socs *StrategicObjectiveCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StrategicObjectiveCategory, error) {
	return socs.repo.GetAll(sort, sortBy, keyword, status)
}

func (socs StrategicObjectiveCategoryService) GetByID(id uint) (*models.StrategicObjectiveCategory, error) {
	strategicObjectiveCategory, err := socs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategicObjectiveCategory, nil
}

func (socs *StrategicObjectiveCategoryService) Create(currentUserID uint, category *models.StrategicObjectiveCategory) (*models.StrategicObjectiveCategory, error) {

	if exists, err := socs.repo.ExistsByTitle(category.TitleEn, category.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdCategory, err := socs.repo.Create(category)
	if err != nil {
		return nil, err
	}

	go func() {
		err := socs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategicObjectiveCategory{}.TableName(),
			createdCategory.ID,
			createdCategory,
			currentUserID,
			"", "", "", nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdCategory, nil
}

func (socs *StrategicObjectiveCategoryService) Update(id uint, currentUserID uint, updatedCategory *models.StrategicObjectiveCategory) (*models.StrategicObjectiveCategory, error) {
	existingCategory, err := socs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedCategory.TitleEn != existingCategory.TitleEn || updatedCategory.TitleAr != existingCategory.TitleAr {
		if exists, err := socs.repo.ExistsByTitleExcludingID(id, updatedCategory.TitleEn, updatedCategory.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := socs.repo.Update(id, updatedCategory)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingCategory, updatedCategory, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := socs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.StrategicObjectiveCategory{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (socs *StrategicObjectiveCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedStrategicObjectiveCategory []map[string]interface{}
	for _, id := range ids {

		strategicObjectiveCategory, err := socs.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = socs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedStrategicObjectiveCategory = append(deletedStrategicObjectiveCategory, map[string]interface{}{
			"id":       id,
			"title_en": strategicObjectiveCategory.TitleEn,
			"title_ar": strategicObjectiveCategory.TitleAr,
		})

	}

	if len(successfulIds) > 0 {

		go func() {
			err := socs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.StrategicObjectiveCategory{}.TableName(),
				successfulIds,
				deletedStrategicObjectiveCategory,
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

func (socs *StrategicObjectiveCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := socs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := socs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.StrategicObjectiveCategory{}.TableName(),
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

func (socs *StrategicObjectiveCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := socs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := socs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.StrategicObjectiveCategory{}.TableName(),
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
