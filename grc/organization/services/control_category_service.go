package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type ControlCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlCategory, error)
	GetByID(id uint) (*models.ControlCategory, error)
	Create(currentUserID uint, controlCategory *models.ControlCategory) (*models.ControlCategory, error)
	Update(id uint, currentUserID uint, controlCategory *models.ControlCategory) (*models.ControlCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ControlCategoryService struct {
	repo           repositories.ControlCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewControlCategoryService(repo repositories.ControlCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ControlCategoryServiceInterface {
	return &ControlCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ccs *ControlCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlCategory, error) {
	return ccs.repo.GetAll(sort, sortBy, keyword, status)
}

func (ccs *ControlCategoryService) GetByID(id uint) (*models.ControlCategory, error) {
	controlCategory, err := ccs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return controlCategory, nil
}

func (ccs *ControlCategoryService) Create(currentUserID uint, controlCategory *models.ControlCategory) (*models.ControlCategory, error) {
	if exists, err := ccs.repo.ExistsByTitle(controlCategory.TitleEn, controlCategory.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdControlCategory, err := ccs.repo.Create(controlCategory)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ccs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ControlCategory{}.TableName(),
			createdControlCategory.ID,
			createdControlCategory,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdControlCategory, nil
}

func (ccs *ControlCategoryService) Update(id uint, currentUserID uint, updatedControlCategory *models.ControlCategory) (*models.ControlCategory, error) {
	existingControlCategory, err := ccs.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedControlCategory.TitleEn != existingControlCategory.TitleEn || updatedControlCategory.TitleAr != existingControlCategory.TitleAr {
		if exists, err := ccs.repo.ExistsByTitleExcludingID(id, updatedControlCategory.TitleEn, updatedControlCategory.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := ccs.repo.Update(id, updatedControlCategory)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingControlCategory, updatedControlCategory, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := ccs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ControlCategory{}.TableName(),
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

func (ccs *ControlCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := ccs.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
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
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ControlCategory{}.TableName(),
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

func (ccs *ControlCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ControlCategory{}.TableName(),
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

func (ccs *ControlCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ControlCategory{}.TableName(),
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
