package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/tprm/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type ThirdPartyQuestionCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyQuestionCategory, error)
	GetByID(id uint) (*models.ThirdPartyQuestionCategory, error)
	Create(currentUserID uint, thirdPartyQuestionCategory *models.ThirdPartyQuestionCategory) (*models.ThirdPartyQuestionCategory, error)
	Update(id uint, currentUserID uint, thirdPartyQuestionCategory *models.ThirdPartyQuestionCategory) (*models.ThirdPartyQuestionCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ThirdPartyQuestionCategoryService struct {
	repo           repositories.ThirdPartyQuestionCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewThirdPartyQuestionCategoryService(repo repositories.ThirdPartyQuestionCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ThirdPartyQuestionCategoryServiceInterface {
	return &ThirdPartyQuestionCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (tpqcs *ThirdPartyQuestionCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyQuestionCategory, error) {
	return tpqcs.repo.GetAll(sort, sortBy, keyword, status)
}

func (tpqcs ThirdPartyQuestionCategoryService) GetByID(id uint) (*models.ThirdPartyQuestionCategory, error) {
	thirdPartyQuestionCategory, err := tpqcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return thirdPartyQuestionCategory, nil

}

func (tpqcs *ThirdPartyQuestionCategoryService) Create(currentUserID uint, thirdPartyQuestionCategory *models.ThirdPartyQuestionCategory) (*models.ThirdPartyQuestionCategory, error) {

	if exists, err := tpqcs.repo.ExistsByTitle(thirdPartyQuestionCategory.TitleEn, thirdPartyQuestionCategory.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdType, err := tpqcs.repo.Create(thirdPartyQuestionCategory)
	if err != nil {
		return nil, err
	}

	go func() {
		err := tpqcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ThirdPartyQuestionCategory{}.TableName(),
			createdType.ID,
			createdType,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdType, nil
}

func (tpqcs *ThirdPartyQuestionCategoryService) Update(id uint, currentUserID uint, updatedType *models.ThirdPartyQuestionCategory) (*models.ThirdPartyQuestionCategory, error) {
	existingType, err := tpqcs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedType.TitleEn != existingType.TitleEn || updatedType.TitleAr != existingType.TitleAr {
		if exists, err := tpqcs.repo.ExistsByTitleExcludingID(id, updatedType.TitleEn, updatedType.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := tpqcs.repo.Update(id, updatedType)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingType, updatedType, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := tpqcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ThirdPartyQuestionCategory{}.TableName(),
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

func (tpqcs *ThirdPartyQuestionCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedThirdPartyQuestionCategory []map[string]interface{}
	for _, id := range ids {

		thirdPartyQuestionCategory, err := tpqcs.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = tpqcs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedThirdPartyQuestionCategory = append(deletedThirdPartyQuestionCategory, map[string]interface{}{
			"id":       id,
			"title_en": thirdPartyQuestionCategory.TitleEn,
			"title_ar": thirdPartyQuestionCategory.TitleAr,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := tpqcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ThirdPartyQuestionCategory{}.TableName(),
				successfulIds,
				deletedThirdPartyQuestionCategory,
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

func (tpqcs *ThirdPartyQuestionCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := tpqcs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := tpqcs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ThirdPartyQuestionCategory{}.TableName(),
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

func (tpqcs *ThirdPartyQuestionCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := tpqcs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := tpqcs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ThirdPartyQuestionCategory{}.TableName(),
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
