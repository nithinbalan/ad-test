package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type ImpactCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ImpactCategory, error)
	GetByID(id uint) (*models.ImpactCategory, error)
	Create(currentUserID uint, category *models.ImpactCategory) (*models.ImpactCategory, error)
	Update(id uint, currentUserID uint, category *models.ImpactCategory) (*models.ImpactCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ImpactCategoryService struct {
	repo           repositories.ImpactCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewImpactCategoryService(repo repositories.ImpactCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ImpactCategoryServiceInterface {
	return &ImpactCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ics *ImpactCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ImpactCategory, error) {
	return ics.repo.GetAll(sort, sortBy, keyword, status)
}

func (ics ImpactCategoryService) GetByID(id uint) (*models.ImpactCategory, error) {
	impactCategory, err := ics.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return impactCategory, nil
}

func (ics *ImpactCategoryService) Create(currentUserID uint, category *models.ImpactCategory) (*models.ImpactCategory, error) {

	if exists, err := ics.repo.ExistsByTitle(category.TitleEn, category.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdCategory, err := ics.repo.Create(category)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ics.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ImpactCategory{}.TableName(),
			createdCategory.ID,
			createdCategory,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdCategory, nil
}

func (ics *ImpactCategoryService) Update(id uint, currentUserID uint, updatedCategory *models.ImpactCategory) (*models.ImpactCategory, error) {
	existingCategory, err := ics.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedCategory.TitleEn != existingCategory.TitleEn || updatedCategory.TitleAr != existingCategory.TitleAr {
		if exists, err := ics.repo.ExistsByTitleExcludingID(id, updatedCategory.TitleEn, updatedCategory.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := ics.repo.Update(id, updatedCategory)

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
		err := ics.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ImpactCategory{}.TableName(),
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

func (ics *ImpactCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedImpactCategory []map[string]interface{}
	for _, id := range ids {

		impactCategory, err := ics.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = ics.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedImpactCategory = append(deletedImpactCategory, map[string]interface{}{
			"id":       id,
			"title_en": impactCategory.TitleEn,
			"title_ar": impactCategory.TitleAr,
		})

	}

	if len(successfulIds) > 0 {

		go func() {
			err := ics.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ImpactCategory{}.TableName(),
				successfulIds,
				deletedImpactCategory,
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

func (ics *ImpactCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ics.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ics.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ImpactCategory{}.TableName(),
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

func (ics *ImpactCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ics.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ics.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ImpactCategory{}.TableName(),
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
