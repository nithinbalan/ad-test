package services

import (
	"fmt"
	"jethur/grc/audit/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"
)

type AuditFindingCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditFindingCategory, error)
	GetByID(id uint) (*models.AuditFindingCategory, error)
	Update(id uint, currentUserID uint, category *models.AuditFindingCategory) (*models.AuditFindingCategory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type AuditFindingCategoryService struct {
	repo           repositories.AuditFindingCategoryRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewAuditFindingCategoryService(repo repositories.AuditFindingCategoryRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) AuditFindingCategoryServiceInterface {
	return &AuditFindingCategoryService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (afcs *AuditFindingCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditFindingCategory, error) {
	return afcs.repo.GetAll(sort, sortBy, keyword, status)
}

func (afcs *AuditFindingCategoryService) GetByID(id uint) (*models.AuditFindingCategory, error) {
	return afcs.repo.GetByID(id)
}

func (afcs *AuditFindingCategoryService) Update(id uint, currentUserID uint, updatedCategory *models.AuditFindingCategory) (*models.AuditFindingCategory, error) {
	existingCategory, err := afcs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedCategory.TitleEn != existingCategory.TitleEn || updatedCategory.TitleAr != existingCategory.TitleAr {
		if exists, err := afcs.repo.ExistsByTitleExcludingID(id, updatedCategory.TitleEn, updatedCategory.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := afcs.repo.Update(id, updatedCategory)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "Status"}
	changes, err := utility.CompareFields(existingCategory, updatedCategory, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := afcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditFindingCategory{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log update activity: %v\n", err)
		}
	}()

	return result, nil
}

func (afcs *AuditFindingCategoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {
		_, err := afcs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = afcs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := afcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Ensure this action constant is defined appropriately
				models.AuditFindingCategory{}.TableName(),
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
