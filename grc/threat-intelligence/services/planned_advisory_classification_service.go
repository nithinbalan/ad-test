package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/threat-intelligence/repositories"
	"jethur/utility"
	"log"
)

type PlannedAdvisoryClassificationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryClassification, error)
	GetByID(id uint) (*models.PlannedAdvisoryClassification, error)
	Create(currentUserID uint, classification *models.PlannedAdvisoryClassification) (*models.PlannedAdvisoryClassification, error)
	Update(id uint, currentUserID uint, classification *models.PlannedAdvisoryClassification) (*models.PlannedAdvisoryClassification, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type PlannedAdvisoryClassificationService struct {
	repo           repositories.PlannedAdvisoryClassificationRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewPlannedAdvisoryClassificationService(repo repositories.PlannedAdvisoryClassificationRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) PlannedAdvisoryClassificationServiceInterface {
	return &PlannedAdvisoryClassificationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (pacs *PlannedAdvisoryClassificationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryClassification, error) {
	return pacs.repo.GetAll(sort, sortBy, keyword, status)
}

func (pacs PlannedAdvisoryClassificationService) GetByID(id uint) (*models.PlannedAdvisoryClassification, error) {
	return pacs.repo.GetByID(id)
}

func (pacs *PlannedAdvisoryClassificationService) Create(currentUserID uint, classification *models.PlannedAdvisoryClassification) (*models.PlannedAdvisoryClassification, error) {

	createdClassification, err := pacs.repo.Create(classification)
	if err != nil {
		return nil, err
	}

	go func() {
		err := pacs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PlannedAdvisoryClassification{}.TableName(),
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

func (pacs *PlannedAdvisoryClassificationService) Update(id uint, currentUserID uint, updatedClassification *models.PlannedAdvisoryClassification) (*models.PlannedAdvisoryClassification, error) {
	existingClassification, err := pacs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	result, err := pacs.repo.Update(id, updatedClassification)

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
		err := pacs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PlannedAdvisoryClassification{}.TableName(),
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

func (pacs *PlannedAdvisoryClassificationService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedPlannedAdvisoryClassifications []map[string]interface{}
	for _, id := range ids {

		plannedAdvisoryClassifications, err := pacs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = pacs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedPlannedAdvisoryClassifications = append(deletedPlannedAdvisoryClassifications, map[string]interface{}{
			"id":    id,
			"title": plannedAdvisoryClassifications.TitleEn,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := pacs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PlannedAdvisoryClassification{}.TableName(),
				successfulIds,
				deletedPlannedAdvisoryClassifications,
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

func (pacs *PlannedAdvisoryClassificationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := pacs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pacs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.PlannedAdvisoryClassification{}.TableName(),
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

func (pacs *PlannedAdvisoryClassificationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := pacs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pacs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.PlannedAdvisoryClassification{}.TableName(),
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
