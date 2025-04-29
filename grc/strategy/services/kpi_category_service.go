package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/strategy/repositories"
	"jethur/utility"
	"log"
)

type KPICategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KPICategory, error)
	GetByID(id uint) (*models.KPICategory, error)
	Create(currentUserID uint, category *models.KPICategory) (*models.KPICategory, error)
	Update(id uint, currentUserID uint, category *models.KPICategory) (*models.KPICategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type KPICategoryService struct {
	repo           repositories.KPICategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewKPICategoryService(repo repositories.KPICategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) KPICategoryServiceInterface {
	return &KPICategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (kcs *KPICategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KPICategory, error) {
	return kcs.repo.GetAll(sort, sortBy, keyword, status)
}

func (kcs KPICategoryService) GetByID(id uint) (*models.KPICategory, error) {
	return kcs.repo.GetByID(id)
}

func (kcs *KPICategoryService) Create(currentUserID uint, category *models.KPICategory) (*models.KPICategory, error) {

	if err := kcs.validateKPICategoryForm(category, nil); err != nil {
		return nil, err
	}

	createdCategory, err := kcs.repo.Create(category)
	if err != nil {
		return nil, err
	}

	go func() {
		err := kcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.KPICategory{}.TableName(),
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

func (kcs *KPICategoryService) Update(id uint, currentUserID uint, updatedCategory *models.KPICategory) (*models.KPICategory, error) {
	existingCategory, err := kcs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if err := kcs.validateKPICategoryForm(updatedCategory, &id); err != nil {
		return nil, err
	}

	result, err := kcs.repo.Update(id, updatedCategory)

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
		err := kcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.KPICategory{}.TableName(),
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

func (kcs *KPICategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedKpiCategory []map[string]interface{}
	for _, id := range ids {

		kriCategory, err := kcs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = kcs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedKpiCategory = append(deletedKpiCategory, map[string]interface{}{
			"id":    id,
			"title": kriCategory.TitleEn,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := kcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.KPICategory{}.TableName(),
				successfulIds,
				deletedKpiCategory,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (kcs *KPICategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := kcs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := kcs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.KPICategory{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (kcs *KPICategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := kcs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := kcs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.KPICategory{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (kcs *KPICategoryService) validateKPICategoryForm(category *models.KPICategory, id *uint) error {

	if id == nil {

		if exists, err := kcs.repo.ExistsByTitle(category.TitleEn, category.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {

		if exists, err := kcs.repo.ExistsByTitleExcludingID(*id, category.TitleEn, category.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	return nil
}
