package services

import (
	"fmt"
	"log"

	"jethur/grc/compliance/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
)

type SecrecyClassificationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SecrecyClassification, error)
	GetByID(id uint) (*models.SecrecyClassification, error)
	Create(currentUserID uint, classification *models.SecrecyClassification) (*models.SecrecyClassification, error)
	Update(id uint, currentUserID uint, classification *models.SecrecyClassification) (*models.SecrecyClassification, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Implementation
type SecrecyClassificationService struct {
	repo           repositories.SecrecyClassificationRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewSecrecyClassificationService(
	repo repositories.SecrecyClassificationRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
) SecrecyClassificationServiceInterface {
	return &SecrecyClassificationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (scs *SecrecyClassificationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SecrecyClassification, error) {
	return scs.repo.GetAll(sort, sortBy, keyword, status)
}

func (scs *SecrecyClassificationService) GetByID(id uint) (*models.SecrecyClassification, error) {
	return scs.repo.GetByID(id)
}

func (scs *SecrecyClassificationService) Create(currentUserID uint, classification *models.SecrecyClassification) (*models.SecrecyClassification, error) {
	if exists, err := scs.repo.ExistsByTitle(classification.TitleEn, classification.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdClassification, err := scs.repo.Create(classification)
	if err != nil {
		return nil, err
	}

	go func() {
		err := scs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SecrecyClassification{}.TableName(),
			classification.ID,
			createdClassification,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdClassification, nil
}

func (scs *SecrecyClassificationService) Update(id uint, currentUserID uint, updatedClassification *models.SecrecyClassification) (*models.SecrecyClassification, error) {
	existing, err := scs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedClassification.TitleEn != existing.TitleEn || updatedClassification.TitleAr != existing.TitleAr {
		if exists, err := scs.repo.ExistsByTitleExcludingID(id, updatedClassification.TitleEn, updatedClassification.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := scs.repo.Update(id, updatedClassification)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existing, updatedClassification, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := scs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.SecrecyClassification{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()
	return result, err
}

func (scs *SecrecyClassificationService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	for _, id := range ids {
		item, err := scs.repo.GetByID(id)
		if err != nil {
			continue
		}

		err = scs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedItems = append(deletedItems, map[string]interface{}{
			"id":       id,
			"title_en": item.TitleEn,
			"title_ar": item.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := scs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.SecrecyClassification{}.TableName(),
				successfulIds,
				deletedItems,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (scs *SecrecyClassificationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := scs.repo.Activate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := scs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.SecrecyClassification{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}
	return err
}

func (scs *SecrecyClassificationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := scs.repo.Deactivate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := scs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.SecrecyClassification{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}
	return err
}
