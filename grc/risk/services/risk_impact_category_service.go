package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"
	"log"
)

type RiskImpactCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskImpactCategory, error)
	GetByID(id uint) (*models.RiskImpactCategory, error)
	Create(currentUserID uint, category *models.RiskImpactCategory) (*models.RiskImpactCategory, error)
	Update(id uint, currentUserID uint, category *models.RiskImpactCategory) (*models.RiskImpactCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type RiskImpactCategoryService struct {
	repo           repositories.RiskImpactCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewRiskImpactCategoryService(repo repositories.RiskImpactCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) RiskImpactCategoryServiceInterface {
	return &RiskImpactCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (rics *RiskImpactCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskImpactCategory, error) {
	return rics.repo.GetAll(sort, sortBy, keyword, status)
}

func (rics RiskImpactCategoryService) GetByID(id uint) (*models.RiskImpactCategory, error) {
	return rics.repo.GetByID(id)
}

func (rics *RiskImpactCategoryService) Create(currentUserID uint, category *models.RiskImpactCategory) (*models.RiskImpactCategory, error) {

	if exists, err := rics.repo.ExistsByTitle(category.TitleEn, category.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdCategory, err := rics.repo.Create(category)
	if err != nil {
		return nil, err
	}

	go func() {
		err := rics.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskImpactCategory{}.TableName(),
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

func (rics *RiskImpactCategoryService) Update(id uint, currentUserID uint, updatedCategory *models.RiskImpactCategory) (*models.RiskImpactCategory, error) {
	existingCategory, err := rics.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedCategory.TitleEn != existingCategory.TitleEn || updatedCategory.TitleAr != existingCategory.TitleAr {
		if exists, err := rics.repo.ExistsByTitleExcludingID(id, updatedCategory.TitleEn, updatedCategory.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := rics.repo.Update(id, updatedCategory)

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
		err := rics.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskImpactCategory{}.TableName(),
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

func (rics *RiskImpactCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedImpactCategory []map[string]interface{}
	for _, id := range ids {

		impactCategory, err := rics.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = rics.repo.Delete(id)
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
			err := rics.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.RiskImpactCategory{}.TableName(),
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

func (rics *RiskImpactCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := rics.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rics.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.RiskImpactCategory{}.TableName(),
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

func (rics *RiskImpactCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := rics.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rics.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.RiskImpactCategory{}.TableName(),
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
