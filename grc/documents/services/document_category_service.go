package services

import (
	"errors"
	"fmt"
	"jethur/grc/documents/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type DocumentCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentCategory, error)
	GetByID(id uint) (*models.DocumentCategory, error)
	Create(currentUserID uint, category *models.DocumentCategory) (*models.DocumentCategory, error)
	Update(id uint, currentUserID uint, category *models.DocumentCategory) (*models.DocumentCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type DocumentCategoryService struct {
	repo           repositories.DocumentCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewDocumentCategoryService(repo repositories.DocumentCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) DocumentCategoryServiceInterface {
	return &DocumentCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (dcs *DocumentCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentCategory, error) {
	return dcs.repo.GetAll(sort, sortBy, keyword, status)
}

func (dcs *DocumentCategoryService) GetByID(id uint) (*models.DocumentCategory, error) {
	documentCategory, err := dcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return documentCategory, nil
}

func (dcs *DocumentCategoryService) Create(currentUserID uint, category *models.DocumentCategory) (*models.DocumentCategory, error) {
	if exists, err := dcs.repo.ExistsByTitle(category.TitleEn, category.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdCategory, err := dcs.repo.Create(category)
	if err != nil {
		return nil, err
	}

	go func() {
		err := dcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.DocumentCategory{}.TableName(),
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

func (dcs *DocumentCategoryService) Update(id uint, currentUserID uint, updatedCategory *models.DocumentCategory) (*models.DocumentCategory, error) {
	existingCategory, err := dcs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedCategory.TitleEn != existingCategory.TitleEn || updatedCategory.TitleAr != existingCategory.TitleAr {
		if exists, err := dcs.repo.ExistsByTitleExcludingID(id, updatedCategory.TitleEn, updatedCategory.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := dcs.repo.Update(id, updatedCategory)
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
		err := dcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.DocumentCategory{}.TableName(),
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

func (dcs *DocumentCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {
		_, err := dcs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = dcs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := dcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.DocumentCategory{}.TableName(),
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

func (dcs *DocumentCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := dcs.repo.Activate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := dcs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.DocumentCategory{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk activate activity: %v\n", err)
			}
		}()
	}
	return err
}

func (dcs *DocumentCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := dcs.repo.Deactivate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := dcs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.DocumentCategory{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk deactivate activity: %v\n", err)
			}
		}()
	}
	return err
}
