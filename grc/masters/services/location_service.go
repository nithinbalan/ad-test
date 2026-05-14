package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
	"jethur/utility"
	"log"
	"time"
)

type LocationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Location, error)
	GetByID(id uint) (*models.Location, error)
	Create(currentUserID uint, location *models.Location) (*models.Location, error)
	Update(id uint, currentUserID uint, location *models.Location) (*models.Location, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type LocationService struct {
	repo           repositories.LocationRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewLocationService(repo repositories.LocationRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) LocationServiceInterface {
	return &LocationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ls *LocationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Location, error) {
	return ls.repo.GetAll(sort, sortBy, keyword, status)
}

func (ls LocationService) GetByID(id uint) (*models.Location, error) {
	return ls.repo.GetByID(id)
}

func (ls *LocationService) Create(currentUserID uint, location *models.Location) (*models.Location, error) {

	location.CreatedBYID = currentUserID

	createdLocation, err := ls.repo.Create(location)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ls.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Location{}.TableName(),
			location.ID,
			createdLocation,
			currentUserID,
			"", "", "", nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdLocation, nil
}

func (ls *LocationService) Update(id uint, currentUserID uint, updatedLocation *models.Location) (*models.Location, error) {
	existingLocation, err := ls.repo.GetByID(id)

	if err != nil {
		return nil, err
	}
	now := time.Now()
	updatedLocation.UpdatedAt = &now
	updatedLocation.UpdatedBYID = &currentUserID

	result, err := ls.repo.Update(id, updatedLocation)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Address", "Status"}
	changes, err := utility.CompareFields(existingLocation, updatedLocation, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := ls.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.Location{}.TableName(),
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

func (ls *LocationService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := ls.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = ls.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := ls.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.Location{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (ls *LocationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ls.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ls.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.Location{}.TableName(),
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

func (ls *LocationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ls.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ls.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.Location{}.TableName(),
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
