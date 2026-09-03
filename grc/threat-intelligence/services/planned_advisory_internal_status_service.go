package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/threat-intelligence/repositories"
	"jethur/utility"
	"log"
)

type PlannedAdvisoryInternalStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryInternalStatus, error)
	GetByID(id uint) (*models.PlannedAdvisoryInternalStatus, error)
	Create(currentUserID uint, internalStatus *models.PlannedAdvisoryInternalStatus) (*models.PlannedAdvisoryInternalStatus, error)
	Update(id uint, currentUserID uint, internalStatus *models.PlannedAdvisoryInternalStatus) (*models.PlannedAdvisoryInternalStatus, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type PlannedAdvisoryInternalStatusService struct {
	repo           repositories.PlannedAdvisoryInternalStatusRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewPlannedAdvisoryInternalStatusService(repo repositories.PlannedAdvisoryInternalStatusRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) PlannedAdvisoryInternalStatusServiceInterface {
	return &PlannedAdvisoryInternalStatusService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (pais *PlannedAdvisoryInternalStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryInternalStatus, error) {
	return pais.repo.GetAll(sort, sortBy, keyword, status)
}

func (pais PlannedAdvisoryInternalStatusService) GetByID(id uint) (*models.PlannedAdvisoryInternalStatus, error) {
	return pais.repo.GetByID(id)
}

func (pais *PlannedAdvisoryInternalStatusService) Create(currentUserID uint, internalStatus *models.PlannedAdvisoryInternalStatus) (*models.PlannedAdvisoryInternalStatus, error) {

	createdInternalStatus, err := pais.repo.Create(internalStatus)
	if err != nil {
		return nil, err
	}

	go func() {
		err := pais.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PlannedAdvisoryInternalStatus{}.TableName(),
			internalStatus.ID,
			createdInternalStatus,
			currentUserID,
			"", "", "", nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdInternalStatus, nil
}

func (pais *PlannedAdvisoryInternalStatusService) Update(id uint, currentUserID uint, updatedInternalStatus *models.PlannedAdvisoryInternalStatus) (*models.PlannedAdvisoryInternalStatus, error) {
	existingInternalStatus, err := pais.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	result, err := pais.repo.Update(id, updatedInternalStatus)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingInternalStatus, updatedInternalStatus, fieldsToCompare)

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
			models.PlannedAdvisoryInternalStatus{}.TableName(),
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

func (pais *PlannedAdvisoryInternalStatusService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedPlannedAdvisoryInternalStatuses []map[string]interface{}
	for _, id := range ids {

		plannedAdvisoryInternalStatus, err := pais.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = pais.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedPlannedAdvisoryInternalStatuses = append(deletedPlannedAdvisoryInternalStatuses, map[string]interface{}{
			"id":    id,
			"title": plannedAdvisoryInternalStatus.TitleEn,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := pais.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PlannedAdvisoryInternalStatus{}.TableName(),
				successfulIds,
				deletedPlannedAdvisoryInternalStatuses,
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

func (pais *PlannedAdvisoryInternalStatusService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := pais.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pais.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.PlannedAdvisoryInternalStatus{}.TableName(),
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

func (pais *PlannedAdvisoryInternalStatusService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := pais.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pais.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.PlannedAdvisoryInternalStatus{}.TableName(),
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
