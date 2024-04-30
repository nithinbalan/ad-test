package services

import (
	"fmt"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"
)

type ComplianceRequirementCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceRequirementCategory, error)
	GetByID(id uint) (*models.ComplianceRequirementCategory, error)
	Create(currentUserID uint, complianceRequirementCategory *models.ComplianceRequirementCategory) (*models.ComplianceRequirementCategory, error)
	Update(id uint, currentUserID uint, complianceRequirementCategory *models.ComplianceRequirementCategory) (*models.ComplianceRequirementCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ComplianceRequirementCategoryService struct {
	repo           repositories.ComplianceRequirementCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewComplianceRequirementCategoryService(repo repositories.ComplianceRequirementCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ComplianceRequirementCategoryServiceInterface {
	return &ComplianceRequirementCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (crcs *ComplianceRequirementCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceRequirementCategory, error) {
	return crcs.repo.GetAll(sort, sortBy, keyword, status)
}

func (crcs ComplianceRequirementCategoryService) GetByID(id uint) (*models.ComplianceRequirementCategory, error) {
	return crcs.repo.GetByID(id)
}

func (crcs *ComplianceRequirementCategoryService) Create(currentUserID uint, complianceRequirementCategory *models.ComplianceRequirementCategory) (*models.ComplianceRequirementCategory, error) {

	if exists, err := crcs.repo.ExistsByTitle(complianceRequirementCategory.TitleEn, complianceRequirementCategory.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdRequirementCategory, err := crcs.repo.Create(complianceRequirementCategory)
	if err != nil {
		return nil, err
	}

	go func() {
		err := crcs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.ComplianceRequirementCategory{}.TableName(),
			createdRequirementCategory.ID,
			createdRequirementCategory,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdRequirementCategory, nil
}

func (crcs *ComplianceRequirementCategoryService) Update(id uint, currentUserID uint, updatedRequirementCategory *models.ComplianceRequirementCategory) (*models.ComplianceRequirementCategory, error) {
	existingRequirementCategory, err := crcs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedRequirementCategory.TitleEn != existingRequirementCategory.TitleEn || updatedRequirementCategory.TitleAr != existingRequirementCategory.TitleAr {
		if exists, err := crcs.repo.ExistsByTitleExcludingID(id, updatedRequirementCategory.TitleEn, updatedRequirementCategory.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := crcs.repo.Update(id, updatedRequirementCategory)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingRequirementCategory, updatedRequirementCategory, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := crcs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.ComplianceRequirementCategory{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk activity: %v\n", err)
		}
	}()
	return result, err
}

func (crcs *ComplianceRequirementCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := crcs.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = crcs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := crcs.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.ComplianceRequirementCategory{}.TableName(),
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

func (crcs *ComplianceRequirementCategoryService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := crcs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := crcs.activityLogSvc.LogBulkActivities(
				logModels.ActivateAction,
				models.ComplianceRequirementCategory{}.TableName(),
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

func (crcs *ComplianceRequirementCategoryService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := crcs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := crcs.activityLogSvc.LogBulkActivities(
				logModels.DeactivateAction,
				models.ComplianceRequirementCategory{}.TableName(),
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
