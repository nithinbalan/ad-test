package services

import (
	"errors"
	"fmt"
	"log"

	awarenessRepositories "jethur/grc/awareness-reporting/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AwarenessReportSubCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessReportSubCategory, error)
	GetByID(id uint) (*models.AwarenessReportSubCategory, error)
	Create(currentUserID uint, subCategory *models.AwarenessReportSubCategory) (*models.AwarenessReportSubCategory, error)
	Update(id uint, currentUserID uint, subCategory *models.AwarenessReportSubCategory) (*models.AwarenessReportSubCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type AwarenessReportSubCategoryService struct {
	repo           awarenessRepositories.AwarenessReportSubCategoryRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewAwarenessReportSubCategoryService(
	repo awarenessRepositories.AwarenessReportSubCategoryRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) AwarenessReportSubCategoryServiceInterface {
	return &AwarenessReportSubCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (arscs *AwarenessReportSubCategoryService) GetAll(sort string,sortBy string,keyword string,status string,) ([]*models.AwarenessReportSubCategory, error) {
	return arscs.repo.GetAll(sort, sortBy, keyword, status)
}

func (arscs *AwarenessReportSubCategoryService) GetByID(id uint) (*models.AwarenessReportSubCategory, error) {
	subCategory, err := arscs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return subCategory, nil
}

func (arscs *AwarenessReportSubCategoryService) Create(currentUserID uint,subCategory *models.AwarenessReportSubCategory,) (*models.AwarenessReportSubCategory, error) {

	if exists, err := arscs.repo.ExistsByTitle(subCategory.TitleEn, subCategory.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(
			fmt.Errorf("the title must be unique"),
			"unique",
			"title",
		)
	}

	createdSubCategory, err := arscs.repo.Create(subCategory)
	if err != nil {
		return nil, err
	}

	go func() {
		err := arscs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AwarenessReportSubCategory{}.TableName(),
			createdSubCategory.ID,
			createdSubCategory,
			currentUserID,
			"", "", "", nil,
		)

		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdSubCategory, nil
}

func (arscs *AwarenessReportSubCategoryService) Update(id uint,currentUserID uint,updatedSubCategory *models.AwarenessReportSubCategory,) (*models.AwarenessReportSubCategory, error) {

	existingSubCategory, err := arscs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedSubCategory.TitleEn != existingSubCategory.TitleEn ||
		updatedSubCategory.TitleAr != existingSubCategory.TitleAr {

		if exists, err := arscs.repo.ExistsByTitleExcludingID(
			id,
			updatedSubCategory.TitleEn,
			updatedSubCategory.TitleAr,
		); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(
				fmt.Errorf("the title must be unique"),
				"unique",
				"title",
			)
		}
	}

	result, err := arscs.repo.Update(id, updatedSubCategory)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{
		"TitleEn",
		"TitleAr",
		"Status",
	}

	changes, err := utility.CompareFields(
		existingSubCategory,
		updatedSubCategory,
		fieldsToCompare,
	)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := arscs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AwarenessReportSubCategory{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)

		if err != nil {
			fmt.Printf("Failed to log update activity: %v\n", err)
		}
	}()

	return result, nil
}

func (arscs *AwarenessReportSubCategoryService) Delete(ids []uint,currentUserID uint,) ([]uint, error) {

	var failedIds []uint
	var successfulIds []uint
	var deletedSubCategories []map[string]interface{}

	for _, id := range ids {
		subCategory, err := arscs.repo.GetByID(id)
		if err != nil {
			continue
		}

		err = arscs.repo.Delete(id)

		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedSubCategories = append(
			deletedSubCategories,
			map[string]interface{}{
				"id":       id,
				"title_en": subCategory.TitleEn,
				"title_ar": subCategory.TitleAr,
			},
		)
	}

	if len(successfulIds) > 0 {
		go func() {
			err := arscs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AwarenessReportSubCategory{}.TableName(),
				successfulIds,
				deletedSubCategories,
				currentUserID,
				nil,
			)

			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (arscs *AwarenessReportSubCategoryService) Activate(ids []uint,currentUserID uint,) error {

	successfulIds, err := arscs.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := arscs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.AwarenessReportSubCategory{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)

			if err != nil {
				fmt.Printf("Failed to log activate activity: %v\n", err)
			}
		}()
	}

	return err
}

func (arscs *AwarenessReportSubCategoryService) Deactivate(ids []uint,currentUserID uint,) error {

	successfulIds, err := arscs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := arscs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.AwarenessReportSubCategory{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)

			if err != nil {
				fmt.Printf("Failed to log deactivate activity: %v\n", err)
			}
		}()
	}

	return err
}