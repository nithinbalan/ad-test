package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
	"jethur/utility"
	"log"
)

type ExpenseCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ExpenseCategory, error)
	GetByID(id uint) (*models.ExpenseCategory, error)
	Create(currentUserID uint, toolType *models.ExpenseCategory) (*models.ExpenseCategory, error)
	Update(id uint, currentUserID uint, toolType *models.ExpenseCategory) (*models.ExpenseCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ExpenseCategoryService struct {
	repo           repositories.ExpenseCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewExpenseCategoryService(repo repositories.ExpenseCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ExpenseCategoryServiceInterface {
	return &ExpenseCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ecs *ExpenseCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ExpenseCategory, error) {
	return ecs.repo.GetAll(sort, sortBy, keyword, status)
}

func (ecs ExpenseCategoryService) GetByID(id uint) (*models.ExpenseCategory, error) {
	return ecs.repo.GetByID(id)
}

func (ecs *ExpenseCategoryService) Create(currentUserID uint, toolType *models.ExpenseCategory) (*models.ExpenseCategory, error) {

	if err := ecs.validateExpenseCategoryForm(toolType, nil); err != nil {
		return nil, err
	}

	createdExpenseCategory, err := ecs.repo.Create(toolType)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ecs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ExpenseCategory{}.TableName(),
			toolType.ID,
			createdExpenseCategory,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdExpenseCategory, nil
}

func (ecs *ExpenseCategoryService) Update(id uint, currentUserID uint, updatedExpenseCategory *models.ExpenseCategory) (*models.ExpenseCategory, error) {
	existingExpenseCategory, err := ecs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if err := ecs.validateExpenseCategoryForm(updatedExpenseCategory, &id); err != nil {
		return nil, err
	}

	result, err := ecs.repo.Update(id, updatedExpenseCategory)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "status"}
	changes, err := utility.CompareFields(existingExpenseCategory, updatedExpenseCategory, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := ecs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ExpenseCategory{}.TableName(),
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

func (ecs *ExpenseCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := ecs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = ecs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := ecs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ExpenseCategory{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (ecs *ExpenseCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ecs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ecs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ExpenseCategory{}.TableName(),
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

func (ecs *ExpenseCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ecs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ecs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ExpenseCategory{}.TableName(),
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

func (ecs *ExpenseCategoryService) validateExpenseCategoryForm(toolTypeForm *models.ExpenseCategory, id *uint) error {

	if id == nil {
		if exists, err := ecs.repo.ExistsByTitle(toolTypeForm.TitleEn, toolTypeForm.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {
		if exists, err := ecs.repo.ExistsByTitleExcludingID(*id, toolTypeForm.TitleEn, toolTypeForm.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	return nil
}
