package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/threat-intelligence/repositories"
	"jethur/utility"
	"log"
)

type PlannedAdvisorySecurityClassificationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisorySecurityClassification, error)
	GetByID(id uint) (*models.PlannedAdvisorySecurityClassification, error)
	Create(currentUserID uint, classification *models.PlannedAdvisorySecurityClassification) (*models.PlannedAdvisorySecurityClassification, error)
	Update(id uint, currentUserID uint, classification *models.PlannedAdvisorySecurityClassification) (*models.PlannedAdvisorySecurityClassification, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type PlannedAdvisorySecurityClassificationService struct {
	repo           repositories.PlannedAdvisorySecurityClassificationRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewPlannedAdvisorySecurityClassificationService(repo repositories.PlannedAdvisorySecurityClassificationRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) PlannedAdvisorySecurityClassificationServiceInterface {
	return &PlannedAdvisorySecurityClassificationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (pasc *PlannedAdvisorySecurityClassificationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisorySecurityClassification, error) {
	return pasc.repo.GetAll(sort, sortBy, keyword, status)
}

func (pasc PlannedAdvisorySecurityClassificationService) GetByID(id uint) (*models.PlannedAdvisorySecurityClassification, error) {
	return pasc.repo.GetByID(id)
}

func (pasc *PlannedAdvisorySecurityClassificationService) Create(currentUserID uint, classification *models.PlannedAdvisorySecurityClassification) (*models.PlannedAdvisorySecurityClassification, error) {

	createdClassification, err := pasc.repo.Create(classification)
	if err != nil {
		return nil, err
	}

	go func() {
		err := pasc.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PlannedAdvisorySecurityClassification{}.TableName(),
			classification.ID,
			createdClassification,
			currentUserID,
			"", "", "", nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdClassification, nil
}

func (pasc *PlannedAdvisorySecurityClassificationService) Update(id uint, currentUserID uint, updatedClassification *models.PlannedAdvisorySecurityClassification) (*models.PlannedAdvisorySecurityClassification, error) {
	existingClassification, err := pasc.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	result, err := pasc.repo.Update(id, updatedClassification)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingClassification, updatedClassification, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := pasc.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PlannedAdvisorySecurityClassification{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()
	return result, err
}

func (pasc *PlannedAdvisorySecurityClassificationService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedPlannedAdvisorySecurityClassifications []map[string]interface{}
	for _, id := range ids {

		plannedAdvisorySecurityClassifications, err := pasc.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = pasc.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedPlannedAdvisorySecurityClassifications = append(deletedPlannedAdvisorySecurityClassifications, map[string]interface{}{
			"id":    id,
			"title": plannedAdvisorySecurityClassifications.TitleEn,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := pasc.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PlannedAdvisorySecurityClassification{}.TableName(),
				successfulIds,
				deletedPlannedAdvisorySecurityClassifications,
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

func (pasc *PlannedAdvisorySecurityClassificationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := pasc.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pasc.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.PlannedAdvisorySecurityClassification{}.TableName(),
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

func (pasc *PlannedAdvisorySecurityClassificationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := pasc.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pasc.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.PlannedAdvisorySecurityClassification{}.TableName(),
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
