package services

import (
	"fmt"
	"jethur/grc/audit/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"
)

type AuditCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditCategory, error)
	GetByID(id uint) (*models.AuditCategory, error)
	Create(currentUserID uint, category *models.AuditCategory) (*models.AuditCategory, error)
	Update(id uint, currentUserID uint, category *models.AuditCategory) (*models.AuditCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type AuditCategoryService struct {
	repo           repositories.AuditCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewAuditCategoryService(repo repositories.AuditCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) AuditCategoryServiceInterface {
	return &AuditCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (acs *AuditCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditCategory, error) {
	return acs.repo.GetAll(sort, sortBy, keyword, status)
}

func (acs AuditCategoryService) GetByID(id uint) (*models.AuditCategory, error) {
	return acs.repo.GetByID(id)
}

func (acs *AuditCategoryService) Create(currentUserID uint, category *models.AuditCategory) (*models.AuditCategory, error) {

	if exists, err := acs.repo.ExistsByTitle(category.TitleEn, category.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdCategory, err := acs.repo.Create(category)
	if err != nil {
		return nil, err
	}

	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditCategory{}.TableName(),
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

func (acs *AuditCategoryService) Update(id uint, currentUserID uint, updatedCategory *models.AuditCategory) (*models.AuditCategory, error) {
	existingCategory, err := acs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedCategory.TitleEn != existingCategory.TitleEn || updatedCategory.TitleAr != existingCategory.TitleAr {
		if exists, err := acs.repo.ExistsByTitleExcludingID(id, updatedCategory.TitleEn, updatedCategory.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := acs.repo.Update(id, updatedCategory)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr","DescriptionEn","DescriptionAr", "Status"}
	changes, err := utility.CompareFields(existingCategory, updatedCategory, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditCategory{}.TableName(),
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

func (acs *AuditCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := acs.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = acs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.AuditCategory{}.TableName(),
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

func (acs *AuditCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := acs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.AuditCategory{}.TableName(),
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

func (acs *AuditCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := acs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.AuditCategory{}.TableName(),
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
