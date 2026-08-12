package services

import (
	"errors"
	"fmt"
	"log"

	"jethur/grc/incident/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentSensitivityClassificationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentSensitivityClassification, error)
	GetByID(id uint) (*models.IncidentSensitivityClassification, error)
	Create(currentUserID uint, classification *models.IncidentSensitivityClassification) (*models.IncidentSensitivityClassification, error)
	Update(id uint, currentUserID uint, classification *models.IncidentSensitivityClassification) (*models.IncidentSensitivityClassification, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type IncidentSensitivityClassificationService struct {
	repo           repositories.IncidentSensitivityClassificationRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewIncidentSensitivityClassificationService(repo repositories.IncidentSensitivityClassificationRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) IncidentSensitivityClassificationServiceInterface {
	return &IncidentSensitivityClassificationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (iscs *IncidentSensitivityClassificationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentSensitivityClassification, error) {
	return iscs.repo.GetAll(sort, sortBy, keyword, status)
}

func (iscs IncidentSensitivityClassificationService) GetByID(id uint) (*models.IncidentSensitivityClassification, error) {
	incidentSensitivityClassification, err := iscs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentSensitivityClassification, nil
}

func (iscs *IncidentSensitivityClassificationService) Create(currentUserID uint, classification *models.IncidentSensitivityClassification) (*models.IncidentSensitivityClassification, error) {
	if exists, err := iscs.repo.ExistsByTitle(classification.TitleEn, classification.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdClassification, err := iscs.repo.Create(classification)
	if err != nil {
		return nil, err
	}

	go func() {
		err := iscs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentSensitivityClassification{}.TableName(),
			createdClassification.ID,
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

func (iscs *IncidentSensitivityClassificationService) Update(id uint, currentUserID uint, updatedClassification *models.IncidentSensitivityClassification) (*models.IncidentSensitivityClassification, error) {
	existingClassification, err := iscs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedClassification.TitleEn != existingClassification.TitleEn || updatedClassification.TitleAr != existingClassification.TitleAr {
		if exists, err := iscs.repo.ExistsByTitleExcludingID(id, updatedClassification.TitleEn, updatedClassification.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := iscs.repo.Update(id, updatedClassification)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Type", "Status"}
	changes, err := utility.CompareFields(existingClassification, updatedClassification, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := iscs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentSensitivityClassification{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (iscs *IncidentSensitivityClassificationService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedIncidentSensitivityClassification []map[string]interface{}
	for _, id := range ids {

		incidentSensitivityClassification, err := iscs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = iscs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedIncidentSensitivityClassification = append(deletedIncidentSensitivityClassification, map[string]interface{}{
			"id":       id,
			"title_en": incidentSensitivityClassification.TitleEn,
			"title_ar": incidentSensitivityClassification.TitleAr,
		})

	}

	if len(successfulIds) > 0 {
		go func() {
			err := iscs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IncidentSensitivityClassification{}.TableName(),
				successfulIds,
				deletedIncidentSensitivityClassification,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (iscs *IncidentSensitivityClassificationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := iscs.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := iscs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.IncidentSensitivityClassification{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (iscs *IncidentSensitivityClassificationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := iscs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := iscs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.IncidentSensitivityClassification{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}
