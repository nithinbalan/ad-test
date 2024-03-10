package services

import (
	"fmt"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	activityLogModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/utility"
	"log"
)

type SlaCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SlaCategory, error)
	GetByID(id uint) (*models.SlaCategory, error)
	Create(currentUserID uint, category *models.SlaCategory) (*models.SlaCategory, error)
	Update(id uint, currentUserID uint, category *models.SlaCategory) (*models.SlaCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type SlaCategoryService struct {
	repo           repositories.SlaCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewSlaCategoryService(repo repositories.SlaCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) SlaCategoryServiceInterface {
	return &SlaCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ccs *SlaCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SlaCategory, error) {
	return ccs.repo.GetAll(sort, sortBy, keyword, status)
}

func (ccs SlaCategoryService) GetByID(id uint) (*models.SlaCategory, error) {
	return ccs.repo.GetByID(id)
}

func (ccs *SlaCategoryService) Create(currentUserID uint, category *models.SlaCategory) (*models.SlaCategory, error) {

	if exists, err := ccs.repo.ExistsByTitle(category.TitleEn, category.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdCategory, err := ccs.repo.Create(category)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ccs.activityLogSvc.LogActivity(
			activityLogModels.CreateAction,
			models.SlaCategory{}.TableName(),
			category.ID,
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

func (ccs *SlaCategoryService) Update(id uint, currentUserID uint, updatedCategory *models.SlaCategory) (*models.SlaCategory, error) {
	existingCategory, err := ccs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedCategory.TitleEn != existingCategory.TitleEn || updatedCategory.TitleAr != existingCategory.TitleAr {
		if exists, err := ccs.repo.ExistsByTitleExcludingID(id, updatedCategory.TitleEn, updatedCategory.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := ccs.repo.Update(id, updatedCategory)

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
		err := ccs.activityLogSvc.LogActivity(
			activityLogModels.UpdateAction,
			models.SlaCategory{}.TableName(),
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

func (ccs *SlaCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := ccs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = ccs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				activityLogModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.SlaCategory{}.TableName(),
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

func (ccs *SlaCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				activityLogModels.ActivateAction,
				models.SlaCategory{}.TableName(),
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

func (ccs *SlaCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				activityLogModels.DeactivateAction,
				models.SlaCategory{}.TableName(),
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
