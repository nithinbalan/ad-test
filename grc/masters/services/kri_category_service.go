package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
	masterModels "jethur/grc/models"
	"jethur/utility"
	"log"
)

type KRICategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*masterModels.KRICategory, error)
	GetByID(id uint) (*masterModels.KRICategory, error)
	Create(currentUserID uint, category *masterModels.KRICategory) (*masterModels.KRICategory, error)
	Update(id uint, currentUserID uint, category *masterModels.KRICategory) (*masterModels.KRICategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type KRICategoryService struct {
	repo           repositories.KRICategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewKRICategoryService(repo repositories.KRICategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) KRICategoryServiceInterface {
	return &KRICategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (kcs *KRICategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*masterModels.KRICategory, error) {
	return kcs.repo.GetAll(sort, sortBy, keyword, status)
}

func (kcs KRICategoryService) GetByID(id uint) (*masterModels.KRICategory, error) {
	return kcs.repo.GetByID(id)
}

func (kcs *KRICategoryService) Create(currentUserID uint, category *masterModels.KRICategory) (*masterModels.KRICategory, error) {

	if err := kcs.validateKRICategoryForm(category, nil); err != nil {
		return nil, err
	}

	createdCategory, err := kcs.repo.Create(category)
	if err != nil {
		return nil, err
	}

	go func() {
		err := kcs.activityLogSvc.LogActivity(
			models.CreateAction,
			masterModels.KRICategory{}.TableName(),
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

func (kcs *KRICategoryService) Update(id uint, currentUserID uint, updatedCategory *masterModels.KRICategory) (*masterModels.KRICategory, error) {
	existingCategory, err := kcs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if err := kcs.validateKRICategoryForm(updatedCategory, &id); err != nil {
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
			masterModels.KRICategory{}.TableName(),
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

func (kcs *KRICategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := kcs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = kcs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := kcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				masterModels.KRICategory{}.TableName(),
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

	return failedIds, nil
}

func (kcs *KRICategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := kcs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := kcs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				masterModels.KRICategory{}.TableName(),
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

func (kcs *KRICategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := kcs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := kcs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				masterModels.KRICategory{}.TableName(),
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

func (kcs *KRICategoryService) validateKRICategoryForm(category *masterModels.KRICategory, id *uint) error {

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
