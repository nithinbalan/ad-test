package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/threat-intelligence/repositories"
	"jethur/utility"
	"log"
)

type PlannedAdvisoryImportanceServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryImportance, error)
	GetByID(id uint) (*models.PlannedAdvisoryImportance, error)
	Create(currentUserID uint, importance *models.PlannedAdvisoryImportance) (*models.PlannedAdvisoryImportance, error)
	Update(id uint, currentUserID uint, importance *models.PlannedAdvisoryImportance) (*models.PlannedAdvisoryImportance, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type PlannedAdvisoryImportanceService struct {
	repo           repositories.PlannedAdvisoryImportanceRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewPlannedAdvisoryImportanceService(repo repositories.PlannedAdvisoryImportanceRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) PlannedAdvisoryImportanceServiceInterface {
	return &PlannedAdvisoryImportanceService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (pais *PlannedAdvisoryImportanceService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryImportance, error) {
	return pais.repo.GetAll(sort, sortBy, keyword, status)
}

func (pais PlannedAdvisoryImportanceService) GetByID(id uint) (*models.PlannedAdvisoryImportance, error) {
	return pais.repo.GetByID(id)
}

func (pais *PlannedAdvisoryImportanceService) Create(currentUserID uint, importance *models.PlannedAdvisoryImportance) (*models.PlannedAdvisoryImportance, error) {

	createdImportance, err := pais.repo.Create(importance)
	if err != nil {
		return nil, err
	}

	go func() {
		err := pais.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PlannedAdvisoryImportance{}.TableName(),
			importance.ID,
			createdImportance,
			currentUserID,
			"", "", "", nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdImportance, nil
}

func (pais *PlannedAdvisoryImportanceService) Update(id uint, currentUserID uint, updatedImportance *models.PlannedAdvisoryImportance) (*models.PlannedAdvisoryImportance, error) {
	existingImportance, err := pais.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	result, err := pais.repo.Update(id, updatedImportance)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingImportance, updatedImportance, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := pais.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PlannedAdvisoryImportance{}.TableName(),
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

func (pais *PlannedAdvisoryImportanceService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedPlannedAdvisoryImportances []map[string]interface{}
	for _, id := range ids {

		plannedAdvisoryImportance, err := pais.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = pais.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedPlannedAdvisoryImportances = append(deletedPlannedAdvisoryImportances, map[string]interface{}{
			"id":    id,
			"title": plannedAdvisoryImportance.TitleEn,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := pais.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PlannedAdvisoryImportance{}.TableName(),
				successfulIds,
				deletedPlannedAdvisoryImportances,
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

func (pais *PlannedAdvisoryImportanceService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := pais.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pais.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.PlannedAdvisoryImportance{}.TableName(),
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

func (pais *PlannedAdvisoryImportanceService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := pais.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pais.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.PlannedAdvisoryImportance{}.TableName(),
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
