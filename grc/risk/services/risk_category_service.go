package services

import (
	"errors"
	"fmt"
	logModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/grc/risk/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type RiskCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskCategory, error)
	GetByID(id uint) (*models.RiskCategory, error)
	Create(currentUserID uint, category *models.RiskCategory) (*models.RiskCategory, error)
	Update(id uint, currentUserID uint, category *models.RiskCategory) (*models.RiskCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type RiskCategoryService struct {
	repo           repositories.RiskCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewRiskCategoryService(repo repositories.RiskCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) RiskCategoryServiceInterface {
	return &RiskCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (rcs *RiskCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskCategory, error) {
	return rcs.repo.GetAll(sort, sortBy, keyword, status)
}

func (rcs RiskCategoryService) GetByID(id uint) (*models.RiskCategory, error) {
	riskCategory, err := rcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskCategory, nil
}

func (rcs *RiskCategoryService) Create(currentUserID uint, category *models.RiskCategory) (*models.RiskCategory, error) {

	if exists, err := rcs.repo.ExistsByTitle(category.TitleEn, category.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdCategory, err := rcs.repo.Create(category)
	if err != nil {
		return nil, err
	}

	go func() {
		err := rcs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.RiskCategory{}.TableName(),
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

func (rcs *RiskCategoryService) Update(id uint, currentUserID uint, updatedCategory *models.RiskCategory) (*models.RiskCategory, error) {
	existingCategory, err := rcs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedCategory.TitleEn != existingCategory.TitleEn || updatedCategory.TitleAr != existingCategory.TitleAr {
		if exists, err := rcs.repo.ExistsByTitleExcludingID(id, updatedCategory.TitleEn, updatedCategory.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := rcs.repo.Update(id, updatedCategory)

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
		err := rcs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.RiskCategory{}.TableName(),
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

func (rcs *RiskCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := rcs.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = rcs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := rcs.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.RiskCategory{}.TableName(),
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

func (rcs *RiskCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := rcs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rcs.activityLogSvc.LogBulkActivities(
				logModels.ActivateAction,
				models.RiskCategory{}.TableName(),
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

func (rcs *RiskCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := rcs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rcs.activityLogSvc.LogBulkActivities(
				logModels.DeactivateAction,
				models.RiskCategory{}.TableName(),
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
