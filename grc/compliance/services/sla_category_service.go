package services

import (
	"fmt"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/utility"
	"log"
)

type SLACategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SLACategory, error)
	GetByID(id uint) (*models.SLACategory, error)
	Create(currentUserID uint, category *models.SLACategory) (*models.SLACategory, error)
	Update(id uint, currentUserID uint, category *models.SLACategory) (*models.SLACategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type SLACategoryService struct {
	repo           repositories.SLACategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewSLACategoryService(repo repositories.SLACategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) SLACategoryServiceInterface {
	return &SLACategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (scs *SLACategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SLACategory, error) {
	return scs.repo.GetAll(sort, sortBy, keyword, status)
}

func (scs SLACategoryService) GetByID(id uint) (*models.SLACategory, error) {
	return scs.repo.GetByID(id)
}

func (scs *SLACategoryService) Create(currentUserID uint, category *models.SLACategory) (*models.SLACategory, error) {

	if exists, err := scs.repo.ExistsByTitle(category.TitleEn, category.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdCategory, err := scs.repo.Create(category)
	if err != nil {
		return nil, err
	}

	go func() {
		err := scs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.SLACategory{}.TableName(),
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

func (scs *SLACategoryService) Update(id uint, currentUserID uint, updatedCategory *models.SLACategory) (*models.SLACategory, error) {
	existingCategory, err := scs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedCategory.TitleEn != existingCategory.TitleEn || updatedCategory.TitleAr != existingCategory.TitleAr {
		if exists, err := scs.repo.ExistsByTitleExcludingID(id, updatedCategory.TitleEn, updatedCategory.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := scs.repo.Update(id, updatedCategory)

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
		err := scs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.SLACategory{}.TableName(),
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

func (scs *SLACategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := scs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = scs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := scs.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction,
				models.SLACategory{}.TableName(),
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

func (scs *SLACategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := scs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := scs.activityLogSvc.LogBulkActivities(
				logModels.ActivateAction,
				models.SLACategory{}.TableName(),
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

func (scs *SLACategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := scs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := scs.activityLogSvc.LogBulkActivities(
				logModels.DeactivateAction,
				models.SLACategory{}.TableName(),
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
