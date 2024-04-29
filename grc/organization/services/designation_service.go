package services

import (
	"errors"
	"fmt"
	logModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type DesignationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Designation, error)
	GetByID(id uint) (*models.Designation, error)
	Create(currentUserID uint, designation *models.Designation) (*models.Designation, error)
	Update(id uint, currentUserID uint, designation *models.Designation) (*models.Designation, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type DesignationService struct {
	repo           repositories.DesignationRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewDesignationService(repo repositories.DesignationRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) DesignationServiceInterface {
	return &DesignationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ds *DesignationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Designation, error) {
	return ds.repo.GetAll(sort, sortBy, keyword, status)
}

func (ds DesignationService) GetByID(id uint) (*models.Designation, error) {
	designation, err := ds.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return designation, nil
}

func (ds *DesignationService) Create(currentUserID uint, designation *models.Designation) (*models.Designation, error) {
	if exists, err := ds.repo.ExistsByTitle(designation.TitleEn, designation.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdDesignation, err := ds.repo.Create(designation)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ds.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.Designation{}.TableName(),
			createdDesignation.ID,
			createdDesignation,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdDesignation, nil
}

func (ds *DesignationService) Update(id uint, currentUserID uint, updatedDesignation *models.Designation) (*models.Designation, error) {
	existingDesignation, err := ds.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedDesignation.TitleEn != existingDesignation.TitleEn || updatedDesignation.TitleAr != existingDesignation.TitleAr {
		if exists, err := ds.repo.ExistsByTitleExcludingID(id, updatedDesignation.TitleEn, updatedDesignation.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := ds.repo.Update(id, updatedDesignation)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingDesignation, updatedDesignation, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := ds.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.Designation{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (ds *DesignationService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := ds.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = ds.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := ds.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.Designation{}.TableName(),
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

func (ds *DesignationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ds.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ds.activityLogSvc.LogBulkActivities(
				logModels.ActivateAction,
				models.Designation{}.TableName(),
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

	return err
}

func (ds *DesignationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ds.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ds.activityLogSvc.LogBulkActivities(
				logModels.DeactivateAction,
				models.Designation{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}
