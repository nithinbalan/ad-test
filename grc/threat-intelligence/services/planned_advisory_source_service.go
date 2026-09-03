package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/threat-intelligence/repositories"
	"jethur/utility"
	"log"
)

type PlannedAdvisorySourceServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisorySource, error)
	GetByID(id uint) (*models.PlannedAdvisorySource, error)
	Create(currentUserID uint, source *models.PlannedAdvisorySource) (*models.PlannedAdvisorySource, error)
	Update(id uint, currentUserID uint, source *models.PlannedAdvisorySource) (*models.PlannedAdvisorySource, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type PlannedAdvisorySourceService struct {
	repo           repositories.PlannedAdvisorySourceRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewPlannedAdvisorySourceService(repo repositories.PlannedAdvisorySourceRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) PlannedAdvisorySourceServiceInterface {
	return &PlannedAdvisorySourceService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (pass *PlannedAdvisorySourceService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisorySource, error) {
	return pass.repo.GetAll(sort, sortBy, keyword, status)
}

func (pass PlannedAdvisorySourceService) GetByID(id uint) (*models.PlannedAdvisorySource, error) {
	return pass.repo.GetByID(id)
}

func (pass *PlannedAdvisorySourceService) Create(currentUserID uint, source *models.PlannedAdvisorySource) (*models.PlannedAdvisorySource, error) {

	createdSource, err := pass.repo.Create(source)
	if err != nil {
		return nil, err
	}

	go func() {
		err := pass.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PlannedAdvisorySource{}.TableName(),
			source.ID,
			createdSource,
			currentUserID,
			"", "", "", nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdSource, nil
}

func (pass *PlannedAdvisorySourceService) Update(id uint, currentUserID uint, updatedSource *models.PlannedAdvisorySource) (*models.PlannedAdvisorySource, error) {
	existingSource, err := pass.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	result, err := pass.repo.Update(id, updatedSource)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingSource, updatedSource, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := pass.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PlannedAdvisorySource{}.TableName(),
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

func (pass *PlannedAdvisorySourceService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedPlannedAdvisorySources []map[string]interface{}
	for _, id := range ids {

		plannedAdvisorySource, err := pass.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = pass.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedPlannedAdvisorySources = append(deletedPlannedAdvisorySources, map[string]interface{}{
			"id":    id,
			"title": plannedAdvisorySource.TitleEn,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := pass.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PlannedAdvisorySource{}.TableName(),
				successfulIds,
				deletedPlannedAdvisorySources,
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

func (pass *PlannedAdvisorySourceService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := pass.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pass.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.PlannedAdvisorySource{}.TableName(),
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

func (pass *PlannedAdvisorySourceService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := pass.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pass.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.PlannedAdvisorySource{}.TableName(),
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
