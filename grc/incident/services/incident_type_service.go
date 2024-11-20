package services

import (
	"errors"
	"fmt"
	"jethur/grc/incident/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type IncidentTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentType, error)
	GetByID(id uint) (*models.IncidentType, error)
	Create(currentUserID uint, incidentType *models.IncidentType) (*models.IncidentType, error)
	Update(id uint, currentUserID uint, incidentType *models.IncidentType) (*models.IncidentType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type IncidentTypeService struct {
	repo           repositories.IncidentTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewIncidentTypeService(repo repositories.IncidentTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) IncidentTypeServiceInterface {
	return &IncidentTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (its *IncidentTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentType, error) {
	return its.repo.GetAll(sort, sortBy, keyword, status)
}

func (its IncidentTypeService) GetByID(id uint) (*models.IncidentType, error) {
	incidentType, err := its.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentType, nil

}

func (its *IncidentTypeService) Create(currentUserID uint, incidentType *models.IncidentType) (*models.IncidentType, error) {

	if exists, err := its.repo.ExistsByTitle(incidentType.TitleEn, incidentType.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdType, err := its.repo.Create(incidentType)
	if err != nil {
		return nil, err
	}

	go func() {
		err := its.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentType{}.TableName(),
			createdType.ID,
			createdType,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdType, nil
}

func (its *IncidentTypeService) Update(id uint, currentUserID uint, updatedType *models.IncidentType) (*models.IncidentType, error) {
	existingType, err := its.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedType.TitleEn != existingType.TitleEn || updatedType.TitleAr != existingType.TitleAr {
		if exists, err := its.repo.ExistsByTitleExcludingID(id, updatedType.TitleEn, updatedType.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := its.repo.Update(id, updatedType)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingType, updatedType, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := its.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentType{}.TableName(),
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

func (its *IncidentTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedIncidentType []map[string]interface{}
	for _, id := range ids {

		incidentType, err := its.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = its.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedIncidentType = append(deletedIncidentType, map[string]interface{}{
			"id":       id,
			"title_en": incidentType.TitleEn,
			"title_ar": incidentType.TitleAr,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := its.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.IncidentType{}.TableName(),
				successfulIds,
				deletedIncidentType,
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

func (its *IncidentTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := its.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := its.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.IncidentType{}.TableName(),
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

func (its *IncidentTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := its.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := its.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.IncidentType{}.TableName(),
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
