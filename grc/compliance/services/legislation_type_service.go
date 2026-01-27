package services

import (
	"fmt"
	"log"

	"jethur/grc/compliance/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
)

type LegislationTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LegislationType, error)
	GetByID(id uint) (*models.LegislationType, error)
	Create(currentUserID uint, body *models.LegislationType) (*models.LegislationType, error)
	Update(id uint, currentUserID uint, body *models.LegislationType) (*models.LegislationType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service struct
type LegislationTypeService struct {
	repo           repositories.LegislationTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

// Constructor
func NewLegislationTypeService(repo repositories.LegislationTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) LegislationTypeServiceInterface {
	return &LegislationTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

// GetAll
func (lts *LegislationTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LegislationType, error) {
	return lts.repo.GetAll(sort, sortBy, keyword, status)
}

// GetByID
func (lts *LegislationTypeService) GetByID(id uint) (*models.LegislationType, error) {
	return lts.repo.GetByID(id)
}

// Create
func (lts *LegislationTypeService) Create(currentUserID uint, body *models.LegislationType) (*models.LegislationType, error) {
	if exists, err := lts.repo.ExistsByTitle(body.TitleEn, body.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdBody, err := lts.repo.Create(body)
	if err != nil {
		return nil, err
	}

	go func() {
		err := lts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.LegislationType{}.TableName(),
			body.ID,
			createdBody,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdBody, nil
}

// Update
func (lts *LegislationTypeService) Update(id uint, currentUserID uint, updatedBody *models.LegislationType) (*models.LegislationType, error) {
	existingBody, err := lts.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedBody.TitleEn != existingBody.TitleEn || updatedBody.TitleAr != existingBody.TitleAr {
		if exists, err := lts.repo.ExistsByTitleExcludingID(id, updatedBody.TitleEn, updatedBody.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := lts.repo.Update(id, updatedBody)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingBody, updatedBody, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := lts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.LegislationType{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()

	return result, nil
}

// Delete
func (lts *LegislationTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedBodies []map[string]interface{}

	for _, id := range ids {
		body, err := lts.repo.GetByID(id)
		if err != nil {
			continue
		}

		err = lts.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedBodies = append(deletedBodies, map[string]interface{}{
			"id":       id,
			"title_en": body.TitleEn,
			"title_ar": body.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := lts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.LegislationType{}.TableName(),
				successfulIds,
				deletedBodies,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

// Activate
func (lts *LegislationTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := lts.repo.Activate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := lts.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.LegislationType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}
	return err
}

// Deactivate
func (lts *LegislationTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := lts.repo.Deactivate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := lts.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.LegislationType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}
	return err
}
