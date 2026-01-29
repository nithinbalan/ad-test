package services

import (
	"fmt"
	"log"

	"jethur/grc/compliance/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
)

type LegislativeBodyServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceLegislativeBody, error)
	GetByID(id uint) (*models.ComplianceLegislativeBody, error)
	Create(currentUserID uint, body *models.ComplianceLegislativeBody) (*models.ComplianceLegislativeBody, error)
	Update(id uint, currentUserID uint, body *models.ComplianceLegislativeBody) (*models.ComplianceLegislativeBody, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service struct
type LegislativeBodyService struct {
	repo           repositories.LegislativeBodyRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

// Constructor
func NewLegislativeBodyService(repo repositories.LegislativeBodyRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) LegislativeBodyServiceInterface {
	return &LegislativeBodyService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

// GetAll
func (lbs *LegislativeBodyService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceLegislativeBody, error) {
	return lbs.repo.GetAll(sort, sortBy, keyword, status)
}

// GetByID
func (lbs *LegislativeBodyService) GetByID(id uint) (*models.ComplianceLegislativeBody, error) {
	return lbs.repo.GetByID(id)
}

// Create
func (lbs *LegislativeBodyService) Create(currentUserID uint, body *models.ComplianceLegislativeBody) (*models.ComplianceLegislativeBody, error) {
	if exists, err := lbs.repo.ExistsByTitle(body.TitleEn, body.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdBody, err := lbs.repo.Create(body)
	if err != nil {
		return nil, err
	}

	go func() {
		err := lbs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceLegislativeBody{}.TableName(),
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
func (lbs *LegislativeBodyService) Update(id uint, currentUserID uint, updatedBody *models.ComplianceLegislativeBody) (*models.ComplianceLegislativeBody, error) {
	existingBody, err := lbs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedBody.TitleEn != existingBody.TitleEn || updatedBody.TitleAr != existingBody.TitleAr {
		if exists, err := lbs.repo.ExistsByTitleExcludingID(id, updatedBody.TitleEn, updatedBody.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := lbs.repo.Update(id, updatedBody)
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
		err := lbs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ComplianceLegislativeBody{}.TableName(),
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
func (lbs *LegislativeBodyService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedBodies []map[string]interface{}

	for _, id := range ids {
		body, err := lbs.repo.GetByID(id)
		if err != nil {
			continue
		}

		err = lbs.repo.Delete(id)
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
			err := lbs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ComplianceLegislativeBody{}.TableName(),
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
func (lbs *LegislativeBodyService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := lbs.repo.Activate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := lbs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ComplianceLegislativeBody{}.TableName(),
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
func (lbs *LegislativeBodyService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := lbs.repo.Deactivate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := lbs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ComplianceLegislativeBody{}.TableName(),
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
