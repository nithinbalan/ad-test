package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/threat-intelligence/repositories"
	"jethur/utility"
	"log"
)

type PlannedAdvisoryTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryType, error)
	GetByID(id uint) (*models.PlannedAdvisoryType, error)
	Create(currentUserID uint, advisoryType *models.PlannedAdvisoryType) (*models.PlannedAdvisoryType, error)
	Update(id uint, currentUserID uint, advisoryType *models.PlannedAdvisoryType) (*models.PlannedAdvisoryType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type PlannedAdvisoryTypeService struct {
	repo           repositories.PlannedAdvisoryTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewPlannedAdvisoryTypeService(repo repositories.PlannedAdvisoryTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) PlannedAdvisoryTypeServiceInterface {
	return &PlannedAdvisoryTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (pats *PlannedAdvisoryTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryType, error) {
	return pats.repo.GetAll(sort, sortBy, keyword, status)
}

func (pats PlannedAdvisoryTypeService) GetByID(id uint) (*models.PlannedAdvisoryType, error) {
	return pats.repo.GetByID(id)
}

func (pats *PlannedAdvisoryTypeService) Create(currentUserID uint, advisoryType *models.PlannedAdvisoryType) (*models.PlannedAdvisoryType, error) {

	createdAdvisoryType, err := pats.repo.Create(advisoryType)
	if err != nil {
		return nil, err
	}

	go func() {
		err := pats.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PlannedAdvisoryType{}.TableName(),
			advisoryType.ID,
			createdAdvisoryType,
			currentUserID,
			"", "", "", nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdAdvisoryType, nil
}

func (pats *PlannedAdvisoryTypeService) Update(id uint, currentUserID uint, updatedAdvisoryType *models.PlannedAdvisoryType) (*models.PlannedAdvisoryType, error) {
	existingAdvisoryType, err := pats.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	result, err := pats.repo.Update(id, updatedAdvisoryType)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingAdvisoryType, updatedAdvisoryType, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := pats.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PlannedAdvisoryType{}.TableName(),
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

func (pats *PlannedAdvisoryTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedPlannedAdvisoryTypes []map[string]interface{}
	for _, id := range ids {

		plannedAdvisoryType, err := pats.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = pats.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedPlannedAdvisoryTypes = append(deletedPlannedAdvisoryTypes, map[string]interface{}{
			"id":    id,
			"title": plannedAdvisoryType.TitleEn,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := pats.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PlannedAdvisoryType{}.TableName(),
				successfulIds,
				deletedPlannedAdvisoryTypes,
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

func (pats *PlannedAdvisoryTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := pats.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pats.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.PlannedAdvisoryType{}.TableName(),
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

func (pats *PlannedAdvisoryTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := pats.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pats.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.PlannedAdvisoryType{}.TableName(),
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
